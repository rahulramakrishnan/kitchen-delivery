package entity

import (
	"time"

	"github.com/kitchen-delivery/entity/exception"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// Order is a kitchen order from a customer.
type Order struct {
	UUID      guuid.UUID
	Name      string    // ex: "Cheeze Pizza"
	Temp      OrderTemp // temperature, enum: ['hot', 'cold', 'frozen']
	ShelfLife int       // shelf life in seconds
	DecayRate float64   // decay rate ex: 0.45
	CreatedAt time.Time // no updated at b/c this is an immutable table
}

// OrderTemp is order temperature enum.
type OrderTemp string

var (
	// OrderTempHot is order temperature hot.
	OrderTempHot = OrderTemp("hot")
	// OrderTempCold is order temperature cold.
	OrderTempCold = OrderTemp("cold")
	// OrderTempFrozen is order temperature frozen.
	OrderTempFrozen = OrderTemp("frozen")
)

// AllOrderTemp holds all order temperatures
// and is used for validation prior to insertion
// and validation after order retrieval.
var AllOrderTemp = map[OrderTemp]bool{
	OrderTempHot:    true,
	OrderTempCold:   true,
	OrderTempFrozen: true,
}

// Validate verifies that an order is valid.
func (o *Order) Validate() error {
	_, ok := AllOrderTemp[o.Temp]
	if !ok {
		return errors.Wrapf(
			exception.ErrInvalidInput, "temp value is invalid, temp: %s", o.Temp)
	}

	return nil
}

// GetShelfType returns shelf name based on order temp.
func (o *Order) GetShelfType() ShelfType {
	switch o.Temp {
	case OrderTempHot:
		return HotShelf
	case OrderTempCold:
		return ColdShelf
	case OrderTempFrozen:
		return FrozenShelf
	default:
		return OverflowShelf
	}
}
