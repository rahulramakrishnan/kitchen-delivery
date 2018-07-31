package entity

import (
	"fmt"
	"strings"
	"time"

	guuid "github.com/satori/go.uuid"
)

// ShelfOrder is an order placed on a shelf entity.
type ShelfOrder struct {
	UUID        guuid.UUID
	OrderUUID   guuid.UUID
	ShelfType   ShelfType
	OrderStatus OrderStatus
	Version     int
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// Used for fetching orders from shelves on a join.
	Order Order
}

// Validate verifies that a shelf order has valid fields.
func (s *ShelfOrder) Validate() error {
	var errorMsgs []string

	// Check shelf type.
	if _, ok := AllShelfTypes[s.ShelfType]; !ok {
		msg := fmt.Sprintf("shelf type %s is invalid", s.ShelfType)
		errorMsgs = append(errorMsgs, msg)
	}

	// Check order status.
	if _, ok := AllOrderStatuses[s.OrderStatus]; !ok {
		msg := fmt.Sprintf("order status %s is invalid", s.OrderStatus)
		errorMsgs = append(errorMsgs, msg)
	}

	// If error msgs exist then we return a combination of them.
	if len(errorMsgs) != 0 {
		// Combine error messages if they exist.
		err := fmt.Errorf(strings.Join(errorMsgs, ", "))
		return err
	}

	return nil
}

// String returns a prettified string representation of an order.
func (s *ShelfOrder) String() string {
	shelfOrderString := fmt.Sprintf("ShelfType: %s, OrderStatus: %s", s.ShelfType, s.OrderStatus)
	return shelfOrderString
}

// ShelfType is the type of shelf to hold the food.
type ShelfType string

var (
	// HotShelf is the hot shelf.
	HotShelf = ShelfType("hot")
	// ColdShelf is the cold shelf.
	ColdShelf = ShelfType("cold")
	// FrozenShelf is the frozen shelf.
	FrozenShelf = ShelfType("frozen")
	// OverflowShelf is the overflow shelf.
	OverflowShelf = ShelfType("overflow")
)

// AllShelfTypes is all shelf types
// and is used to verify if a shelf type is valid or not.
// We use a hashmap for O(1) look up.
var AllShelfTypes = map[ShelfType]bool{
	HotShelf:      true,
	ColdShelf:     true,
	FrozenShelf:   true,
	OverflowShelf: true,
}

// OrderStatus is order status enum.
type OrderStatus string

var (
	// OrderStatusReadyForPickup is for when an order is ready for picked up.
	OrderStatusReadyForPickup = OrderStatus("ready_for_pickup")
	// OrderStatusPickedUp is for when an order is picked up.
	OrderStatusPickedUp = OrderStatus("picked_up")
	// OrderStatusWasted is for when an order is dropped as waste after TTL has expired.
	OrderStatusWasted = OrderStatus("wasted")
)

// AllOrderStatuses holds all order statuses
// and is used for validation prior to insertion.
// We use a hashmap for O(1) look up.
var AllOrderStatuses = map[OrderStatus]bool{
	OrderStatusReadyForPickup: true,
	OrderStatusWasted:         true,
	OrderStatusPickedUp:       true,
}
