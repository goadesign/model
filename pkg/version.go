package structurizr

import (
	"fmt"
	"regexp"
)

const (
	// Major version number
	Major = 0
	// Minor version number
	Minor = 0
	// Build number
	Build = 11
	// Suffix - set to empty string in release tag commits.
	Suffix = ""
)

var (
	// Version format
	versionFormat = regexp.MustCompile(`v(\d+?)\.(\d+?)\.(\d+?)(?:-.+)?`)
)

// Version returns the complete version number.
func Version() string {
	if Suffix != "" {
		return fmt.Sprintf("v%d.%d.%d-%s", Major, Minor, Build, Suffix)
	}
	return fmt.Sprintf("v%d.%d.%d", Major, Minor, Build)
}
