package order

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity/endpoint"
	"github.com/kitchen-delivery/mapper"
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
	// Parse form so we can access key value pairs of post request.
	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("failed to parse form err: %+v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// We map key, value pair http request to a request entity.
	// We check if there are any errors in the submission and return an error if there is.
	formData := endpoint.FormData(r.PostForm)
	fieldsToExtract := endpoint.FieldsToExtract{
		RequiredFields: []string{"name", "temp", "shelfLife", "decayRate"},
		OptionalFields: []string{"uuid"},
	}
	createOrderRequest := endpoint.CreateOrderRequest{}
	err = endpoint.ExtractRequest(formData, fieldsToExtract, &createOrderRequest)
	if err != nil {
		msg := fmt.Sprintf("failed to handle create order request %+v", err.Error())
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// Map a HTTP create order request to an order entity.
	order, err := mapper.CreateOrderRequestToOrder(createOrderRequest)
	if err != nil {
		msg := fmt.Sprintf("failed to map create order request to order %+v", err.Error())
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// TODO: put on a channel that get's pulled off of by workers spawned at start time.
	err = o.services.Order.Create(*order)
	if err != nil {
		msg := fmt.Sprintf("failed to store order %+v", err.Error())
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	// Send back order uuid to client on success.
	// This will support client-polling and allow for idempotency.
	orderUUID := order.UUID

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(orderUUID.String()))
}
