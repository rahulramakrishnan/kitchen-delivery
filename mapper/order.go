package mapper

import (
	"strconv"

	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/endpoint"
	"github.com/kitchen-delivery/service/repository/record"

	"github.com/pkg/errors"
	guuid "github.com/satori/go.uuid"
)

// CreateOrderRequestToOrder maps a HTTP create order request to an order entity.
func CreateOrderRequestToOrder(createOrderRequest endpoint.CreateOrderRequest) (*entity.Order, error) {
	shelfLife, err := strconv.ParseInt(createOrderRequest.ShelfLife, 0, 32)
	if err != nil {
		return nil, errors.Wrapf(
			err, "failed to parse int shelf life %s", createOrderRequest.ShelfLife)
	}

	decayRate, err := strconv.ParseFloat(createOrderRequest.DecayRate, 64)
	if err != nil {
		return nil, errors.Wrapf(
			err, "failed to parse float decay rate %s", createOrderRequest.DecayRate)
	}

	// We support idempotency by checking if an order UUID is passed.
	// If it is, we convert it to a UUID.
	var orderUUID guuid.UUID

	if createOrderRequest.UUID != "" {
		orderUUID, err = guuid.FromString(createOrderRequest.UUID)
		if err != nil {
			return nil, errors.Wrapf(
				err, "create order request uuid is invalid - uuid: %s", createOrderRequest.UUID)
		}
	} else {
		// If not uuid is passed, we generate a new order uuid.
		orderUUID = guuid.NewV4()
	}

	order := entity.Order{
		UUID:      orderUUID,
		Name:      createOrderRequest.Name,
		Temp:      entity.OrderTemp(createOrderRequest.Temp),
		ShelfLife: int(shelfLife),
		DecayRate: decayRate,
	}

	err = order.Validate()
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// OrderToRecord maps an order entity to an order record.
func OrderToRecord(order entity.Order) (*record.Order, error) {
	record := record.Order{
		UUID:      order.UUID.String(),
		Name:      order.Name,
		Temp:      string(order.Temp),
		ShelfLife: order.ShelfLife,
		DecayRate: order.DecayRate,
		CreatedAt: order.CreatedAt,
	}

	// We set a random uuid for order if there is not one passed in.
	nullUUID := guuid.NullUUID{}
	if nullUUID.UUID == order.UUID {
		record.UUID = guuid.NewV4().String()
	}

	return &record, nil
}

// RecordToOrder maps an order record to an order entity.
func RecordToOrder(record record.Order) (*entity.Order, error) {
	orderUUID, err := guuid.FromString(record.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "uuid is not valid, uuid: %s", record.UUID)
	}

	order := entity.Order{
		UUID:      orderUUID,
		Name:      record.Name,
		Temp:      entity.OrderTemp(record.Temp),
		ShelfLife: record.ShelfLife,
		DecayRate: record.DecayRate,
		CreatedAt: record.CreatedAt,
	}

	err = order.Validate()
	if err != nil {
		return nil, err
	}

	return &order, nil
}
