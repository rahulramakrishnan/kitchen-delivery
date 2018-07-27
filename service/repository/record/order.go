package record

import "time"

// Order is an order record.
type Order struct {
	UUID      string    `gorm:"column:uuid;primary_key"`
	Name      string    `gorm:"column:name"`
	Temp      string    `gorm:"column:temp"`
	ShelfLife int       `gorm:"column:shelf_life"`
	DecayRate float64   `gorm:"column:decay_rate"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
