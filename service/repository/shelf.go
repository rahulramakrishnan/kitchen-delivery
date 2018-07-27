package repository

import (
	"log"
	"math"

	"github.com/gomodule/redigo/redis"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/entity/exception"
	"github.com/pkg/errors"
)

// ShelfRepository is the shelf repository interface.
type ShelfRepository interface {
	AddOrder(order entity.Order) error
	PullOrder() (*entity.Order, error)
}

type shelfRepository struct {
	redis      redis.Conn
	shelfSpace map[string]int
}

// NewShelfRepository is a new order repository.
func NewShelfRepository(redisClient redis.Conn) ShelfRepository {
	shelfSpace := map[string]int{
		"hot":      15,
		"cold":     15,
		"frozen":   15,
		"overflow": 20,
	}
	return &shelfRepository{
		redis:      redisClient,
		shelfSpace: shelfSpace,
	}
}

// AddOrder adds an order to the right shelf.
func (s *shelfRepository) AddOrder(order entity.Order) error {
	// Check corresponding shelf space before attempting to add to shelf.
	shelfName := string(order.Temp) // ex: "hot", "cold", "frozen"
	count, err := s.redis.Do("SCARD", shelfName)
	if err != nil {
		return errors.Wrapf(
			exception.ErrDatabase, "redis failed to add order %+v err: %+v", order, err)
	}

	// If there is space on the corresponding shelf we add it w/ an expiration time.
	// if count < s.shelfSpace[shelfName] {
	// 	log.Printf("")
	// }

	log.Printf("count: %d, %T", count, count)
	return nil
}

// PullOrder pulls an order from the shelf.
func (s *shelfRepository) PullOrder() (*entity.Order, error) {
	return nil, nil
}

func (s *shelfRepository) getTTL(order entity.Order) int {
	// Calculate time to live in seconds based on formula.
	// Remember an order is waste after the "value" becomes zero.
	// This leads the formula to be reduced to:
	// => orderAge = shelfLife / (1 + decayRate)
	// We're given shelfLife and decayRate so we can solve for
	// how old an order can get before we consider it as waste.
	expirationTime := float64(order.ShelfLife) / (1.0 + order.DecayRate)
	ttl := int(math.Floor(expirationTime))
	return ttl
}
