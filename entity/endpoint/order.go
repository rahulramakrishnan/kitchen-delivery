package endpoint

// CreateOrderRequest is create order request.
type CreateOrderRequest struct {
	Name      string  `json:"name"`
	Temp      string  `json:"temp"`
	ShelfLife int     `json:"shelfLife"`
	DecayRate float64 `json:"decayRate"`
}
