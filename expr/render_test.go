package expr

import (
	"fmt"
	"testing"
)

func makeSimModel() (*Model, *SoftwareSystem) {
	mPerson := Person{
		Element: &Element{Name: "Person"},
	}
	mPeople := make([]*Person, 1)
	mPeople[0] = &mPerson
	mContainers := make([]*Container, 1)
	mContainers[0] = &Container{
		Element: &Element{Name: "SubContainer"},
	}
	mSoftwareSystem := SoftwareSystem{
		Element:    &Element{Name: "SoftwareSystem"},
		Containers: mContainers,
	}
	Identify(&mSoftwareSystem)
	mSystems := make([]*SoftwareSystem, 1)
	mSystems[0] = &mSoftwareSystem

	mDeploymentNodes := make([]*DeploymentNode, 1)
	mDeploymentNodes[0] = &DeploymentNode{
		Element:     &Element{Name: "Fields"},
		Environment: "GreenField",
	}
	mNewModel := Model{
		Systems:         mSystems,
		People:          mPeople,
		DeploymentNodes: mDeploymentNodes,
	}
	return &mNewModel, &mSoftwareSystem
}

func makeViewProps() ViewProps {
	mElementViews := make([]*ElementView, 1)
	mElementViews[0] = &ElementView{
		Element: &Element{ID: "1", Name: "Element"},
	}
	mViewProps := ViewProps{
		ElementViews: mElementViews,
	}
	return mViewProps
}

func makeComponents(mSoftwareSystem *SoftwareSystem) Components {
	mRelationships := make([]*Relationship, 1)
	mSource := Container{
		Element: &Element{Name: "Source"},
		System:  mSoftwareSystem,
	}
	Identify(&mSource)
	mDestination := Container{
		Element: &Element{Name: "Destination"},
		System:  mSoftwareSystem,
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
	}
	mContainer2 := Container{
		Element: &Element{Name: "ComponentContainer"},
		System:  mSoftwareSystem,
	}
	mComponent := Component{
		Element: &Element{
			ID:            "1",
			Name:          "Component",
			Relationships: mRelationships,
		},
		Container: &mContainer2,
	}
	mComponents := make([]*Component, 1)
	mComponents[0] = &mComponent
	Identify(&mComponent)
	return mComponents
}

func Test_AddAllElements(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel()
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps()
	mLandscapeView := LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := ContextView{
		ViewProps: &mViewProps,
	}

	mContainer := Container{
		Element: &Element{Name: "Container"},
		System:  mSoftwareSystem,
	}
	Identify(&mContainer)
	mComponentView := ComponentView{
		ViewProps:   &mViewProps,
		ContainerID: mContainer.ID,
	}
	mDeploymentView := DeploymentView{
		ViewProps:   &mViewProps,
		Environment: "GreenField",
	}
	tests := []struct {
		in View
	}{
		{in: &mLandscapeView},
		{in: &mContextView},
		{in: &mComponentView},
		{in: &mDeploymentView},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			addAllElements(tt.in)
		})
	}
}

func Test_AddDefaultElements(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel()
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps()
	mLandscapeView := LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := ContextView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
	}
	mComponents := makeComponents(mSoftwareSystem)
	mContainer := Container{
		Element:    &Element{Name: "Container"},
		Components: mComponents,
		System:     mSoftwareSystem,
	}
	Identify(&mContainer)
	mComponentView := ComponentView{
		ViewProps:   &mViewProps,
		ContainerID: mContainer.ID,
	}
	mDeploymentView := DeploymentView{
		ViewProps:   &mViewProps,
		Environment: "GreenField",
	}
	mContainerView := ContainerView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
		AddInfluencers:   true,
	}
	tests := []struct {
		in View
	}{
		{in: &mLandscapeView},
		{in: &mContextView},
		{in: &mComponentView},
		{in: &mContainerView},
		{in: &mDeploymentView},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			addDefaultElements(tt.in)
		})
	}
}

