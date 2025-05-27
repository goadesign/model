package dsl

import (
	"regexp"

	"goa.design/goa/v3/eval"
	"goa.design/model/expr"
)

type (
	// ShapeKind is the enum used to represent element shapes.
	ShapeKind int

	// BorderKind is the enum used to represent element border styles.
	BorderKind int
)

const (
	// ShapeBox represents a box shape for elements.
	ShapeBox ShapeKind = iota + 1
	// ShapeCircle represents a circle shape for elements.
	ShapeCircle
	// ShapeCylinder represents a cylinder shape for elements.
	ShapeCylinder
	// ShapeEllipse represents an ellipse shape for elements.
	ShapeEllipse
	// ShapeHexagon represents a hexagon shape for elements.
	ShapeHexagon
	// ShapeRoundedBox represents a rounded box shape for elements.
	ShapeRoundedBox
	// ShapeComponent represents a component shape for elements.
	ShapeComponent
	// ShapeFolder represents a folder shape for elements.
	ShapeFolder
	// ShapeMobileDeviceLandscape represents a mobile device landscape shape for elements.
	ShapeMobileDeviceLandscape
	// ShapeMobileDevicePortrait represents a mobile device portrait shape for elements.
	ShapeMobileDevicePortrait
	// ShapePerson represents a person shape for elements.
	ShapePerson
	// ShapePipe represents a pipe shape for elements.
	ShapePipe
	// ShapeRobot represents a robot shape for elements.
	ShapeRobot
	// ShapeWebBrowser represents a web browser shape for elements.
	ShapeWebBrowser
)

const (
	// BorderSolid represents a solid border style.
	BorderSolid BorderKind = iota + 1
	// BorderDashed represents a dashed border style.
	BorderDashed
	// BorderDotted represents a dotted border style.
	BorderDotted
)

// Styles is a wrapper for one or more element/relationship styles,
// which are used when rendering diagrams.
//
// Styles must appear in Views.
//
// Styles accepts a single argument: a function that defines the styles.
//
// Example:
//
//	var _ = Design(func() {
//	    var System = SoftwareSystem("Software System", "Great system.", func() {
//	        Tag("blue")
//	    })
//
//	    var User = Person("User", "A user of my software system.", func() {
//	        Tag("blue", "person")
//	        Uses(System, "Uses", func() {
//	            Tag("client")
//	        })
//	    })
//
//	    Views(func() {
//	        SystemContext(MySystem, "SystemContext", "Context diagram.", func() {
//	            AddAll()
//	            AutoLayout(RankTopBottom)
//	        })
//	        Styles(func() {
//	            ElementStyle("blue", func() {
//	                Background("#1168bd")
//	                Color("#3333ee")
//	             })
//	            ElementStyle("person", func() {
//	                Shape("ShapePerson")
//	            })
//	            RelationshipStyle("client", func() {
//	                Routing(RoutingCurved)
//	                Thickness(42)
//	            })
//	        })
//	    })
//	})
func Styles(dsl func()) {
	vs, ok := eval.Current().(*expr.Views)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	styles := &expr.Styles{}
	eval.Execute(dsl, styles)
	vs.Styles = styles
}

// ElementStyle defines element styles.
//
// ElementStyle must appear in Styles.
//
// ElementStyle accepts two arguments: the tag that identifies the elements that
// the style should be applied to and a function describing the style
// properties.
//
// Example:
//
//	var _ = Design(func() {
//	    // ...
//	    Views(func() {
//	        // ...
//	        Styles(func() {
//	            ElementStyle("default", func() {
//	                Shape(ShapeBox)
//	                Icon("https://goa.design/goa-logo.png")
//	                Background("#dddddd")
//	                Color("#000000")
//	                Stroke("#000000")
//	                ShowMetadata()
//	                ShowDescription()
//	            })
//	        })
//	    })
//	})
func ElementStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	es := &expr.ElementStyle{Tag: tag}
	eval.Execute(dsl, es)
	cfg.Elements = append(cfg.Elements, es)
}

// RelationshipStyle defines relationship styles.
//
// RelationshipStyle must appear in Styles.
//
// RelationshipStyle accepts two arguments: the tag that identifies the
// relationships that the style should be applied to and a function describing
// the style properties.
//
// Example:
//
//	var _ = Design(func() {
//	    // ...
//	    Views(func() {
//	        // ...
//	        Styles(func() {
//	            RelationshipStyle("default", func() {
//	                Thick()
//	                Color("#000000")
//	                Stroke("#000000")
//	                Solid()
//	                Routing(RoutingOrthogonal)
//	            })
//	        })
//	    })
//	})
func RelationshipStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	rs := &expr.RelationshipStyle{Tag: tag}
	eval.Execute(dsl, rs)
	cfg.Relationships = append(cfg.Relationships, rs)
}

// Shape defines element shapes, default is ShapeBox.
//
// Shape must apear in ElementStyle.
//
// Shape accepts one argument, one of: ShapeBox, ShapeRoundedBox, ShapeCircle,
// ShapeEllipse, ShapeHexagon or ShapeCylinder, ShapePipe, ShapePerson
// ShapeRobot, ShapeFolder, ShapeWebBrowser, ShapeMobileDevicePortrait,
// ShapeMobileDeviceLandscape or ShapeComponent.
func Shape(kind ShapeKind) {
	switch es := eval.Current().(type) {
	case *expr.ElementStyle:
		es.Shape = expr.ShapeKind(kind)
	default:
		eval.IncompatibleDSL()
	}
}

