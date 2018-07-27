package entity

import (
	"time"

	"github.com/kitchen-delivery/entity/exception"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// OrderLog is an order history record.
type OrderLog struct {
	UUID        guuid.UUID
	OrderUUID   guuid.UUID
	OrderStatus OrderStatus // ex: "hot_shelf", "overflow_shelf", "dropped", "wasted", "picked_up"
	Description string      // ex: TTL w/ doubled decay rate is 211
	CreatedAt   time.Time   // no updated at b/c this is an immutable table.
}

// OrderStatus is order status enum.
type OrderStatus string

var (
	// OrderStatusHotShelf is for when an order is placed on the hot shelf.
	OrderStatusHotShelf = OrderStatus("hot_shelf")
	// OrderStatusColdShelf is for when an order is placed on the cold shelf.
	OrderStatusColdShelf = OrderStatus("hot_shelf")
	// OrderStatusFrozenShelf is for when an order is placed on the frozen shelf.
	OrderStatusFrozenShelf = OrderStatus("frozen_shelf")
	// OrderStatusOverflowShelf is for when an order is placed on the overflow shelf.
	OrderStatusOverflowShelf = OrderStatus("overflow_shelf")
	// OrderStatusWasted is for when an order is dropped as waste after TTL has expired.
	OrderStatusWasted = OrderStatus("wasted")
	// OrderStatusDropped is for when an order is dropped if there is not any space on any shelf.
	OrderStatusDropped = OrderStatus("dropped")
	// OrderStatusPickedUp is for when an order is picked up.
	OrderStatusPickedUp = OrderStatus("picked_up")
)

// AllOrderStatuses holds all order statuses
// and is used for validation prior to insertion.
var AllOrderStatuses = map[OrderStatus]bool{
	OrderStatusHotShelf:      true,
	OrderStatusColdShelf:     true,
	OrderStatusFrozenShelf:   true,
	OrderStatusOverflowShelf: true,
	OrderStatusWasted:        true,
	OrderStatusDropped:       true,
	OrderStatusPickedUp:      true,
}

// Validate verifies that an order is valid.
func (o *OrderLog) Validate() error {
	_, ok := AllOrderStatuses[o.OrderStatus]
	if !ok {
		return errors.Wrapf(
			exception.ErrInvalidInput, "status value is invalid %s", o.OrderStatus)
	}

	return nil
}
