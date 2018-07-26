package endpoint

// CreateOrderRequest holds an HTTP create order request
// with url encoded values.
type CreateOrderRequest struct {
	UUID      string `json:"uuid"` // optional and used for idempotency on creation endpoint
	Name      string `json:"name"`
	Temp      string `json:"temp"`
	ShelfLife string `json:"shelfLife"`
	DecayRate string `json:"decayRate"`
}

// OrderJSON holds the order json from input.json.
type OrderJSON struct {
	Name      string  `json:"name"`
	Temp      string  `json:"temp"`
	ShelfLife int     `json:"shelfLife"`
	DecayRate float64 `json:"decayRate"`
}
