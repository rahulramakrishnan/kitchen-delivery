package job

import (
	"log"
	"time"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/lib"
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
	queues   *entity.Queues
}

// NewOrderJob returns a new order job.
func NewOrderJob(cfg config.AppConfig, services service.Services, queues *entity.Queues) OrderJob {
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
		go o.handleIncomingOrder(i)
	}
}

func (o *orderJob) handleIncomingOrder(workerNum int) {
	// Poll redis queue until we stop service.
	for {
		// Sleep for 1s before polling Redis queue.
		time.Sleep(1)

		// Fetch redis connection from redis pool.
		redisConn := o.queues.Order.Pool.Get() // Fetch redis connection from redis pool.
		switch redisConn.Err() {
		case nil:
			orderUUIDObj, err := redisConn.Do("RPOP", o.queues.Order.Name)
			redisConn.Close()
			if err != nil {
				log.Printf("worker %d failed to fetch order uuid from order queue - err: %+v", workerNum, err)
				continue
			}
			if orderUUIDObj == nil {
				// Nothing in the queue to pull and work on.
				continue
			}

			orderUUIDStr := string(orderUUIDObj.([]uint8))
			orderUUID, err := guuid.FromString(orderUUIDStr)
			if err != nil {
				log.Printf("worker %d order uuid got corrupted - err: %s", workerNum, err.Error())
				continue
			}

			log.Printf("worker %d pulled orderUUID %s from order queue", workerNum, orderUUID.String())

			o.placeOrderOnShelf(orderUUID)
		default:
			redisConn.Close()
			err := redisConn.Err()
			if err != nil {
				log.Printf("worker %d failed to connect to redis queue, err: %+v", workerNum, err)
			}
		}
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

	// Fetch and print the contents of the shelves as a best effort.
	allShelfOrders, err := o.services.Order.GetAllOrdersOnShelves()
	shelves := lib.StringifyShelves(allShelfOrders)
	log.Printf("\n ----- Shelf Contents ------ \n%s", shelves)
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