// Icon sets elements icon. Icon accepts the URL or data URI
// (https://css-tricks.com/data-uris/) of the icon.
//
// Tip: Generating icons programatically can be done using the "image" package
// (to draw the image), "image/png" to render the image and "encoding/base64" to
// encode the result into a data URI.
//
// Icon must appear in ElementStyle.
//
// Icon accepts URL to the icon image or a data URI.
func Icon(icon string) {
	switch es := eval.Current().(type) {
	case *expr.ElementStyle:
		es.Icon = icon
	default:
		eval.IncompatibleDSL()
	}
}

// Width sets elements or a relationships width, default is 450.
//
// Width must appear in ElementStyle or RelationshipStyle.
//
// Width accepts a single argument: the width in pixel.
func Width(width int) {
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Width = &width
	case *expr.RelationshipStyle:
		a.Width = &width
	default:
		eval.IncompatibleDSL()
	}
}

// Height sets elements height, default is 300.
//
// Height must appear in ElementStyle.
//
// Height accepts a single argument: the height in pixel.
func Height(height int) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Height = &height
		return
	}
	eval.IncompatibleDSL()
}

// colorRegex is used to validate strings that represent colors.
var colorRegex = regexp.MustCompile("#[A-Fa-f0-9]{6}")

// Background sets elements background color, default is #dddddd.
//
// Background must appear in ElementStyle.
//
// Background accepts a single argument: the background color encoded as HTML
// hex value (e.g. "#ffffff").
func Background(color string) {
	if !colorRegex.MatchString(color) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, color)
	}
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Background = color
		return
	}
	eval.IncompatibleDSL()
}

// Color sets elements text color, default is #000000.
//
// Color must appear in ElementStyle or RelationshipStyle.
//
// Color accepts a single argument: the color encoded as HTML hex value (e.g.
// "#ffffff").
func Color(color string) {
	if !colorRegex.MatchString(color) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, color)
	}
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Color = color
	case *expr.RelationshipStyle:
		a.Color = color
	default:
		eval.IncompatibleDSL()
	}
}

// Stroke sets elements stroke color.
//
// Stroke must appear in ElementStyle.
//
// Stroke accepts a single argument: the background color encoded as HTML
// hex value (e.g. "#ffffff").
func Stroke(color string) {
	if !colorRegex.MatchString(color) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, color)
	}
	switch es := eval.Current().(type) {
	case *expr.ElementStyle:
		es.Stroke = color
	case *expr.RelationshipStyle:
		es.Stroke = color
	default:
		eval.IncompatibleDSL()
	}
}

// FontSize sets elements or relationships text font size, default is 24.
//
// FontSize must appear in ElementStyle or RelationshipStyle.
//
// FontSize accepts a single argument: the size of the font in pixels.
func FontSize(pixels int) {
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.FontSize = &pixels
	case *expr.RelationshipStyle:
		a.FontSize = &pixels
	default:
		eval.IncompatibleDSL()
	}
}

// Border sets elements border style, default is BorderSolid.
//
// Border must appear in ElementStyle.
//
// Border takes a single argument: one of BorderSolid, BorderDashed or
// BorderDotted.
func Border(kind BorderKind) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Border = expr.BorderKind(kind)
		return
	}
	eval.IncompatibleDSL()
}

// Opacity sets elements or relationships opacity, default is 100.
//
// Opacity must appear in ElementStyle or RelationshipStyle.
//
// Opacity accepts a single argument: the opacity value between 0 (transparent)
// and 100 (opaque).
func Opacity(percent int) {
	if percent < 0 || percent > 100 {
		eval.InvalidArgError("value between 0 and 100", percent)
	}
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Opacity = &percent
	case *expr.RelationshipStyle:
		a.Opacity = &percent
	default:
		eval.IncompatibleDSL()
	}
}

// ShowMetadata shows the elements metadata.
//
// ShowMetadata must appear in ElementStyle.
//
// ShowMetadata takes no argument.
func ShowMetadata() {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		t := true
		es.Metadata = &t
		return
	}
	eval.IncompatibleDSL()
}

// ShowDescription shows the elements description.
//
// ShowDescription must appear in ElementStyle.
//
// ShowDescription takes no argument.
func ShowDescription() {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		t := true
		es.Description = &t
		return
	}
	eval.IncompatibleDSL()
}

// Thickness sets relationships thickness.
//
// Thickness must appear in RelationshipStyle.
//
// Thickness takes one argument: the thickness in pixels.
func Thickness(pixels int) {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		rs.Thickness = &pixels
		return
	}
	eval.IncompatibleDSL()
}

// Solid makes relationship lines solid (non-dashed).
//
// Solid must appear in RelationshipStyle.
//
// Solid takes no argument.
func Solid() {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		f := false
		rs.Dashed = &f
		return
	}
	eval.IncompatibleDSL()
}

// Dashed makes relationship lines dashed.
//
// Dashed must appear in RelationshipStyle.
//
// Dashed takes no argument.
func Dashed() {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		f := true
		rs.Dashed = &f
		return
	}
	eval.IncompatibleDSL()
}
