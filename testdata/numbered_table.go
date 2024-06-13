package testdata

// NumberedStruct is a test struct
type NumberedStruct struct {
	Header1 string `json:"Header 1" seltabl:"Header 1" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
	Header2 string `json:"Header 2" seltabl:"Header 2" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
	Header3 string `json:"Header 3" seltabl:"Header 3" hSel:"tr:nth-child(1) td:nth-child(3)" dSel:"tr:not(:first-child) td:nth-child(3)" cSel:"$text"`
}

// NumberedTable is the html for the numbered table
var NumberedTable = `
<table border="1">
	<tr>
		<td>Header 1</td>
		<td>Header 2</td>
		<td>Header 3</td>
	</tr>
	<tr>
		<td>Row 1, Cell 1</td>
		<td>Row 1, Cell 2</td>
		<td>Row 1, Cell 3</td>
	</tr>
	<tr>
		<td>Row 2, Cell 1</td>
		<td>Row 2, Cell 2</td>
		<td>Row 2, Cell 3</td>
	</tr>
	<tr>
		<td>Row 3, Cell 1</td>
		<td>Row 3, Cell 2</td>
		<td>Row 3, Cell 3</td>
	</tr>
</table>
`

// NumberedTableResult is the expected result of parsing the numbered table
var NumberedTableResult = []NumberedStruct{
	{
		Header1: "Row 1, Cell 1",
		Header2: "Row 1, Cell 2",
		Header3: "Row 1, Cell 3",
	},
	{
		Header1: "Row 2, Cell 1",
		Header2: "Row 2, Cell 2",
		Header3: "Row 2, Cell 3",
	},
	{
		Header1: "Row 3, Cell 1",
		Header2: "Row 3, Cell 2",
		Header3: "Row 3, Cell 3",
	},
}
