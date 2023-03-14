package expr

import (
	"fmt"
	"testing"
)

func Test_AddAllElements(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel(t)
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps(t)
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
	mModel, mSoftwareSystem := makeSimModel(t)
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps(t)
	mLandscapeView := LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := ContextView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
	}
	mComponents := makeComponents(mSoftwareSystem, t)
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
	mModel, mSoftwareSystem := makeSimModel(t)
	Root = &Design{
		Model: mModel,
	}
	mViewProps := makeViewProps(t)
	mLandscapeView := LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := ContextView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
	}
	mComponents := makeComponents(mSoftwareSystem, t)
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

func Test_RelatedPeople(t *testing.T) {
	mSource, mDestination := makeRelatedPeople(t)
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

func Test_RelatedSoftwareSystems(t *testing.T) {
	mSource, mDestination := makeRelatedSystems(t)
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

func Test_RelatedContainers(t *testing.T) {
	mSource, mDestination := makeRelatedContainers(t)
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

func Test_RelatedComponents(t *testing.T) {
	mSource, mDestination := makeRelatedComponents(t)
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

func Test_RelatedInfrastructureNodes(t *testing.T) {
	mSource, mDestination := makeRelatedInfrastructureNodes(t)
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

func Test_RelatedContainerInstances(t *testing.T) {
	mSource, mDestination := makeRelatedContainerInstances(t)
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

func Test_AddInfluencers(t *testing.T) {
	mModel, mSoftwareSystem := makeSimModel(t)

	mContainerView, mRelationships := makeContainerView(*mSoftwareSystem, t)
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
		})
	}

}
