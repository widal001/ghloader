package csvloader

import (
	"path/filepath"
	"testing"
)

// CompareMaps compares two maps and returns true if they are equivalent
func AssertMapsEqual(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}
	return true
}

var wanted = []map[string]string{
	{
		"URL":          "https://github.com/widal001/project-demo/issues/52",
		"Story Points": "3",
		"Status":       "Done",
	},
	{
		"URL":          "https://github.com/widal001/project-demo/issues/53",
		"Story Points": "5",
		"Status":       "In Progress",
	},
}

func TestLoadCSV(t *testing.T) {
	// Arrange
	path := filepath.Join("testdata", "test.csv")
	// Act
	got, err := LoadCSVFromPath(path)
	// Assert
	if err != nil {
		t.Fatalf("unexpected error loading csv: %v", err)
	}
	if !AssertMapsEqual(wanted[0], got[0]) {
		t.Fatalf("\n\nGot: %v\n Wanted: %v\n\n", got, wanted)
	}

}

func TestLoadTSV(t *testing.T) {
	// Arrange
	path := filepath.Join("testdata", "test.tsv")
	// Act
	got, err := LoadCSVFromPath(path)
	// Assert
	if err != nil {
		t.Fatalf("unexpected error loading csv: %v", err)
	}
	if !AssertMapsEqual(wanted[0], got[0]) {
		t.Fatalf("\n\nGot: %v\n Wanted: %v\n\n", got, wanted)
	}

}
