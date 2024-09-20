package project

import (
	"fmt"
	"strconv"
	"time"
)

// Struct to represent a field in the project
type ProjectV2Field struct {
	ID      string
	Name    string
	Type    string
	Options []struct {
		Id   string
		Name string
	}
	Iterations []struct {
		Id    string
		Title string
	}
}

func (f *ProjectV2Field) FormatUpdateValue(strVal string) (map[string]interface{}, error) {
	switch f.Type {
	case "NUMBER":
		return f.FormatNumberInput(strVal)
	case "TITLE":
		return map[string]interface{}{"text": strVal}, nil
	case "DATE":
		return f.FormatDateInput(strVal)
	case "SINGLE_SELECT":
		return f.FormatSingleSelectInput(strVal)
	case "ITERATION":
		return f.FormatIterationInput(strVal)
	default:
		// if its not one of the field types above, skip the update
		return make(map[string]interface{}), nil
	}
}

// Format the input value for a SingleSelect field
func (f *ProjectV2Field) FormatSingleSelectInput(input string) (map[string]interface{}, error) {
	for _, opt := range f.Options {
		if opt.Name == input {
			return map[string]interface{}{"singleSelectOptionId": opt.Id}, nil
		}
	}
	return make(map[string]interface{}), fmt.Errorf("%s isn't a valid option for field %s", input, f.Name)
}

// Format the input value for an Iteration field
func (f *ProjectV2Field) FormatIterationInput(name string) (map[string]interface{}, error) {
	for _, iteration := range f.Iterations {
		if iteration.Title == name {
			return map[string]interface{}{"iterationId": iteration.Id}, nil
		}
	}
	return make(map[string]interface{}), fmt.Errorf("%s isn't a valid iteration for this field", name)
}

// Format the input value for a Date field
func (f *ProjectV2Field) FormatDateInput(input string) (map[string]interface{}, error) {
	date, err := time.Parse("2006-01-02", input)
	if err != nil {
		return make(map[string]interface{}), fmt.Errorf("%s is not a valid date", input)
	}
	return map[string]interface{}{"date": date}, nil
}

// Format the input value for a Number field
func (f *ProjectV2Field) FormatNumberInput(input string) (map[string]interface{}, error) {
	floatVal, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return make(map[string]interface{}), fmt.Errorf("%s is not a valid number", input)
	}
	return map[string]interface{}{"number": floatVal}, nil
}
