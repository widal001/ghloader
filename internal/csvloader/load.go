package csvloader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Loads a CSV or TSV from any io.Reader and converts it into a list of map[string]string entries
func loadCSV(reader io.Reader, fileName string) ([]map[string]string, error) {
	// Make an empty result array for records parsed from the CSV
	var result []map[string]string

	// Determine the delimiter based on file type
	delimiter, err := detectDelimiter(fileName)
	if err != nil {
		return result, err
	}

	// Create a new CSV reader
	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimiter

	// Read all records from the CSV
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Check that the CSV contains records, then treat the first row as headers
	if len(records) < 1 {
		return nil, errors.New("no records found")
	}
	headers := records[0]

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
func LoadCSVFromPath(path string) ([]map[string]string, error) {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Load and return the CSV contents
	return loadCSV(file, path)
}

// Parse CSV/TSV from an io.Reader (e.g., for multipart uploads)
func LoadCSVFromReader(reader io.Reader, fileName string) ([]map[string]string, error) {
	// Load and return the CSV contents
	return loadCSV(reader, fileName)
}

// Detect delimiter based on file extension (csv or tsv)
func detectDelimiter(path string) (rune, error) {
	// Get the file extension
	ext := filepath.Ext(path)
	switch ext {
	case ".csv":
		return ',', nil
	case ".tsv":
		return '\t', nil
	default:
		return 0, fmt.Errorf("only CSV or TSV files are allowed")
	}
}
