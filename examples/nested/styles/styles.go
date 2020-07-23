/*
Package styles provide shared styles used by multiple models.
*/
package styles

import . "goa.design/structurizr/dsl"

// DefineAll defines all the styles described in this package.
func DefineAll() {
	SystemStyle()
	PersonStyle()
}

// SystemStyle defines the style used to render software systems. All elements tagged
// with "system" inherit the style.
func SystemStyle() {
	Styles(func() {
		ElementStyle("person", func() {
			Shape(ShapePerson)
			Background("#08427b")
			Color("#ffffff")
		})
	})
}

// PersonStyle defines the style used to render people. All elements tagged with
// "person" inherit the style.
func PersonStyle() {
	Styles(func() {
		ElementStyle("person", func() {
			Shape(ShapePerson)
			Background("#08427b")
			Color("#ffffff")
		})
	})
}
