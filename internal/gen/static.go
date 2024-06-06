package gen

var testTempl = `
package {{.PackageName}}

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"
)

// fixtureFile is a test fixture file to scrape.
//
//go:embed fixture.html
var fixtureFile string

// honestFile is a single line manually extracted from the fixture
//
//go:embed honest.json
var honestFile string

// jsonToMap converts a JSON string to a map[string]interface{}.
func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

func jsonArrayToMap(jsonStr string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

// isEqual checks if two JSON objects are equal.
func isEqual(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// contains checks if the larger JSON array contains the smaller JSON object.
func contains(larger []map[string]interface{}, smaller map[string]interface{}) bool {
	for _, obj := range larger {
		if isEqual(obj, smaller) {
			return true
		}
	}
	return false
}

// TestScrape tests the scrape function
func TestScrape(t *testing.T) {
	output, err := Scrape(fixtureFile)
	if err != nil {
		t.Fatalf("failed to scrape: %v", err)
	}
	jsonMap := jsonArrayToMap(output)
	outputMap := jsonToMap(honestFile)
	if !contains(jsonMap, outputMap) {
		t.Errorf("expected %v, got %v", jsonToMap(honestFile), jsonMap)
	}
}
`
