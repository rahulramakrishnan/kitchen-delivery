package service

import (
	"testing"
	"time"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/service/repository"
	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	guuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Load app config.
	cfg := config.AppConfig{}
	cfg.LoadConfig("../config/development.yaml")

	orderRepository := repository.NewMockOrderRepository(ctrl)
	shelfOrderRepository := repository.NewMockShelfOrderRepository(ctrl)

	expected := &orderService{
		cfg:                  cfg,
		orderRepository:      orderRepository,
		shelfOrderRepository: shelfOrderRepository,
		shelfSpace: map[entity.ShelfType]int{
			entity.HotShelf:      cfg.ShelfSpace.Hot,
			entity.ColdShelf:     cfg.ShelfSpace.Cold,
			entity.FrozenShelf:   cfg.ShelfSpace.Frozen,
			entity.OverflowShelf: cfg.ShelfSpace.Overflow,
		},
	}

	orderService := NewOrderService(cfg, orderRepository, shelfOrderRepository)
	assert.Equal(t, expected, orderService)
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Load app config.
	cfg := config.AppConfig{}
	cfg.LoadConfig("../config/development.yaml")

	orderRepository := repository.NewMockOrderRepository(ctrl)
	shelfOrderRepository := repository.NewMockShelfOrderRepository(ctrl)
	orderService := NewOrderService(cfg, orderRepository, shelfOrderRepository)

	order := entity.Order{
		UUID:      guuid.NewV4(),
		Name:      "Cheeze Pizza",
		Temp:      entity.OrderTempHot,
		ShelfLife: 300,
		DecayRate: 0.45,
	}

	orderRepository.EXPECT().CreateOrder(order)

	err := orderService.CreateOrder(order)
	assert.Nil(t, err)
}

func TestPlaceOrderOnShelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Load app config.
	cfg := config.AppConfig{}
	cfg.LoadConfig("../config/development.yaml")

	orderRepository := repository.NewMockOrderRepository(ctrl)
	shelfOrderRepository := repository.NewMockShelfOrderRepository(ctrl)
	orderService := NewOrderService(cfg, orderRepository, shelfOrderRepository)

	order := entity.Order{
		UUID:      guuid.NewV4(),
		Name:      "Cheeze Pizza",
		Temp:      entity.OrderTempHot,
		ShelfLife: 300,
		DecayRate: 0.45,
	}
	shelfType := order.GetShelfType()

	// Prepare expected shelf order.
	ttl := order.GetTTL()
	now := time.Now()
	expirationDate := now.Add(time.Second * time.Duration(ttl))

	expectedShelfOrder := entity.ShelfOrder{
		OrderUUID:   order.UUID,
		ShelfType:   order.GetShelfType(),
		OrderStatus: entity.OrderStatusReadyForPickup,
		Version:     0,
		ExpiresAt:   expirationDate,
	}

	gomock.InOrder(
		shelfOrderRepository.EXPECT().CountOrdersOnShelf(shelfType).Return(3, nil),
		shelfOrderRepository.EXPECT().AddOrderToShelf(&shelfOrderMatcher{expectedShelfOrder}).Return(nil),
	)

	err := orderService.PlaceOrderOnShelf(order)
	assert.Nil(t, err)
}

func TestPlaceOrderOnShelf_CountOrdersOnShelfError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Load app config.
	cfg := config.AppConfig{}
	cfg.LoadConfig("../config/development.yaml")

	orderRepository := repository.NewMockOrderRepository(ctrl)
	shelfOrderRepository := repository.NewMockShelfOrderRepository(ctrl)
	orderService := NewOrderService(cfg, orderRepository, shelfOrderRepository)

	order := entity.Order{
		UUID:      guuid.NewV4(),
		Name:      "Cheeze Pizza",
		Temp:      entity.OrderTempHot,
		ShelfLife: 300,
		DecayRate: 0.45,
	}
	shelfType := order.GetShelfType()

	shelfOrderRepository.EXPECT().CountOrdersOnShelf(shelfType).Return(0, exception.ErrDatabase)

	err := orderService.PlaceOrderOnShelf(order)
	assert.Equal(t, exception.ErrDatabase, errors.Cause(err))
}

func TestPlaceOrderOnShelf_ShelfSpaceIsFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Load app config.
	cfg := config.AppConfig{}
	cfg.LoadConfig("../config/development.yaml")

	orderRepository := repository.NewMockOrderRepository(ctrl)
	shelfOrderRepository := repository.NewMockShelfOrderRepository(ctrl)
	orderService := NewOrderService(cfg, orderRepository, shelfOrderRepository)

	order := entity.Order{
		UUID:      guuid.NewV4(),
		Name:      "Cheeze Pizza",
		Temp:      entity.OrderTempHot,
		ShelfLife: 300,
		DecayRate: 0.45,
	}
	shelfType := order.GetShelfType()

	hotLimit := cfg.ShelfSpace.Hot
	overflowLimit := cfg.ShelfSpace.Overflow

	gomock.InOrder(
		shelfOrderRepository.EXPECT().CountOrdersOnShelf(shelfType).Return(hotLimit, nil),
		shelfOrderRepository.EXPECT().CountOrdersOnShelf(entity.OverflowShelf).Return(overflowLimit, nil),
	)

	err := orderService.PlaceOrderOnShelf(order)
	assert.Equal(t, exception.ErrFullShelf, errors.Cause(err))
}

// shelfOrderMatcher holds shelf order matchers.
type shelfOrderMatcher struct {
	ShelfOrder entity.ShelfOrder
}

func (s *shelfOrderMatcher) String() string {
	return "shelf order matches expected parameters"
}

func (s *shelfOrderMatcher) Matches(x interface{}) bool {
	shelfOrder := x.(entity.ShelfOrder)
	doesMatch := (s.ShelfOrder.OrderUUID == shelfOrder.OrderUUID &&
		s.ShelfOrder.ShelfType == shelfOrder.ShelfType &&
		s.ShelfOrder.OrderStatus == shelfOrder.OrderStatus &&
		s.ShelfOrder.Version == shelfOrder.Version)

	now := time.Now()
	isExpiresAtInFuture := shelfOrder.ExpiresAt.After(now)
	return doesMatch && isExpiresAtInFuture
}
