Extract the information that can be extracted from two given urls.

You want to extract the information that can be extracted from two given urls.

Identify the information that can be extracted from the two given urls contents.
Notice differences between the two url's contents.

Content of {{.URL1}}:
```html
{{.URL1Content}}
```

Content of {{.URL2}}:
```html
{{.URL2Content}}
```

Your task is to extract the information that can be extracted from the two given urls.

You should respond with a list of the information that can be extracted from the two given urls.

The format of your response should be a list of strings.
Each string should be a single piece of information that can be extracted from the two given urls.

For example, if the two given urls are:

Content of https://example.com/team/1234:
```
<div class="team-1234">
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
</div>
```

Content of https://example.com/team/4321:
```
<div class="team-4321" id="team-4321">
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
```

an example response and the information that can be extracted from the two examples is:

```response
player individual scores
team indiviual leaders 
team scores
```

If there is no information that can be extracted from the two given urls, respond with an empty list. 
