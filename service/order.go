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
	repository repository.OrderRepository
}

// NewOrderService returns a new user service.
// switch to userRepositories
func NewOrderService(repository repository.OrderRepository) OrderService {
	return &orderService{
		repository: repository,
	}
}

// Create stores an order in the orders table.
func (o *orderService) CreateOrder(order entity.Order) error {
	err := o.repository.CreateOrder(order)
	if err != nil {
		return errors.Wrapf(err, "failed to create order, order: %+v", order)
	}

	return nil
}

// CreateOrderLog creates an order history event.
func (o *orderService) CreateOrderLog(orderLog entity.OrderLog) error {
	err := o.repository.CreateOrderLog(orderLog)
	if err != nil {
		return errors.Wrapf(err, "failed to create order history %+v", orderLog)
	}

	return nil
}
