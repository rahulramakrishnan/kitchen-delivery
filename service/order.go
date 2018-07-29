package service

import (
	"math"
	"time"

	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/service/repository"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// OrderService is order serivce interface.
type OrderService interface {
	CreateOrder(order entity.Order) error
}

type orderService struct {
	orderRepository      repository.OrderRepository
	shelfOrderRepository repository.ShelfOrderRepository
	shelfSpace           map[entity.ShelfType]int
}

// NewOrderService returns a new user service.
// switch to userRepositories
func NewOrderService(orderRepository repository.OrderRepository, shelfOrderRepository repository.ShelfOrderRepository) OrderService {
	// TODO: move to configuration.
	// Holds how many items each type of shelf can hold
	// at any given time.
	shelfSpace := map[entity.ShelfType]int{
		entity.HotShelf:      15,
		entity.ColdShelf:     15,
		entity.FrozenShelf:   15,
		entity.OverflowShelf: 20,
	}
	return &orderService{
		orderRepository:      orderRepository,
		shelfOrderRepository: shelfOrderRepository,
		shelfSpace:           shelfSpace,
	}
}

// Create stores an order in the orders table.
func (o *orderService) CreateOrder(order entity.Order) error {
	err := o.orderRepository.CreateOrder(order)
	if err != nil {
		return errors.Wrapf(err, "failed to create order, order: %+v", order)
	}

	// Check count of orders in "hot" w/ status of ready for pick up.
	shelfType := order.GetShelfType()
	numOfOrders, err := o.shelfOrderRepository.CountOrdersOnShelf(shelfType)
	if err != nil {
		return errors.Wrapf(err, "failed to create order %+v", order)
	}

	// If num of orders on shelf is greater than allowed limit
	// then we return a retriable error to the consumer.
	if numOfOrders > o.shelfSpace[shelfType] {
		return errors.Wrap(
			exception.ErrServiceUnavailable, "all shelves are filled, please retry again later")
	}

	// Otherwise we:
	//    a. calculate ttl and expiration date
	//    b. form a shelf order w/ version 0
	//    c. store order in mysql table shelf_orders

	ttl := o.getTTL(order)
	expirationDate := time.Now().Add(time.Second * time.Duration(ttl)).UTC()

	shelfOrder := entity.ShelfOrder{
		UUID:        guuid.NewV4(),
		OrderUUID:   order.UUID,
		ShelfType:   order.GetShelfType(),
		OrderStatus: entity.OrderStatusReadyForPickup,
		Version:     0,
		ExpiresAt:   expirationDate,
	}

	err = o.shelfOrderRepository.AddOrderToShelf(shelfOrder)
	if err != nil {
		return errors.Wrapf(err, "failed to add order, order: %+v", order)
	}

	return nil
}

func (o *orderService) getTTL(order entity.Order) int {
	// Calculate time to live in seconds based on formula.
	// Remember an order is waste after the "value" becomes zero.
	// This leads the formula to be reduced to:
	// => orderAge = shelfLife / (1 + decayRate)
	// We're given shelfLife and decayRate so we can solve for
	// how old an order can get before we consider it as waste.
	expirationTime := float64(order.ShelfLife) / (1.0 + order.DecayRate)
	ttl := int(math.Floor(expirationTime))
	return ttl
}
