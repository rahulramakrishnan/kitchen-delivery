package health

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	guuid "github.com/satori/go.uuid"
)

// CheckHealth checks service health and returns 200 OK if reachable.
func (h *healthHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: read from input.json
	orderUUID, err := h.submitOrderRequest()
	if err != nil {
		msg := fmt.Sprintf("submit order request %+v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	// Send back order uuid.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(orderUUID.String()))
}

// submitOrderRequest submits an order creation HTTP request.
func (h *healthHandler) submitOrderRequest() (*guuid.UUID, error) {
	formData := url.Values{
		"name":      {"Pepporinni Pizza"},
		"temp":      {"hot"},
		"shelfLife": {"300"},
		"decayRate": {"0.45"},
	}
	resp, err := http.PostForm("http://localhost:8080/order", formData)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	orderUUID, err := guuid.FromString(string(content))
	if err != nil {
		return nil, err
	}

	return &orderUUID, nil
}
