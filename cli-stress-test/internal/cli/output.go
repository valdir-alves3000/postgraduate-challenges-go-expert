package cli

import (
	"fmt"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/cli-stress-test/internal/entity"
)

func PrintReport(report entity.Report) {
	fmt.Println("\nLoad Test Report")
	fmt.Println("================")
	fmt.Printf("Total Time: %v\n", report.TotalTime)
	fmt.Printf("Total Requests: %d\n", report.TotalRequests)
	fmt.Printf("Successful Requests (200): %d\n", report.SuccessCount)
	fmt.Println("Status Code Distribution:")
	for code, count := range report.StatusCodeCount {
		fmt.Printf("  %d: %d\n", code, count)
	}
}
