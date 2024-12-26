package ratelimiter

type Store interface {
	Get(key string) (int, error)
	Increment(key string, expiration int) (int, error)
	Block(key string, duration int) error
	IsBlocked(key string) (bool, error)
	ListKeys(pattern string) ([]string, error)
}
