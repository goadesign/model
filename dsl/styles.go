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
	// Shapes allowed in ElementStyle
	ShapeBox ShapeKind = iota + 1
	ShapeCircle
	ShapeCylinder
	ShapeEllipse
	ShapeHexagon
	ShapeRoundedBox

	// Shapes allowed in StructurizrElementStyle
	ShapeComponent
	ShapeFolder
	ShapeMobileDeviceLandscape
	ShapeMobileDevicePortrait
	ShapePerson
	ShapePipe
	ShapeRobot
	ShapeWebBrowser
)

const (
	BorderSolid BorderKind = iota + 1
	BorderDashed
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
//     var _ = Design(func() {
//         var System = SoftwareSystem("Software System", "Great system.", func() {
//             Tag("blue")
//         })
//
//         var User = Person("User", "A user of my software system.", func() {
//             Tag("blue", "person")
//             Uses(System, "Uses", func() {
//                 Tag("client")
//             })
//         })
//
//         Views(func() {
//             SystemContext(MySystem, "SystemContext", "Context diagram.", func() {
//                 AddAll()
//                 AutoLayout(RankTopBottom)
//             })
//             Styles(func() {
//                 ElementStyle("blue", func() {
//                     Background("#1168bd")
//                     Color("#3333ee")
//                  })
//                 ElementStyle("person", func() {
//                     Shape("ShapePerson")
//                 })
//                 RelationshipStyle("client", func() {
//                     Routing(RoutingCurved)
//                     Thickness(42)
//                 })
//             })
//         })
//     })
//
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
//     var _ = Design(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 ElementStyle("default", func() {
//                     Shape(ShapeBox)
//                     Icon("https://goa.design/goa-logo.png")
//                     Background("#dddddd")
//                     Color("#000000")
//                     Stroke("#000000")
//                     ShowMetadata()
//                     ShowDescription()
//                 })
//             })
//         })
//     })
//
func ElementStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	es := &expr.ElementStyle{Tag: tag}
	eval.Execute(dsl, es)
	cfg.Elements = append(cfg.Elements, es)
}

// StructurizrElementStyle defines additional element styles used for views
// rendered in the Structurizr service. Shape accepts additional values when
// used in StructurizrElementStyle.
//
// StructurizrElementStyle must appear in Styles.
//
// StructurizrElementStyle accepts two arguments: the tag that identifies the
// elements that the style should be applied to and a function describing the
// style properties.
//
// Example:
//
//     var _ = Design(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 StructurizrElementStyle("default", func() {
//                     Shape(ShapeRobot)
//                     Icon("https://goa.design/goa-logo.png")
//                     Width(300)
//                     Height(450)
//                     FontSize(24)
//                     Border(BorderSolid)
//                     Opacity(100)
//                 })
//             })
//         })
//     })
//
func StructurizrElementStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	es := &expr.StructurizrElementStyle{Tag: tag}
	eval.Execute(dsl, es)
	cfg.StructurizrElements = append(cfg.StructurizrElements, es)
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
//     var _ = Design(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 RelationshipStyle("default", func() {
//                     Thick()
//                     Color("#000000")
//                     Stroke("#000000")
//                     Solid()
//                     Routing(RoutingOrthogonal)
//                 })
//             })
//         })
//     })
//
func RelationshipStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	rs := &expr.RelationshipStyle{Tag: tag}
	eval.Execute(dsl, rs)
	cfg.Relationships = append(cfg.Relationships, rs)
}

// StructurizrRelationshipStyle defines additional relationship styles that
// apply to views rendered in the Structurizr service.
//
// StructurizrRelationshipStyle must appear in Styles.
//
// StructurizrRelationshipStyle accepts two arguments: the tag that identifies
// the relationships that the style should be applied to and a function
// describing the style properties.
//
// Example:
//
//     var _ = Design(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 StructurizrRelationshipStyle("default", func() {
//                     Thickness(2)
//                     FontSize(24)
//                     Width(300)
//                     Position(50)
//                     Opacity(100)
//                 })
//             })
//         })
//     })
//
func StructurizrRelationshipStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Styles)
	if !ok {
		eval.IncompatibleDSL()
	}
	rs := &expr.StructurizrRelationshipStyle{Tag: tag}
	eval.Execute(dsl, rs)
	cfg.StructurizrRelationships = append(cfg.StructurizrRelationships, rs)
}

// Shape defines element shapes, default is ShapeBox.
//
// Shape must apear in ElementStyle or StructurizrElementStyle.
//
// Shape accepts one argument, one of: ShapeBox, ShapeRoundedBox, ShapeCircle,
// ShapeEllipse, ShapeHexagon or ShapeCylinder. Additionally when used
// in StructurizrElementStyle Shape also accepts one of ShapePipe, ShapePerson
// ShapeRobot, ShapeFolder, ShapeWebBrowser, ShapeMobileDevicePortrait,
// ShapeMobileDeviceLandscape or ShapeComponent.
func Shape(kind ShapeKind) {
	switch es := eval.Current().(type) {
	case *expr.ElementStyle:
		if int(kind) >= int(expr.ShapeLast) {
			eval.ReportError("Shape: value can only be used in StructurizrElementStyle")
		}
		es.Shape = expr.ShapeKind(kind)
	case *expr.StructurizrElementStyle:
		es.Shape = expr.ExtendedShapeKind(kind)
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
// Icon accepts URL to the icon image or a data URI (StructurizrElementStyle).
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
// Width must appear in StructurizrElementStyle or StructurizrRelationshipStyle.
//
// Width accepts a single argument: the width in pixel.
func Width(width int) {
	switch a := eval.Current().(type) {
	case *expr.StructurizrElementStyle:
		a.Width = &width
	case *expr.StructurizrRelationshipStyle:
		a.Width = &width
	default:
		eval.IncompatibleDSL()
	}
}

// Height sets elements height, default is 300.
//
// Height must appear in StructurizrElementStyle.
//
// Height accepts a single argument: the height in pixel.
func Height(height int) {
	if es, ok := eval.Current().(*expr.StructurizrElementStyle); ok {
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
// FontSize must appear in StructurizrElementStyle or
// StructurizrRelationshipStyle.
//
// FontSize accepts a single argument: the size of the font in pixels.
func FontSize(pixels int) {
	switch a := eval.Current().(type) {
	case *expr.StructurizrElementStyle:
		a.FontSize = &pixels
	case *expr.StructurizrRelationshipStyle:
		a.FontSize = &pixels
	default:
		eval.IncompatibleDSL()
	}
}

// Border sets elements border style, default is BorderSolid.
//
// Border must appear in StructurizrElementStyle.
//
// Border takes a single argument: one of BorderSolid, BorderDashed or
// BorderDotted.
func Border(kind BorderKind) {
	if es, ok := eval.Current().(*expr.StructurizrElementStyle); ok {
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
	if percent < 0 || 0 > 100 {
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

// Thick renders a thick line to represent the relationship.
//
// Thick must appear in RelationshipStyle.
//
// Thick takes no argument.
func Thick() {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		t := true
		rs.Thick = &t
		return
	}
	eval.IncompatibleDSL()
}

// Thickness sets relationships thickness.
//
// Thickness must appear in StructurizrRelationshipStyle.
//
// Thickness takes one argument: the thickness in pixels.
func Thickness(pixels int) {
	if rs, ok := eval.Current().(*expr.StructurizrRelationshipStyle); ok {
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
