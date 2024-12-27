package cli

import (
	"flag"
	"fmt"
	"os"
)

func ParseFlags() (string, int, int) {
	url := flag.String("url", "", "URL of the service to test")
	totalRequests := flag.Int("requests", 1, "Total number of requests")
	concurrency := flag.Int("concurrency", 1, "Number of simultaneous calls")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: URL must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	if *totalRequests <= 0 || *concurrency <= 0 {
		fmt.Println("Error: Requests and concurrency must be positive integers.")
		os.Exit(1)
	}

	return *url, *totalRequests, *concurrency
}
