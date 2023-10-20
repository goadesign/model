package design

import (
	"fmt"

	. "goa.design/model/dsl"
	"goa.design/model/mdl"
)

var shapes = []ShapeKind{
	ShapeBox,
	ShapeCircle,
	ShapeCylinder,
	ShapeEllipse,
	ShapeHexagon,
	ShapeRoundedBox,
	ShapeComponent,
	ShapeFolder,
	ShapeMobileDeviceLandscape,
	ShapeMobileDevicePortrait,
	//ShapePerson,
	ShapePipe,
	ShapeRobot,
	ShapeWebBrowser,
}

var _ = Design("Getting Started", "This is a model of my software system.", func() {
	for i, sh := range shapes {
		func(i int) {
			SoftwareSystem(fmt.Sprintf("System %d", i+1), fmt.Sprintf("Shape: %s.", shapeName(sh)), func() {
				Tag(fmt.Sprintf("system%d", i+1))
			})
		}(i)
	}

	Person("User", "A person using shapes.", func() {
		for i := range shapes {
			Uses(fmt.Sprintf("System %d", i+1), "Uses")
		}
		Tag("person")
	})

	Views(func() {
		SystemContextView("System 1", "SystemContext", "An example of a System Context diagram.", func() {
			AddAll()
			AutoLayout(ImplementationDagre, RankLeftRight)
		})
		Styles(func() {
			ElementStyle("person", func() {
				Shape(ShapePerson)
				Background("#f0f7ff")
			})
			for i, shape := range shapes {
				ElementStyle(fmt.Sprintf("system%d", i+1), func() {
					Shape(shape)
					Background("#f0f7ff")
				})
			}
		})
	})
})

func shapeName(sh ShapeKind) string {
	b, _ := mdl.ShapeKind(sh).MarshalJSON()
	return string(b)
}
