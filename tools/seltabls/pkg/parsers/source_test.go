package parsers_test

import (
	"os"
	"path/filepath"
	"testing"

	_ "embed"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
)

//go:embed testdata/test.go
var testFile string

func TestParseSingleFile(t *testing.T) {
	source := `
package example

import "fmt"

// Person is a struct for a person
type Person struct {
	Name string
	Age  int
}

func (p *Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

var GlobalVar = "Test"

const GlobalConst = 42
`
	// Write the source to a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "example.go")
	err := os.WriteFile(tmpFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Parse the file
	goFile, err := parsers.ParseSingleFile(tmpFile, true)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Check the package name
	if goFile.Package != "example" {
		t.Errorf("Expected package name to be 'example', got '%s'", goFile.Package)
	}

	// Check the number of structs
	if len(goFile.Structs) != 1 {
		t.Errorf("Expected 1 struct, got %d", len(goFile.Structs))
	}

	// Check the struct name
	if goFile.Structs[0].Name != "Person" {
		t.Errorf("Expected struct name to be 'Person', got '%s'", goFile.Structs[0].Name)
	}
	// Check the struct comment
	if goFile.Structs[0].Comments != "Person is a struct for a person" {
		t.Errorf("Expected struct comment to be 'Person is a struct for a person', got '%s'", goFile.Structs[0].Comments)
	}

	// Check the number of methods
	if len(goFile.StructMethods) != 1 {
		t.Errorf("Expected 1 method, got %d", len(goFile.StructMethods))
	}

	// Check the method name
	if goFile.StructMethods[0].Name != "Greet" {
		t.Errorf("Expected method name to be 'Greet', got '%s'", goFile.StructMethods[0].Name)
	}

	// Check the number of imports
	if len(goFile.Imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(goFile.Imports))
	}

	// Check the import path
	if goFile.Imports[0].Path != `"fmt"` {
		t.Errorf("Expected import path to be '\"fmt\"', got '%s'", goFile.Imports[0].Path)
	}

	// Test ImportPath method
	importPath, isExternal, err := goFile.ImportPath()
	if err != nil {
		t.Fatalf("ImportPath method failed: %v", err)
	}
	if isExternal {
		t.Errorf("Expected isExternalPackage to be false, got true")
	}
	if importPath == "" {
		t.Errorf("Expected importPath to be non-empty, got an empty string")
	}

	// Test Prefix method on the import
	if goFile.Imports[0].Prefix() != "fmt" {
		t.Errorf("Expected Prefix to be 'fmt', got '%s'", goFile.Imports[0].Prefix())
	}
}

func TestGoTag_Get(t *testing.T) {
	tag := &parsers.GoTag{
		Value: "`json:\"name\" xml:\"name\"`",
	}

	jsonTag := tag.Get("json")
	if jsonTag != "name" {
		t.Errorf("Expected json tag to be 'name', got '%s'", jsonTag)
	}

	xmlTag := tag.Get("xml")
	if xmlTag != "name" {
		t.Errorf("Expected xml tag to be 'name', got '%s'", xmlTag)
	}
}

func TestParseTestFile(t *testing.T) {
	goFile, err := parsers.ParseSource(testFile, "testdata", true)
	if err != nil {
		t.Fatalf("Failed to parse test file: %v", err)
	}

	// Check the package name
	if goFile.Package != "testdata" {
		t.Errorf("Expected package name to be 'testdata', got '%s'", goFile.Package)
	}

	// Check the number of structs
	if len(goFile.Structs) != 4 {
		t.Errorf("Expected 4 structs, got %d", len(goFile.Structs))
	}

	// Check struct names
	expectedStructs := []string{"Structure", "Field", "Tags", "Tag"}
	for i, expectedName := range expectedStructs {
		if goFile.Structs[i].Name != expectedName {
			t.Errorf("Expected struct name to be '%s', got '%s'", expectedName, goFile.Structs[i].Name)
		}
	}

	// Check the number of interfaces
	if len(goFile.Interfaces) != 1 {
		t.Errorf("Expected 1 interface, got %d", len(goFile.Interfaces))
	}

	// Check the interface name
	if goFile.Interfaces[0].Name != "Tagger" {
		t.Errorf("Expected interface name to be 'Tagger', got '%s'", goFile.Interfaces[0].Name)
	}

	// Check the number of methods in the Tagger interface
	if len(goFile.Interfaces[0].Methods) != 3 {
		t.Errorf("Expected 3 methods in the Tagger interface, got %d", len(goFile.Interfaces[0].Methods))
	}

	// Check method names in the Tagger interface
	expectedMethods := []string{"Get", "Set", "AddOptions"}
	for i, expectedName := range expectedMethods {
		if goFile.Interfaces[0].Methods[i].Name != expectedName {
			t.Errorf("Expected method name to be '%s', got '%s'", expectedName, goFile.Interfaces[0].Methods[i].Name)
		}
	}
}
