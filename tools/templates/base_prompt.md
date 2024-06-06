Respond a json for a scraper to parse a html table into a struct.

The json should be a list of objects with the following fields:
{{ .Fields }}

YOU MUST Output just the json without any other text, '```' or any other formatting.

An short extract of an example of the json file is:
```json
[
	{
		"name": "HomeTeamScore",
		"json_name": "home_team_score",
		"type": "int",
		"header_name": "Home Team Score",
		"header_selector": "tr:nth-child(1) th:nth-child(1)",
		"data_selector": "tr:nth-child(1) td:nth-child(1)",
		"cell_selector": "$text"
	},
	/...
]
```

You will need to refine the json file to operate properly with further prompts if an error is returned while using your generated json.

Your given html:
```html
{{ .HTML }}
```

One exmaple result column of the json after executing the program is:

```json
{{ .Input }}
```
