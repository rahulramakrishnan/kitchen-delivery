package mapper

import (
	"fmt"
	"testing"
	"time"

	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/endpoint"
	"github.com/kitchen-delivery/service/repository/record"

	guuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderRequestToOrder(t *testing.T) {
	orderUUID := guuid.NewV4()
	shelfLife := 300
	decayRate := 0.45
	createOrderRequests := []endpoint.CreateOrderRequest{
		{
			UUID:      orderUUID.String(),
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: fmt.Sprintf("%d", shelfLife),
			DecayRate: fmt.Sprintf("%f", decayRate),
		},
		{
			// No UUID passed in, so we ensure that we generate one.
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: fmt.Sprintf("%d", shelfLife),
			DecayRate: fmt.Sprintf("%f", decayRate),
		},
	}

	expected := &entity.Order{
		UUID:      orderUUID,
		Name:      createOrderRequests[0].Name,
		Temp:      entity.OrderTempHot,
		ShelfLife: shelfLife,
		DecayRate: decayRate,
	}

	// Verify we can map with an idempotency uuid.
	order1, err := CreateOrderRequestToOrder(createOrderRequests[0])
	assert.Nil(t, err, "no error mapping create order request to order")
	assert.Equal(t, expected, order1)

	// Verify we generate a new uuid when one is not given.
	order2, err := CreateOrderRequestToOrder(createOrderRequests[1])
	assert.Nil(t, err, "no error mapping create order request to order")
	assert.NotEqual(t, expected.UUID, guuid.NullUUID{}.UUID)
	assert.Equal(t, expected.Name, order2.Name)
	assert.Equal(t, expected.Temp, order2.Temp)
	assert.Equal(t, expected.ShelfLife, order2.ShelfLife)
	assert.Equal(t, expected.DecayRate, order2.DecayRate)
}

func TestCreateOrderRequestToOrder_InvalidRequests(t *testing.T) {
	orderUUID := guuid.NewV4()
	shelfLife := 300
	decayRate := 0.45
	invalidRequests := []endpoint.CreateOrderRequest{
		{
			UUID:      orderUUID.String(),
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: "not a valid integer", // invalid integer shelf life
			DecayRate: fmt.Sprintf("%f", decayRate),
		},
		{
			UUID:      orderUUID.String(),
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: fmt.Sprintf("%d", shelfLife),
			DecayRate: "not a valid float", // invalid float decay rate
		},
		{
			UUID:      "invalid order uuid", // invalid order uuid
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: fmt.Sprintf("%d", shelfLife),
			DecayRate: fmt.Sprintf("%f", decayRate),
		},
		{
			UUID:      orderUUID.String(),
			Name:      "Cheeze Pizza",
			Temp:      "invalid order temp", // invalid order temperature
			ShelfLife: fmt.Sprintf("%d", shelfLife),
			DecayRate: fmt.Sprintf("%f", decayRate),
		},
	}

	for _, invalidRequest := range invalidRequests {
		_, err := CreateOrderRequestToOrder(invalidRequest)
		assert.Error(t, err, "failed mapping create order request to order")
	}
}

func TestOrderToRecord(t *testing.T) {
	now := time.Now()
	orders := []entity.Order{
		{
			UUID:      guuid.NewV4(), // uuid passed in
			Name:      "Cheeze Pizza",
			Temp:      entity.OrderTempHot,
			ShelfLife: 300,
			DecayRate: 0.45,
			CreatedAt: now,
		},
		{
			// No order uuid passed in
			// so a new one should be created.
			Name:      "Cheeze Pizza",
			Temp:      entity.OrderTempHot,
			ShelfLife: 300,
			DecayRate: 0.45,
			CreatedAt: now,
		},
	}

	// Expected record for order 1 w/ uuid.
	expected := &record.Order{
		UUID:      orders[0].UUID.String(),
		Name:      orders[0].Name,
		Temp:      string(orders[0].Temp),
		ShelfLife: orders[0].ShelfLife,
		DecayRate: orders[0].DecayRate,
		CreatedAt: orders[0].CreatedAt,
	}

	record1, err := OrderToRecord(orders[0])
	assert.Nil(t, err, "no error mapping order 1 to record")
	assert.Equal(t, expected, record1)

	record2, err := OrderToRecord(orders[1])
	assert.Nil(t, err, "no error mapping order 2 to record")
	// Verify that we record 2's uuid is not a null uuid.
	assert.NotEqual(t, record2.UUID, guuid.NullUUID{}.UUID)
	assert.Equal(t, expected.Name, record2.Name)
	assert.Equal(t, expected.Temp, record2.Temp)
	assert.Equal(t, expected.ShelfLife, record2.ShelfLife)
	assert.Equal(t, expected.DecayRate, record2.DecayRate)
	assert.Equal(t, expected.CreatedAt, record2.CreatedAt)
}

func TestRecordToOrder(t *testing.T) {
	orderUUID := guuid.NewV4()
	record := record.Order{
		UUID:      orderUUID.String(),
		Name:      "Cheeze Pizza",
		Temp:      string(entity.OrderTempHot),
		ShelfLife: 300,
		DecayRate: 0.45,
		CreatedAt: time.Now(),
	}

	expected := &entity.Order{
		UUID:      orderUUID,
		Name:      record.Name,
		Temp:      entity.OrderTemp(record.Temp),
		ShelfLife: record.ShelfLife,
		DecayRate: record.DecayRate,
		CreatedAt: record.CreatedAt,
	}

	order, err := RecordToOrder(record)
	assert.Nil(t, err, "no error mapping record to order")
	assert.Equal(t, expected, order)
}

func TestRecordToOrder_InvalidRecords(t *testing.T) {
	records := []record.Order{
		{
			UUID:      "invalid uuid", // invalid uuid
			Name:      "Cheeze Pizza",
			Temp:      string(entity.OrderTempHot),
			ShelfLife: 300,
			DecayRate: 0.45,
			CreatedAt: time.Now(),
		},
		{
			UUID:      guuid.NewV4().String(),
			Name:      "Cheeze Pizza",
			Temp:      "not a valid temp enum", // invalid order temp enum
			ShelfLife: 300,
			DecayRate: 0.45,
			CreatedAt: time.Now(),
		},
	}

	for _, record := range records {
		_, err := RecordToOrder(record)
		assert.Error(t, err, "error mapping record to order")
	}
}
