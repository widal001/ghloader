package project

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/widal001/ghloader/internal/graphql"
	"github.com/widal001/ghloader/internal/utils"
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

// Generate a list of `ItemData` entries from the contents of a CSV file
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

// Upsert multiple project items concurrently
func (p *ProjectV2) BatchUpsertItems(items []ItemData) ([]string, bool) {
	// Create a wait group and channels to synchronize goroutine output
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(items))
	updatedChan := make(chan string, len(items))

	// Launch goroutines to update items
	for _, item := range items {
		wg.Add(1)
		go func(i ItemData) {
			defer wg.Done()
			err := p.UpsertItem("templates", item)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to update field %s: %w", i.ItemURL, err)
				return
			}
			updatedChan <- i.ItemURL
		}(item)
	}

	// Close channels when all goroutines have finished
	go func() {
		wg.Wait()
		close(errorsChan)
		close(updatedChan)
	}()

	// Collect errors and updated items
	var updated []string
	var hasError bool
	for err := range errorsChan {
		hasError = true
		fmt.Println(err)
	}
	for item := range updatedChan {
		updated = append(updated, item)
	}

	// Return updated items and success flag
	return updated, !hasError && len(updated) == len(items)
}

// =================================================
// Add or update project item
// =================================================

// Add an issue to a project or update its metadata
func (p *ProjectV2) UpsertItem(templDir string, data ItemData) error {
	// Add the item to the project and get its ID
	itemId, err := p.AddItemByURL(data.ItemURL)
	if err != nil {
		return fmt.Errorf("failed to add item with URL %s: %w", data.ItemURL, err)
	}
	// Create the query string from the input fields
	queryStr, err := p.GenerateUpdateQuery(templDir, data.Fields)
	if err != nil {
		return fmt.Errorf("failed to generate update query: %w", err)
	}
	// Generate the query
	query := graphql.Query{
		QueryStr: queryStr,
		Vars: map[string]interface{}{
			"projectId": p.Id,
			"itemId":    itemId,
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

// =================================================
// Generate update query
// =================================================

// Generate the GraphQL query to update multiple fields in a single API call
func (p *ProjectV2) GenerateUpdateQuery(
	tmplDir string,
	fields []FieldData,
) (string, error) {
	// Load the query template
	tmpl, err := utils.LoadTemplate(tmplDir, "updateProjectField.graphql")
	if err != nil {
		return "", err
	}
	// Format the field data for the template
	var tmplFields []fieldProps
	for _, f := range fields {
		// Retrieve the field from the project using its name
		field, ok := p.Fields[f.Name]
		if !ok {
			fmt.Printf("field %s not found", f.Name)
			continue
		}
		// Format the field for the query template
		fieldData, err := field.formatFieldProps(f.Value)
		if err != nil {
			fmt.Printf("error formatting field %s: %s", f.Name, err)
			continue
		}
		tmplFields = append(tmplFields, fieldData)
	}
	// Render the query template and return it as a string
	var query bytes.Buffer
	tmpl.Execute(&query, tmplFields)
	return query.String(), nil
}
