package dsl

import (
	"regexp"

	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"
)

const (
	SymbolSquareBrackets      = expr.SymbolSquareBrackets
	SymbolRoundBrackets       = expr.SymbolRoundBrackets
	SymbolCurlyBrackets       = expr.SymbolCurlyBrackets
	SymbolAngleBrackets       = expr.SymbolAngleBrackets
	SymbolDoubleAngleBrackets = expr.SymbolDoubleAngleBrackets
	SymbolNone                = expr.SymbolNone
)

const (
	ShapeBox                   = expr.ShapeBox
	ShapeRoundedBox            = expr.ShapeRoundedBox
	ShapeComponent             = expr.ShapeComponent
	ShapeCircle                = expr.ShapeCircle
	ShapeEllipse               = expr.ShapeEllipse
	ShapeHexagon               = expr.ShapeHexagon
	ShapeFolder                = expr.ShapeFolder
	ShapeCylinder              = expr.ShapeCylinder
	ShapePipe                  = expr.ShapePipe
	ShapeWebBrowser            = expr.ShapeWebBrowser
	ShapeMobileDevicePortrait  = expr.ShapeMobileDevicePortrait
	ShapeMobileDeviceLandscape = expr.ShapeMobileDeviceLandscape
	ShapePerson                = expr.ShapePerson
	ShapeRobot                 = expr.ShapeRobot
)

const (
	BorderSolid  = expr.BorderSolid
	BorderDashed = expr.BorderDashed
	BorderDotted = expr.BorderDotted
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
//     var _ = Workspace(func() {
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
	cfg := &expr.Configuration{Styles: &expr.Styles{}}
	eval.Execute(dsl, cfg)
	vs.Configuration = cfg
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
//     var _ = Workspace(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 ElementStyle("default", func() {
//                     Shape(ShapeBox)
//                     Icon("https://goa.design/goa-logo.png")
//                     Width(300)
//                     Height(450)
//                     Background("#dddddd")
//                     Color("#000000")
//                     Stroke("#000000")
//                     FontSize(24)
//                     Boder(BorderSolid)
//                     Opacity(100)
//                     ShowMetadata()
//                     ShowDescription()
//                 })
//             })
//         })
//     })
//
func ElementStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Configuration)
	if !ok {
		eval.IncompatibleDSL()
	}
	es := &expr.ElementStyle{Tag: tag}
	eval.Execute(dsl, es)
	cfg.Styles.Elements = append(cfg.Styles.Elements, es)
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
//     var _ = Workspace(func() {
//         // ...
//         Views(func() {
//             // ...
//             Styles(func() {
//                 RelationshipStyle("default", func() {
//                     Thickness(2)
//                     Color("#000000")
//                     Dashed(true)
//                     Routing(RoutingOrthogonal)
//                     FontSize(24)
//                     Width(300)
//                     Position(50)
//                     Opacity(100)
//                 })
//             })
//         })
//     })
//
func RelationshipStyle(tag string, dsl func()) {
	cfg, ok := eval.Current().(*expr.Configuration)
	if !ok {
		eval.IncompatibleDSL()
	}
	rs := &expr.RelationshipStyle{Tag: tag}
	eval.Execute(dsl, rs)
	cfg.Styles.Relationships = append(cfg.Styles.Relationships, rs)
}

// Shape defines element shapes, default is ShapeBox.
//
// Shape must apear in ElementStyle.
//
// Shape accepts one argument, one of: ShapeBox, ShapeRoundedBox, ShapeCircle,
// ShapeEllipse, ShapeHexagon, ShapeCylinder, ShapePipe, ShapePerson ShapeRobot,
// ShapeFolder, ShapeWebBrowser, ShapeMobileDevicePortrait,
// ShapeMobileDeviceLandscape or ShapeComponent.
//
// See https://structurizr.com/help/shapes.
func Shape(s expr.ShapeKind) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Shape = s
		return
	}
	eval.IncompatibleDSL()
}

// Icon sets elements icon by URL or data URI
// (https://css-tricks.com/data-uris/).
//
// Tip: Generating icons programatically can be done using the "image" package
// (to draw the image), "image/png" to render the image and "encoding/base64" to
// encode the result into a data URI.
//
// See https://structurizr.com/help/icons.
func Icon(file string) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Icon = file
		return
	}
	eval.IncompatibleDSL()
}

// Width sets elements or a relationships width, default is 450.
//
// Width must appear in ElementStyle or RelationshipStyle.
//
// Width accepts a single argument: the width in pixel.
func Width(w int) {
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Width = w
	case *expr.RelationshipStyle:
		a.Width = w
	default:
		eval.IncompatibleDSL()
	}
}

// Height sets elements height, default is 300.
//
// Height must appear in ElementStyle.
//
// Height accepts a single argument: the height in pixel.
func Height(h int) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Height = h
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
func Background(c string) {
	if !colorRegex.MatchString(c) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, c)
	}
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Background = c
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
func Color(c string) {
	if !colorRegex.MatchString(c) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, c)
	}
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Color = c
	case *expr.RelationshipStyle:
		a.Color = c
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
func Stroke(c string) {
	if !colorRegex.MatchString(c) {
		eval.InvalidArgError(`color hex value (e.g. "#ffffff")`, c)
	}
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Stroke = c
		return
	}
	eval.IncompatibleDSL()
}

// FontSize sets elements or relationships text font size, default is 24.
//
// FontSize must appear in ElementStyle or RelationshipStyle.
//
// FontSize accepts a single argument: the size of the font in pixels.
func FontSize(s int) {
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.FontSize = s
	case *expr.RelationshipStyle:
		a.FontSize = s
	default:
		eval.IncompatibleDSL()
	}
}

// Boder sets elements border style, default is BorderSolid.
//
// Border must appear in ElementStyle.
//
// Border takes a single argument: one of BorderSolid, BorderDashed or
// BorderDotted.
func Boder(b expr.BorderKind) {
	if es, ok := eval.Current().(*expr.ElementStyle); ok {
		es.Border = b
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
func Opacity(o int) {
	if o < 0 || 0 > 100 {
		eval.InvalidArgError("value between 0 and 100", o)
	}
	switch a := eval.Current().(type) {
	case *expr.ElementStyle:
		a.Opacity = o
	case *expr.RelationshipStyle:
		a.Opacity = o
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
		es.Metadata = true
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
		es.Description = true
		return
	}
	eval.IncompatibleDSL()
}

// Thickness sets relationships thickness.
//
// Thickness must appear in RelationshipStyle.
//
// Thickness takes one argument: the thickness in pixels.
func Thickness(t int) {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		rs.Thickness = t
		return
	}
	eval.IncompatibleDSL()
}

// Dashed dashes relationship lines.
//
// Dashed must appear in RelationshipStyle.
//
// Dashed takes no argument.
func Dashed() {
	if rs, ok := eval.Current().(*expr.RelationshipStyle); ok {
		rs.Dashed = true
		return
	}
	eval.IncompatibleDSL()
}
