package main

// SeltablConfig is the configuration for the seltabl command
type SeltablConfig struct {
	InputFile   string `json:"inputFile"`
	OutputFile  string `json:"outputFile"`
	URL         string `json:"url"`
	PackageName string `json:"packageName"`
	StructName  string `json:"structName"`
}
