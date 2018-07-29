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
	PlaceOrderOnShelf(order entity.Order) error
	PickupOrder() (*entity.Order, error)
	GetExpiredOrdersOnShelf() ([]*entity.ShelfOrder, error)
	MarkOrderAsWasted(entity.ShelfOrder) error
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
	// Store an immutable record of incoming orders.
	err := o.orderRepository.CreateOrder(order)
	if err != nil {
		return errors.Wrapf(err, "failed to create order, order: %+v", order)
	}

	return nil
}

// PlaceOrderOnShelf places an order on the shelf.
func (o *orderService) PlaceOrderOnShelf(order entity.Order) error {
	// Check count of orders in "hot" w/ status of ready for pick up.
	correspondingShelfType := order.GetShelfType()
	numOfOrders, err := o.shelfOrderRepository.CountOrdersOnShelf(correspondingShelfType)
	if err != nil {
		return errors.Wrapf(err, "failed to create order %+v", order)
	}

	// We try to place the order on the right shelf but we first
	// have to check the current num of items on both the corresponding shelf
	// and the overflow shelf.
	isCorrespondingShelfFull := numOfOrders > o.shelfSpace[correspondingShelfType]
	isOverflowShelfFull := numOfOrders > o.shelfSpace[entity.OverflowShelf]

	// First, if both the corresponding shelf and the overflow shelf are full
	// we throw a retriable service full shelf exception so a caller can handle it explictly.
	areAllShelvesFull := isCorrespondingShelfFull && isOverflowShelfFull
	if areAllShelvesFull {
		return errors.Wrap(
			exception.ErrFullShelf, "all shelves are filled, please retry again later")
	}

	var shelfType entity.ShelfType

	// Second, we check if we can place the order on the corresponding shelf.
	if !isCorrespondingShelfFull {
		shelfType = order.GetShelfType()
	} else { // If we cannot, we know we can place it on the overflow shelf is not full.
		shelfType = entity.OverflowShelf
	}

	// Next:
	//    a. Calculate ttl and expiration date
	//    b. Form a shelf order w/ version 0
	//    c. Place shelf order on a queue that the kitchen pulls off of.

	ttl := o.getTTL(order)
	now := time.Now()
	expirationDate := now.Add(time.Second * time.Duration(ttl))

	shelfOrder := entity.ShelfOrder{
		UUID:        guuid.NewV4(),
		OrderUUID:   order.UUID,
		ShelfType:   shelfType,
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

func (o *orderService) PickupOrder() (*entity.Order, error) {
	// Get order that is ready for pickup from shelf that
	// has an expiration date that is the most soon.
	// We do this to minimize waste.
	shelfOrder, err := o.shelfOrderRepository.GetOpenOrder()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch open order")
	}

	// Update shelf order status to be "picked_up".
	err = o.shelfOrderRepository.UpdateOrderStatus(*shelfOrder, entity.OrderStatusPickedUp)
	if err != nil {
		return nil, errors.Wrapf(
			err, "failed to update status of shelf order %+v", shelfOrder)
	}

	// Fetch the corresponding order so the consumer (driver) has all the details.
	order, err := o.orderRepository.GetOrder(shelfOrder.OrderUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order")
	}

	return order, nil
}

func (o *orderService) GetExpiredOrdersOnShelf() ([]*entity.ShelfOrder, error) {
	shelfOrders, err := o.shelfOrderRepository.GetExpiredOrders()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch expired orders on shelf")
	}

	return shelfOrders, nil
}

func (o *orderService) MarkOrderAsWasted(shelfOrder entity.ShelfOrder) error {
	newOrderStatus := entity.OrderStatusWasted
	err := o.shelfOrderRepository.UpdateOrderStatus(shelfOrder, newOrderStatus)
	if err != nil {
		return errors.Wrapf(err, "faield to mark order as wasted %s", err.Error())
	}

	return nil
}