func Test_AddNeighbors(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel()
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps()
	mLandscapeView := LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := ContextView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
	}
	mComponents := makeComponents(mSoftwareSystem)
	mContainer := Container{
		Element:    &Element{Name: "Container"},
		Components: mComponents,
		System:     mSoftwareSystem,
	}
	Identify(&mContainer)
	mComponentView := ComponentView{
		ViewProps:   &mViewProps,
		ContainerID: mContainer.ID,
	}
	mDeploymentView := DeploymentView{
		ViewProps:   &mViewProps,
		Environment: "GreenField",
	}
	mContainerView := ContainerView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
		AddInfluencers:   true,
	}
	tests := []struct {
		el *Element
		in View
	}{
		{el: mContainer.Element, in: &mLandscapeView},
		{el: mContainer.Element, in: &mContextView},
		{el: mContainer.Element, in: &mComponentView},
		{el: mContainer.Element, in: &mContainerView},
		{el: mContainer.Element, in: &mDeploymentView},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			addNeighbors(tt.el, tt.in)
		})
	}
}

func makeRelatedPeople() (Person, Person) {
	mRelationships := make([]*Relationship, 1)
	mSource := Person{
		Element: &Element{Name: "SourcePerson"},
	}
	Identify(&mSource)
	mDestination := Person{
		Element: &Element{Name: "DestinationPerson"},
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Family",
	}
	Identify(mRelationships[0])
	return mSource, mDestination
}

