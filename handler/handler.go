package handler

import (
	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/handler/health"
)

// Handlers holds HTTP handlers.
type Handlers struct {
	Health health.Handler
}

// NewHandlers returns new HTTP handlers.
func NewHandlers(cfg config.AppConfig) (*Handlers, error) {
	healthHandler := health.NewHandler(cfg)

	return &Handlers{
		Health: healthHandler,
	}, nil
}
