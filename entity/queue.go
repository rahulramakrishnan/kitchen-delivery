package entity

import "github.com/gomodule/redigo/redis"

// Queues holds Redis queues.
type Queues struct {
	// Kitchen categorize and store incoming orders.
	Order Queue
	// We can extend this to include more queues
	// as our Kitchen Delivery system expands.
}

// Queue holds queue name and redis connection.
type Queue struct {
	Name string
	Pool *redis.Pool
}
