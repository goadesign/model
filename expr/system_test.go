package expr

import (
	"fmt"
	"testing"
)

func Test_AddContainer(t *testing.T) {

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{ID: "1", Name: "Software system"},
	}
	Identify(&mSoftwareSystem)
	mContainer := Container{
		Element: &Element{ID: "1", Name: "Container", Tags: "One,Two"},
		System:  &mSoftwareSystem,
	}

	mNewContainer := Container{
		Element: &Element{ID: "1",
			Name:        "Container",
			Description: "Uncertain",
			Technology:  "CottonOnString",
			URL:         "microsoft.com",
			Tags:        "Three,Four",
		},
		System: &mSoftwareSystem,
	}

	tests := []struct {
		in   *Container
		want *Container
	}{
		{in: &mContainer, want: &mNewContainer},
		{in: &mNewContainer, want: &mNewContainer},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := mSoftwareSystem.AddContainer(tt.in)
			if i == 0 {
				if got.Name != mContainer.Name {
					t.Errorf("Got %s, wanted %s", got.Name, tt.want.Name)
				}
			} else {
				if got.Description != mContainer.Description {
					t.Errorf("Got %s, wanted %s", got.Description, tt.want.Description)
				}
				if got.Technology != mContainer.Technology {
					t.Errorf("Got %s, wanted %s", got.Technology, tt.want.Technology)
				}
			}
		})
	}
}
