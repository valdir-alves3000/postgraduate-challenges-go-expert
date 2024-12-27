package usecase

import (
	"net/http"
	"sync"
	"time"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/cli-stress-test/internal/entity"
)

func LoadTestReport(url string, totalRequests, concurrency int) entity.Report {
	var wg sync.WaitGroup
	requestChannel := make(chan struct{}, concurrency)
	statusCodeCount := make(map[int]int)
	mu := sync.Mutex{}
	startTime := time.Now()
	var successCount int

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		requestChannel <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-requestChannel }()

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			defer mu.Unlock()

			statusCodeCount[resp.StatusCode]++
			if resp.StatusCode == 200 {
				successCount++
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)

	return entity.Report{
		TotalTime:       totalTime,
		TotalRequests:   totalRequests,
		SuccessCount:    successCount,
		StatusCodeCount: statusCodeCount,
	}
}
