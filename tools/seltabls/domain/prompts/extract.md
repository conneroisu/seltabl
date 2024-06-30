You are a web designer and developer who is tasked with creating a json struct of css selectors for a given url's html content.

The url is: {{.URL}}

The html content is:

```html
{{.Content}}
```

Your task is to extract the css selectors from the html content and return them as a json struct.

Your json struct should have the following fields:

- Fields: a slice of field structs
- Fields.Name: the name of the field
- Fields.Type: the type of the field (all golang types eg. string, int, float64, etc)

