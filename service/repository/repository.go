package repository

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// Repositories stores MySQL DB drivers.
type Repositories struct {
	Order      OrderRepository
	ShelfOrder ShelfOrderRepository
}

// InitializeRepositories initializes repositories.
func InitializeRepositories(db *gorm.DB, redisConn redis.Conn) Repositories {
	orderRepository := NewOrderRepository(db)
	shelfOrderRepository := NewShelfOrderRepository(db, redisConn)

	repositories := Repositories{
		Order:      orderRepository,
		ShelfOrder: shelfOrderRepository,
	}

	return repositories
}
