package main

import (
	"fmt"
	"log"

	"github.com/widal001/ghloader/internal/project"
)

func main() {
	// Example GitHub project details
	projURL := "https://github.com/users/widal001/projects/3"
	issueURL := "https://github.com/widal001/project-demo/issues/52"

	// Load the fields from the ProjectV2
	proj, err := project.FromURL(projURL)
	if err != nil {
		log.Fatalf("Error loading project fields: %v", err)
	}

	// Add an issue to the project and get its ID
	itemId, err := proj.AddItemByURL(issueURL)
	if err != nil {
		log.Fatalf("Error adding issue to project: %v", err)
	}
	fmt.Printf("The project item ID is: %s\n\n", itemId)

	// Update the value of each field
	data := project.UpdateProps{
		ItemId: itemId,
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
