package project

import (
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

// Shared project and field metadata returned by the GraphQL API
type projectV2Fragment struct {
	Id     string
	Title  string
	Fields struct {
		Nodes []struct {
			Id      string
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

// LoadProjectFields fetches all fields of a GitHub ProjectV2
func LoadProjectMetadata(login string, projectNumber int, projectType string) (*ProjectV2, error) {
	// Find and load the correct .graphql file based on project type
	var queryFile string
	switch projectType {
	case "orgs":
		queryFile = "queryOrg.graphql"
	default:
		queryFile = "queryUser.graphql"
	}
	// Create the query
	query := graphql.Query{
		Options: graphql.QueryOptions{
			QueryDir:      "queries/projectV2Metadata",
			QueryPath:     queryFile,
			FragmentPaths: []string{"fragments.graphql"},
		},
		Vars: map[string]interface{}{
			"login":         login,
			"projectNumber": projectNumber,
		},
	}
	// Create a struct to parse the response body
	var response struct {
		Data struct {
			Login struct {
				ProjectV2 projectV2Fragment
			}
		}
	}
	// Post the query
	err := query.Post(&response)
	if err != nil {
		fmt.Printf("Failed to post query: %s\n", err)
		return nil, err
	}
	// Map each field into a ProjectV2Field struct
	fieldsMap := make(map[string]ProjectV2Field)
	for _, node := range response.Data.Login.ProjectV2.Fields.Nodes {
		fieldsMap[node.Name] = ProjectV2Field{
			Id:         node.Id,
			Name:       node.Name,
			Type:       node.Type,
			Options:    node.Options,
			Iterations: node.Configuration.Iterations,
		}
	}

	// Return the ProjectV2 struct with its parsed metadata and fields
	return &ProjectV2{
		Id:     response.Data.Login.ProjectV2.Id,
		Title:  response.Data.Login.ProjectV2.Title,
		Fields: fieldsMap,
	}, nil
}
