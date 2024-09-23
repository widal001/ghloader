package project

import (
	"fmt"
	"sync"

	"github.com/widal001/ghloader/internal/graphql"
)

// =================================================
// Types
// =================================================

type FieldData struct {
	Name  string
	Value string
}

type ItemData struct {
	ItemURL string
	Fields  []FieldData
}

// =================================================
// Add or update project item
// =================================================

func (p *ProjectV2) UpsertItem(data ItemData) error {
	// Add the item to the project and get its ID
	itemId, err := p.AddItemByURL(data.ItemURL)
	if err != nil {
		return fmt.Errorf("failed to add item with URL %s: %w", data.ItemURL, err)
	}

	// Create the wait group and the error channel to synchronize errors
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(data.Fields)) // To capture errors from goroutines

	// Update fields
	for _, field := range data.Fields {
		wg.Add(1) // Increment the WaitGroup counter for each field
		go func(f FieldData) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			err := p.UpdateItemField(itemId, f)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to update field %s: %w", f.Name, err)
			}
		}(field) // Pass field as an argument to avoid data races
	}

	// Wait for all goroutines to finish and close the channel
	wg.Wait()
	close(errorsChan)

	// Check if any errors were returned
	for err := range errorsChan {
		if err != nil {
			fmt.Println(err) // Handle or log the error
		}
	}
	return nil
}

// =================================================
// Update project item
// =================================================

func (proj *ProjectV2) UpdateItemField(
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
	// Catch or return the errors
	if err != nil {
		return err
	} else {
		fmt.Printf("%v\n", response)
	}
	return nil
}
