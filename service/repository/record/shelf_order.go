package record

import "time"

// ShelfOrder is an order on a shelf record.
type ShelfOrder struct {
	UUID        string    `gorm:"column:uuid;primary_key"`
	OrderUUID   string    `gorm:"column:order_uuid"`   // FK on Orders
	ShelfType   string    `gorm:"column:shelf_type"`   // "hot", "cold", "frozen", "overflow"
	OrderStatus string    `gorm:"column:order_status"` // "ready_for_pickup", "picked_up", "wasted"
	Version     int       `gorm:"column:version"`      // Used for optimistic locking.
	ExpiresAt   time.Time `gorm:"column:expires_at"`   // time when order expires
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
