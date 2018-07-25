package health

import (
	"net/http"

	"github.com/kitchen-delivery/config"
)

// Handler is Health handler interface.
type Handler interface {
	// CheckHealth verifies that the service is running and reachable.
	CheckHealth(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
	cfg config.AppConfig
}

// NewHandler creates a new HTTP health handler instance.
func NewHandler(appConfig config.AppConfig) Handler {
	return &healthHandler{
		cfg: appConfig,
	}
}

// CheckHealth checks service health and returns 200 OK if reachable.
func (h *healthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
