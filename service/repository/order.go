package repository

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service/mapper"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// OrderRepository is the user repository interface.
type OrderRepository interface {
	Create(order entity.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository is a new order repository.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Create stores an order into the orders table.
func (o *orderRepository) Create(order entity.Order) error {
	record := mapper.OrderToRecord(order)

	tx := o.db.Begin()
	err := tx.Create(&record).Error

	// We ensure idempotency on creation using order UUID.
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == mysqlerr.ER_DUP_ENTRY {
			tx.Rollback()
			return nil
		}
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
