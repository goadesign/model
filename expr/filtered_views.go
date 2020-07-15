package expr

import (
	"fmt"
)

type (
	// FilteredView describes a filtered view on top of a specified view.
	FilteredView struct {
		// Title of the view.
		Title string `json:"title,omitempty"`
		// Description of view.
		Description string `json:"description,omitempty"`
		// Key used to refer to the view.
		Key string `json:"key"`
		// BaseKey is the key of the view on which this filtered view is based.
		BaseKey string `json:"baseViewKey"`
		// Whether elements/relationships are being included ("Include") or
		// excluded ("Exclude") based upon the set of tags.
		Mode string `json:"mode"`
		// The set of tags to include/exclude elements/relationships when
		// rendering this filtered view.
		Tags []string `json:"tags,omitempty"`
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
