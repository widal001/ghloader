package main

import (
	"fmt"
	"log"

	"github.com/widal001/ghloader/internal/project"
)

func main() {
	// Example GitHub project details
	url := "https://github.com/users/widal001/projects/3"

	// Load the fields from the ProjectV2
	proj, err := project.FromURL(url)
	if err != nil {
		log.Fatalf("Error loading project fields: %v", err)
	}

	// Print the type of each field
	for fieldName := range proj.Fields {
		field, ok := proj.Fields[fieldName]
		if !ok {
			log.Fatalf("Error getting field type: %v", err)
		}
		fmt.Printf("The ID of '%s' field is: %s\n", fieldName, field.ID)
	}

	// Update the value of each field
	data := project.UpdateProps{
		ItemId: "PVTI_lAHOAUXHu84AR8K0zgJRro4",
		Fields: []project.FieldData{
			{
				Name:  "Status",
				Value: "✅ Done",
			},
			{
				Name:  "Sprint",
				Value: "Sprint 12",
			},
			{
				Name:  "Story Points",
				Value: "3",
			},
			{
				Name:  "Due date",
				Value: "2024-05-15",
			},
		},
	}
	proj.UpsertProjectItem(data)
}
