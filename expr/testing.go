package expr

import (
	"testing"
)

func makeViews(mSystem SoftwareSystem, t *testing.T) *Views {
	t.Helper()
	mViews := &Views{}
	mContainers := make([]*Container, 1)
	mContainers[0] = &Container{
		Element: &Element{Name: "SubContainer"},
	}
	return mViews
}

func makeContainerView(mSystem SoftwareSystem, t *testing.T) (ContainerView, []*Relationship) {
	t.Helper()
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

func makeSimModel(t *testing.T) (*Model, *SoftwareSystem) {
	t.Helper()
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

func makeViewProps(t *testing.T) ViewProps {
	t.Helper()
	mElementViews := make([]*ElementView, 1)
	mElementViews[0] = &ElementView{
		Element: &Element{ID: "1", Name: "Element"},
	}
	mViewProps := ViewProps{
		ElementViews: mElementViews,
	}
	return mViewProps
}

func makeComponents(mSoftwareSystem *SoftwareSystem, t *testing.T) Components {
	t.Helper()
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

func makeRelatedPeople(t *testing.T) (Person, Person) {
	t.Helper()
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

func makeRelatedSystems(t *testing.T) (SoftwareSystem, SoftwareSystem) {
	t.Helper()
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

func makeRelatedContainers(t *testing.T) (Container, Container) {
	t.Helper()
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

func makeRelatedComponents(t *testing.T) (Component, Component) {
	t.Helper()
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

func makeRelatedInfrastructureNodes(t *testing.T) (InfrastructureNode, InfrastructureNode) {
	t.Helper()
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

func makeRelatedContainerInstances(t *testing.T) (ContainerInstance, ContainerInstance) {
	t.Helper()
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
