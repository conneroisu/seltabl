package cmds

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"text/template"

	"github.com/charmbracelet/huh"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/spf13/cobra"
)

var (
	uuri        string
	packageName string
	fileName    string

	body string

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the URL.").
				Value(&uuri).
				Validate(parsers.ValidateURL),
			huh.NewInput().
				Title("Enter the package name.").
				Value(&packageName).
				Validate(parsers.ValidatePackageName),
			huh.NewInput().
				Title("Enter the file name.").
				Value(&fileName).
				Validate(parsers.ValidateFileName),
		),
	)
)

func NewStaticCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "static",
		Short: "Statically define html given a url.",
		Long: `
Statically define html given a url.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.PersistentFlags().StringVarP(
				&uuri,
				"url",
				"u",
				"",
				"url to scrape",
			)
			cmd.PersistentFlags().StringVarP(
				&packageName,
				"package-name",
				"p",
				"",
				"package name to use",
			)
			cmd.PersistentFlags().StringVarP(
				&fileName,
				"file-name",
				"f",
				"",
				"file name to test (without .go extension or _test.go) ex: koala",
			)
			err := form.Run()
			if err != nil {
				return fmt.Errorf("failed to run form: %w", err)

			}
			var staticFileName string
			staticFileName, err = getURLFileName(uuri)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			println("wrote html to:", staticFileName)
			content, err := TestFileString(testfileTemplate{
				PackageName: packageName,
				FileName:    staticFileName,
			})
			err = os.WriteFile(fileName+"_test.go", []byte(content), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			println("wrote file:", fileName+"_test.go")
			return nil
		},
	}
}

// testfileTemplateString is a string for a test file template.
const testfileTemplateString = `
package {{.PackageName}}

import (
	"fmt"
	"github.com/conneroisu/seltabl"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "embed"
)

//go:embed {{.FileName}}
var Fixture string

func Test{{.FileName}}(t *testing.T) {
	a := assert.New(t)
	res, err := seltabl.NewFromString[TableStruct](Fixture)
	a.NoError(err)
	a.NotNil(res)
}
`

var (
	tt = template.Must(template.New("test").Parse(testfileTemplateString))
)

// testfileTemplate is a struct for a test file template.
type testfileTemplate struct {
	PackageName string
	FileName    string
}

// TestFileString returns a string for a test file.
func TestFileString(tmpl testfileTemplate) (string, error) {
	var buf bytes.Buffer
	err := tt.Execute(&buf, tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// sanitizeFileName ensures that the filename is safe to use.
func sanitizeFileName(filename string) string {
	// Define a regex to match invalid filename characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	// Replace invalid characters with an underscore
	return invalidChars.ReplaceAllString(filename, "_")
}

// getURLFileName gets the file name from a url
func getURLFileName(fileURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	// Extract the base filename from the URL path
	baseFileName := path.Base(parsedURL.Path)

	// Sanitize the filename to ensure it's safe to use
	safeFileName := sanitizeFileName(baseFileName)

	if safeFileName == "" {
		return "", fmt.Errorf("filename is empty after sanitization")
	}

	client := &http.Client{}
	resp, err := client.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	f, err := os.Create(safeFileName + ".html")
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	f.Write(body)

	return safeFileName + ".html", nil
}
