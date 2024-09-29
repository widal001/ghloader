package project

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Option struct {
	Id   string
	Name string
}

type Iteration struct {
	Id    string
	Title string
}

// Struct to represent a field in the project
type ProjectV2Field struct {
	Id         string
	Name       string
	Type       string
	Options    []Option
	Iterations []Iteration
}

type fieldProps struct {
	Key   string
	Id    string
	Value string
}

// format a field's name
func (f *ProjectV2Field) formatUpdateKey() string {
	key := strings.ReplaceAll(f.Name, " ", "")
	return strings.ToLower(key)
}

// Format the fieldProps used to generate the GraphQL mutation query for a given field
func (f *ProjectV2Field) formatFieldProps(strVal string) (fieldProps, error) {
	// Format the update value
	value, err := f.FormatUpdateValue(strVal)
	if err != nil {
		return fieldProps{}, err
	}
	// Return the fieldProps
	return fieldProps{
		Id:    f.Id,
		Key:   f.formatUpdateKey(),
		Value: value,
	}, nil
}

// Format the `value` param for the input field in the updateProjectV2ItemFieldValue GraphQL mutation
func (f *ProjectV2Field) FormatUpdateValue(strVal string) (string, error) {
	switch f.Type {
	case "NUMBER":
		return f.formatNumberInput(strVal)
	case "TITLE":
		return fmt.Sprintf("{text: \"%s\"}", strVal), nil
	case "DATE":
		return f.formatDateInput(strVal)
	case "SINGLE_SELECT":
		return f.formatSingleSelectInput(strVal)
	case "ITERATION":
		return f.formatIterationInput(strVal)
	default:
		// if its not one of the field types above, skip the update
		return "", nil
	}
}

// Format the input value for a SingleSelect field
func (f *ProjectV2Field) formatSingleSelectInput(input string) (string, error) {
	for _, opt := range f.Options {
		if opt.Name == input {
			return fmt.Sprintf("{singleSelectOptionId: \"%s\"}", opt.Id), nil
		}
	}
	return "", fmt.Errorf("%s isn't a valid option for field %s", input, f.Name)
}

// Format the input value for an Iteration field
func (f *ProjectV2Field) formatIterationInput(name string) (string, error) {
	for _, iteration := range f.Iterations {
		if iteration.Title == name {
			return fmt.Sprintf("{iterationId: \"%s\"}", iteration.Id), nil
		}
	}
	return "", fmt.Errorf("%s isn't a valid iteration for this field", name)
}

// Format the input value for a Date field
func (f *ProjectV2Field) formatDateInput(input string) (string, error) {
	date, err := time.Parse("2006-01-02", input)
	if err != nil {
		return "", fmt.Errorf("%s is not a valid date", input)
	}
	return fmt.Sprintf("{date: \"%s\"}", date.Format("2006-01-02")), nil
}

// Format the input value for a Number field
func (f *ProjectV2Field) formatNumberInput(input string) (string, error) {
	floatVal, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "", fmt.Errorf("%s is not a valid number", input)
	}
	return fmt.Sprintf("{number: %f}", floatVal), nil
}
