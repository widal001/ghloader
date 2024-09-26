package main

import (
	"fmt"
	"log"
	"os"

	"github.com/widal001/ghloader/internal/csvloader"
	"github.com/widal001/ghloader/internal/project"
)

func main() {
	// Parse the command line arguments
	if len(os.Args) != 3 {
		log.Fatal(`
This script expects exactly two positional args:
1. The URL of the project to update
2. A path to a source CSV or TSV with data to use during the update
`)
	}
	projURL := os.Args[1]
	updateCSV := os.Args[2]

	// Load the fields from the ProjectV2
	proj, err := project.FromURL(projURL)
	if err != nil {
		log.Fatalf("Error loading project fields: %v", err)
	}
	fmt.Printf("Project has %d fields\n", len(proj.Fields))

	// Load the contents of the csv
	contents, err := csvloader.LoadCSV(updateCSV)
	if err != nil {
		log.Fatalf("Error loading CSV: %v", updateCSV)
	}

	// Parse those contents into items to insert or update on the project
	items, err := project.ItemsFromCSV(contents, "URL")
	if err != nil {
		log.Fatal("Error parsing contents of the CSV")
	}

	// Update or insert the parsed items
	proj.BatchUpsertItems(items)

}
