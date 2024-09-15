package graphql

import (
	"os"
	"path/filepath"
)

// LoadQuery loads a GraphQL query from a file
func LoadQuery(filename string) (string, error) {
	// Load the query from the path provided
	queryPath := filepath.Join("queries", filename)
	content, err := os.ReadFile(queryPath)
	if err != nil {
		return "", err
	}
	// Load any shared query fragments from the fragments file in the same directory
	fragmentPath := filepath.Join(filepath.Dir(queryPath), "fragments.graphql")
	fragments, err := os.ReadFile(fragmentPath)
	if err != nil {
		return "", err
	}
	content = append(content, fragments...)
	return string(content), nil
}
