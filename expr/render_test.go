package expr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupTestViewsWithContainer creates common view types with a container for testing view operations
func setupTestViewsWithContainer(t *testing.T) (mModel *Model, mSoftwareSystem *SoftwareSystem, mContainer *Container, mViewProps ViewProps, allViews []View) {
	t.Helper()

	mModel, mSoftwareSystem = makeSimModel(t)
	Root = &Design{
		Model: mModel,
	}
	mViewProps = makeViewProps(t)

	// Create container
	mComponents := makeComponents(mSoftwareSystem, t)
	mContainer = &Container{
		Element:    &Element{Name: "Container"},
		Components: mComponents,
		System:     mSoftwareSystem,
	}
	Identify(mContainer)

	// Create all view types
	mLandscapeView := &LandscapeView{
		ViewProps: &mViewProps,
	}
	mContextView := &ContextView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
	}
	mComponentView := &ComponentView{
		ViewProps:   &mViewProps,
		ContainerID: mContainer.ID,
	}
	mDeploymentView := &DeploymentView{
		ViewProps:   &mViewProps,
		Environment: "GreenField",
	}
	mContainerView := &ContainerView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
		AddInfluencers:   true,
	}

	allViews = []View{mLandscapeView, mContextView, mComponentView, mContainerView, mDeploymentView}
	return
}

func Test_AddAllElements(t *testing.T) {
	_, _, _, _, allViews := setupTestViewsWithContainer(t)

	tests := []struct {
		in View
	}{
		{in: allViews[0]}, // LandscapeView
		{in: allViews[1]}, // ContextView
		{in: allViews[2]}, // ComponentView
		{in: allViews[4]}, // DeploymentView (skip ContainerView as it's not supported)
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
			addAllElements(tt.in)
		})
	}
}

func Test_AddDefaultElements(t *testing.T) {
	_, _, _, _, allViews := setupTestViewsWithContainer(t)

	tests := []struct {
		in View
	}{
		{in: allViews[0]}, // LandscapeView
		{in: allViews[1]}, // ContextView
		{in: allViews[2]}, // ComponentView
		{in: allViews[3]}, // ContainerView
		{in: allViews[4]}, // DeploymentView
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
			addDefaultElements(tt.in)
		})
	}
}

func Test_AddNeighbors(t *testing.T) {
	_, _, mContainer, _, allViews := setupTestViewsWithContainer(t)

	tests := []struct {
		el *Element
		in View
	}{
		{el: mContainer.Element, in: allViews[0]}, // LandscapeView
		{el: mContainer.Element, in: allViews[1]}, // ContextView
		{el: mContainer.Element, in: allViews[2]}, // ComponentView
		{el: mContainer.Element, in: allViews[3]}, // ContainerView
		{el: mContainer.Element, in: allViews[4]}, // DeploymentView
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
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
		t.Run(fmt.Sprint(i), func(_ *testing.T) {
			addInfluencers(tt.el)
		})
	}
}