func Test_RelatedPeople(t *testing.T) {
	mSource, mDestination := makeRelatedPeople()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedPeople(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeRelatedSystems() (SoftwareSystem, SoftwareSystem) {
	mRelationships := make([]*Relationship, 1)
	mSource := SoftwareSystem{
		Element: &Element{Name: "SourceSoftwareSystem"},
	}
	Identify(&mSource)
	mDestination := SoftwareSystem{
		Element: &Element{Name: "DestinationSoftwareSystem"},
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Empire",
	}
	Identify(mRelationships[0])
	return mSource, mDestination
}
func Test_RelatedSoftwareSystems(t *testing.T) {
	mSource, mDestination := makeRelatedSystems()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedSoftwareSystems(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeRelatedContainers() (Container, Container) {
	mRelationships := make([]*Relationship, 1)
	mSystem := SoftwareSystem{
		Element: &Element{Name: "SourceSoftwareSystem"},
	}
	Identify(&mSystem)
	mSource := Container{
		Element: &Element{Name: "SourceContainer"},
		System:  &mSystem,
	}
	Identify(&mSource)
	mDestination := Container{
		Element: &Element{Name: "DestinationContainer"},
		System:  &mSystem,
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Group",
	}
	Identify(mRelationships[0])
	return mSource, mDestination
}
func Test_RelatedContainers(t *testing.T) {
	mSource, mDestination := makeRelatedContainers()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedContainers(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeRelatedComponents() (Component, Component) {
	mRelationships := make([]*Relationship, 1)
	mSystem := SoftwareSystem{
		Element: &Element{Name: "SourceSoftwareSystem"},
	}
	Identify(&mSystem)

	mSrcContainer := Container{
		Element: &Element{Name: "SourceContainer"},
		System:  &mSystem,
	}
	Identify(&mSrcContainer)
	mDstContainer := Container{
		Element: &Element{Name: "SourceContainer"},
		System:  &mSystem,
	}
	Identify(&mDstContainer)
	mSource := Component{
		Element:   &Element{Name: "SourceComponent"},
		Container: &mSrcContainer,
	}
	Identify(&mSource)
	mDestination := Component{
		Element:   &Element{Name: "DestinationComponent"},
		Container: &mDstContainer,
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Group",
	}
	Identify(mRelationships[0])
	return mSource, mDestination
}

func Test_RelatedComponents(t *testing.T) {
	mSource, mDestination := makeRelatedComponents()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedComponents(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeRelatedInfrastructureNodes() (InfrastructureNode, InfrastructureNode) {
	mRelationships := make([]*Relationship, 1)

	mDeploymentNode := DeploymentNode{
		Element:     &Element{Name: "DemploymentNode"},
		Environment: "BrownField",
	}
	Identify(&mDeploymentNode)
	mSource := InfrastructureNode{
		Element:     &Element{Name: "SourceInfrastructureNode"},
		Parent:      &mDeploymentNode,
		Environment: "BrownField",
	}
	Identify(&mSource)
	mDestination := InfrastructureNode{
		Element:     &Element{Name: "DestinationInfrastructureNode"},
		Parent:      &mDeploymentNode,
		Environment: "BrownField",
	}
	Identify(&mDestination)
	mRelationships[0] = &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Group",
	}
	Identify(mRelationships[0])
	return mSource, mDestination
}
func Test_RelatedInfrastructureNodes(t *testing.T) {
	mSource, mDestination := makeRelatedInfrastructureNodes()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedInfrastructureNodes(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeRelatedContainerInstances() (ContainerInstance, ContainerInstance) {
	mRelationships := make([]*Relationship, 1)
	mSystem := SoftwareSystem{
		Element: &Element{Name: "SourceSoftwareSystem"},
	}
	Identify(&mSystem)

	mSrcContainer := Container{
		Element: &Element{Name: "SourceContainer"},
		System:  &mSystem,
	}
	Identify(&mSrcContainer)
	mDstContainer := Container{
		Element: &Element{Name: "SourceContainer"},
		System:  &mSystem,
	}
	Identify(&mDstContainer)

	mSrcDeploymentNode := DeploymentNode{
		Element:     &Element{Name: "SourceContainer"},
		Environment: "mSrcDeploymentNode",
	}
	Identify(&mSrcDeploymentNode)
	mDestDeploymentNode := DeploymentNode{
		Element:     &Element{Name: "DestinationContainer"},
		Environment: "mDestDeploymentNode",
	}
	Identify(&mDestDeploymentNode)

	mSrcContainerInstance := ContainerInstance{
		Element:     &Element{Name: "SourceContainerInstance"},
		Parent:      &mSrcDeploymentNode,
		ContainerID: mSrcContainer.ID,
	}
	Identify(&mSrcContainerInstance)
	mDstContainerInstance := ContainerInstance{
		Element:     &Element{Name: "SourceContainerInstance"},
		Parent:      &mDestDeploymentNode,
		ContainerID: mDstContainer.ID,
	}
	Identify(&mDstContainerInstance)
	mRelationships[0] = &Relationship{
		Source:      mSrcContainerInstance.Element,
		Destination: mDstContainerInstance.Element,
		Description: "Group",
	}
	Identify(mRelationships[0])
	return mSrcContainerInstance, mDstContainerInstance
}

func Test_RelatedContainerInstances(t *testing.T) {
	mSource, mDestination := makeRelatedContainerInstances()
	tests := []struct {
		el   *Element
		want *Element
	}{
		{el: mSource.Element, want: mDestination.Element},
		{el: mDestination.Element, want: mSource.Element},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := relatedContainerInstances(tt.el)
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}
		})
	}

}

func makeContainerView(mSystem SoftwareSystem) (ContainerView, []*Relationship) {
	mContainers := make([]*Container, 1)
	mContainers[0] = &Container{
		Element: &Element{Name: "SubContainer"},
	}
	mOtherSoftwareSystem := SoftwareSystem{
		Element:    &Element{Name: "SoftwareSystem"},
		Containers: mContainers,
	}
	Identify(&mOtherSoftwareSystem)
	mRelationshipView := RelationshipView{
		Source:      mOtherSoftwareSystem.Element,
		Destination: mSystem.Element,
	}
	mRelationshipViews := make([]*RelationshipView, 1)
	mRelationshipViews[0] = &mRelationshipView
	mViewProps := ViewProps{
		RelationshipViews: mRelationshipViews,
	}
	mContainerView := ContainerView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSystem.ID,
	}
	mRelationship := Relationship{
		Source:      mOtherSoftwareSystem.Element,
		Destination: mSystem.Element,
	}
	mRelationships := make([]*Relationship, 1)
	mRelationships[0] = &mRelationship

	return mContainerView, mRelationships
}

func Test_AddInfluencers(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel()

	mContainerView, mRelationships := makeContainerView(*mSoftwareSystem)
	mModel.Systems[0].Relationships = mRelationships
	mSoftwareSystem.Relationships = mRelationships

	tests := []struct {
		el *ContainerView
	}{
		{el: &mContainerView},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			addInfluencers(tt.el)
			/*got := cv
			if got != nil {
				if got[0].Name != tt.want.Name {
					t.Errorf("Got %s, wanted %s", got[0].Name, tt.want.Name)
				}
			}*/
		})
	}

}
