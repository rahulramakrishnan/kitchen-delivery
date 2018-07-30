package entity

import (
	"fmt"
	"math"
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

// String returns a prettified string representation of an order.
func (o *Order) String() string {
	orderString := fmt.Sprintf("Name: %s, Temp: %s", o.Name, o.Temp)
	return orderString
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

// GetTTL returns the ttl for the order.
func (o *Order) GetTTL() int {
	// Calculate time to live in seconds based on formula.
	// Remember an order is waste after the "value" becomes zero.
	// This leads the formula to be reduced to:
	// => orderAge = shelfLife / (1 + decayRate)
	// We're given shelfLife and decayRate so we can solve for
	// how old an order can get before we consider it as waste.
	expirationTime := float64(o.ShelfLife) / (1.0 + o.DecayRate)
	ttl := int(math.Floor(expirationTime))
	return ttl
}
