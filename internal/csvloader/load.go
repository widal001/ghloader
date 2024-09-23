package csvloader

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

// Loads a CSV or TSV file and converts it into a list of map[string]string entries
func loadCSV(filepath string, delimiter rune) ([]map[string]string, error) {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.Comma = delimiter // Set the delimiter (e.g., ',' for CSV or '\t' for TSV)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Ensure there are records and the first row contains the headers
	if len(records) < 1 {
		return nil, errors.New("no records found")
	}

	headers := records[0] // First row is treated as headers
	var result []map[string]string

	// Iterate over each subsequent row and map headers to values
	for _, row := range records[1:] {
		if len(row) != len(headers) {
			return nil, errors.New("row length does not match header length")
		}
		entry := make(map[string]string)
		for i, header := range headers {
			entry[header] = row[i]
		}
		result = append(result, entry)
	}

	return result, nil
}

// Detect file type based on the extension and loads CSV/TSV accordingly
func LoadCSV(filepath string) ([]map[string]string, error) {
	// Detect delimiter based on file extension
	delimiter := ','
	if strings.HasSuffix(filepath, ".tsv") {
		delimiter = '\t'
	}

	return loadCSV(filepath, delimiter)
}
