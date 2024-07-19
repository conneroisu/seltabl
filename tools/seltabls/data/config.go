package data

// ConfigFile is a struct for a config file.
type ConfigFile struct {
	URLs struct {
		URL string `json:"url" yaml:"url"`
	} `json:"urls" yaml:"urls"`
}
