package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/config"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/middleware"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/ratelimiter"
)

func TestRateLimitByIP(t *testing.T) {
	redisAddr := "localhost:6379"
	cfg := config.Config{
		RequestLimitIP: 2,
		BlockTime:      5,
	}
	store := ratelimiter.NewRedisStore(redisAddr)
	rl := ratelimiter.NewRateLimiter(store, cfg)

	handler := middleware.RateLimiterMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request OK"))
	}))

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	// Create the test server
	rr := httptest.NewRecorder()

	// First request
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Third request must be blocked
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Wait for the blocking time and try again
	time.Sleep(6 * time.Second)

	// After blocking, it should be allowed again
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRateLimitByAPIKey(t *testing.T) {
	redisAddr := "localhost:6379"
	store := ratelimiter.NewRedisStore(redisAddr)
	cfg := config.Config{
		RequestLimitToken: 3,
		BlockTime:         5,
	}
	rl := ratelimiter.NewRateLimiter(store, cfg)

	handler := middleware.RateLimiterMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request OK"))
	}))

	// Create a test request with API_KEY
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("API_KEY", "token123")

	// Create the test server
	rr := httptest.NewRecorder()

	// First request
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Third request
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Fourth request must be blocked
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Wait for the blocking time and try again
	time.Sleep(6 * time.Second)

	// After blocking, it should be allowed again
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
