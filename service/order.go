package service

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service/repository"

	"github.com/pkg/errors"
)

// OrderService is order serivce interface.
type OrderService interface {
	Create(order entity.Order) error
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

// Create stores a user in a user table.
func (o *orderService) Create(order entity.Order) error {
	err := o.repository.Create(order)
	if err != nil {
		return errors.Wrapf(err, "failed to create order %+v", order)
	}

	return nil
}
