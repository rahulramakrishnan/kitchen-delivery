package repository

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/mapper"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// OrderRepository is the order repository interface.
type OrderRepository interface {
	CreateOrder(order entity.Order) error
	CreateOrderLog(orderLog entity.OrderLog) error
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

// CreateOrder stores an order into the orders table.
func (o *orderRepository) CreateOrder(order entity.Order) error {
	// Map order entity to order record.
	record, err := mapper.OrderToRecord(order)
	if err != nil {
		return errors.Wrapf(
			exception.ErrInvalidInput, "failed to map order to record, err: %+v", err)
	}

	// Begin DB transaction.
	tx := o.db.Begin()
	err = tx.Create(&record).Error

	// We ensure idempotency on creation using order UUID.
	// If the same order already exists we rollback transaction.
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == mysqlerr.ER_DUP_ENTRY {
			tx.Rollback()
			return nil
		}
	}

	if err != nil {
		tx.Rollback()
		return errors.Wrapf(
			exception.ErrDatabase, "failed to store order, err: %+v", err)
	}

	// Commit DB transaction.
	tx.Commit()
	return nil
}

// CreateOrderLog stores an order history entry into the table.
func (o *orderRepository) CreateOrderLog(orderLog entity.OrderLog) error {
	// Map order history entity to order history record.
	record, err := mapper.OrderLogToRecord(orderLog)
	if err != nil {
		return errors.Wrapf(
			exception.ErrInvalidInput, "failed to map order to record, err: %+v", err)
	}

	// Begin DB transaction.
	tx := o.db.Begin()
	err = tx.Create(&record).Error

	// We ensure idempotency on creation using order history UUID.
	// If the same order already exists we rollback transaction.
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == mysqlerr.ER_DUP_ENTRY {
			tx.Rollback()
			return nil
		}
	}

	if err != nil {
		tx.Rollback()
		return errors.Wrapf(
			exception.ErrDatabase, "failed to store order, err: %+v", err)
	}

	// Commit DB transaction.
	tx.Commit()
	return nil
}
