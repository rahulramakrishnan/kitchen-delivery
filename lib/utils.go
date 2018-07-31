package lib

import (
	"fmt"
	"strings"

	"github.com/kitchen-delivery/entity"
)

// StringifyShelves normalizes shelves.
func StringifyShelves(shelves map[entity.ShelfType][]*entity.ShelfOrder) string {
	var prettifiedShelves []string

	for shelfType, shelfOrders := range shelves {
		var orderStrs []string
		for _, shelfOrder := range shelfOrders {
			orderStrs = append(orderStrs, shelfOrder.Order.Name)
		}

		ordersOnShelf := strings.Join(orderStrs, ", ")
		shelf := fmt.Sprintf("%s: %s", shelfType, ordersOnShelf)

		prettifiedShelves = append(prettifiedShelves, shelf)
	}

	return fmt.Sprintf("%s", strings.Join(prettifiedShelves, "\n\n"))
}
