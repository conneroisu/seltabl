Create a Golang struct similar to the one found in the example below.
The struct should be used to parse an HTML table into a Go struct with proper field tags for `json`, `seltabl`, `hSel`, `dSel`, and `cSel`.
These tags are necessary to correctly map the HTML table cells to the struct fields.
The struct will be used in conjunction with the `seltabl` package to parse the table.

Here is an example of the usage and the struct definition:

```go
package main

import (
	"fmt"
	"github.com/conneroisu/seltabl"
	"github.com/conneroisu/seltabl/testdata"
)

// TableStruct defines a struct to represent the HTML table structure
//
// td:nth-child(1) (the first table data element)
//
// td:nth-child(2) (the second table data element)
//
// tr:nth-child(1) (the first table row)
//
// tr:nth-child(2) (the second table row)
//
// tr (a table row)
type TableStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

// Example HTML table to be parsed
var fixture = `
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
	<tr>
		<td>5</td>
		<td>6</td>
	</tr>
	<tr>
		<td>7</td>
		<td>8</td>
	</tr>
</table>
`
```

Make sure the struct has the following field tags:
- `json`: for JSON serialization/deserialization.
- `seltabl`: for the seltabl library to map the field to the table header.
- `hSel`: to specify the CSS selector for the header cells.
- `dSel`: to specify the CSS selector for the data cells.
- `cSel`: to specify the content selection method (e.g., `$text` to extract text content).

The HTML table contains headers and multiple rows of data that need to be parsed into the struct fields.

YOU MUST Output just the struct definition without any other text, '```' or any other formatting.

e.g.
```go
type TableStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}
```

Once you have the struct definition, we will run `go test` to verify that the struct is properly defined.

The test will fail if the struct is not properly defined giving you a chance to fix the issue.
When/if you must fix the issue, you must output just the struct definition without any other text, '```' or any other formatting.


---