func Test_CoalesceRelationships(t *testing.T) {
	t.Run("auto_description_and_technology", func(t *testing.T) {
		_, _, mViewProps := setupCoalescingTest(t, "1", &CoalescedRelationship{
			// No Description or Technology - should auto-concatenate
		}, false)

		assert.Equal(t, 2, len(mViewProps.RelationshipViews), "Should start with 2 relationship views")

		coalesceRelationships(&mViewProps)

		assert.Equal(t, 1, len(mViewProps.RelationshipViews), "Should have 1 relationship view after coalescing")

		coalescedView := mViewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify auto-concatenated description and technology
		assert.Equal(t, "First relationship. Second relationship", coalescedRel.Description, "Should auto-concatenate descriptions")
		assert.Equal(t, "HTTP, HTTPS", coalescedRel.Technology, "Should auto-concatenate technologies")
	})

	t.Run("override_description_only", func(t *testing.T) {
		_, _, mViewProps := setupCoalescingTest(t, "2", &CoalescedRelationship{
			Description: "Custom description override",
			// No Technology - should auto-concatenate
		}, false)

		assert.Equal(t, 2, len(mViewProps.RelationshipViews), "Should start with 2 relationship views")

		coalesceRelationships(&mViewProps)

		assert.Equal(t, 1, len(mViewProps.RelationshipViews), "Should have 1 relationship view after coalescing")

		coalescedView := mViewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify explicit description and auto-concatenated technology
		assert.Equal(t, "Custom description override", coalescedRel.Description, "Should use explicit description")
		assert.Equal(t, "HTTP, HTTPS", coalescedRel.Technology, "Should auto-concatenate technologies")
	})

	t.Run("override_technology_only", func(t *testing.T) {
		_, _, mViewProps := setupCoalescingTest(t, "3", &CoalescedRelationship{
			// No Description - should auto-concatenate
			Technology: "Custom tech override",
		}, true) // Changed to true to ensure explicit coalescing works

		assert.Equal(t, 2, len(mViewProps.RelationshipViews), "Should start with 2 relationship views")

		coalesceRelationships(&mViewProps)

		assert.Equal(t, 1, len(mViewProps.RelationshipViews), "Should have 1 relationship view after coalescing")

		coalescedView := mViewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify auto-concatenated description and explicit technology
		assert.Equal(t, "First relationship. Second relationship", coalescedRel.Description, "Should auto-concatenate descriptions")
		assert.Equal(t, "Custom tech override", coalescedRel.Technology, "Should use explicit technology")
	})

	t.Run("override_both_description_and_technology", func(t *testing.T) {
		_, _, mViewProps := setupCoalescingTest(t, "4", &CoalescedRelationship{
			Description: "Custom description override",
			Technology:  "Custom tech override",
		}, true) // Changed to true to ensure explicit coalescing works

		assert.Equal(t, 2, len(mViewProps.RelationshipViews), "Should start with 2 relationship views")

		coalesceRelationships(&mViewProps)

		assert.Equal(t, 1, len(mViewProps.RelationshipViews), "Should have 1 relationship view after coalescing")

		coalescedView := mViewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify both explicit description and technology
		assert.Equal(t, "Custom description override", coalescedRel.Description, "Should use explicit description")
		assert.Equal(t, "Custom tech override", coalescedRel.Technology, "Should use explicit technology")
	})

	t.Run("verify_common_properties", func(t *testing.T) {
		mSource, mDestination, mViewProps := setupCoalescingTest(t, "5", &CoalescedRelationship{
			Description: "Test description",
			Technology:  "Test tech",
		}, false)

		coalesceRelationships(&mViewProps)

		coalescedView := mViewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify source and destination are correct
		assert.Equal(t, mSource.ID, coalescedRel.Source.ID, "Source should match")
		assert.Equal(t, mDestination.ID, coalescedRel.Destination.ID, "Destination should match")
	})
}

func Test_CoalesceAllRelationshipsInView(t *testing.T) {
	// Use setupCoalescingTest with no explicit coalesced relationship and coalesceAll=false
	mSource, mDestination, baseViewProps := setupCoalescingTest(t, "CoalesceAll", nil, false)

	// Test 1: No explicit pairs - should coalesce automatically
	t.Run("auto_coalesce_with_no_explicit_pairs", func(t *testing.T) {
		// Make a copy to avoid modifying original
		viewProps := baseViewProps
		viewProps.RelationshipViews = make([]*RelationshipView, len(baseViewProps.RelationshipViews))
		copy(viewProps.RelationshipViews, baseViewProps.RelationshipViews)

		explicitPairs := make(map[string]bool)

		initialCount := len(viewProps.RelationshipViews)
		assert.Equal(t, 2, initialCount, "Should start with 2 relationship views")

		coalesceAllRelationshipsInView(&viewProps, explicitPairs)

		// Should have one coalesced relationship view
		assert.Equal(t, 1, len(viewProps.RelationshipViews), "Should have 1 relationship view after auto-coalescing")

		// Get the coalesced relationship
		coalescedView := viewProps.RelationshipViews[0]
		coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

		// Verify auto-concatenated description and technology
		assert.Equal(t, "First relationship. Second relationship", coalescedRel.Description, "Should auto-concatenate descriptions")
		assert.Equal(t, "HTTP, HTTPS", coalescedRel.Technology, "Should auto-concatenate technologies")
	})

	// Test 2: With explicit pairs - should skip the explicit pair
	t.Run("skip_explicit_pairs", func(t *testing.T) {
		// Make a copy to avoid modifying original
		viewProps := baseViewProps
		viewProps.RelationshipViews = make([]*RelationshipView, len(baseViewProps.RelationshipViews))
		copy(viewProps.RelationshipViews, baseViewProps.RelationshipViews)

		explicitPairs := make(map[string]bool)
		explicitPairs[mSource.ID+"->"+mDestination.ID] = true

		initialCount := len(viewProps.RelationshipViews)
		assert.Equal(t, 2, initialCount, "Should start with 2 relationship views")

		coalesceAllRelationshipsInView(&viewProps, explicitPairs)

		// Should still have 2 relationship views since the pair was explicitly handled
		assert.Equal(t, 2, len(viewProps.RelationshipViews), "Should still have 2 relationship views when pair is explicit")
	})
}

