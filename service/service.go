package service

import (
	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/service/repository"
)

// Services contains service layer.
type Services struct {
	Order OrderService
}

// InitializeServices initializes service layer.
func InitializeServices(cfg config.AppConfig, repositories repository.Repositories) Services {
	orderService := NewOrderService(cfg, repositories.Order, repositories.ShelfOrder)

	return Services{
		Order: orderService,
	}
}
