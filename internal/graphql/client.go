package graphql

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const URL = "https://api.github.com/graphql"

// Client is a struct that represents the HTTP client used to interact with the GitHub GraphQL API.
type Client struct {
	token string
	URL   string
}

// NewClient creates a new GitHub GraphQL API client with the provided token.
func NewClient() *Client {
	// Load variables from a .env file in the working directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create a new client with that token
	return &Client{
		token: os.Getenv("GITHUB_TOKEN"),
		URL:   URL,
	}
}

// Post sends a POST request to the GitHub GraphQL API with the provided requestBody.
func (c *Client) Post(requestBody []byte) ([]byte, error) {
	// Create the HTTP POST request
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers for the request
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}
