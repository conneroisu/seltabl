package parsers

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// TestSelectors tests the selectors function
func TestSelectors(t *testing.T) {
	t.Run("TestSelectors", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  []string
		}{
			{
				name: "TestSelectors",
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
					"html body table",
					"html body table tbody tr",
				},
			},
			{
				name: "TestSelectors",
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
					"html body",
					"html body table",
					"html body table tbody tr",
					"html head title",
					"html body div h1",
				},
			},
			{
				name: "TestSelectors",
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
					"head",
					"body",
					"table",
					"tbody",
					"tr",
					"td",
					"html",
					"div",
					"h1",
				},
			},
			{
				name: "TestSelectors",
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
					"head",
					"body",
					"table",
					"tbody",
					"tr",
					"td",
					"html",
					"div",
					"h1",
					"p",
				},
			},
			{
				name: "TestSelectors",
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
					"html head",
					"html body",
					"html",
					"html body table",
					"html body table tbody",
					"html body table tr",
					"html body table tr td",
					"html body table tr td a[href=https://example.com]",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.input),
				)
				if err != nil {
					t.Fatalf("failed to create document: %v", err)
				}
				got := GetAllSelectors(doc)
				for _, wa := range tt.want {
					// TODO: Need a better message here
					assert.NotEmpty(t, doc.Find(wa))
					// TODO: Need a better message here
					assert.Contains(t, got, wa)
				}
			})
		}
	})
}
