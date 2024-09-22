package project

import (
	"encoding/json"
	"fmt"

	"github.com/widal001/ghloader/internal/graphql"
)

// =================================================
// Types
// =================================================

type FieldData struct {
	Name  string
	Value string
}

type UpdateProps struct {
	ItemId string
	Fields []FieldData
}

// =================================================
// Add or update project item
// =================================================

func (p *ProjectV2) UpsertProjectItem(data UpdateProps) error {
	// get the client
	client := graphql.NewClient()
	for _, field := range data.Fields {
		fmt.Printf("%s\n", field.Name)
		err := p.UpdateProjectItemField(*client, data.ItemId, field)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// =================================================
// Update project item
// =================================================

func (proj *ProjectV2) UpdateProjectItemField(
	client graphql.Client,
	itemId string,
	fieldData FieldData,
) error {
	// Load the query from the GraphQL file
	query, err := graphql.LoadQuery("projectV2Update/mutation.graphql")
	if err != nil {
		return fmt.Errorf("failed to load GraphQL query: %v", err)
	}

	// Retrieve the field from the project using its name
	field, ok := proj.Fields[fieldData.Name]
	if !ok {
		return fmt.Errorf("field %s not found", fieldData.Name)
	}

	// Set the value
	value, err := field.FormatUpdateValue(fieldData.Value)
	if err != nil {
		return err
	}

	// Create a request body for the GraphQL query
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"projectId": proj.Id,
			"itemId":    itemId,
			"fieldId":   field.ID,
			"value":     value,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make the update
	response, err := client.Post(requestBody)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to update field: %s", field.Name)
	} else {
		fmt.Printf("%v\n\n", string(response))
	}
	return nil
}
