package project

import (
	"strings"
	"testing"
)

var tmplDir string = "../../templates"

var testProj = ProjectV2{
	Id: "barbaz",
	Fields: map[string]ProjectV2Field{
		"Bar Name": {
			Id:   "barId",
			Name: "Bar Name",
			Type: "NUMBER",
		},
		"Baz field": {
			Id:   "bazId",
			Name: "Baz field",
			Type: "SINGLE_SELECT",
			Options: []Option{
				{
					Id:   "baz1",
					Name: "Option 1",
				},
				{
					Id:   "baz2",
					Name: "Option 2",
				},
			},
		},
	},
}

var testData = []FieldData{
	{
		Name:  "Bar Name",
		Value: "3",
	},
	{
		Name:  "Baz field",
		Value: "Option 2",
	},
}

func assertContains(t *testing.T, got, wanted string) {
	t.Helper()
	// Assert - Bar formatted correctly
	if !strings.Contains(got, wanted) {
		t.Fatalf("Query is missing `%s`:\n %s", wanted, got)
	}
}

func TestGenerateUpdateQuery(t *testing.T) {
	// Act
	got, err := testProj.GenerateUpdateQuery(tmplDir, testData)
	// Assert - rendered correctly
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	// Assert - contains key elements
	assertContains(t, got, "barname: update")
	assertContains(t, got, "bazfield: update")
	assertContains(t, got, "fieldId: \"barId\"")
	assertContains(t, got, "fieldId: \"bazId\"")
	assertContains(t, got, "value: {number: 3.000000}")
	assertContains(t, got, "value: {singleSelectOptionId: \"baz2\"}")
}
