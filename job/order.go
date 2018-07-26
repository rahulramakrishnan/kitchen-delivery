package job

import (
	"log"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service"
)

// OrderJob is order job interface.
type OrderJob interface {
	HandleOrders()
}

type orderJob struct {
	cfg      config.AppConfig
	services service.Services
	queue    <-chan *entity.Order // only pulls off of Order queue
}

// NewOrderJob returns a new order job.
func NewOrderJob(cfg config.AppConfig, services service.Services, queue <-chan *entity.Order) OrderJob {
	return &orderJob{
		cfg:      cfg,
		services: services,
		queue:    queue,
	}
}

// HandleOrders pulls orders off of order queue.
// TODO: Store job status in a table so we can track when what was run.
func (o *orderJob) HandleOrders() {
	// Pull order of queue and spawn a go-routine to handle order.
	for order := range o.queue {
		// TODO: Use a worker pool.
		go o.handleOrder(*order)
	}
}

// handleOrder pulls an order off of an order queue and stores it.
func (o *orderJob) handleOrder(order entity.Order) {
	log.Printf("pulled order off of queue order: %+v", order)

	err := o.services.Order.Create(order)
	if err != nil {
		// TODO: store job status in a table.
		log.Printf("failed to create an order err: %+v", err)
	}

	log.Printf("successfull stored order %+v", order)
}
