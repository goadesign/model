package expr

import (
	"fmt"
)

type (
	// FilteredView describes a filtered view on top of a specified view.
	FilteredView struct {
		Title       string
		Description string
		Key         string `json:"key"`
		BaseKey     string
		Exclude     bool
		FilterTags  []string
	}
)

// EvalName returns the generic expression name used in error messages.
func (fv *FilteredView) EvalName() string {
	var suffix string
	if fv.Key != "" {
		suffix = fmt.Sprintf(" key %q and", fv.Key)
	}
	return fmt.Sprintf("filtered view with%s base key %q", suffix, fv.BaseKey)
}
