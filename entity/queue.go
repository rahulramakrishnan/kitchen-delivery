package entity

// Queues holds local queues.
type Queues struct {
	// Kitchen categorize and store incoming orders.
	OrderQueue chan *Order
	// Kitchen pulls off finished orders and places them on shelves
	// when there is space.
	ShelfQueue chan *Order
}
