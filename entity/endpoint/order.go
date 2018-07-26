package endpoint

// CreateOrderRequest is create order request.
type CreateOrderRequest struct {
	UUID      string `json:"uuid"` // optional and used for idempotency
	Name      string `json:"name"`
	Temp      string `json:"temp"`
	ShelfLife string `json:"shelfLife"`
	DecayRate string `json:"decayRate"`
}
