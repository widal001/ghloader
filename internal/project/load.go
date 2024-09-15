package project

import (
	"encoding/json"
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

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

	type projectV2Fragment struct {
		ProjectV2 struct {
			Id     string
			Title  string
			Fields struct {
				Nodes []struct {
					ID   string
					Name string
					Type string
				}
			}
		}
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

	var projectData projectV2Fragment
	switch projectType {
	case "users":
		projectData = response.Data.User.projectV2Fragment
	case "orgs":
		projectData = response.Data.Org.projectV2Fragment
	}

	// Map fields into a ProjectV2Info struct
	fieldsMap := make(map[string]ProjectV2Field)
	for _, node := range projectData.ProjectV2.Fields.Nodes {
		fieldsMap[node.Name] = ProjectV2Field{
			ID:   node.ID,
			Name: node.Name,
			Type: node.Type,
		}
	}

	return &ProjectV2{
		Id:     response.Data.User.ProjectV2.Id,
		Title:  response.Data.User.ProjectV2.Title,
		Fields: fieldsMap,
	}, nil
}
