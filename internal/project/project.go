package project

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Struct to represent a field in the project
type ProjectV2Field struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Struct to hold all fields from the ProjectV2
type ProjectV2 struct {
	Login  string
	Number int
	Type   string
	Id     string
	Title  string
	Fields map[string]ProjectV2Field
}

// Method to get the type of a field by name
func (p *ProjectV2) GetFieldType(fieldName string) (string, error) {
	field, ok := p.Fields[fieldName]
	if !ok {
		return "", fmt.Errorf("field %s not found", fieldName)
	}
	return field.Type, nil
}

func FromURL(projectURL string) (*ProjectV2, error) {
	// Validate the URL
	re := regexp.MustCompile(`https:\/\/github.com\/(users|orgs)\/\w+\/projects\/\d+`)
	ok := re.Match([]byte(projectURL))
	if !ok {
		return nil, fmt.Errorf("not a valid GitHub project URL: %s", projectURL)
	}

	// Parse the URL
	u, err := url.Parse(projectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Split the URL path params into segments
	// and extract login, type, and project number
	segments := strings.Split(u.Path, "/")
	login := segments[2]
	projectType := segments[1]
	projectNumber, err := strconv.Atoi(segments[4])
	if err != nil {
		return nil, fmt.Errorf("invalid project number: %v", err)
	}

	// Populate the struct
	proj, err := LoadProjectMetadata(login, projectNumber, projectType)
	if err != nil {
		return nil, fmt.Errorf("failed to load project metadata: %v", err)
	}
	return proj, nil
}
