package repository

import "github.com/jinzhu/gorm"

// Repositories stores MySQL DB drivers.
type Repositories struct {
	Order OrderRepository
}

// InitializeRepositories initializes repositories.
func InitializeRepositories(db *gorm.DB) Repositories {
	orderRepository := NewOrderRepository(db)

	repositories := Repositories{
		Order: orderRepository,
	}

	return repositories
}
