package mapper

import (
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/endpoint"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// CreateOrderRequestToOrder maps a HTTP create order request to an order entity.
func CreateOrderRequestToOrder(createOrderRequest endpoint.CreateOrderRequest) (*entity.Order, error) {
	order := entity.Order{
		Name:      createOrderRequest.Name,
		Temp:      entity.OrderTemp(createOrderRequest.Temp),
		ShelfLife: createOrderRequest.ShelfLife,
		DecayRate: createOrderRequest.DecayRate,
	}

	err := order.Validate()
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// OrderToRecord maps an order entity to an order record.
func OrderToRecord(order entity.Order) record.Order {
	record := record.Order{
		UUID:      order.UUID.String(),
		Name:      order.Name,
		Temp:      string(order.Temp),
		ShelfLife: order.ShelfLife,
		CreatedAt: order.CreatedAt,
	}

	return record
}

// RecordToOrder maps an order record to an order entity.
func RecordToOrder(record record.Order) (*entity.Order, error) {
	orderUUID, err := guuid.FromString(record.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "uuid is not valid %s", record.UUID)
	}

	order := entity.Order{
		UUID:      orderUUID,
		Name:      record.Name,
		Temp:      entity.OrderTemp(record.Temp),
		ShelfLife: record.ShelfLife,
		CreatedAt: record.CreatedAt,
	}

	err = order.Validate()
	if err != nil {
		return nil, err
	}

	return &order, nil
}
