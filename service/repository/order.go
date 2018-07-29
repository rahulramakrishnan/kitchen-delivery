package repository

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/mapper"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// OrderRepository is the order repository interface.
type OrderRepository interface {
	CreateOrder(order entity.Order) error
	GetOrder(orderUUID guuid.UUID) (*entity.Order, error)
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
			exception.ErrInvalidInput, "failed to map order to record, err: %s", err)
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
			exception.ErrDatabase, "failed to store order, err: %s", err)
	}

	// Commit DB transaction.
	tx.Commit()
	return nil
}

// GetOrder returns a specific order.
func (o *orderRepository) GetOrder(orderUUID guuid.UUID) (*entity.Order, error) {
	var orderRecord record.Order

	err := o.db.
		Where("uuid = ?", orderUUID.String()).
		First(&orderRecord).Error
	if err == gorm.ErrRecordNotFound {
		return nil, exception.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(exception.ErrDatabase, err.Error())
	}

	order, err := mapper.RecordToOrder(orderRecord)
	if err != nil {
		return nil, errors.Wrapf(
			exception.ErrDataCorrupted, "failed to map record to order %+v, err: %s", orderRecord, err)
	}

	return order, nil
}
