package web

import (
	"fmt"
	"net/http"
)

// RunServer starts the web server on the specified port
func RunServer(port string) {
	http.HandleFunc("/", UploadForm)
	http.HandleFunc("/upload", UploadFile)

	fmt.Printf("Starting web server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting web server: %v\n", err)
	}
}
