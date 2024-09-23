package graphql

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Query struct {
	Options QueryOptions
	Vars    map[string]interface{}
}

type QueryOptions struct {
	QueryDir      string
	QueryPath     string
	FragmentPaths []string
}

func (q *Query) LoadFromFile() (string, error) {
	// Load the main query file
	queryPath := filepath.Join(q.Options.QueryDir, q.Options.QueryPath)
	content, err := os.ReadFile(queryPath)
	if err != nil {
		return "", fmt.Errorf("failed to read query file: %v", err)
	}

	// Load additional fragments if provided
	for _, fragmentFile := range q.Options.FragmentPaths {
		fragmentPath := filepath.Join(q.Options.QueryDir, fragmentFile)
		fragmentContent, err := os.ReadFile(fragmentPath)
		if err != nil {
			return "", fmt.Errorf("failed to read fragment file: %v", err)
		}
		content = append(content, fragmentContent...)
	}

	return string(content), nil
}

func (q *Query) Post(response interface{}) error {
	// Get the client
	client := NewClient()
	// Load the query
	query, err := q.LoadFromFile()
	if err != nil {
		return err
	}
	// Create the requestBody
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": q.Vars,
	})
	if err != nil {
		return nil
	}
	// Make the API request
	responseBody, err := client.Post(requestBody)
	if err != nil {
		return err
	}
	// Check for errors
	var errors struct {
		Errors interface{}
	}
	json.Unmarshal([]byte(responseBody), &errors)
	if errors.Errors != nil {
		return fmt.Errorf("GitHub API raised the following errors: %s", errors.Errors)
	}
	// Unmarshal the JSON response into the struct
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}
	return nil
}
