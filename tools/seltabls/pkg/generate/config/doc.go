// Package config contains the config file generation logic.
//
// It includes functions for generating the contents config file.
// More specifically, it includes functions for generating the sections
// that are used within the config file.
//
// A config file is any file that has the name "{name}_seltabl.yaml".
// It is used to configure the seltabl package for a given struct.
// It ensures that each struct has a unique config file, allows for
// repeatedly regenerating of the struct, and provides a way to
// document the history of the struct in code version control.
package config
