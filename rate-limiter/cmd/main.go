package main

import (
	"fmt"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/config"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/ratelimiter"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/web"
)

func main() {
	cfg := config.LoadConfig()
	store := ratelimiter.NewRedisStore(fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
	rl := ratelimiter.NewRateLimiter(store, cfg)

	web.StartServer(rl)
}
