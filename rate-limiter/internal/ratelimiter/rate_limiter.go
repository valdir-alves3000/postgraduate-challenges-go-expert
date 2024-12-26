package ratelimiter

import (
	"fmt"
	"sync"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/config"
)

type RateLimiter struct {
	Store  Store
	Config config.Config
	mu     sync.Mutex
}

func NewRateLimiter(store Store, config config.Config) *RateLimiter {
	return &RateLimiter{Store: store, Config: config}
}

func (rl *RateLimiter) CheckLimit(key string, limit, blockTime int) (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	blocked, err := rl.Store.IsBlocked(key)
	if err != nil {
		return false, fmt.Errorf("error checking block status: %v", err)
	}
	if blocked {
		return false, fmt.Errorf("you are blocked")
	}

	count, err := rl.Store.Increment(key, 1)
	if err != nil {
		return false, fmt.Errorf("error incrementing request count: %v", err)
	}

	if count > limit {
		err = rl.Store.Block(key, blockTime)

		if err != nil {
			return false, fmt.Errorf("error blocking the key: %v", err)
		}
		return false, fmt.Errorf("rate limit exceeded")
	}
	return true, nil
}
