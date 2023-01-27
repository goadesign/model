package expr

import (
	"fmt"
	"testing"
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
	var res string = ""
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
				mLandscapeView.AddElements(tt.ehs1)
			} else {
				mLandscapeView.AddElements(tt.ehs1, tt.ehs2)
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
				mContextView.AddElements(tt.ehs1)
			} else {
				mContextView.AddElements(tt.ehs1, tt.ehs2)
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
				mContainerView.AddElements(tt.ehs1)
			} else {
				mContainerView.AddElements(tt.ehs1, tt.ehs2)
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
				mComponentView.AddElements(tt.ehs1)
			} else {
				mComponentView.AddElements(tt.ehs1, tt.ehs2)
			}

			got := len(mComponentView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}

/*
func Test_AddElements2DeploymentView(t *testing.T) {
	mLViewProp := ViewProps{Key: "LV"}
	mDeploymentView := DeploymentView{ViewProps: &mLViewProp}
	mDeploymentNode := DeploymentNode{
		Element: &Element{Name: "SoftwareSystem"},
	}
	mSoftwareSystem := SoftwareSystem{
		Element: &Element{ID: "2", Name: "SoftwareSystem"},
	}

	mContainer := Container{
		Element: &Element{ID: "7", Name: "Container"},
		System:  &mSoftwareSystem,
	}
	Identify(&mContainer)
	mContainerInstance := ContainerInstance{
		Element:     &Element{Name: "ContainerInstance"},
		ContainerID: "7:Container",
		Container:   &mContainer,
	}
	mDeploymentView.SoftwareSystemID = "1"
	/*mComponent := Component{
		Element: &Element{Name: "Component"},
	}*/
/*	tests := []struct {
		ehs1 ElementHolder
		ehs2 ElementHolder
		want int
	}{
		{ehs1: &mDeploymentNode, want: 1},
		{ehs1: &mContainerInstance, want: 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if i == 0 {
				mDeploymentView.AddElements(tt.ehs1)
			} else {
				mDeploymentView.AddElements(tt.ehs1)
			}

			got := len(mDeploymentView.ElementViews)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}
*/
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
				mLandscapeView.AddAnimationStep(&tt.ehs1)
			} else {
				mLandscapeView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mLandscapeView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
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
				mContextView.AddAnimationStep(&tt.ehs1)
			} else {
				mContextView.AddAnimationStep(&tt.ehs1)
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
				mContainerView.AddAnimationStep(&tt.ehs1)
			} else {
				mContainerView.AddAnimationStep(&tt.ehs1)
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
				mComponentView.AddAnimationStep(&tt.ehs1)
			} else {
				mComponentView.AddAnimationStep(&tt.ehs1)
			}

			got := len(mComponentView.AnimationSteps)
			if got != tt.want {
				t.Errorf("Got %d, wanted %d", got, tt.want)
			}
		})
	}

}
