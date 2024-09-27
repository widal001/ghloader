package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/widal001/ghloader/internal/csvloader"
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

	// Parse the URL provided
	url := r.FormValue("url")
	fmt.Fprintf(w, "Project URL: %s\n", url)

	// Parse the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Load and process the file content
	records, err := csvloader.LoadCSVFromReader(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	// Respond with the processed data (for demonstration)
	for i, record := range records {
		fmt.Fprintf(w, "Record %d: %+v\n", i+1, record)
	}
}
