package record

import "time"

// OrderLog is an order history record.
type OrderLog struct {
	UUID        string    `gorm:"column:uuid;primary_key"`
	OrderUUID   string    `gorm:"column:order_uuid"` // FK on Orders
	OrderStatus string    `gorm:"column:order_status"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
