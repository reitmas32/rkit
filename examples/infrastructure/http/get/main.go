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

// Character represents a character from the Rick and Morty API
type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Origin  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Image   string `json:"image"`
	URL     string `json:"url"`
	Created string `json:"created"`
}

func PrintFormattedResponse(resp *corehttp.TypedResponse[Character]) {
	// use json.Marshal to print the response body in a formatted way
	json, err := json.MarshalIndent(resp.Body, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling response body: %v\n", err)
	}
	fmt.Printf("Response Body (typed):\n")
	fmt.Printf("%s\n", json)
}

func main() {
	fmt.Println("=== HTTP Client GET Example ===")

	// Create context
	ctx := customctx.New(context.Background())

	// Create HTTP client with default config
	config := http.DefaultConfig()
	client := http.NewClient(config)

	// GET request with typed JSON response
	fmt.Println("GET request with typed JSON response:")
	fmt.Println("GET https://rickandmortyapi.com/api/character/1")
	fmt.Println()

	// Make GET request with typed response
	resp, err := corehttp.GetTyped[Character](
		client,
		ctx,
		"https://rickandmortyapi.com/api/character/1",
		corehttp.WithHeader("Accept", "application/json"),
	)

	if err != nil {
		log.Fatalf("Error making GET request: %v\n", err)
	}

	// Access typed response
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Request Time: %s\n", resp.RequestTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Response Time: %s\n", resp.ResponseTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Printf("Is Success (default 2xx): %v\n", resp.IsSuccess())
	fmt.Println()

	// Example: Configure expected status code
	resp.SetExpectedStatusCode(200)
	fmt.Printf("Is Success (expected 200): %v\n", resp.IsSuccess())
	fmt.Println()

	// Example: Configure status code range
	resp.SetSuccessStatusCodeRange(200, 299)
	resp.ExpectedStatusCode = nil // Clear expected code to use range
	fmt.Printf("Is Success (range 200-299): %v\n", resp.IsSuccess())
	fmt.Println()

	PrintFormattedResponse(resp)

	fmt.Println("✓ Example completed!")
}
