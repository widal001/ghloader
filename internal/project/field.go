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

// Get the ID of an option for a single select field by its name
func (f *ProjectV2Field) GetSelectOptionId(name string) (string, error) {
	for _, opt := range f.Options {
		if opt.Name == name {
			return opt.Id, nil
		}
	}
	return "", fmt.Errorf("%s isn't a valid option for this field", name)
}

// Get the ID of an option for a single select field by its name
func (f *ProjectV2Field) GetIterationId(name string) (string, error) {
	for _, opt := range f.Iterations {
		if opt.Title == name {
			return opt.Id, nil
		}
	}
	return "", fmt.Errorf("%s isn't a valid iteration for this field", name)
}

func (f *ProjectV2Field) FormatUpdateValue(strVal string) (map[string]interface{}, error) {
	var value map[string]interface{}
	switch f.Type {
	case "NUMBER":
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return value, fmt.Errorf("%s is not a valid number", strVal)
		}
		value = map[string]interface{}{"number": floatVal}
	case "TITLE":
		value = map[string]interface{}{"text": strVal}
	case "DATE":
		date, err := time.Parse(time.RFC3339, strVal)
		if err != nil {
			return value, fmt.Errorf("%s is not a valid date", strVal)
		}
		value = map[string]interface{}{"text": date}
	case "SINGLE_SELECT":
		optId, err := f.GetSelectOptionId(strVal)
		if err != nil {
			return value, fmt.Errorf("couldn't update field %s", f.Name)
		}
		value = map[string]interface{}{"singleSelectOptionId": optId}
	case "ITERATION":
		iterationId, err := f.GetIterationId(strVal)
		if err != nil {
			return value, fmt.Errorf("couldn't update field %s", f.Name)
		}
		value = map[string]interface{}{"iterationId": iterationId}
	default:
		// if its not one of the field types above, skip the update
		return value, nil
	}
	return value, nil
}
