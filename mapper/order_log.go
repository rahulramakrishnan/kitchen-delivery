package mapper

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// OrderLogToRecord maps an order entity to an order record.
func OrderLogToRecord(orderLog entity.OrderLog) (*record.OrderLog, error) {
	record := record.OrderLog{
		UUID:        orderLog.UUID.String(),
		OrderUUID:   orderLog.OrderUUID.String(),
		OrderStatus: string(orderLog.OrderStatus),
		Description: orderLog.Description,
		CreatedAt:   orderLog.CreatedAt,
	}

	// We set a random uuid for order history
	// if there is not one passed in.
	nullUUID := guuid.NullUUID{}
	if nullUUID.UUID == orderLog.UUID {
		newUUID, err := guuid.NewV4()
		if err != nil {
			return nil, err
		}

		record.UUID = newUUID.String()
	}

	return &record, nil
}

// RecordToOrderLog maps an order record to an order history entity.
func RecordToOrderLog(record record.OrderLog) (*entity.OrderLog, error) {
	uuid, err := guuid.FromString(record.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "uuid is not valid, uuid: %s", record.UUID)
	}

	orderUUID, err := guuid.FromString(record.OrderUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "order uuid is not valid, uuid: %s", record.OrderUUID)
	}

	orderLog := entity.OrderLog{
		UUID:        uuid,
		OrderUUID:   orderUUID,
		OrderStatus: entity.OrderStatus(record.OrderStatus),
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
	}

	err = orderLog.Validate()
	if err != nil {
		return nil, err
	}

	return &orderLog, nil
}
