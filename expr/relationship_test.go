package expr

import (
	"fmt"
	"testing"
)

// Don't use t.parallel. Crashes when file tests are run as concurrent map read & write.

func Test_Relationship_EvalName(t *testing.T) {
	mRelationship := Relationship{
		Source:      &Element{Name: "Source"},
		Destination: &Element{Name: "Destination"},
		Description: "Sample",
	}

	tests := []struct {
		want string
	}{
		{want: "relationship \"Sample\" [Source -> Destination]"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got := mRelationship.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}

}

func Test_Relationship_Finalize(t *testing.T) {
	mRelationship := Relationship{
		Description: "Sample",
		Tags:        "Tag0",
	}

	tests := []struct {
		InteractionStyle InteractionStyleKind
		want             string
	}{
		{InteractionStyle: InteractionAsynchronous, want: "Relationship,Asynchronous"},
		{InteractionStyle: InteractionSynchronous, want: "Relationship"},
	}
	for i, tt := range tests {
		mRelationship.InteractionStyle = tt.InteractionStyle
		mRelationship.Tags = ""
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			mRelationship.Finalize()
			if mRelationship.Tags != tt.want {
				t.Errorf("received %s, wanted %s", mRelationship.Tags, tt.want)
			}
		})
	}
}

func Test_Relationship_MergeTags(t *testing.T) {
	mRelationship := Relationship{
		Source:      &Element{Name: "Source"},
		Destination: &Element{Name: "Destination"},
		Description: "Sample",
		Tags:        "Tag0",
	}

	tests := []struct {
		tag1 string
		tag2 string
	}{
		{tag1: "Tag1", tag2: "Tag2"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			mRelationship.MergeTags(tt.tag1, tt.tag2)
			if mRelationship.Tags != "Tag0,Tag1,Tag2" {
				t.Errorf("had %s, wanted %s", mRelationship.Tags, "Tag0,Tag1,Tag2")
			}
		})
	}
}
