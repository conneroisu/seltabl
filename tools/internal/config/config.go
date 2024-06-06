package config

// Config is the configuration for the seltabl command
type Config struct {
	InputFilePath  string `json:"inputFile"`
	OutputFilePath string `json:"outputFile"`
	URL            string `json:"url"`
	PackageName    string `json:"packageName"`
	StructName     string `json:"structName"`
}
