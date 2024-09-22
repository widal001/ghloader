package project

import (
	"encoding/json"
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
	// Get the client
	client := graphql.NewClient()
	// Load the query
	query, err := graphql.LoadQuery("projectV2AddItem/mutation.graphql")
	if err != nil {
		return "", err
	}
	// Create the requestBody
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"projectId": proj.Id,
			"issueId":   nodeId,
		},
	})
	if err != nil {
		return "", nil
	}
	// Make the API request
	responseBody, err := client.Post(requestBody)
	if err != nil {
		return "", err
	}
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
	// Unmarshal the JSON response into the struct
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", err
	}
	return response.Data.Update.Item.Id, nil

}

func GetItemByURL(url string) (string, error) {
	// Get the client
	client := graphql.NewClient()
	// Load the query
	query, err := graphql.LoadQuery("projectV2AddItem/getItemByURL.graphql")
	if err != nil {
		return "", err
	}
	// Create the requestBody
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": map[string]interface{}{"issueUrl": url},
	})
	if err != nil {
		return "", err
	}
	// Make the API request
	responseBody, err := client.Post(requestBody)
	if err != nil {
		return "", err
	}
	// Declare a struct that matches the expected response JSON
	var response struct {
		Data struct {
			Resource struct {
				Id string
			}
		}
	}
	// Unmarshal the JSON response into the struct
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", err
	}
	return response.Data.Resource.Id, nil
}
