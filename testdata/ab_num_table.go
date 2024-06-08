package testdata

// FixtureStruct is a test struct
type FixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr:nth-child(1+) td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr:nth-child(1+) td:nth-child(2)" cSel:"$text"`
}

var FixtureABNumTable string = `
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
		<td>4 </td>
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

// FixtureABNumTableResult is the expected result of parsing the ab_num_table
var FixtureABNumTableResult = []FixtureStruct{
	{A: "1", B: "2"},
	{A: "3", B: "4"},
	{A: "5", B: "6"},
	{A: "7", B: "8"},
}
