package expr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_EvalName(t *testing.T) {
	var mDeployNode DeploymentNode
	mDeployNode = DeploymentNode{}
	mDeploymentNodes := make([]*DeploymentNode, 1)
	mDeploymentNodes[0] = &mDeployNode
	m := Model{}
	tests := []struct {
		want string
	}{
		{want: "model"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got := m.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_Parent(t *testing.T) {
	mSoftwareSystem := SoftwareSystem{}
	mContainer := Container{}
	mComponent := Component{}
	mPerson := Person{}
	t.Parallel()
	tests := []struct {
		eh   ElementHolder
		want ElementHolder
	}{
		{eh: &mSoftwareSystem, want: nil},
		{eh: &mPerson, want: nil},
		{eh: &mContainer, want: &mSoftwareSystem},
		{eh: &mComponent, want: &mContainer},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := Parent(tt.eh)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_Person(t *testing.T) {
	var mNilPerson Person
	mPersonFoo := Person{
		Element: &Element{Name: "JohnDoe"},
	}
	mPeople := make([]*Person, 1)
	mPeople[0] = &mPersonFoo
	m := Model{
		People: mPeople,
	}
	t.Parallel()
	tests := []struct {
		name string
		want Person
	}{
		{name: "JohnDoe", want: mPersonFoo},
		{name: "JulieSmith", want: mNilPerson},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.Person(tt.name)
			if i == 0 {
				if *got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			}
			if i == 1 {
				if got == &tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_SoftwareSystem(t *testing.T) {
	var mNilSystem *SoftwareSystem
	mBigBankSystem := SoftwareSystem{
		Element: &Element{Name: "BigBank"},
	}
	mSystem := make([]*SoftwareSystem, 1)
	mSystem[0] = &mBigBankSystem
	m := Model{
		Systems: mSystem,
	}
	t.Parallel()
	tests := []struct {
		name string
		want *SoftwareSystem
	}{
		{name: "BigBank", want: &mBigBankSystem},
		{name: "CornerShop", want: mNilSystem},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.SoftwareSystem(tt.name)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DeploymentNode(t *testing.T) {
	var mNilDeploymentNode *DeploymentNode
	mDeploymentNode := DeploymentNode{
		Element:     &Element{Name: "MainServerBank"},
		Environment: "Google Cloud",
	}
	mBackupDeploymentNode := DeploymentNode{
		Element:     &Element{Name: "BackupServerBank"},
		Environment: "Google Cloud",
	}
	mDPNode := make([]*DeploymentNode, 2)
	mDPNode[0] = &mDeploymentNode
	mDPNode[1] = &mBackupDeploymentNode
	m := Model{
		DeploymentNodes: mDPNode,
	}
	t.Parallel()
	tests := []struct {
		name        string
		environment string
		want        *DeploymentNode
	}{
		{name: "MainServerBank", environment: "Google Cloud", want: &mDeploymentNode},
		{name: "BackupServerBank", environment: "AWS", want: mNilDeploymentNode},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.DeploymentNode(tt.environment, tt.name)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Validate_DuplicatePeople(t *testing.T) {

	mPeople := make([]*Person, 2)
	mPeople[0] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}
	mPeople[1] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}
	mBigBankSystem := SoftwareSystem{
		Element: &Element{Name: "BigBank"},
	}
	mSystem := make([]*SoftwareSystem, 1)
	mSystem[0] = &mBigBankSystem
	//mSystem[1] = &mBigBankSystem

	mDuplicatePeople := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	/*mPeople1 := make([]*Person, 1)
	mPeople1[0] = &Person{Element: &Element{Name: "Julie"}, Location: LocationExternal}*/
	/*mDuplicateSystems := Model{
		People:  mPeople1,
		Systems: mSystem,
	}*/
	duplicate_person_verr := errors.New("person \"Brian\": name already in use")
	t.Parallel()
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownPerson", model: mDuplicatePeople, want: duplicate_person_verr},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tt.model.Validate()
			// the types are awkward - arrange to pick the underlying strings out as these are the important
			// issues to check for equality
			mgot := string((got.Error()))
			mwant := string((tt.want.Error()))
			if mgot != mwant {
				t.Errorf("got %v, want %v", mgot, mwant)
			}
		})
	}
}

func Test_Validate_DuplicateSystems(t *testing.T) {

	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}
	mBigBankSystem := SoftwareSystem{
		Element: &Element{Name: "BigBank"},
	}
	mSystem := make([]*SoftwareSystem, 2)
	mSystem[0] = &mBigBankSystem
	mSystem[1] = &mBigBankSystem

	mDuplicateSystems := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	duplicate_person_verr := errors.New("software system \"BigBank\": name already in use")
	t.Parallel()
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownSystem", model: mDuplicateSystems, want: duplicate_person_verr},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tt.model.Validate()
			// the types are awkward - arrange to pick the underlying strings out as these are the important
			// issues to check for equality
			mgot := string((got.Error()))
			mwant := string((tt.want.Error()))
			if mgot != mwant {
				t.Errorf("got %v, want %v", mgot, mwant)
			}
		})
	}
}

/*
Yet to complete - commented out to push interim to github
func Test_Validate_DuplicateContainers(t *testing.T) {

	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}

	mComponents := Components{} //make([]Component,1)
	mContainer := make([]Container, 2)
	mContainer[0] = Container{
		Element:    &Element{Name: "Box"},
		Components: mComponents,
	}
	mContainer[1] = Container{
		Element:    &Element{Name: "Box"},
		Components: mComponents,
	}

	mBigBankSystem := SoftwareSystem{
		Element:    &Element{Name: "BigBank"},
		Containers: mContainer,
	}
	mSystem := make([]*SoftwareSystem, 2)
	mSystem[0] = &mBigBankSystem

	mDuplicateSystems := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	duplicate_person_verr := errors.New("container \"BigBank\": name already in use")
	t.Parallel()
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownSystem", model: mDuplicateSystems, want: duplicate_person_verr},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tt.model.Validate()
			// the types are awkward - arrange to pick the underlying strings out as these are the important
			// issues to check for equality
			mgot := string((got.Error()))
			mwant := string((tt.want.Error()))
			if mgot != mwant {
				t.Errorf("got %v, want %v", mgot, mwant)
			}
		})
	}
}
*/
