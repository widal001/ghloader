package project

import (
	"encoding/json"
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

// GraphQL query to fetch all fields from the GitHub ProjectV2
const projectV2MetadataQuery = `
query($login: String!, $projectNumber: Int!) {
  user(login: $login) {
    projectV2(number: $projectNumber) {
	  id
	  title
      fields(first: 100) {
        nodes {
          ... on ProjectV2FieldCommon {
            id
            name
            type: dataType
          }
        }
      }
    }
  }
}`

// LoadProjectFields fetches all fields of a GitHub ProjectV2
func LoadProjectMetadata(login string, projectNumber int, projectType string) (*ProjectV2, error) {
	// Create a request body for the GraphQL query
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": projectV2MetadataQuery,
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
		}
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Map fields into a ProjectV2Info struct
	fieldsMap := make(map[string]ProjectV2Field)
	for _, node := range response.Data.User.ProjectV2.Fields.Nodes {
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
