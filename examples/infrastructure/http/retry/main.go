package main

import (
	"context"
	"fmt"
	"log"

	"github.com/reitmas32/rkit/core/customctx"
	corehttp "github.com/reitmas32/rkit/core/http"
	"github.com/reitmas32/rkit/core/logger"
	"github.com/reitmas32/rkit/infrastructure/http"
)

func main() {
	fmt.Println("=== HTTP Client Retry Example ===")

	// Create context
	ctx := customctx.New(context.Background())

	// Create a simple logger to see retry attempts
	simpleLogger := logger.NewSimpleLogger("debug")

	// Create HTTP client with retry configuration
	config := http.DefaultConfig()
	config.Logger = simpleLogger
	config.MaxRetries = 3                                                  // Retry up to 3 times (4 total attempts)
	config.RetryDelay = 200                                                // 200ms delay between retries
	config.RetryableStatusCodes = []int{429, 500, 502, 503, 504}           // Retry on these status codes
	config.RetryableMethods = []string{"GET", "HEAD", "OPTIONS", "DELETE"} // Only retry idempotent methods

	client := http.NewClient(config)

	fmt.Println("Making GET request with retries configured...")
	fmt.Println("MaxRetries: 3")
	fmt.Println("RetryDelay: 200ms")
	fmt.Println("RetryableStatusCodes: 429, 500, 502, 503, 504")
	fmt.Println("RetryableMethods: GET, HEAD, OPTIONS, DELETE")
	fmt.Println()

	// Make GET request - this will retry if it gets a retryable status code or network error
	// Note: In real scenarios, you'd use your actual API endpoint
	// This example shows the retry configuration, but the endpoint should work on first try
	resp, err := corehttp.GetTyped[map[string]interface{}](
		client,
		ctx,
		"https://jsonplaceholder.typicode.com/posts/1", // This endpoint should work on first try
		corehttp.WithHeader("Accept", "application/json"),
	)

	if err != nil {
		log.Printf("Error after retries: %v\n", err)
		return
	}

	fmt.Printf("\nFinal Response:\n")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Total Duration: %v\n", resp.Duration)
	fmt.Println()

	fmt.Println("Note: Check the logs above to see retry attempts.")
	fmt.Println("      The request will retry up to 3 times if it receives a retryable status code.")
}
