As an expert web programmer, I need you to create a scraping profile containing the css selectors for each bit of content on the given url that I want to scrape.
I want you to create a scraping profile for a given url in the form of a json struct.
The url is: {{.URL}} and the html content is:

```html
{{.Content}}
```

Your task is to create a scraping profile for the given url in the form of a json struct.

Information that I need to include in the scraping profile:

{{.Info}}

Your json struct should have the following fields:

- Fields: a slice of field structs
- Fields.Name: the name of the field
- Fields.Type: the type of the field (all golang types eg.
