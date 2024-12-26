package middleware

import (
	"log"
	"net/http"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/ratelimiter"
)

func RateLimiterMiddleware(rl *ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			token := r.Header.Get("API_KEY")

			limit := rl.Config.RequestLimitIP
			if token != "" {
				key = token
				limit = rl.Config.RequestLimitToken
			}

			if ok, err := rl.CheckLimit(key, limit, rl.Config.BlockTime); !ok {
				if err != nil {
					log.Printf("Error while checking rate limit: %v", err)
				}

				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
