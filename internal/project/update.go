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
// Load item data
// =================================================

func ItemsFromCSV(data []map[string]string, urlField string) ([]ItemData, error) {
	var items []ItemData
	for _, row := range data {
		// Get the URL for the item to add
		url, ok := row[urlField]
		if !ok {
			return items, fmt.Errorf("missing a URL field")
		}
		// Create the item
		item := ItemData{ItemURL: url}
		// Populate the item's fields from the other columns
		for field, val := range row {
			if field == urlField {
				continue
			}
			item.Fields = append(item.Fields, FieldData{Name: field, Value: val})
		}
		items = append(items, item)
	}
	return items, nil
}

// =================================================
// Add or update multiple project items
// =================================================

func (p *ProjectV2) BatchUpsertItems(items []ItemData) ([]string, bool) {
	var updated []string
	var ok = true
	for _, item := range items {
		// Increment the WaitGroup counter for each item
		err := p.UpsertItem(item)
		if err != nil {
			fmt.Printf("failed to update item %s: %s", item.ItemURL, err.Error())
			ok = false
		}
		updated = append(updated, item.ItemURL)
	}
	return updated, ok
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

	// Create the wait group and the error channel to capture errors from goroutine
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(data.Fields))

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
// Update project item metadata
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
	}
	return nil
}
