package parsers

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// TestSelectors tests the selectors function
func TestSelectors(t *testing.T) {
	type args struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}
	tests := []args{
		{
			name: "Test Get Selectors for HTML with single table",
			input: `
				<table>
					<tr>
						<td>a</td>
						<td>b</td>
					</tr>
					<tr>
						<td>1</td>
						<td>2</td>
					</tr>
					<tr>
						<td>3</td>
						<td>4</td>
					</tr>
				</table>
				`,
			want: []string{
				"html",
				"html > body > table",
				"html > body > table > tbody > tr",
			},
		},
		{
			name: "Test Get Selectors for HTML with multiple tables",
			input: `
				<html>
					<head>
						<title>Test</title>
					</head>
					<body>
						<div>
							<h1>Test</h1>
							<p>Test</p>
						</div>
						<table>
							<tr>
								<td>a</td>
								<td>b</td>
							</tr>
							<tr>
								<td>1</td>
								<td>2</td>
							</tr>
							<tr>
								<td>3</td>
								<td>4</td>
							</tr>
						</table>
					</body>
				</html>
				`,
			want: []string{
				"html",
				"html > head",
				"html > body",
				"html > body > table",
				"html > body > table > tbody",
				"html > body > table > tbody > tr",
				"html > body > table > tbody > tr > td",
			},
		},
		{
			name: "Test Get Selectors for HTML with multiple tables and multiple selectors",
			input: `
				<html>
					<head>
						<title>Test</title>
					</head>
					<body>
						<div>
							<h2>Test</h2>
						</div>
					</body>
				</html>
				`,
			want: []string{
				"html",
				"html > head",
				"html > body",
				"html > body > div",
				"html > body > div > h2",
			},
		},
		{
			name: "Test Get Selectors for HTML with multiple tables and multiple selectors",
			input: `
				<html>
					<head>
						<title>Test</title>
					</head>
					<body>
						<div>
							<h1>Test</h1>
							<p>Test</p>
						</div>
						<table>
							<tr>
								<td>
									<a href="https://example.com">
										Test
									</a>
								</td>
								<td>b</td>
							</tr>
							<tr>
								<td>1</td>
								<td>2</td>
							</tr>
							<tr>
								<td>3</td>
								<td>4</td>
							</tr>
						</table>
					</body>
				</html>
				`,
			want: []string{
				"html",
				"html > head",
				"html > body",
				"html > body > table",
				"html > body > table > tbody",
				"html > body > table > tbody > tr",
				"html > body > table > tbody > tr > td",
				"html > body > table > tbody > tr > td > a[href]",
			},
		},
		{
			name: "Test Get Selectors for HTML with multiple tables and multiple selectors and multiple tags",
			input: `
				<html>
					<head>
						<title>Test</title>
					</head>
					<body>
						<div>
							<h1>Test</h1>
							<p>Test</p>
						</div>
						<table>
							<tr>
								<td>
									<a href="https://example.com">
										Test
									</a>
								</td>
								<td>b</td>
							</tr>
							<tr>
								<td>1</td>
								<td>2</td>
							</tr>
							<tr>
								<td>3</td>
								<td>4</td>
							</tr>
						</table>
					</body>
				</html>
				`,
			want: []string{
				"html",
				"html > head",
				"html > body",
				"html > body > table",
				"html > body > table > tbody",
				"html > body > table > tbody > tr",
				"html > body > table > tbody > tr > td",
				"html > body > table > tbody > tr > td > a[href]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			doc, err := goquery.NewDocumentFromReader(
				strings.NewReader(tt.input),
			)
			if err != nil {
				t.Fatalf("failed to create document: %v", err)
			}
			got, err := GetAllSelectors(doc)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("failed to get selectors: %v", err)
			}
			for _, wa := range got {
				t.Logf("\"%s\",", wa)
				// verify that the selector can be found in the document
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.input),
				)
				if err != nil {
					t.Fatalf("failed to create document: %v", err)
				}
				sel := doc.Find(wa)
				assert.NotEqual(t, sel.Length(), 0)
			}
			for _, wa := range tt.want {
				assert.Contains(
					t,
					got,
					wa,
					"selector %s not found in got",
					wa,
				)
			}
		})
	}
}
