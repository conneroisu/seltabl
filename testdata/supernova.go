package testdata

import (
	_ "embed"
)

// SuperNovaStruct is a test struct
type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year"      seltabl:"Year"      hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type"      seltabl:"Type"      hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `json:"Distance"  seltabl:"Distance"  hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)"`
	Notes     string `json:"Notes"     seltabl:"Notes"     hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

// SuperNovaTable is the html for the supernova table
var SuperNovaTable = `
<table>
	<tr>
		<th>Supernova</th>
		<th>Year</th>
		<th>Type</th>
		<th>Distance (light-years)</th>
		<th>Notes</th>
	</tr>
	<tr>
		<td>SN 1006</td>
		<td>1006</td>
		<td>Type Ia</td>
		<td>7,200</td>
		<td>Brightest recorded supernova in history</td>
	</tr>
	<tr>
		<td>SN 1054 (Crab Nebula)</td>
		<td>1054</td>
		<td>Type II</td>
		<td>6,500</td>
		<td>Formed the Crab Nebula and pulsar</td>
	</tr>
	<tr>
		<td>SN 1572 (Tycho's Supernova)</td>
		<td>1572</td>
		<td>Type Ia</td>
		<td>8,000-10,000</td>
		<td>Observed by Tycho Brahe</td>
	</tr>
	<tr>
		<td>SN 1604 (Kepler's Supernova)</td>
		<td>1604</td>
		<td>Type Ia</td>
		<td>20,000</td>
		<td>Last observed supernova in the Milky Way</td>
	</tr>
	<tr>
		<td>SN 1987A</td>
		<td>1987</td>
		<td>Type II</td>
		<td>168,000</td>
		<td>Closest observed supernova since 1604</td>
	</tr>
	<tr>
		<td>SN 1993J</td>
		<td>1993</td>
		<td>Type IIb</td>
		<td>11,000,000</td>
		<td>In the galaxy M81</td>
	</tr>
</table>
`

// SuperNovaTableResult is the expected result of parsing the supernova table
var SuperNovaTableResult = []SuperNovaStruct{
	{
		Supernova: "SN 1006",
		Year:      "1006",
		Type:      "Type Ia",
		Distance:  "7,200",
		Notes:     "Brightest recorded supernova in history",
	},
	{
		Supernova: "SN 1054 (Crab Nebula)",
		Year:      "1054",
		Type:      "Type II",
		Distance:  "6,500",
		Notes:     "Formed the Crab Nebula and pulsar",
	},
	{
		Supernova: "SN 1572 (Tycho's Supernova)",
		Year:      "1572",
		Type:      "Type Ia",
		Distance:  "8,000-10,000",
		Notes:     "Observed by Tycho Brahe",
	},
	{
		Supernova: "SN 1604 (Kepler's Supernova)",
		Year:      "1604",
		Type:      "Type Ia",
		Distance:  "20,000",
		Notes:     "Last observed supernova in the Milky Way",
	},
	{
		Supernova: "SN 1987A",
		Year:      "1987",
		Type:      "Type II",
		Distance:  "168,000",
		Notes:     "Closest observed supernova since 1604",
	},
	{
		Supernova: "SN 1993J",
		Year:      "1993",
		Type:      "Type IIb",
		Distance:  "11,000,000",
		Notes:     "In the galaxy M81",
	},
}
