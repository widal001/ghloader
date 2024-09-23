package project

import (
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

func (proj *ProjectV2) AddItemByURL(url string) (string, error) {
	// Get the Node ID of the item from its url
	nodeId, err := GetItemByURL(url)
	if err != nil {
		return "", err
	}
	// Add the item to the project
	itemId, err := proj.AddItemByNodeId(nodeId)
	if err != nil {
		return "", err
	}
	return itemId, nil

}

func (proj *ProjectV2) AddItemByNodeId(nodeId string) (string, error) {
	// Declare a struct that matches the expected response JSON
	var response struct {
		Data struct {
			Update struct {
				Item struct {
					Id string
				}
			} `json:"addProjectV2ItemById"`
		}
	}
	// Create the query
	query := graphql.Query{
		Options: graphql.QueryOptions{
			QueryDir:  "queries/projectV2AddItem",
			QueryPath: "mutation.graphql",
		},
		Vars: map[string]interface{}{
			"issueId":   nodeId,
			"projectId": proj.Id,
		},
	}
	// Post the query
	err := query.Post(&response)
	if err != nil {
		fmt.Printf("Failed to post query: %s\n", err)
		return "", err
	}
	return response.Data.Update.Item.Id, nil

}

func GetItemByURL(url string) (string, error) {
	// Declare a struct that matches the expected response JSON
	var response struct {
		Data struct {
			Resource struct {
				Id string
			}
		}
	}
	// Create the query
	query := graphql.Query{
		Options: graphql.QueryOptions{
			QueryDir:  "queries/projectV2AddItem",
			QueryPath: "getItemByUrl.graphql",
		},
		Vars: map[string]interface{}{"issueUrl": url},
	}
	// Post the query
	err := query.Post(&response)
	if err != nil {
		fmt.Printf("Failed to post query: %s\n", err)
		return "", err
	}
	return response.Data.Resource.Id, nil
}
