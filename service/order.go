package service

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service/repository"

	"github.com/pkg/errors"
)

// OrderService is order serivce interface.
type OrderService interface {
	CreateOrder(order entity.Order) error
	CreateOrderLog(orderLog entity.OrderLog) error
}

type orderService struct {
	orderRepository repository.OrderRepository
	shelfRepository repository.ShelfRepository
}

// NewOrderService returns a new user service.
// switch to userRepositories
func NewOrderService(orderRepository repository.OrderRepository, shelfRepository repository.ShelfRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
		shelfRepository: shelfRepository,
	}
}

// Create stores an order in the orders table.
func (o *orderService) CreateOrder(order entity.Order) error {
	err := o.orderRepository.CreateOrder(order)
	if err != nil {
		return errors.Wrapf(err, "failed to create order, order: %+v", order)
	}

	err = o.shelfRepository.AddOrder(order)
	if err != nil {
		return errors.Wrapf(err, "failed to add order, order: %+v", order)
	}

	return nil
}

// CreateOrderLog creates an order history event.
func (o *orderService) CreateOrderLog(orderLog entity.OrderLog) error {
	err := o.orderRepository.CreateOrderLog(orderLog)
	if err != nil {
		return errors.Wrapf(err, "failed to create order history %+v", orderLog)
	}

	return nil
}
