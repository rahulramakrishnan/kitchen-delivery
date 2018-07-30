package order

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/endpoint"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/kitchen-delivery/mapper"
	"github.com/kitchen-delivery/service"

	"github.com/pkg/errors"
)

// Handler is Health handler interface.
type Handler interface {
	HandleOrder(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	cfg      config.AppConfig
	services service.Services
	queue    chan<- *entity.Order // we only send to this channel
}

// NewHandler creates a new HTTP order handler instance.
func NewHandler(appConfig config.AppConfig, services service.Services, queue chan<- *entity.Order) Handler {
	return &orderHandler{
		cfg:      appConfig,
		services: services,
		queue:    queue,
	}
}

// HandleOrder either creates an order, or sends an order back to a driver.
func (o *orderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	// If request is a HTTP GET then we send back an order.
	if r.Method == http.MethodGet {
		o.pickupOrder(w, r)
		return
	}

	o.createOrder(w, r)
}

func (o *orderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	// Otherwise we handle an order creation w/ an HTTP POST request.
	// Parse form so we can access key value pairs of post request.
	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("failed to parse form - err: %s", err)
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
		msg := fmt.Sprintf("failed to handle create order request - err: %s", err)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// Map a HTTP create order request to an order entity.
	order, err := mapper.CreateOrderRequestToOrder(createOrderRequest)
	if err != nil {
		msg := fmt.Sprintf("failed to map create order request to order - err: %s", err)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	// Persist order to DB, before returning success to client.
	err = o.services.Order.CreateOrder(*order)
	if err != nil {
		if errors.Cause(err) == exception.ErrFullShelf {
			msg := "shelf is full"
			log.Println(msg)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(msg))
			return
		}

		msg := fmt.Sprintf("failed to store order - err: %s", err)
		log.Println(msg)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(msg))
		return
	}

	// Place order on queue which multiple worker threads pull off
	// concurrently. This is increases the throughput that our API can handle.
	// We purposesfully do not close this channel because we want to keep it open
	// for workers to continue pulling indefinitely.
	o.queue <- order

	// Send back order uuid to client on success.
	// This will support client-polling and allow for idempotency.
	orderUUID := order.UUID

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(orderUUID.String()))
}

// pickupOrder picks up an order.
func (o *orderHandler) pickupOrder(w http.ResponseWriter, r *http.Request) {
	order, err := o.services.Order.PickupOrder()
	if err != nil {
		switch errors.Cause(err) {
		case exception.ErrNotFound:
			msg := fmt.Sprintf("no more orders - err: %s", err)
			log.Println(msg)

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(msg))
			return
		default:
			msg := fmt.Sprintf("failed to pickup order - err: %s", err)
			log.Println(msg)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}
	}

	// Stringify the contents of the order.
	orderContents := order.String()

	log.Printf("driver picked up order successfully - %s", orderContents)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(orderContents))
}
