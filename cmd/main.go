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

	// Get the type of a specific field
	for fieldName := range proj.Fields {
		fieldType, err := proj.GetFieldType(fieldName)
		if err != nil {
			log.Fatalf("Error getting field type: %v", err)
		}
		fmt.Printf("The type of '%s' field is: %s\n", fieldName, fieldType)
	}
	sprint := proj.Fields["Sprint"]
	fmt.Printf("%v\n", sprint.Iterations)
	status := proj.Fields["Status"]
	fmt.Printf("%v\n", status.Options)
}
