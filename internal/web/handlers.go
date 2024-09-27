package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/widal001/ghloader/internal/csvloader"
	"github.com/widal001/ghloader/internal/project"
)

// UploadForm renders the HTML form that asks for the URL and CSV data to update a GH project
func UploadForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "upload.html")))
	tmpl.Execute(w, nil)
}

// UpdateProject uses the URL and CSV data to update the GitHub project after submission
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the URL provided and use it to fetch GitHub project data
	url := r.FormValue("url")
	proj, err := project.FromURL(url)
	if err != nil {
		http.Error(w, "Failed to load project data", http.StatusInternalServerError)
		return
	}

	// Parse the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Load and process the file content
	contents, err := csvloader.LoadCSVFromReader(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	// Parse those contents into items to insert or update on the project
	items, err := project.ItemsFromCSV(contents, "URL")
	if err != nil {
		http.Error(w, "Error parsing contents of the CSV", http.StatusInternalServerError)
	}

	// Update or insert the parsed items
	updated, _ := proj.BatchUpsertItems(items)
	for _, itemURL := range updated {
		fmt.Fprintf(w, "Updated project item with URL: %s\n", itemURL)
	}

}
