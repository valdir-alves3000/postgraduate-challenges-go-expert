package main

import (
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/cli-stress-test/internal/cli"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/cli-stress-test/internal/usecase"
)

func main() {
	url, totalRequests, concurrency := cli.ParseFlags()
	if url == "" || totalRequests == 0 || concurrency == 0 {
		return
	}

	report := usecase.LoadTestReport(url, totalRequests, concurrency)
	cli.PrintReport(report)
}
