package repository

import (
	"github.com/jinzhu/gorm"
)

// Repositories stores MySQL DB drivers.
type Repositories struct {
	Order      OrderRepository
	ShelfOrder ShelfOrderRepository
}

// InitializeRepositories initializes repositories.
func InitializeRepositories(db *gorm.DB) Repositories {
	orderRepository := NewOrderRepository(db)
	shelfOrderRepository := NewShelfOrderRepository(db)

	repositories := Repositories{
		Order:      orderRepository,
		ShelfOrder: shelfOrderRepository,
	}

	return repositories
}
