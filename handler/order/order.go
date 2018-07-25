package order

import (
	"net/http"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/service"
)

// Handler is Health handler interface.
type Handler interface {
	HandleOrder(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	cfg      config.AppConfig
	services service.Services
}

// NewHandler creates a new HTTP order handler instance.
func NewHandler(appConfig config.AppConfig, services service.Services) Handler {
	return &orderHandler{
		cfg:      appConfig,
		services: services,
	}
}

// HandleOrder either creates an order, or sends an order back to a driver.
func (o *orderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
