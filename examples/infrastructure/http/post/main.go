package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/reitmas32/rkit/core/customctx"
	corehttp "github.com/reitmas32/rkit/core/http"
	"github.com/reitmas32/rkit/infrastructure/http"
)

// Post represents a post from JSONPlaceholder API
type PostResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (p *PostResponse) PrintFormattedResponse() {
	// use json.Marshal to print the response body in a formatted way
	json, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling response body: %v\n", err)
	}
	fmt.Printf("%s\n", json)
	fmt.Println()
}

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	fmt.Println("=== HTTP Client POST Example ===")

	// Create context
	ctx := customctx.New(context.Background())

	// Create HTTP client with default config
	config := http.DefaultConfig()
	client := http.NewClient(config)

	// Prepare request body (directly as object, no manual marshaling needed)
	createReq := CreatePostRequest{
		UserID: 1,
		Title:  "Mi primer post con base-kit",
		Body:   "Este es el contenido del post creado usando el cliente HTTP del base-kit.",
	}

	// POST request with typed JSON response
	fmt.Println("POST request to create a new post:")
	fmt.Println("POST https://jsonplaceholder.typicode.com/posts")
	fmt.Println("Note: The request body is automatically JSON-marshaled")
	fmt.Println()

	// Make POST request with typed response
	// TRequest = CreatePostRequest, TResponse = Post
	resp, err := corehttp.PostTyped[CreatePostRequest, PostResponse](
		client,
		ctx,
		"https://jsonplaceholder.typicode.com/posts",
		createReq, // Pass the object directly, not bytes
		corehttp.WithContentType("application/json"),
		corehttp.WithHeader("Accept", "application/json"),
	)

	if err != nil {
		log.Fatalf("Error making POST request: %v\n", err)
	}

	// Access typed response
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Request Time: %s\n", resp.RequestTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Response Time: %s\n", resp.ResponseTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Printf("Is Success: %v\n", resp.IsSuccess())
	fmt.Println()

	// Print formatted response
	fmt.Println("Response Body (typed):")
	resp.Body.PrintFormattedResponse()

	fmt.Println("✓ Example completed!")
	fmt.Println("\nNote: JSONPlaceholder is a fake REST API for testing.")
	fmt.Println("      The post was not actually created, but the API returns a simulated response.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
