package expr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Don't use t.parallel. Crashes when file tests are run as concurrent map read & write.
func Test_EvalName(t *testing.T) {
	var mDeployNode = DeploymentNode{}
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

	mDuplicatePeople := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	duplicatePersonVerr := errors.New("person \"Brian\": name already in use")
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownPerson", model: mDuplicatePeople, want: duplicatePersonVerr},
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
	duplicatePersonVerr := errors.New("software system \"BigBank\": name already in use")
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownSystem", model: mDuplicateSystems, want: duplicatePersonVerr},
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

func Test_Validate_DuplicateContainers(t *testing.T) {

	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}

	mComponents := Components{}
	mContainers := make([]*Container, 2)
	mContainers[0] = &Container{
		Element:    &Element{Name: "Box"},
		Components: mComponents,
	}
	mContainers[1] = &Container{
		Element:    &Element{Name: "Box"},
		Components: mComponents,
	}

	mBigBankSystem := SoftwareSystem{
		Element:    &Element{Name: "BigBank"},
		Containers: mContainers,
	}
	mSystem := make([]*SoftwareSystem, 1)
	mSystem[0] = &mBigBankSystem

	mDuplicateSystems := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	duplicatePersonVerr := errors.New("container \"Box\": name already in use")
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownSystem", model: mDuplicateSystems, want: duplicatePersonVerr},
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

func Test_Validate_DuplicateComponents(t *testing.T) {

	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{Element: &Element{Name: "Brian"}, Location: LocationExternal}

	mComponents := make([]*Component, 2)
	mComponents[0] = &Component{
		Element: &Element{Name: "Widget"},
	}
	mComponents[1] = &Component{
		Element: &Element{Name: "Widget"},
	}

	mContainers := make([]*Container, 1)
	mContainers[0] = &Container{
		Element:    &Element{Name: "Box"},
		Components: mComponents,
	}

	mBigBankSystem := SoftwareSystem{
		Element:    &Element{Name: "BigBank"},
		Containers: mContainers,
	}
	mSystem := make([]*SoftwareSystem, 1)
	mSystem[0] = &mBigBankSystem

	mDuplicateSystems := Model{
		People:  mPeople,
		Systems: mSystem,
	}
	duplicatePersonVerr := errors.New("component \"Widget\": name already in use")
	tests := []struct {
		name  string
		model Model
		want  error
	}{
		{name: "KnownSystem", model: mDuplicateSystems, want: duplicatePersonVerr},
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

func Test_FindElement(t *testing.T) {
	mComponents := make([]*Component, 2)
	mComponents[0] = &Component{
		Element: &Element{Name: "Widget"},
	}
	mComponents[1] = &Component{
		Element: &Element{Name: "Truit"},
	}
	mContainers := make([]*Container, 1)
	mContainers[0] = &Container{
		Element:    &Element{Name: "Mainframe"},
		Components: mComponents,
	}

	mBigBankSystem := SoftwareSystem{
		Element:    &Element{Name: "BigBank"},
		Containers: mContainers,
	}
	mSystem := make([]*SoftwareSystem, 1)
	mSystem[0] = &mBigBankSystem
	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{
		Element: &Element{Name: "Brian"},
	}
	m := Model{
		Systems: mSystem,
		People:  mPeople,
	}

	tests := []struct {
		eh    ElementHolder
		path  string
		want  ElementHolder
		want2 error
	}{
		{eh: mContainers[0], path: "Widget", want: mComponents[0], want2: nil},
		{eh: &mBigBankSystem, path: "Mainframe", want: mContainers[0], want2: nil},
		{eh: mPeople[0], path: "Brian", want: mPeople[0], want2: nil},
		{eh: mComponents[0], path: "Widget", want: mComponents[0], want2: nil},
		{eh: mComponents[0], path: "BigBank", want: &mBigBankSystem, want2: nil},
		{eh: mContainers[0], path: "BigBank/Mainframe", want: mContainers[0], want2: nil},
		{eh: &mBigBankSystem, path: "BigBank/Widget", want: mComponents[0], want2: nil},
		{eh: mBigBankSystem, path: "BigBank/Mainframe", want: mContainers[0], want2: nil},
		{eh: &mBigBankSystem, path: "Mainframe/Widget", want: mComponents[0], want2: nil},
		{eh: nil, path: "Mainframe/Widget", want: mComponents[0], want2: nil},
		{eh: mBigBankSystem, path: "BigBank/Mainframe/Widget", want: mComponents[0], want2: nil},
		{eh: nil, path: "SmallBank/Mainframe/Widget", want: mComponents[0], want2: nil},
		{eh: nil, path: "BigBank/Mainframe/Widget/Item", want: mComponents[0], want2: nil},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := m.FindElement(tt.eh, tt.path)
			if got != tt.want && err == nil {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_AddPerson(t *testing.T) {
	mPeople := make([]*Person, 1)
	mPeople[0] = &Person{
		Element: &Element{Name: "Brian"},
	}
	BrianFull := &Person{
		Element: &Element{Name: "Brian", Description: "Stevedore"},
	}
	Steve := &Person{
		Element: &Element{Name: "Steve"},
	}
	m := Model{
		People: mPeople,
	}
	tests := []struct {
		Person *Person
		want   *Person
	}{
		{Person: mPeople[0], want: mPeople[0]},
		{Person: BrianFull, want: mPeople[0]},
		{Person: Steve, want: Steve},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.AddPerson(tt.Person)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_AddSystem(t *testing.T) {
	mSoftwareSystems := make([]*SoftwareSystem, 1)
	mSoftwareSystems[0] = &SoftwareSystem{
		Element: &Element{Name: "Mainframe", DSLFunc: func() {
			name := "hello"
			_ = name
		}},
	}
	Mainframe := &SoftwareSystem{
		Element: &Element{Name: "Mainframe", Description: "CDC64", DSLFunc: func() {
			name := "world"
			_ = name
		}},
	}

	Tablet := &SoftwareSystem{
		Element: &Element{Name: "iPad"},
	}
	m := Model{
		Systems: mSoftwareSystems,
	}

	tests := []struct {
		System *SoftwareSystem
		want   *SoftwareSystem
	}{
		{System: Tablet, want: Tablet},
		{System: Mainframe, want: mSoftwareSystems[0]},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.AddSystem(tt.System)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_AddDeployment(t *testing.T) {
	mDeploymentNodes := make([]*DeploymentNode, 1)
	mDeploymentNodes[0] = &DeploymentNode{
		Element: &Element{Name: "Corner Shop", DSLFunc: func() {
			name := "hello"
			_ = name
		}},
	}
	Tesco := &DeploymentNode{
		Element: &Element{Name: "Corner Shop", Description: "Tesco Express", Technology: "WAN Connection", DSLFunc: func() {
			name := "world"
			_ = name
		}},
	}

	Aldi := &DeploymentNode{
		Element: &Element{Name: "Parade"},
	}
	m := Model{
		DeploymentNodes: mDeploymentNodes,
	}

	tests := []struct {
		DPNode *DeploymentNode
		want   *DeploymentNode
	}{
		{DPNode: Aldi, want: Aldi},
		{DPNode: Tesco, want: mDeploymentNodes[0]},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := m.AddDeploymentNode(tt.DPNode)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_AddImpliedRelationship(t *testing.T) {
	// define lineage A
	mAComponents := make([]*Component, 1)
	mAComponents[0] = &Component{
		Element: &Element{ID: "1", Name: "WidgetA"},
	}
	mAContainers := make([]*Container, 1)
	mAContainers[0] = &Container{
		Element:    &Element{ID: "2", Name: "MainframeA"},
		Components: mAComponents,
	}

	mBigBankSystem := &SoftwareSystem{
		Element:    &Element{ID: "3", Name: "BigBankA"},
		Containers: mAContainers,
	}
	mASystem := make([]*SoftwareSystem, 1)
	mASystem[0] = mBigBankSystem
	mAContainers[0].System = mBigBankSystem
	mAComponents[0].Container = mAContainers[0]
	// define lineage B
	mBComponents := make([]*Component, 1)
	mBComponents[0] = &Component{
		Element: &Element{ID: "4", Name: "WidgetB"},
	}
	mBContainers := make([]*Container, 1)
	mBContainers[0] = &Container{
		Element:    &Element{ID: "5", Name: "MainframeB"},
		Components: mAComponents,
	}

	mBigBankB := &SoftwareSystem{
		Element:    &Element{ID: "6", Name: "BigBankB"},
		Containers: mBContainers,
	}
	mBSystem := make([]*SoftwareSystem, 1)
	mBSystem[0] = mBigBankB
	mBContainers[0].System = mBigBankB
	mBComponents[0].Container = mBContainers[0]

	// Register the pieces
	Identify(mASystem[0])
	Identify(mAContainers[0])
	Identify(mAComponents[0])
	Identify(mBSystem[0])
	Identify(mBContainers[0])
	Identify(mBComponents[0])
	nilRelationship := Relationship{}

	tests := []struct {
		src      ElementHolder
		dst      *Element
		existing *Relationship
	}{
		{src: mAContainers[0], dst: mBComponents[0].Element, existing: &nilRelationship},
		{src: mASystem[0], dst: mBContainers[0].Element, existing: &nilRelationship},
		{src: mASystem[0], dst: mAContainers[0].Element, existing: &nilRelationship},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
			addImpliedRelationships(tt.src, tt.dst, tt.existing)
		})
	}
}

// Note the vast majority of Finalize is actually exercised from registry.go. So it's the registry_test.go that should cause the code coverage...
func TestModelFinalize(t *testing.T) {
	// define lineage A
	mAComponents := make([]*Component, 1)
	mAComponents[0] = &Component{
		Element: &Element{ID: "1", Name: "WidgetA"},
	}
	mAContainers := make([]*Container, 1)
	mAContainers[0] = &Container{
		Element:    &Element{ID: "2", Name: "MainframeA"},
		Components: mAComponents,
	}

	mBigBankSystem := &SoftwareSystem{
		Element:    &Element{ID: "3", Name: "BigBankA"},
		Containers: mAContainers,
	}
	mASystem := make([]*SoftwareSystem, 1)
	mASystem[0] = mBigBankSystem
	mAContainers[0].System = mBigBankSystem
	mAComponents[0].Container = mAContainers[0]
	// define lineage B
	mBComponents := make([]*Component, 1)
	mBComponents[0] = &Component{
		Element: &Element{ID: "4", Name: "WidgetB"},
	}
	mBContainers := make([]*Container, 1)
	mBContainers[0] = &Container{
		Element:    &Element{ID: "5", Name: "MainframeB"},
		Components: mAComponents,
	}

	mBigBankB := &SoftwareSystem{
		Element:    &Element{ID: "6", Name: "BigBankB"},
		Containers: mBContainers,
	}
	mBSystem := make([]*SoftwareSystem, 2)
	mBSystem[0] = mBigBankB
	mBSystem[1] = mBigBankSystem
	mBContainers[0].System = mBigBankB
	mBComponents[0].Container = mBContainers[0]

	mRelationshipA1 := Relationship{
		ID:          "R123",
		Source:      mAContainers[0].Element,
		Destination: mBContainers[0].Element,
	}
	mRelationshipA2 := Relationship{
		ID:          "R123",
		Source:      mBContainers[0].Element,
		Destination: mAContainers[0].Element,
	}

	mRelationshipA3 := Relationship{
		ID:          "R123",
		Source:      mAComponents[0].Element,
		Destination: mBComponents[0].Element,
	}
	mRelationshipA4 := Relationship{
		ID:          "R123",
		Source:      mBComponents[0].Element,
		Destination: mAComponents[0].Element,
	}

	mAContainers[0].Relationships = make([]*Relationship, 1)
	mAContainers[0].Relationships[0] = &mRelationshipA1
	mBContainers[0].Relationships = make([]*Relationship, 1)
	mBContainers[0].Relationships[0] = &mRelationshipA2
	mAComponents[0].Relationships = make([]*Relationship, 1)
	mAComponents[0].Relationships[0] = &mRelationshipA3
	mBComponents[0].Relationships = make([]*Relationship, 1)
	mBComponents[0].Relationships[0] = &mRelationshipA4
	// Register the pieces as if the registry is empty finalize does nothing.
	Identify(mASystem[0])
	Identify(mAContainers[0])
	Identify(mAComponents[0])
	Identify(mBSystem[0])
	Identify(mBContainers[0])
	Identify(mBComponents[0])
	Identify(&mRelationshipA1)
	Identify(&mRelationshipA2)
	Identify(&mRelationshipA3)
	Identify(&mRelationshipA4)

	mSrcDeploymentNode := DeploymentNode{
		Element: &Element{ID: "123SD", Name: "SrcDeploymentNode"},
	}
	mDstDeploymentNode := DeploymentNode{
		Element: &Element{ID: "123DD", Name: "DstDeploymentNode"},
	}

	mContainerInstanceA := ContainerInstance{
		Element:     &Element{ID: "123", Name: "ContainerInstanceA"},
		Parent:      &mSrcDeploymentNode,
		ContainerID: mAContainers[0].ID,
		Environment: "London",
	}

	mContainerInstanceB := ContainerInstance{
		Element:     &Element{ID: "456", Name: "ContainerInstanceB"},
		Parent:      &mDstDeploymentNode,
		ContainerID: mBContainers[0].ID,
		Environment: "Manchester",
	}
	Identify(&mContainerInstanceA)
	Identify(&mContainerInstanceB)

	m := Model{
		Systems:                 mBSystem,
		AddImpliedRelationships: true,
	}

	tests := []struct {
	}{
		{},
	}
	for i := range tests {
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
			m.Finalize()
		})
	}
}
