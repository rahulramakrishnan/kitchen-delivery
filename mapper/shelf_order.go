package mapper

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// ShelfOrderToRecord maps an order entity to an order record.
func ShelfOrderToRecord(shelfOrder entity.ShelfOrder) record.ShelfOrder {
	record := record.ShelfOrder{
		UUID:        shelfOrder.UUID.String(),
		OrderUUID:   shelfOrder.OrderUUID.String(),
		ShelfType:   string(shelfOrder.ShelfType),
		OrderStatus: string(shelfOrder.OrderStatus),
		Version:     shelfOrder.Version,
		ExpiresAt:   shelfOrder.ExpiresAt,
		CreatedAt:   shelfOrder.CreatedAt,
		UpdatedAt:   shelfOrder.UpdatedAt,
	}

	// We set a random uuid for shelf order if there is not one passed in.
	nullUUID := guuid.NullUUID{}
	if nullUUID.UUID == shelfOrder.UUID {
		record.UUID = guuid.NewV4().String()
	}

	return record
}

// RecordToShelfOrder maps a shelf order record to an shelf order entity.
func RecordToShelfOrder(record record.ShelfOrder) (*entity.ShelfOrder, error) {
	uuid, err := guuid.FromString(record.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "uuid is not valid, uuid: %s", record.UUID)
	}

	orderUUID, err := guuid.FromString(record.OrderUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "order uuid is not valid, uuid: %s", record.OrderUUID)
	}

	shelfOrder := entity.ShelfOrder{
		UUID:        uuid,
		OrderUUID:   orderUUID,
		ShelfType:   entity.ShelfType(record.ShelfType),
		OrderStatus: entity.OrderStatus(record.OrderStatus),
		Version:     record.Version,
		ExpiresAt:   record.ExpiresAt,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}

	err = shelfOrder.Validate()
	if err != nil {
		return nil, err
	}

	return &shelfOrder, nil
}
