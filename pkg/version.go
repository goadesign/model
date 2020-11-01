package model

import (
	"fmt"
)

const (
	// Major version number
	Major = 1
	// Minor version number
	Minor = 7
	// Build number
	Build = 0
	// Suffix - set to empty string in release tag commits.
	Suffix = ""
)

// Version returns the complete version number.
func Version() string {
	if Suffix != "" {
		return fmt.Sprintf("v%d.%d.%d-%s", Major, Minor, Build, Suffix)
	}
	return fmt.Sprintf("v%d.%d.%d", Major, Minor, Build)
}
