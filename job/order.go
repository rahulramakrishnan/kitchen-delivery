package job

import (
	"log"
	"time"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/service"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
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
func (o *orderJob) HandleIncomingOrders() {
	// Pull order off of shelf queue and spawn a go-routine to retry placing order
	// on the right shelf.
	for i := 0; i < o.cfg.WorkerPool.MaxWorkers; i++ {
		go o.handleIncomingOrder()
	}
}

func (o *orderJob) handleIncomingOrder() {
	// Pull order off of shelf queue and spawn a go-routine to retry placing order
	// on the right shelf.
	orderQueue := o.queues.Order

	for {
		// Sleep for 1s before polling Redis queue.
		time.Sleep(1)

		orderUUIDObj, err := orderQueue.Conn.Do("RPOP", orderQueue.Name)
		if err != nil {
			log.Printf("failed to fetch order uuid from order queue - err: %+v", err)
			continue
		}
		if orderUUIDObj == nil {
			// Nothing in the queue to pull and work on.
			continue
		}

		orderUUIDStr := string(orderUUIDObj.([]uint8))
		orderUUID, err := guuid.FromString(orderUUIDStr)
		if err != nil {
			log.Printf("order uuid got correupted - err: %s", err.Error())
			continue
		}

		log.Printf("pulled orderUUID %s from order queue", orderUUID)

		o.placeOrderOnShelf(orderUUID)
	}
}

// placeOrderOnShelf pulls an order off of an order queue and stores it.
func (o *orderJob) placeOrderOnShelf(orderUUID guuid.UUID) {
	order, err := o.services.Order.GetOrder(orderUUID)
	if err != nil {
		log.Printf("worker | failed to fetch order - orderUUID: %s", orderUUID.String())
		return
	}

	err = o.services.Order.PlaceOrderOnShelf(*order)
	if err != nil {
		if errors.Cause(err) == exception.ErrFullShelf {
			log.Printf("worker | kitchen is over capacity - dropping order: %s", order.String())
			return
		}

		log.Printf("worker | failed to place order on shelf, err: %s", err.Error())
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
