package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kitchen-delivery/entity/endpoint"
)

// CheckHealth checks service health and returns 200 OK if reachable.
func (h *healthHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Load order data from input json file, and make requests to
	// kitchen-delivery's order endpoint.
	jsonFile, err := os.Open("data/input.json")
	if err != nil {
		msg := fmt.Sprintf("failed to open input file err: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}
	defer jsonFile.Close()

	// Read order inputs and unmarshall into order requests.
	// Note: We read the json file at once, because that is the only way we
	// can parse a json file, a more scalable solution is paginating
	// from entries that prepopulate a database, but this is just for
	// demo purposes.
	var orders []endpoint.OrderJSON
	ordersByteArray, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		msg := fmt.Sprintf("failed to read json file, err: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	err = json.Unmarshal(ordersByteArray, &orders)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal json, err: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	// Spawn thread to submit order requests asynchronously.
	go h.submitOrderRequests(orders)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// submitOrderRequests submits orders to an orders endpoint.
func (h *healthHandler) submitOrderRequests(orders []endpoint.OrderJSON) {
	// TODO: make a request w/ the rate provided by the assignment.
	// Iterate over order requests and submit request.
	for _, order := range orders {
		time.Sleep(250 * time.Millisecond) // rate of submitting an order is 1/4th a second

		err := h.submitOrderRequest(order)
		if err != nil {
			// We fail open here as we don't want an error in
			// the creation of one order to stop the creation of subsequent ones.
			msg := fmt.Sprintf("failed to submit order request err: %+v", err)
			log.Println(msg)
		}
	}
}

// submitOrderRequest submits an order creation HTTP request.
func (h *healthHandler) submitOrderRequest(order endpoint.OrderJSON) error {
	// Prepare HTTP Post request by stringifying URL values.
	shelfLife := fmt.Sprintf("%d", order.ShelfLife)   // safely convert int to str int
	decayRate := fmt.Sprintf("%.2f", order.DecayRate) // keep float to 2 decimal places ex: 2.56

	formData := url.Values{
		"name":      {order.Name},
		"temp":      {order.Temp},
		"shelfLife": {shelfLife},
		"decayRate": {decayRate},
	}
	resp, err := http.PostForm("http://localhost:8080/order", formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If request was not successful then return the content
	// of the response as the error.
	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body, err: %+v", err)
		}

		return fmt.Errorf("%+v", content)
	}

	return nil
}
