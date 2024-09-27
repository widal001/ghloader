package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/widal001/ghloader/internal/csvloader"
	"github.com/widal001/ghloader/internal/project"
	"github.com/widal001/ghloader/internal/web"
)

func main() {
	// Define flags for web and CLI modes
	webMode := flag.Bool("web", false, "Run the application in web server mode")
	port := flag.String("port", "8080", "Port to run the web server on")
	flag.Parse()

	if *webMode {
		web.RunServer(*port)
	} else {
		runCLI()
	}
}

func runCLI() {
	// Parse the command line arguments
	projURL := flag.String("url", "", "URL of the project to update")
	updateCSV := flag.String("file", "", "Path to the file with data used to update the project")
	if *projURL != "" {
		log.Fatal("To use the CLI, you must pass the URL a GitHub project to the -url flag.")
	}
	if *updateCSV != "" {
		log.Fatal("To use the CLI, you must pass the path of a CSV or TSV file -file flag.")
	}

	// Load the contents of the csv
	contents, err := csvloader.LoadCSVFromPath(*updateCSV)
	if err != nil {
		log.Fatalf("Error loading CSV: %v", updateCSV)
	}

	// Load the project metadata from the URL
	proj, err := project.FromURL(*projURL)
	if err != nil {
		log.Fatalf("Error loading project fields: %v", err)
	}
	fmt.Printf("Project has %d fields\n", len(proj.Fields))

	// Parse those contents into items to insert or update on the project
	items, err := project.ItemsFromCSV(contents, "URL")
	if err != nil {
		log.Fatal("Error parsing contents of the CSV")
	}

	// Update or insert the parsed items
	updated, _ := proj.BatchUpsertItems(items)
	for _, item := range updated {
		fmt.Printf("Updated item with URL: %s", item)
	}

}
