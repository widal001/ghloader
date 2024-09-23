package main

import (
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

	// Update the value of each field
	data := project.ItemData{
		ItemURL: issueURL,
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
	proj.UpsertItem(data)
}
