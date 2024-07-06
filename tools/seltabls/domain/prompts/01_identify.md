Extract the information that can be extracted from two given urls.

You want to extract the information that can be extracted from two given urls.

Identify the information that can be extracted from the a given urls html response.

Content of {{.URL1}}:

```html
{{.URL1Content}}
```

Your task is to extract the information that can be extracted from the two given urls.

You should respond with a list of the information that can be extracted from the two given urls.

The format of your response should be a list of strings.

Each string should be a single piece of information that can be extracted from the two given urls.

For example, if the two given urls are:

Content of https://example.com/team/1234:
```
<div class="team-1234" id="player_scores">
	<table>
		<tr><td>Player</td><td>Score</td></tr>
		<tr><td>John Doe</td><td>100</td></tr>
		<tr><td>Tomas Jefferson</td><td>109</td></tr>
	</table>
	<table border="1" class="dataframe team_leaders">
		<h1>Team Leaders</h1>
		<tr><td>Name</td><td>Stat</td><td>Value</td></tr>
		<tr><td>JohnDoe</td><td>Strikeouts</td><td>100</td></tr>
		<tr><td>NikolaTesla</td><td>Hits</td><td>90</td></tr>
	</table>
	<table border="1" class="dataframe" id="team_scores">
		<tr><td>Team</td><td>Score</td></tr>
		<tr><td>TeamA</td><td>100</td></tr>
		<tr><td>TeamB</td><td>90</td></tr>
	</table>
</div>
```
an example response and the information that can be extracted from the two examples is:

```response
player individual scores (.player_scores)
team indiviual leaders (#team_leaders)
team scores (#dataframe.team_scores)
```

If there is no information that can be extracted from the given url, respond with an empty list. 
