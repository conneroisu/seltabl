package prompts

import (
	"testing"
)

// TestNewIdentifyPrompt tests the NewIdentifyPrompt function
func TestNewIdentifyPrompt(t *testing.T) {
	url1 := "https://example.com/team/1234"
	url2 := "https://example.com/team/4321"
	content1 := `<div class="team-1234">
	<table>
		<tr>
			<td>Player</td>
			<td>Score</td>
		</tr>
		<tr>
			<td>John Doe</td>
			<td>100</td>				
		</tr>
		<tr>
			<td>Jane Doe</td>
			<td>90</td>
		</tr>
	</table>
	<table border="1" class="dataframe" id="team_leaders">
		<h1>Team Leaders</h1>
		<tr>
			<td>Name</td>												
			<td>Stat</td>
			<td>Value</td>
		</tr>
		<tr>
			<td>John Doe</td>
			<td>Strikeouts</td>
			<td>100</td>
		</tr>
		<tr>
			<td>Jane Doe</td>
			<td>Hits</td>
			<td>90</td>
		</tr>
	</table>
	<table border="1" class="dataframe" id="team_scores">
		<tr>
			<td>Team</td>
			<td>Score</td>
		</tr>							
			
		<tr>
			<td>Team A</td>
			<td>100</td>
		</tr>
		<tr>
			<td>Team B</td>
			<td>90</td>
		</tr>
	</table>
</div>`

	content2 := `<div class="team-4321" id="team-4321">
	<h1>Team 4321</h1>
	<table>
		<tr>
			<td>Player</td>
			<td>Score</td>
		</tr>
		<tr>
			<td>John Doe</td>
			<td>100</td>
		</tr>
		<tr>
			<td>Jane Doe</td>
			<td>90</td>
		</tr>
	</table>
	<table border="1" class="dataframe" id="team_leaders">
		<h1>Team Leaders</h1>
		<tr>
			<td>Name</td>												
			<td>Stat1</td>
			<td>Value</td>												
		</tr>
		<tr>
			<td>John Doe</td>
			<td>11</td>
			<td>100</td>
		</tr>
		<tr>
			<td>Jane Doe</td>
			<td>12</td>
			<td>90</td>
		</tr>
	</table>
	<table border="1" class="dataframe" id="team_scores">
		<tr>
			<td>Team</td>
			<td>Score</td>
		</tr>
		<tr>
			<td>Team A</td>
			<td>100</td>
		</tr>
		<tr>
			<td>Team B</td>
			<td>90</td>
		</tr>
	</table>
	</div>
	`

	prompt, err := NewIdentifyPrompt(url1, url2, content1, content2)
	if err != nil {
		t.Errorf("NewIdentifyPrompt() error = %v", err)
		return
	}

	t.Logf("prompt: %s", prompt)
}
