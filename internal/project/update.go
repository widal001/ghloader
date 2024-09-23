package project

import (
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
	// Create the query
	query := graphql.Query{
		Options: graphql.QueryOptions{
			QueryDir:  "queries/projectV2Update",
			QueryPath: "mutation.graphql",
		},
		Vars: map[string]interface{}{
			"projectId": proj.Id,
			"itemId":    itemId,
			"fieldId":   field.ID,
			"value":     value,
		},
	}
	// Post the query
	var response map[string]interface{}
	err = query.Post(&response)
	clientErrors, queryErr := response["errors"]
	// Catch or return the errors
	if err != nil {
		fmt.Printf("Failed to post query: %s\n", err)
		return err
	} else if queryErr {
		return fmt.Errorf("graphql client return the following errors: %v", clientErrors)
	} else {
		fmt.Printf("%v", response)
	}
	return nil
}
