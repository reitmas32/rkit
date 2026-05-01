package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/reitmas32/rkit/core/customctx"
	corehttp "github.com/reitmas32/rkit/core/http"
	"github.com/reitmas32/rkit/core/logger"
	"github.com/reitmas32/rkit/infrastructure/http"
)

// Post represents a post from JSONPlaceholder API
type PostResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	fmt.Println("=== HTTP Client POST Example with Logger ===")

	// Create context
	ctx := customctx.New(context.Background())

	// Create a simple logger
	simpleLogger := logger.NewSimpleLogger("debug")

	// Create HTTP client with logger enabled
	config := http.DefaultConfig()
	config.Logger = simpleLogger
	// config.DisableLogging = true // Uncomment to disable logs even with logger set
	client := http.NewClient(config)

	// Prepare request body
	createReq := CreatePostRequest{
		UserID: 1,
		Title:  "Mi primer post con base-kit",
		Body:   "Este es el contenido del post creado usando el cliente HTTP del base-kit.",
	}

	fmt.Println("Making POST request with logging enabled...")
	fmt.Println()

	// Make POST request with typed response
	resp, err := corehttp.PostTyped[CreatePostRequest, PostResponse](
		client,
		ctx,
		"https://jsonplaceholder.typicode.com/posts",
		createReq,
		corehttp.WithContentType("application/json"),
		corehttp.WithHeader("Accept", "application/json"),
	)

	if err != nil {
		log.Fatalf("Error making POST request: %v\n", err)
	}

	// Access typed response
	fmt.Printf("\nResponse Details:\n")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Println()

	// Print formatted response
	fmt.Println("Response Body (typed):")
	jsonResp, err := json.MarshalIndent(resp.Body, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling response body: %v\n", err)
	}
	fmt.Printf("%s\n", jsonResp)
	fmt.Println()

	fmt.Println("✓ Example completed!")
	fmt.Println("\nNote: Check the logs above to see HTTP request/response logging.")
}
