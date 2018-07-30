package service

import (
	"time"

	"github.com/kitchen-delivery/config"
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
	GetOrder(orderUUID guuid.UUID) (*entity.Order, error)
	PickupOrder() (*entity.Order, error)
	GetExpiredOrdersOnShelf() ([]*entity.ShelfOrder, error)
	MarkOrderAsWasted(entity.ShelfOrder) error
}

type orderService struct {
	cfg                  config.AppConfig
	orderRepository      repository.OrderRepository
	shelfOrderRepository repository.ShelfOrderRepository
	shelfSpace           map[entity.ShelfType]int
}

// NewOrderService returns a new user service.
// switch to userRepositories
func NewOrderService(cfg config.AppConfig, orderRepository repository.OrderRepository, shelfOrderRepository repository.ShelfOrderRepository) OrderService {
	// Holds how many items each type of shelf can hold at any given time.
	shelfSpace := map[entity.ShelfType]int{
		entity.HotShelf:      cfg.ShelfSpace.Hot,
		entity.ColdShelf:     cfg.ShelfSpace.Cold,
		entity.FrozenShelf:   cfg.ShelfSpace.Frozen,
		entity.OverflowShelf: cfg.ShelfSpace.Overflow,
	}

	return &orderService{
		cfg:                  cfg,
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

	// Determine if the corresponding shelf has space.
	var isOverflowShelfFull bool
	isCorrespondingShelfFull := numOfOrders >= o.shelfSpace[correspondingShelfType]
	if isCorrespondingShelfFull {
		// If the corresponding shelf is full we count orders on overflow shelf.
		numOfOrders, err = o.shelfOrderRepository.CountOrdersOnShelf(entity.OverflowShelf)
		if err != nil {
			return errors.Wrapf(err, "failed to count orders on shelf %+v", order)
		}

		isOverflowShelfFull = numOfOrders >= o.shelfSpace[entity.OverflowShelf]
	}

	// First, if both the corresponding shelf and the overflow shelf are full
	// we throw a retriable service full shelf exception so a caller can handle it explictly.
	areAllShelvesFull := isCorrespondingShelfFull && isOverflowShelfFull
	if areAllShelvesFull {
		return errors.Wrap(
			exception.ErrFullShelf, "all shelves are filled, please retry again later")
	}

	// Now we know that either the corresponding shelf has
	// space or the overflow shelf has space.
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

	ttl := order.GetTTL()
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

func (o *orderService) GetOrder(orderUUID guuid.UUID) (*entity.Order, error) {
	// Fetch the corresponding order so the consumer (driver) has all the details.
	order, err := o.orderRepository.GetOrder(orderUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order")
	}

	return order, nil
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
