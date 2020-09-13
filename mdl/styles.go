package mdl

import (
	"fmt"
	"strings"

	"goa.design/model/expr"
)

// stroke computes the stroke color for the given element.
func stroke(data *elementData) string {
	s, bg := data.Stroke, data.Background
	if s != "" {
		return s
	}
	if bg == "" {
		bg = "#FFFFFF"
	}
	// Darken background by 50%
	var r, g, b int
	switch len(bg) {
	case 7:
		fmt.Sscanf(bg, "#%02x%02x%02x", &r, &g, &b)
	case 4:
		fmt.Sscanf(bg, "#%1x%1x%1x", &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
	}
	r /= 2
	g /= 2
	b /= 2
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// interpolate returns the D3 curve shape for the given relationship view if not
// "linear" (https://github.com/d3/d3-shape/blob/master/README.md#curves).
func interpolate(rs *expr.RelationshipStyle) string {
	switch rs.Routing {
	case expr.RoutingCurved:
		return "basis"
	case expr.RoutingOrthogonal:
		return "step"
	default:
		return ""
	}
}

// elementStyle compute the style of the given element view. It does that by
// merging all the styling information from all styles that apply (i.e. that
// apply to a tag of the corresponding element).
func elemStyle(ev *expr.ElementView) (style *expr.ElementStyle) {
	style = &expr.ElementStyle{}
	styles := expr.Root.Views.Styles
	if styles == nil {
		return
	}
loop:
	for _, tag := range strings.Split(ev.Element.Tags, ",") {
		for _, es := range styles.Elements {
			if tag == es.Tag {
				if es.Background != "" {
					style.Background = es.Background
				}
				if es.Stroke != "" {
					style.Stroke = es.Stroke
				}
				if es.Color != "" {
					style.Color = es.Color
				}
				if es.Shape != expr.ShapeUndefined {
					style.Shape = es.Shape
				}
				if es.Icon != "" {
					style.Icon = es.Icon
				}
				if es.Opacity != nil {
					style.Opacity = es.Opacity
				}
				if es.Metadata != nil {
					style.Metadata = es.Metadata
				}
				if es.Description != nil {
					style.Description = es.Description
				}
				if es.Border != expr.BorderUndefined {
					style.Border = es.Border
				}
				continue loop
			}
		}
	}
	return
}

// relationshipStyle compute the style of the given relationship view. It does that by
// merging all the styling information from all styles that apply (i.e. that
// apply to a tag of the corresponding relationship).
func relStyle(rv *expr.RelationshipView) (style *expr.RelationshipStyle) {
	style = &expr.RelationshipStyle{}
	styles := expr.Root.Views.Styles
	if styles == nil {
		return
	}
	rel := expr.Registry[rv.RelationshipID].(*expr.Relationship)
loop:
	for _, tag := range strings.Split(rel.Tags, ",") {
		for _, rs := range styles.Relationships {
			if tag == rs.Tag {
				if rs.Thick != nil {
					style.Thick = rs.Thick
				}
				if rs.Color != "" {
					style.Color = rs.Color
				}
				if rs.Stroke != "" {
					style.Stroke = rs.Stroke
				}
				if rs.Dashed != nil {
					style.Dashed = rs.Dashed
				}
				if rs.Routing != expr.RoutingUndefined {
					style.Routing = rs.Routing
				}
				if rs.Opacity != nil {
					style.Opacity = rs.Opacity
				}
				continue loop
			}
		}
	}
	return
}
