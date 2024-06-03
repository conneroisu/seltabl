package test_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/conneroisu/seltabl"
)

// Sample HTML fixture
var fixture = `
<table>
	<tr>
		<td>a</td>
		<td>b</td>
	</tr>
	<tr>
		<td> 1 </td>
		<td>2</td>
	</tr>
	<tr>
		<td>3 </td>
		<td> 4</td>
	</tr>
	<tr>
		<td> 5 </td>
		<td> 6</td>
	</tr>
	<tr>
		<td>7 </td>
		<td> 8</td>
	</tr>
</table>
`

// Struct for the test
type fixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

func TestNewFromString(t *testing.T) {
	expected := []fixtureStruct{
		{A: "1", B: "2"},
		{A: "3", B: "4"},
		{A: "5", B: "6"},
		{A: "7", B: "8"},
	}

	result, err := seltabl.NewFromString[fixtureStruct](fixture)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, result)
	}
}

func TestNewFromString_InvalidHTML(t *testing.T) {
	invalidHTML := `<table><tr><td>a</td></tr>`

	_, err := seltabl.NewFromString[fixtureStruct](invalidHTML)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestNewFromString_MissingHeaderSelector(t *testing.T) {
	type invalidStruct struct {
		A string `json:"a" seltabl:"a" dSel:"tr td:nth-child(1)" cSel:"$text"`
		B string `json:"b" seltabl:"b" dSel:"tr td:nth-child(2)" cSel:"$text"`
	}

	_, err := seltabl.NewFromString[invalidStruct](fixture)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestNewFromString_MissingDataSelector(t *testing.T) {
	type invalidStruct struct {
		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" cSel:"$text"`
		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" cSel:"$text"`
	}

	_, err := seltabl.NewFromString[invalidStruct](fixture)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestNewFromString_EmptyHTML(t *testing.T) {
	emptyHTML := ``

	_, err := seltabl.NewFromString[fixtureStruct](emptyHTML)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

// TestNewFromURL tests the NewFromURL function
func TestNewFromURL(t *testing.T) {
	// Mock HTTP server for testing
	server := httpTestServer(t, fixture)
	defer server.Close()

	expected := []fixtureStruct{
		{A: "1", B: "2"},
		{A: "3", B: "4"},
		{A: "5", B: "6"},
		{A: "7", B: "8"},
	}

	result, err := seltabl.NewFromURL[fixtureStruct](server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected %v,\n got %v", expected, result)
	}
}

// httpTestServer sets up a test HTTP server
func httpTestServer(t *testing.T, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(body))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
}
