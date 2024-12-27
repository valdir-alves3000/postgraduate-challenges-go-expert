package entity

import "time"

type Report struct {
	TotalTime       time.Duration
	TotalRequests   int
	SuccessCount    int
	StatusCodeCount map[int]int
}
