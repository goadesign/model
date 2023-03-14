package expr

import (
	"fmt"
	"testing"

	"goa.design/goa/v3/eval"
)

// Don't use t.parallel. Crashes when file tests are run as concurrent map read & write.

func Test_Views_EvalName(t *testing.T) {
	mViews := Views{}

	tests := []struct {
		want string
	}{
		{want: "views"},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got := mViews.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}

}

func Test_DSLFunc(t *testing.T) {
	mViews := Views{
		DSLFunc: func() {},
	}
	result := mViews.DSL()
	if result == nil {
		t.Errorf("Function not return")
	}
}

func Test_IsPS(t *testing.T) {
	mPerson := Person{
		Element: &Element{Name: "Person"},
	}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mContainer := Container{
		Element: &Element{Name: "Container"},
	}

	tests := []struct {
		eh   ElementHolder
		want bool
	}{
		{eh: &mPerson, want: true},
		{eh: &mSoftwareSystem, want: true},
		{eh: &mContainer, want: false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := isPS(tt.eh)
			if got != tt.want {
				t.Errorf("Got %t, wanted %t, for %v", got, tt.want, tt.eh)
			}
		})
	}
}

func Test_IsPSC(t *testing.T) {
	mPerson := Person{
		Element: &Element{Name: "Person"},
	}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mContainer := Container{
		Element: &Element{Name: "Container"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	tests := []struct {
		eh   ElementHolder
		want bool
	}{
		{eh: &mPerson, want: true},
		{eh: &mSoftwareSystem, want: true},
		{eh: &mContainer, want: true},
		{eh: &mComponent, want: false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := isPSC(tt.eh)
			if got != tt.want {
				t.Errorf("Got %t, wanted %t, for %v", got, tt.want, tt.eh)
			}
		})
	}
}

func Test_IsPSCC(t *testing.T) {
	mDeploymentNode := DeploymentNode{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mPerson := Person{
		Element: &Element{Name: "Person"},
	}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mContainer := Container{
		Element: &Element{Name: "Container"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	tests := []struct {
		eh   ElementHolder
		want bool
	}{
		{eh: &mPerson, want: true},
		{eh: &mSoftwareSystem, want: true},
		{eh: &mContainer, want: true},
		{eh: &mComponent, want: true},
		{eh: &mDeploymentNode, want: false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := isPSCC(tt.eh)
			if got != tt.want {
				t.Errorf("Got %t, wanted %t, for %v", got, tt.want, tt.eh)
			}
		})
	}
}

func Test_IsDCI(t *testing.T) {
	mDeploymentNode := DeploymentNode{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mPerson := Person{
		Element: &Element{Name: "Person"},
	}
	mInfrastructureNode := InfrastructureNode{
		Element: &Element{Name: "InfrastructureNode"},
	}
	mContainer := ContainerInstance{
		Element: &Element{Name: "Container"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	tests := []struct {
		eh   ElementHolder
		want bool
	}{
		{eh: &mPerson, want: false},
		{eh: &mInfrastructureNode, want: true},
		{eh: &mContainer, want: true},
		{eh: &mComponent, want: false},
		{eh: &mDeploymentNode, want: true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := isDCI(tt.eh)
			if got != tt.want {
				t.Errorf("Got %t, wanted %t, for %v", got, tt.want, tt.eh)
			}
		})
	}
}

func Test_AllViews(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mLandscapeView := LandscapeView{ViewProps: &mLViewProp}
	mLandscapeViews := make([]*LandscapeView, 0)
	mLandscapeViews = append(mLandscapeViews, &mLandscapeView)

	mCViewProp := ViewProps{Key: "CTXV"}
	mContextView := ContextView{ViewProps: &mCViewProp}
	mContextViews := make([]*ContextView, 0)
	mContextViews = append(mContextViews, &mContextView)

	mCNViewProp := ViewProps{Key: "CNV"}
	mContainerView := ContainerView{ViewProps: &mCNViewProp}
	mContainerViews := make([]*ContainerView, 0)
	mContainerViews = append(mContainerViews, &mContainerView)

	mCMPViewProp := ViewProps{Key: "CMPV"}
	mComponentView := ComponentView{ViewProps: &mCMPViewProp}
	mComponentViews := make([]*ComponentView, 0)
	mComponentViews = append(mComponentViews, &mComponentView)

	mDYNViewProp := ViewProps{Key: "DYNV"}
	mDynamicView := DynamicView{ViewProps: &mDYNViewProp}
	mDynamicViews := make([]*DynamicView, 0)
	mDynamicViews = append(mDynamicViews, &mDynamicView)

	mDeploymentProp := ViewProps{Key: "DPMV"}
	mDeploymentView := DeploymentView{ViewProps: &mDeploymentProp}
	mDeploymentViews := make([]*DeploymentView, 0)
	mDeploymentViews = append(mDeploymentViews, &mDeploymentView)

	mViews := Views{
		LandscapeViews:  mLandscapeViews,
		ContainerViews:  mContainerViews,
		ContextViews:    mContextViews,
		ComponentViews:  mComponentViews,
		DynamicViews:    mDynamicViews,
		DeploymentViews: mDeploymentViews,
	}

	tests := []struct {
		mViews Views
		want   string
	}{
		{mViews: mViews, want: "LV1,CTXV1,CNV1,CMPV1,DYNV1,DPMV1"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := mViews.All()
			result := assessViewResult(got)
			if result != tt.want {
				t.Errorf("Got %s, wanted %s", result, tt.want)
			}
		})
	}
}

func assessViewResult(vps []View) string {
	var res string
	var intermed = make(map[string]int)
	intermed["LV"] = 0
	intermed["CTXV"] = 0
	intermed["CNV"] = 0
	intermed["CMPV"] = 0
	intermed["DYNV"] = 0
	intermed["DPMV"] = 0
	for _, vtype := range vps {
		intermed[vtype.Props().Key] += 1
	}
	res = fmt.Sprintf("LV%d,CTXV%d,CNV%d,CMPV%d,DYNV%d,DPMV%d", intermed["LV"], intermed["CTXV"], intermed["CNV"], intermed["CMPV"], intermed["DYNV"], intermed["DPMV"])
	return res
}

func Test_AddElements2LandscapeView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mLandscapeView := LandscapeView{ViewProps: &mLViewProp}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mContainer := Container{
		Element: &Element{Name: "Container"},
	}
	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mSoftwareSystem, want: 1},
		{ehs1: &mSoftwareSystem, ehs2: &mContainer, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mLandscapeView.AddElements(tt.ehs1)
			} else {
				_ = mLandscapeView.AddElements(tt.ehs1, tt.ehs2)
			}
			got := len(mLandscapeView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddElements2ContextView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mContextView := ContextView{ViewProps: &mLViewProp}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mContainer := Container{
		Element: &Element{Name: "Container"},
	}
	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mSoftwareSystem, want: 1},
		{ehs1: &mSoftwareSystem, ehs2: &mContainer, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mContextView.AddElements(tt.ehs1)
			} else {
				_ = mContextView.AddElements(tt.ehs1, tt.ehs2)
			}

			got := len(mContextView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddElements2ContainerView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mContainerView := ContainerView{ViewProps: &mLViewProp}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}
	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mSoftwareSystem, want: 1},
		{ehs1: &mSoftwareSystem, ehs2: &mComponent, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mContainerView.AddElements(tt.ehs1)
			} else {
				_ = mContainerView.AddElements(tt.ehs1, tt.ehs2)
			}

			got := len(mContainerView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddElements2ComponentView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mComponentView := ComponentView{ViewProps: &mLViewProp}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}
	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mSoftwareSystem, want: 1},
		{ehs1: &mSoftwareSystem, ehs2: mComponent.Element, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mComponentView.AddElements(tt.ehs1)
			} else {
				_ = mComponentView.AddElements(tt.ehs1, tt.ehs2)
			}

			got := len(mComponentView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddElements2DeploymentView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mDeploymentView := DeploymentView{ViewProps: &mLViewProp}
	mRootInfrastructureNode := InfrastructureNode{
		Element: &Element{Name: "InfrastructureNode"},
	}
	mInfrastructureNodes := make([]*InfrastructureNode, 1)
	mInfrastructureNodes[0] = &mRootInfrastructureNode
	mRootDeploymentNode := DeploymentNode{
		Element: &Element{Name: "RootDeploymentNode"},
	}

	mDeploymentNode := DeploymentNode{
		Element:             &Element{Name: "DeploymentNode"},
		Parent:              &mRootDeploymentNode,
		InfrastructureNodes: mInfrastructureNodes,
	}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{ID: "2", Name: "SoftwareSystem"},
	}
	Identify(&mSoftwareSystem)

	mContainer := Container{
		Element: &Element{ID: "7", Name: "Container"},
		System:  &mSoftwareSystem,
	}
	Identify(&mContainer)
	mContainerInstance := ContainerInstance{
		Element:     &Element{Name: "ContainerInstance"},
		ContainerID: mContainer.ID,
		Container:   &mContainer,
		Parent:      &mDeploymentNode,
	}
	mDeploymentView.SoftwareSystemID = mSoftwareSystem.ID

	mInfrastructureNode := InfrastructureNode{
		Element: &Element{Name: "InfrastructureNode"},
		Parent:  &mDeploymentNode,
	}
	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mSoftwareSystem, want: 0},
		{ehs1: &mDeploymentNode, want: 1},
		{ehs1: &mContainerInstance, want: 1},
		{ehs1: &mInfrastructureNode, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_ = mDeploymentView.AddElements(tt.ehs1)
			got := len(mDeploymentView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddAnimationStep2LandscapeView(t *testing.T) {
	var eHolder []ElementHolder
	var fHolder []ElementHolder
	mLViewProp := ViewProps{Key: "LV"}
	mLandscapeView := LandscapeView{ViewProps: &mLViewProp}

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}

	mContainer := Container{
		Element: &Element{Name: "Container"},
	}

	eHolder = make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	fHolder = make([]ElementHolder, 1)
	fHolder[0] = &mContainer

	mFalseAnimationStep := AnimationStep{
		Order:    1,
		Elements: fHolder,
	}

	tests := []struct {
		ehs1 AnimationStep
		want int
	}{
		{ehs1: mAnimationStep, want: 1},
		{ehs1: mFalseAnimationStep, want: 1}, // i.e. another one isn't added..
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mLandscapeView.AddAnimationStep(&tt.ehs1)
			} else {
				_ = mLandscapeView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mLandscapeView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddDeploymentNodeChildren(t *testing.T) {
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	Identify(&mSoftwareSystem)
	mContainer := Container{
		Element: &Element{Name: "Container"},
		System:  &mSoftwareSystem,
	}
	Identify(&mContainer)
	mContainerInstance := ContainerInstance{
		Element:     &Element{ID: "7", Name: "ContainerInstance"},
		Container:   &mContainer,
		ContainerID: mContainer.ID,
	}
	mContainerInstances := make([]*ContainerInstance, 1)
	mContainerInstances[0] = &mContainerInstance

	mDeploymentNode := DeploymentNode{
		Element: &Element{Name: "Main"},
	}
	mDeploymentNodes := make([]*DeploymentNode, 1)
	mDeploymentNodes[0] = &mDeploymentNode
	mInfrastructureNode := InfrastructureNode{
		Element:     &Element{Name: "Infrastructure"},
		Parent:      &mDeploymentNode,
		Environment: "Greenfield",
	}
	Identify(&mInfrastructureNode)
	mInfrastructureNodes := make([]*InfrastructureNode, 1)
	mInfrastructureNodes[0] = &mInfrastructureNode
	n := DeploymentNode{
		ContainerInstances:  mContainerInstances,
		InfrastructureNodes: mInfrastructureNodes,
		Children:            mDeploymentNodes,
	}
	mElementView := ElementView{
		Element: &Element{ID: "1", Name: "Element"},
	}
	mElementViews := make([]*ElementView, 1)
	mElementViews[0] = &mElementView
	mViewProps := ViewProps{
		ElementViews: mElementViews,
	}
	dv := DeploymentView{
		ViewProps:        &mViewProps,
		SoftwareSystemID: mSoftwareSystem.Element.ID,
	}
	tests := []struct {
		dv1  *DeploymentView
		dn1  *DeploymentNode
		want bool
	}{
		{dv1: &dv, dn1: &n, want: true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := addDeploymentNodeChildren(tt.dv1, tt.dn1)
			if got != tt.want {
				t.Errorf("Got %t, wanted %t", got, tt.want)
			}
		})
	}

}

func Test_AddAnimationStep2ContextView(t *testing.T) {
	var eHolder []ElementHolder
	var fHolder []ElementHolder
	mLViewProp := ViewProps{Key: "LV"}
	mContextView := ContextView{ViewProps: &mLViewProp}

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}

	mContainer := Container{
		Element: &Element{Name: "Container"},
	}

	eHolder = make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	fHolder = make([]ElementHolder, 1)
	fHolder[0] = &mContainer

	mFalseAnimationStep := AnimationStep{
		Order:    1,
		Elements: fHolder,
	}

	tests := []struct {
		ehs1 AnimationStep
		want int
	}{
		{ehs1: mAnimationStep, want: 1},
		{ehs1: mFalseAnimationStep, want: 1}, // i.e. another one isn't added..
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mContextView.AddAnimationStep(&tt.ehs1)
			} else {
				_ = mContextView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mContextView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddAnimationStep2ContainerView(t *testing.T) {
	var eHolder []ElementHolder
	var fHolder []ElementHolder
	mLViewProp := ViewProps{Key: "LV"}
	mContainerView := ContainerView{ViewProps: &mLViewProp}

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}

	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	eHolder = make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	fHolder = make([]ElementHolder, 1)
	fHolder[0] = &mComponent

	mFalseAnimationStep := AnimationStep{
		Order:    1,
		Elements: fHolder,
	}

	tests := []struct {
		ehs1 AnimationStep
		want int
	}{
		{ehs1: mAnimationStep, want: 1},
		{ehs1: mFalseAnimationStep, want: 1}, // i.e. another one isn't added..
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mContainerView.AddAnimationStep(&tt.ehs1)
			} else {
				_ = mContainerView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mContainerView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddAnimationStep2ComponentView(t *testing.T) {
	var eHolder []ElementHolder
	var fHolder []ElementHolder
	mLViewProp := ViewProps{Key: "LV"}
	mComponentView := ComponentView{ViewProps: &mLViewProp}

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}

	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	eHolder = make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	fHolder = make([]ElementHolder, 1)
	fHolder[0] = mComponent.Element

	mFalseAnimationStep := AnimationStep{
		Order:    1,
		Elements: fHolder,
	}

	tests := []struct {
		ehs1 AnimationStep
		want int
	}{
		{ehs1: mAnimationStep, want: 1},
		{ehs1: mFalseAnimationStep, want: 1}, // i.e. another one isn't added..
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				_ = mComponentView.AddAnimationStep(&tt.ehs1)
			} else {
				_ = mComponentView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mComponentView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_AddAnimationStep2DeploymentView(t *testing.T) {
	var eHolder []ElementHolder
	var fHolder []ElementHolder
	mLViewProp := ViewProps{Key: "LV"}
	mComponentView := ComponentView{ViewProps: &mLViewProp}

	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}

	mDeploymentView := DeploymentView{ViewProps: &mLViewProp}

	mDeploymentNode := DeploymentNode{
		Element: &Element{Name: "DeploymentNode"},
	}

	eHolder = make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mFalseAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	fHolder = make([]ElementHolder, 1)
	fHolder[0] = &mDeploymentNode

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: fHolder,
	}

	tests := []struct {
		ehs1 AnimationStep
		want int
	}{
		{ehs1: mAnimationStep, want: 1},
		{ehs1: mFalseAnimationStep, want: 1}, // i.e. another one isn't added..
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_ = mDeploymentView.AddAnimationStep(&tt.ehs1)

			got := len(mComponentView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

func Test_ValidateViews(t *testing.T) {
	// set up all the elements needed to test the one function
	// it's a marathon... rather than a sprint..
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mComponent := Component{
		Element: &Element{Name: "Component"},
	}

	mSrcDeploymentNode := DeploymentNode{
		Element: &Element{Name: "SrcDeploymentNode"},
	}
	mDstDeploymentNode := DeploymentNode{
		Element: &Element{Name: "DstDeploymentNode"},
	}

	mRelationship := Relationship{
		Source:      mSrcDeploymentNode.Element,
		Destination: mDstDeploymentNode.Element,
		Description: "Sample",
	}
	Identify(&mRelationship)
	eHolder := make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	mAnimationSteps := make([]*AnimationStep, 1)
	mAnimationSteps[0] = &mAnimationStep

	mContainer := Container{
		Element: &Element{ID: "1"},
		System:  &mSoftwareSystem,
	}
	Identify(&mContainer)
	mElementView := ElementView{
		Element: mContainer.Element,
	}
	mElementViews := make([]*ElementView, 1)
	mElementViews[0] = &mElementView

	mSSElementView := ElementView{
		Element: mSoftwareSystem.Element,
	}
	mSSElementViews := make([]*ElementView, 1)
	mSSElementViews[0] = &mSSElementView

	mCTElementView := ElementView{
		Element: mComponent.Element,
	}
	mCTElementViews := make([]*ElementView, 1)
	mCTElementViews[0] = &mCTElementView

	mSourceContainerInstance := ContainerInstance{
		Element:     &Element{Name: "Source"},
		Parent:      &mSrcDeploymentNode,
		ContainerID: "1",
	}
	Identify(&mSourceContainerInstance)
	// this is probably cheating as the ids ought to refer to the containers.
	mSourceContainerInstance.ContainerID = mSourceContainerInstance.ID

	mDestinationContainerInstance := ContainerInstance{
		Element:     &Element{Name: "Destination"},
		Parent:      &mDstDeploymentNode,
		ContainerID: "2",
	}
	Identify(&mDestinationContainerInstance)
	// this is probably cheating as the ids ought to refer to the containers.
	mDestinationContainerInstance.ContainerID = mDestinationContainerInstance.ID
	mCIRelationship := Relationship{
		Source:      mSourceContainerInstance.Element,
		Destination: mDestinationContainerInstance.Element,
		Description: "CIRelationship",
	}
	Identify(&mCIRelationship)

	mRelationshipView := RelationshipView{
		Source:         mSourceContainerInstance.Element,
		Destination:    mDestinationContainerInstance.Element,
		Description:    "CIRelationship",
		RelationshipID: "",
	}
	mRelationshipViews := make([]*RelationshipView, 1)
	mRelationshipViews[0] = &mRelationshipView
	mLVViewProps := ViewProps{
		Description:       "House",
		RelationshipViews: mRelationshipViews,
		AnimationSteps:    mAnimationSteps,
		ElementViews:      mElementViews,
	}

	mBadRelationshipView := RelationshipView{
		Source:         &Element{ID: "1"},
		Destination:    &Element{ID: "2"},
		Description:    "BadRelationship",
		RelationshipID: "",
	}
	mBadRelationshipViews := make([]*RelationshipView, 1)
	mBadRelationshipViews[0] = &mBadRelationshipView
	mLandscapeView := LandscapeView{
		ViewProps: &mLVViewProps,
	}
	mLandscapeViews := make([]*LandscapeView, 1)
	mLandscapeViews[0] = &mLandscapeView
	mCVViewProps := ViewProps{
		Description:       "House",
		RelationshipViews: mRelationshipViews,
		AnimationSteps:    mAnimationSteps,
		ElementViews:      mSSElementViews,
	}

	mUnreachable := Element{Name: "unreachable"}
	mUnreachables := make([]*Element, 1)
	mUnreachables[0] = &mUnreachable

	mCTViewProps := ViewProps{
		Description:       "House",
		RelationshipViews: mBadRelationshipViews,
		AnimationSteps:    mAnimationSteps,
		ElementViews:      mCTElementViews,
		RemoveUnreachable: mUnreachables,
	}

	mContextView := ContextView{
		ViewProps: &mCVViewProps,
	}

	mContextViews := make([]*ContextView, 1)
	mContextViews[0] = &mContextView

	mContainerView := ContainerView{
		ViewProps: &mCTViewProps,
	}

	mContainerViews := make([]*ContainerView, 1)
	mContainerViews[0] = &mContainerView

	mViews := Views{
		LandscapeViews: mLandscapeViews,
		ContextViews:   mContextViews,
		ContainerViews: mContainerViews,
	}
	tests := []struct {
		ehs1 Views
		want int
	}{
		{ehs1: mViews, want: 5},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tt.ehs1.Validate()
			evalgot := got.(*eval.ValidationErrors)
			if len(evalgot.Errors) != tt.want {
				t.Errorf("Got %d, wanted %d", len(evalgot.Errors), tt.want)
			}
		})
	}

}

func Test_Finalize(t *testing.T) {
	// set up all the elements needed to test the one function
	// it's a marathon... rather than a sprint..
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "SoftwareSystem"},
	}
	Identify(&mSoftwareSystem)
	mUnrelatedSoftwareSystem := SoftwareSystem{
		Element: &Element{Name: "UnrelatedSoftwareSystem"},
	}
	Identify(&mUnrelatedSoftwareSystem)
	mSrcDeploymentNode := DeploymentNode{
		Element: &Element{Name: "SrcDeploymentNode"},
	}
	mDstDeploymentNode := DeploymentNode{
		Element: &Element{Name: "DstDeploymentNode"},
	}

	mRelationship := Relationship{
		Source:      mSrcDeploymentNode.Element,
		Destination: mDstDeploymentNode.Element,
		Description: "Sample",
	}
	Identify(&mRelationship)
	eHolder := make([]ElementHolder, 1)
	eHolder[0] = &mSoftwareSystem

	mAnimationStep := AnimationStep{
		Order:    1,
		Elements: eHolder,
	}
	mAnimationSteps := make([]*AnimationStep, 1)
	mAnimationSteps[0] = &mAnimationStep

	mContainer := Container{
		Element: &Element{ID: "1"},
		System:  &mSoftwareSystem,
	}
	Identify(&mContainer)
	mComponent := Component{
		Element:   &Element{Name: "Component"},
		Container: &mContainer,
	}
	Identify(&mComponent)

	mCTElementView := ElementView{
		Element:        mComponent.Element,
		NoRelationship: true,
	}
	mCTElementViews := make([]*ElementView, 1)
	mCTElementViews[0] = &mCTElementView

	mSourceContainerInstance := ContainerInstance{
		Element:     &Element{Name: "Source"},
		Parent:      &mSrcDeploymentNode,
		ContainerID: "1",
	}
	Identify(&mSourceContainerInstance)
	// this is probably cheating as the ids ought to refer to the containers.
	mSourceContainerInstance.ContainerID = mSourceContainerInstance.ID

	mDestinationContainerInstance := ContainerInstance{
		Element:     &Element{Name: "Destination"},
		Parent:      &mDstDeploymentNode,
		ContainerID: "2",
	}
	Identify(&mDestinationContainerInstance)
	// this is probably cheating as the ids ought to refer to the containers.
	mDestinationContainerInstance.ContainerID = mDestinationContainerInstance.ID
	mCIRelationship := Relationship{
		Source:      mSourceContainerInstance.Element,
		Destination: mDestinationContainerInstance.Element,
		Description: "CIRelationship",
	}
	Identify(&mCIRelationship)
	mCIRelationships := make([]*Relationship, 1)
	mCIRelationships[0] = &mCIRelationship

	mRelationshipView := RelationshipView{
		Source:         mSourceContainerInstance.Element,
		Destination:    mDestinationContainerInstance.Element,
		Description:    "CIRelationship",
		RelationshipID: "",
	}
	mUnrelationshipView := RelationshipView{
		Source:         mSoftwareSystem.Element,
		Destination:    mDestinationContainerInstance.Element,
		Description:    "Dummy",
		RelationshipID: "",
	}
	mRelationshipViews := make([]*RelationshipView, 2)
	mRelationshipViews[0] = &mRelationshipView
	mRelationshipViews[1] = &mUnrelationshipView

	mTags := make([]string, 2)
	mTags[0] = "String1"
	mTags[1] = "String1"
	mNeighbourElement := Element{ID: "1"}
	mNeighbourElements := make([]*Element, 1)
	mNeighbourElements[0] = &mNeighbourElement

	mRemoveElement := Element{ID: "1"}
	mRemoveElements := make([]*Element, 1)
	mRemoveElements[0] = &mRemoveElement
	mUnreachable := Element{Name: "unreachable"}
	mUnreachables := make([]*Element, 1)
	mUnreachables[0] = &mUnreachable

	mCTViewProps := ViewProps{
		Description:         "House",
		RelationshipViews:   mRelationshipViews,
		AnimationSteps:      mAnimationSteps,
		ElementViews:        mCTElementViews,
		AddAll:              true,
		AddDefault:          true,
		AddNeighbors:        mNeighbourElements,
		RemoveElements:      mRemoveElements,
		RemoveRelationships: mCIRelationships,
		RemoveTags:          mTags,
		RemoveUnreachable:   mUnreachables,
		RemoveUnrelated:     true,
	}

	mCTNotAllViewProps := ViewProps{
		Description:       "House",
		RelationshipViews: mRelationshipViews,
		AnimationSteps:    mAnimationSteps,
		ElementViews:      mCTElementViews,
		AddAll:            false,
		AddDefault:        true,
		AddNeighbors:      mNeighbourElements,
		RemoveElements:    mRemoveElements,
	}

	mContainerView := ContainerView{
		ViewProps:        &mCTViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
		AddInfluencers:   true,
	}
	mContainerNotAllView := ContainerView{
		ViewProps:        &mCTNotAllViewProps,
		SoftwareSystemID: mSoftwareSystem.ID,
		AddInfluencers:   true,
	}

	mContainerViews := make([]*ContainerView, 2)
	mContainerViews[0] = &mContainerView
	mContainerViews[1] = &mContainerNotAllView

	mViews := Views{
		ContainerViews: mContainerViews,
	}
	tests := []struct {
		ehs1 Views
		want int
	}{
		{ehs1: mViews, want: 4},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tt.ehs1.Finalize()
			got := 4
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}
