package domain

// Prompt is a struct for the Prompt.
type Prompt struct{}

// NewStructContentArgs is a struct for the NewStructContentArgs.
//
//	type {{ .Name }} struct {
//	        {{ range $i, $field := .Fields }}
//	                // {{ $field.Name }} - {{ $field.Description }}
//	                {{ $field.Name }} {{ $field.Type }} `json:"{{ $field.Name }}" hSel:"{{ $field.HeaderSelector }}" dSel:"{{ $field.DataSelector }}" ctl:"{{ $field.ControlSelector }}" {{ if $field.MustBePresent }}must:"{{ $field.MustBePresent }}"{{ end }}`
//	        {{ end }}
//	}}
//
// template name: schema_struct
type NewStructContentArgs struct {
	Prompt         `prompt:"struct_content"`
	URL            string   // required
	Name           string   // required
	IgnoreElements []string // required
	Fields         []Field  // required
}

// NewAggregateStuctPromptArgs is a struct for the NewAggregateStuctPromptArgs.
//
// - header selector: used to find the header row and column for the field in the given struct.
// - data-selector: used to find the data column for the field in the given struct.
// - query-selector: used to query for the inner text or attribute of the cell.
// - control-selector: used to control what to query for the inner text or attribute of the cell.
// - must-be-present: used to ensure that the field is present in the given url.
//
// template name: schema_aggregate
type NewAggregateStuctPromptArgs struct {
	Prompt      `prompt:"schema_aggregate"`
	URL         string   // required
	HTMLContent string   // required
	Selectors   []string // required
	Schemas     []string // required
}

// NewStructFileArgs is a struct for the NewStructFileArgs.
type NewStructFileArgs struct {
	Prompt      `prompt:"struct_file"`
	PackageName string // required
	URL         string // required
}

// NewSelectorPromptArgs is a struct for the NewSelectorPromptArgs.
//
// The different selector name and descriptions are:
// - header-selector: used to find the header row and column for the field in the given struct.
// - data-selector: used to find the data column for the field in the given struct.
// - query-selector: used to query for the inner text or attribute of the cell.
// - control-selector: used to control what to query for the inner text or attribute of the cell.
// - must-be-present: used to ensure that the field is present in the given url.
//
// template name: schema_prompt
type NewSelectorPromptArgs struct {
	Prompt              `prompt:"schema_prompt"`
	URL                 string   // required
	SelectorName        string   // required
	SelectorDescription string   // required
	Content             string   // required
	Selectors           []string // required
}

// NewErrorPromptArgs is a struct for the NewErrorPromptArgs.
type NewErrorPromptArgs struct {
	Prompt `prompt:"error_prompt"`
	Error  error  // required
	Out    string // required
}

// NewErrorAggregatePromptArgs is a struct for the NewErrorAggregatePromptArgs.
type NewErrorAggregatePromptArgs struct {
	Prompt `prompt:"error_aggregate"`
	Errors []error // required
	Out    string  // required
}
