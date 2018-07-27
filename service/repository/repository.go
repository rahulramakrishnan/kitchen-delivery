package repository

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// Repositories stores MySQL DB drivers.
type Repositories struct {
	Order OrderRepository
	Shelf ShelfRepository
}

// InitializeRepositories initializes repositories.
func InitializeRepositories(db *gorm.DB, redisConn redis.Conn) Repositories {
	orderRepository := NewOrderRepository(db)
	shelfRepository := NewShelfRepository(redisConn)

	repositories := Repositories{
		Order: orderRepository,
		Shelf: shelfRepository,
	}

	return repositories
}
