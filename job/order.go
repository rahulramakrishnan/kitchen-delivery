package job

import (
	"log"
	"time"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/service"
	"github.com/pkg/errors"
)

// OrderJob is order job interface.
type OrderJob interface {
	HandleIncomingOrders()
	RemoveExpiredOrders()
}

type orderJob struct {
	cfg      config.AppConfig
	services service.Services
	queues   entity.Queues
}

// NewOrderJob returns a new order job.
func NewOrderJob(cfg config.AppConfig, services service.Services, queues entity.Queues) OrderJob {
	return &orderJob{
		cfg:      cfg,
		services: services,
		queues:   queues,
	}
}

// HandleOrders pulls orders off of order queue and shelf queue.
// TODO: Store job status in a table so we can track when what was run.
func (o *orderJob) HandleIncomingOrders() {
	// Spawn thread to handle incoming orders.
	go o.handleIncomingOrders()
	// Spawn thread to handle placing orders on the right shelves.
	go o.handlePlacingOrdersOnShelves()
}

func (o *orderJob) handleIncomingOrders() {
	// Pull order of queue and spawn a go-routine to handle order.
	for order := range o.queues.OrderQueue {
		// TODO: Use a worker pool.
		go o.handleOrder(*order)
	}
}

func (o *orderJob) handlePlacingOrdersOnShelves() {
	// Pull order off of shelf queue and spawn a go-routine to retry placing order
	// on the right shelf.
	for order := range o.queues.ShelfQueue {
		// TODO: Use a worker pool.
		go o.handleOrder(*order)
	}
}

// handleOrder pulls an order off of an order queue and stores it.
func (o *orderJob) handleOrder(order entity.Order) {
	log.Printf("worker | pulled order off of queue - order: %s", order.String())

	err := o.services.Order.PlaceOrderOnShelf(order)
	if err != nil {
		log.Printf("worker | failed to place order on shelf so putting it on the shelf queue- err: %s", err)
		// Just because we can't place an order on a shelf right now
		// doesn't mean we should fail the request and throw the food away.
		// We can keep the order on the stove and the cook can pull it off and put it on
		// a shelf when there is space.
		o.queues.ShelfQueue <- &order
		return
	}

	log.Printf("worker | placed order on correct shelf - %s", order.String())
}

func (o *orderJob) RemoveExpiredOrders() {
	// Every 5 seconds, find all food that is wasted and status is "ready_for_pickup".
	// and update it's status to "wasted".
	for {
		time.Sleep(5 * time.Second)

		expiredOrdersOnShelf, err := o.services.Order.GetExpiredOrdersOnShelf()
		if err != nil {
			continue
		}

		for _, shelfOrder := range expiredOrdersOnShelf {
			err := o.removeExpiredOrder(*shelfOrder)
			if err != nil {
				log.Printf("Failed to mark order as waste %s - err: %s", shelfOrder.String(), err.Error())
				continue
			}

		}
	}
}

func (o *orderJob) removeExpiredOrder(shelfOrder entity.ShelfOrder) error {
	err := o.services.Order.MarkOrderAsWasted(shelfOrder)
	if err != nil {
		switch errors.Cause(err) {
		case exception.ErrVersionInvalid:
			// This is not an exceptional case.
			// We do not want to retry the operation b/c the order status has
			// changed to being picked up.
			return nil
		case exception.ErrDatabase:
			// We should retry the operation if it's a database error b/c
			// if we don't mark the food as waste then a customer might get
			// food that might get them sick and CSS might get sued.
			return err
		}

		return err
	}

	log.Printf("Marked order on shelf wasted %s", shelfOrder.String())
	return nil
}
