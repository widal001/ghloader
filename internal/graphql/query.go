package graphql

import (
	"fmt"
	"os"
	"path/filepath"
)

type QueryOptions struct {
	FragmentPaths []string
	QueryDir      string
}

// Option is a functional option for configuring the query loader
type Option func(*QueryOptions)

// WithFragment adds a fragment file to be loaded with the query
func WithFragment(fragmentPath string) Option {
	return func(opts *QueryOptions) {
		opts.FragmentPaths = append(opts.FragmentPaths, fragmentPath)
	}
}

// WithQueryDir allows specifying a custom query directory
func WithQueryDir(dir string) Option {
	return func(opts *QueryOptions) {
		opts.QueryDir = dir
	}
}

// LoadQuery loads a GraphQL query from a file with the provided options
func LoadQuery(filename string, options ...Option) (string, error) {
	// Default options
	opts := &QueryOptions{
		QueryDir: "queries", // default directory
	}

	// Apply functional options
	for _, option := range options {
		option(opts)
	}

	// Load the main query file
	queryPath := filepath.Join(opts.QueryDir, filename)
	content, err := os.ReadFile(queryPath)
	if err != nil {
		return "", fmt.Errorf("failed to read query file: %v", err)
	}

	// Load additional fragments if provided
	for _, fragmentFile := range opts.FragmentPaths {
		fragmentPath := filepath.Join(opts.QueryDir, fragmentFile)
		fragmentContent, err := os.ReadFile(fragmentPath)
		if err != nil {
			return "", fmt.Errorf("failed to read fragment file: %v", err)
		}
		content = append(content, fragmentContent...)
	}

	return string(content), nil
}
