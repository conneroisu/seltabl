package domain

// Config is a struct for the config
type Config struct {
	Cfgs []StructConfig `json:"configs" yaml:"configs" toml:"configs"`
}

// StructConfig is a struct for the config
// url1 and url2 are the urls to compare for differences
// information is the information that can be extracted from the two urls
type StructConfig struct {
	// FileName is the name of the file to write the struct to
	FileName string `json:"file_name"   yaml:"file_name"   toml:"file_name"`
	// URL1 is the first url to compare
	URL1 string `json:"url1"        yaml:"url1"        toml:"url1"`
	// URL2 is the second url to compare
	URL2 string `json:"url2"        yaml:"url2"        toml:"url2"`
	// Information is the information that can be extracted from the two urls
	Information []string `json:"information" yaml:"information" toml:"information"`
	// Prompt is the prompt to use for the struct
	Prompt string `json:"prompt"      yaml:"prompt"      toml:"prompt"`
}
