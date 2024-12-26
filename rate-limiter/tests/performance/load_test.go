package performance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/config"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/middleware"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/ratelimiter"
)

func setupTest(requestLimitIP, requestLimitToken, blockTime int) (*ratelimiter.RateLimiter, http.Handler) {
	// Configuração do Redis para o teste
	redisStore := ratelimiter.NewRedisStore("localhost:6379")
	cfg := config.Config{
		RequestLimitIP:    requestLimitIP,
		RequestLimitToken: requestLimitToken,
		BlockTime:         60,
	}

	// Criar o RateLimiter
	rl := ratelimiter.NewRateLimiter(redisStore, cfg)

	// Criar o handler com o middleware de rate limiting
	handler := middleware.RateLimiterMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return rl, handler
}

func makeRequest(handler http.Handler, ip, token string) *httptest.ResponseRecorder {
	// Criar uma solicitação de teste
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	if ip != "" {
		req.Header.Set("X-Real-IP", ip) // Definindo o IP no cabeçalho
	}
	if token != "" {
		req.Header.Set("API_KEY", token) // Definindo a chave de API no cabeçalho
	}

	// Gravar a resposta
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func BenchmarkRateLimiterMiddleware(b *testing.B) {
	// redisAddr := "localhost:6379"
	// store := ratelimiter.NewRedisStore(redisAddr)

	// cfg := config.Config{
	// 	RequestLimitIP:    5,
	// 	RequestLimitToken: 10,
	// 	BlockTime:         10,
	// }
	// rl := ratelimiter.NewRateLimiter(store, cfg)

	// handler := middleware.RateLimiterMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Request OK"))
	// }))

	_, handler := setupTest(5, 10, 10)

	server := http.Server{
		Addr:    ":8081",
		Handler: handler,
	}
	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(2 * time.Second)

	// Simulated Concurrent Requests
	const totalClients = 1000
	const requestsPerClient = 20
	var successCount int32
	var blockedCount int32

	var wg sync.WaitGroup
	wg.Add(totalClients)

	client := http.Client{}
	for i := 0; i < totalClients; i++ {
		go func(clientID int) {
			defer wg.Done()
			for j := 0; j < requestsPerClient; j++ {
				req, err := http.NewRequest("GET", "http://localhost:8081/test", nil)
				assert.NoError(b, err)

				// Add a unique token for half of the customers
				if clientID%2 == 0 {
					req.Header.Set("API_KEY", "TOKEN"+fmt.Sprint(clientID))
				}

				resp, err := client.Do(req)
				assert.NoError(b, err)

				if resp.StatusCode == http.StatusOK {
					successCount++
				} else if resp.StatusCode == http.StatusTooManyRequests {
					blockedCount++
				}

				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	b.Logf("Total successful requests: %d", successCount)
	b.Logf("Total blocked requests: %d", blockedCount)
}

func TestStreesMonitoring(t *testing.T) {
	rl, handler := setupTest(5, 5, 60)

	getMemStats := func() runTime.MemStats {
		var m runTime.MemStats
		runTime.ReadMemStats(&m)
		return m
	}

	getGoroutineCount := func() int {
		return runtime.NumGoroutine()
	}

	getCPUUsage := func() float64 {
		start := time.Now()
		for i := 0; i < 100000; i++ {
			_ = i * i
		}
		elapsed := time.Since(start).Seconds()
		return elapsed
	}

	getRedisPoolStats := func() *redis.PoolStats {
		client := rl.Store.(*rateLimiter.RedisStore).client
		return client.PoolStats()
	}

	initialMemStats := getMemStats()
	initialGoroutines := getGoroutineCount()
	initialCPUUsage := getCPUUsage()
	initialRedisPoolStats := getRedisPoolStats()

	for i := 0; i < 1000; i++ {
		rr := makeRequest(handler, "127.0.0.1", "")
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	finalMemStats := getMemStats()
	finalGoroutines := getGoroutineCount()
	finalCPUUsage := getCPUUsage()
	finalRedisPoolStats := getRedisPoolStats()

	assert.True(t, finalMemStats.Alloc <= initialMemStats.Alloc+10*1024*1024, "Memory usage increased too much")
	assert.True(t, finalGoroutines <= initialGoroutines+50, "Too many goroutines created")
	assert.True(t, finalCPUUsage <= initialCPUUsage*1.5, "CPU usage increased too much")
	assert.True(t, finalRedisPoolStats.TotalConns <= initialRedisPoolStats.TotalConns+100, "Redis pool connections exceeded limit")

}