func Test_CoalesceRelationships_ExplicitTakesPrecedence(t *testing.T) {
	_, _, mViewProps := setupCoalescingTest(t, "ExplicitPrecedence", &CoalescedRelationship{
		Description: "Explicit custom description", // This should take precedence
		Technology:  "Explicit custom tech",
	}, true) // Both CoalesceAllRelationships AND explicit CoalescedRelationships

	initialCount := len(mViewProps.RelationshipViews)
	assert.Equal(t, 2, initialCount, "Should start with 2 relationship views")

	// Call main coalescing function (which handles precedence)
	coalesceRelationships(&mViewProps)

	// Should have one coalesced relationship view
	assert.Equal(t, 1, len(mViewProps.RelationshipViews), "Should have 1 relationship view after coalescing")

	// Get the coalesced relationship
	coalescedView := mViewProps.RelationshipViews[0]
	coalescedRel := Registry[coalescedView.RelationshipID].(*Relationship)

	// Verify explicit description and technology were used (not auto-concatenated)
	assert.Equal(t, "Explicit custom description", coalescedRel.Description, "Should use explicit description, not auto-concatenated")
	assert.Equal(t, "Explicit custom tech", coalescedRel.Technology, "Should use explicit technology, not auto-concatenated")

	// Verify it's NOT the auto-concatenated version
	assert.NotEqual(t, "First relationship. Second relationship", coalescedRel.Description, "Should NOT be auto-concatenated description")
	assert.NotEqual(t, "HTTP, HTTPS", coalescedRel.Technology, "Should NOT be auto-concatenated technology")
}

// setupCoalescingTest creates independent test data for coalescing tests
func setupCoalescingTest(t *testing.T, testSuffix string, coalescedRel *CoalescedRelationship, coalesceAll bool) (*SoftwareSystem, *SoftwareSystem, ViewProps) {
	t.Helper()

	// Create source and destination systems
	mSource, mDestination := createTestSystems(t, testSuffix)

	// Create fresh relationships
	rel1 := &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "First relationship",
		Technology:  "HTTP",
	}
	Identify(rel1)

	rel2 := &Relationship{
		Source:      mSource.Element,
		Destination: mDestination.Element,
		Description: "Second relationship",
		Technology:  "HTTPS",
	}
	Identify(rel2)

	// Set the source and destination for the coalesced relationship
	if coalescedRel != nil {
		coalescedRel.Source = mSource.Element
		coalescedRel.Destination = mDestination.Element
	}

	// Build the ViewProps
	viewProps := ViewProps{
		RelationshipViews: []*RelationshipView{
			{
				Source:         mSource.Element,
				Destination:    mDestination.Element,
				RelationshipID: rel1.ID,
			},
			{
				Source:         mSource.Element,
				Destination:    mDestination.Element,
				RelationshipID: rel2.ID,
			},
		},
		CoalesceAllRelationships: coalesceAll,
	}

	if coalescedRel != nil {
		viewProps.CoalescedRelationships = []*CoalescedRelationship{coalescedRel}
	}

	return mSource, mDestination, viewProps
}

// createTestSystems creates a pair of test software systems with unique names
func createTestSystems(t *testing.T, suffix string) (*SoftwareSystem, *SoftwareSystem) {
	t.Helper()

	source := &SoftwareSystem{
		Element: &Element{
			Name: "TestSource" + suffix,
		},
	}
	Identify(source)

	destination := &SoftwareSystem{
		Element: &Element{
			Name: "TestDestination" + suffix,
		},
	}
	Identify(destination)

	return source, destination
}
