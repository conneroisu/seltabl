You are to generate a json schema for a given url, {{ .URL }}.

The json schema must have a number of attributes for each field in the struct.

- header selector (hSel): used to find the header row and column for the field in the given struct.
- data selector (dSel): used to find the data column for the field in the given struct.
- query selector (qSel): used to query for the inner text or attribute of the cell.
- control selector (cSel): used to control what to query for the inner text or attribute of the cell.

Example Output:

```go
package main

// @url: https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
	B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
}
```

Your task is to generate the golang struct for the given url to capture all the data from the webpage not just the first table.

You must use the given url's html content to generate the golang struct.

Your html content:

{{ .Content }}

Your selectors with attriubutes (meaning without `:nth-child(1)`, `:not(:first-child)`, etc):

```go
{{ .Selectors }}
```
