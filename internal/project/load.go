package project

import (
	"encoding/json"
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

// Shared project and field metadata returned by the GraphQL API
type projectV2Fragment struct {
	ProjectV2 struct {
		Id     string
		Title  string
		Fields struct {
			Nodes []struct {
				ID      string
				Name    string
				Type    string
				Options []struct {
					Id   string
					Name string
				}
				Configuration struct {
					Iterations []struct {
						Id    string
						Title string
					}
				}
			}
		}
	}
}

// LoadProjectFields fetches all fields of a GitHub ProjectV2
func LoadProjectMetadata(login string, projectNumber int, projectType string) (*ProjectV2, error) {
	// Find and load the correct .graphql file based on project type
	var queryFile string
	switch projectType {
	case "orgs":
		queryFile = "projectV2Metadata/queryOrg.graphql"
	case "users":
		queryFile = "ProjectV2Metadata/queryUser.graphql"
	}
	query, err := graphql.LoadQuery(queryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load GraphQL query: %v", err)
	}

	// Create a request body for the GraphQL query
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"login":         login,
			"projectNumber": projectNumber,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send request to GitHub GraphQL API
	client := graphql.NewClient()
	responseBody, err := client.Post(requestBody)
	if err != nil {
		return nil, err
	}

	// Parse the response JSON
	var response struct {
		Data struct {
			User struct {
				projectV2Fragment
			}
			Org struct {
				projectV2Fragment
			}
		}
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Pluck the project metadata based on the project type
	var projectMetadata projectV2Fragment
	switch projectType {
	case "users":
		projectMetadata = response.Data.User.projectV2Fragment
	case "orgs":
		projectMetadata = response.Data.Org.projectV2Fragment
	}

	// Map each field into a ProjectV2Field struct
	fieldsMap := make(map[string]ProjectV2Field)
	for _, node := range projectMetadata.ProjectV2.Fields.Nodes {
		fieldsMap[node.Name] = ProjectV2Field{
			ID:         node.ID,
			Name:       node.Name,
			Type:       node.Type,
			Options:    node.Options,
			Iterations: node.Configuration.Iterations,
		}
	}

	// Return the ProjectV2 struct with its parsed metadata and fields
	return &ProjectV2{
		Id:     response.Data.User.ProjectV2.Id,
		Title:  response.Data.User.ProjectV2.Title,
		Fields: fieldsMap,
	}, nil
}
