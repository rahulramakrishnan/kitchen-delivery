package handler

import (
	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/handler/health"
	"github.com/kitchen-delivery/handler/order"
)

// Handlers holds HTTP handlers.
type Handlers struct {
	Health health.Handler
	Order  order.Handler
}

// NewHandlers returns new HTTP handlers.
func NewHandlers(cfg config.AppConfig) (*Handlers, error) {
	healthHandler := health.NewHandler(cfg)
	orderHandler := order.NewHandler(cfg)

	return &Handlers{
		Health: healthHandler,
		Order:  orderHandler,
	}, nil
}
