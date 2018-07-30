package repository

import (
	"time"

	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/mapper"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// ShelfOrderRepository is the shelf order repository interface.
type ShelfOrderRepository interface {
	AddOrderToShelf(shelfOrder entity.ShelfOrder) error
	CountOrdersOnShelf(shelfType entity.ShelfType) (int, error)
	UpdateOrderStatus(shelfOrder entity.ShelfOrder, orderStatus entity.OrderStatus) error
	GetOpenOrder() (*entity.ShelfOrder, error)
	GetExpiredOrders() ([]*entity.ShelfOrder, error)
}

type shelfRepository struct {
	db *gorm.DB
}

// NewShelfOrderRepository is a new order repository.
func NewShelfOrderRepository(db *gorm.DB) ShelfOrderRepository {
	return &shelfRepository{
		db: db,
	}
}

// AddOrderToShelf adds an order to a designated shelf.
func (s *shelfRepository) AddOrderToShelf(shelfOrder entity.ShelfOrder) error {
	record := mapper.ShelfOrderToRecord(shelfOrder)

	// Begin DB transaction.
	tx := s.db.Begin()
	err := tx.Create(&record).Error

	// We ensure idempotency on DB create
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == mysqlerr.ER_DUP_ENTRY {
			tx.Rollback()
			return nil
		}
	}

	if err != nil {
		tx.Rollback()
		return errors.Wrapf(exception.ErrDatabase, "failed to add order to shelf - err: %s", err)
	}

	tx.Commit()
	return nil
}

// CountOrdersOnShelf counts shelf orders.
func (s *shelfRepository) CountOrdersOnShelf(shelfType entity.ShelfType) (int, error) {
	// Check count of orders in "hot" w/ status of ready for pick up.
	count := 0

	orderStatus := entity.OrderStatusReadyForPickup
	err := s.db.Model(&record.ShelfOrder{}).
		Where("shelf_type = ?", string(shelfType)).
		Where("order_status = ?", string(orderStatus)).
		Count(&count).
		Error

	if err != nil {
		return 0, errors.Wrapf(
			exception.ErrDatabase, "failed to count shelf orders - err: %s", err.Error())
	}

	return count, nil
}

// UpdateOrderStatus updates a shelf order's status.
func (s *shelfRepository) UpdateOrderStatus(shelfOrder entity.ShelfOrder, orderStatus entity.OrderStatus) error {
	// We set up map of conditions to update a request with.
	newVersion := shelfOrder.Version + 1 // increment version number - optimistic locking

	conditions := make(map[string]interface{})
	conditions["order_status"] = string(orderStatus)
	conditions["version"] = newVersion

	// We map user entity to user record.
	record := mapper.ShelfOrderToRecord(shelfOrder)

	// We start db transaction master instance.
	tx := s.db.Begin()

	// We update the row that matches the set of conditions.
	updateOperation := tx.Model(&record).
		Where("uuid = ?", shelfOrder.UUID.String()).
		Where("version = ?", shelfOrder.Version).
		Updates(conditions)

	if updateOperation.Error != nil {
		tx.Rollback()
		return updateOperation.Error
	}

	// If we do not update anything then we do not commit the transaction.
	// We also do not return an error to support idempotency.
	// If an update operation fails because of a database issue
	// we would have caught it in the error check above.
	if updateOperation.RowsAffected == 0 {
		tx.Rollback()
		return exception.ErrVersionInvalid
	}

	// If everything is successful we commit the txn.
	tx.Commit()
	return nil
}

// GetOpenOrder returns an order ready for pickup w/ the most soon expiration date.
func (s *shelfRepository) GetOpenOrder() (*entity.ShelfOrder, error) {
	var shelfOrderRecord record.ShelfOrder

	err := s.db.
		// Only return orders ready for pick up.
		Where("order_status = ?", string(entity.OrderStatusReadyForPickup)).
		// We want to optimize for minimizing waste.
		Order("expires_at asc").
		First(&shelfOrderRecord).Error
	if err == gorm.ErrRecordNotFound {
		return nil, exception.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(exception.ErrDatabase, err.Error())
	}

	shelfOrder, err := mapper.RecordToShelfOrder(shelfOrderRecord)
	if err != nil {
		return nil, errors.Wrapf(
			exception.ErrDataCorrupted, "failed to map record to shelf order %+v - err: %s", shelfOrderRecord, err.Error())
	}

	return shelfOrder, nil
}

// GetExpiredOrders returns orders that have expired.
func (s *shelfRepository) GetExpiredOrders() ([]*entity.ShelfOrder, error) {
	var shelfOrderRecords []*record.ShelfOrder
	now := time.Now()

	err := s.db.
		// Only return orders ready for pick up.
		Where("order_status = ?", string(entity.OrderStatusReadyForPickup)).
		Where("expires_at < ?", now.String()). // records that have already expired
		Find(&shelfOrderRecords).Error
	if err == gorm.ErrRecordNotFound {
		// Finding expired orders .
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(exception.ErrDatabase, err.Error())
	}

	// Map records to shelf orders.
	shelfOrders, err := mapper.RecordsToShelfOrders(shelfOrderRecords)
	if err != nil {
		return nil, errors.Wrapf(
			exception.ErrDataCorrupted, "failed to map record to shelf order - err: %s", err.Error())
	}

	return shelfOrders, nil
}
