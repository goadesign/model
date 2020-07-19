package expr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	// ViewProps contains common properties of a view as well as helper
	// methods to fetch them.
	ViewProps struct {
		// Title of the view.
		Title string `json:"title,omitempty"`
		// Description of view.
		Description string `json:"description,omitempty"`
		// Key used to refer to the view.
		Key string `json:"key"`
		// PaperSize is the paper size that should be used to render this view.
		PaperSize PaperSizeKind `json:"paperSize,omitempty"`
		// Layout describes the automatic layout mode for the diagram if
		// defined.
		Layout *Layout `json:"automaticLayout,omitempty"`
		// ElementViews list the elements included in the view.
		ElementViews []*ElementView `json:"elements,omitempty"`
		// RelationshipViews list the relationships included in the view.
		RelationshipViews []*RelationshipView `json:"relationships,omitempty"`
		// Animations describes the animation steps if any.
		Animations []*Animation `json:"animationSteps,omitempty"`
	}

	// ElementView describes an instance of a model element (Person,
	// Software System, Container or Component) in a View.
	ElementView struct {
		// ID of element.
		ID string `json:"id"`
		// Horizontal position of element when rendered
		X *int `json:"x,omitempty"`
		// Vertical position of element when rendered.
		Y *int `json:"y,omitempty"`
		// Correpsonding model element.
		Element *Element `json:"-"`
		// Remove relationships before rendering.
		NoRelationship bool `json:"-"`
	}

	// RelationshipView describes an instance of a model relationship in a
	// view.
	RelationshipView struct {
		// ID of relationship.
		ID string `json:"id"`
		// Description of relationship used in dynamic views.
		Description string `json:"description,omitempty"`
		// Order of relationship in dynamic views.
		Order string `json:"order"`
		// Set of vertices used to render relationship
		Vertices []*Vertex `json:"vertices,omitempty"`
		// Routing algorithm used to render relationship.
		Routing RoutingKind `json:"routing"`
		// Position of annotation along line; 0 (start) to 100 (end).
		Position *int `json:"position,omitempty"`
		// Corresponding relationship.
		Relationship *Relationship `json:"-"`
	}

	// Vertex describes the x and y coordinate of a bend in a line.
	Vertex struct {
		// Horizontal position of vertex when rendered.
		X int `json:"x"`
		// Vertical position of vertex when rendered.
		Y int `json:"y"`
	}

	// Animation represents an animation step.
	Animation struct {
		// Order of animation step.
		Order int `json:"order"`
		// Set of element IDs that should be included.
		ElementIDs []string `json:"elements,omitempty"`
		// Set of relationship IDs tat should be included.
		Relationships []string `json:"relationships,omitempty"`
		// Set of element that should be included.
		Elements []ElementHolder `json:"-"`
	}

	// Layout describes an automatic layout.
	Layout struct {
		// Algorithm rank direction.
		RankDirection RankDirectionKind `json:"rankDirection"`
		// RankSep defines the separation between ranks in pixels.
		RankSep *int `json:"rankSeparation,omitempty"`
		// NodeSep defines the separation between nodes in pixels.
		NodeSep *int `json:"nodeSeparation,omitempty"`
		// EdgeSep defines the separation between edges in pixels.
		EdgeSep *int `json:"edgeSeparation,omitempty"`
		// Render vertices if true.
		Vertices *bool `json:"vertices,omitempty"`
	}
)

type (
	// PaperSizeKind is the enum for possible paper kinds.
	PaperSizeKind int

	// RoutingKind is the enum for possible routing algorithms.
	RoutingKind int

	// RankDirectionKind is the enum for possible automatic layout rank
	// directions.
	RankDirectionKind int
)

const (
	SizeA0Landscape PaperSizeKind = iota + 1
	SizeA0Portrait
	SizeA1Landscape
	SizeA1Portrait
	SizeA2Landscape
	SizeA2Portrait
	SizeA3Landscape
	SizeA3Portrait
	SizeA4Landscape
	SizeA4Portrait
	SizeA5Landscape
	SizeA5Portrait
	SizeA6Landscape
	SizeA6Portrait
	SizeLegalLandscape
	SizeLegalPortrait
	SizeLetterLandscape
	SizeLetterPortrait
	SizeSlide16X10
	SizeSlide16X9
	SizeSlide4X3
)

const (
	RoutingCurved RoutingKind = iota + 1
	RoutingDirect
	RoutingOrthogonal
)

const (
	RankTopBottom RankDirectionKind = iota + 1
	RankBottomTop
	RankLeftRight
	RankRightLeft
)

// ElementView returns the element view with the given ID if any.
func (v *ViewProps) ElementView(id string) *ElementView {
	for _, e := range v.ElementViews {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// RelationshipView returns the relationship view with the given ID if any.
func (v *ViewProps) RelationshipView(id string) *RelationshipView {
	for _, r := range v.RelationshipViews {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// AllTagged returns all elements with the given tag in the view.
func (v *ViewProps) AllTagged(tag string) (elts []*Element) {
	for _, ev := range v.ElementViews {
		vals := strings.Split(ev.Element.Tags, ",")
		for _, val := range vals {
			if strings.Trim(val, " ") == tag {
				elts = append(elts, ev.Element)
				break
			}
		}
	}
	return
}

// AllUnreachable fetches all elements in view related to the element (directly
// or not).
func (v *ViewProps) AllUnreachable(eh ElementHolder) (elts []*Element) {
	e := eh.GetElement()
	if v.index(e.ID) == -1 {
		return
	}
	ids := e.Reachable()
loop:
	for _, e := range v.ElementViews {
		for _, id := range ids {
			if id == e.ID {
				continue loop
			}
		}
		elts = append(elts, e.Element)
	}
	return
}

// AllUnrelated fetches all elements that have no relationship to other elements
// in the view.
func (v *ViewProps) AllUnrelated() (elts []*Element) {
	for _, ev := range v.ElementViews {
		related := false
		for _, r := range v.RelationshipViews {
			if r.Relationship.SourceID == ev.ID {
				related = true
				break
			}
			if r.Relationship.FindDestination().ID == ev.ID {
				related = true
				break
			}
		}
		if !related {
			elts = append(elts, ev.Element)
		}
	}
	return
}

// Props returns the underlying properties object.
func (v *ViewProps) Props() *ViewProps { return v }

// index returns the index of the element with the given ID, -1 if not found.
func (v *ViewProps) index(id string) int {
	for i, e := range v.ElementViews {
		if e.ID == id {
			return i
		}
	}
	return -1
}

// EvalName returns the generic expression name used in error messages.
func (v *ViewProps) EvalName() string {
	var suf string
	switch {
	case v.Title != "":
		suf = fmt.Sprintf(" with title %q", v.Title)
	case v.Title != "":
		suf = fmt.Sprintf(" with  key %q", v.Key)
	case v.Title != "" && v.Key != "":
		suf = fmt.Sprintf(" with title %q and key %q", v.Title, v.Key)
	}
	return fmt.Sprintf("view%s", suf)
}

// Finalize computes the relationships used in the view.
func (v *ViewProps) Finalize() {
	var rels []*Relationship
	for _, x := range Registry {
		r, ok := x.(*Relationship)
		if !ok {
			continue
		}
		for _, ev := range v.ElementViews {
			if r.SourceID == ev.ID {
				if v.ElementView(r.FindDestination().ID) != nil {
					rels = append(rels, r)
				}
			}
		}
	}
	addRelationships(v, rels)
}

// EvalName returns the generic expression name used in error messages.
func (l *Layout) EvalName() string { return "automatic layout" }

// EvalName returns the generic expression name used in error messages.
func (v *ElementView) EvalName() string { return "element view" }

// EvalName returns the generic expression name used in error messages.
func (v *RelationshipView) EvalName() string { return "relationship view" }

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (p PaperSizeKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch p {
	case SizeA6Portrait:
		buf.WriteString("A6_Portrait")
	case SizeA6Landscape:
		buf.WriteString("A6_Landscape")
	case SizeA5Portrait:
		buf.WriteString("A5_Portrait")
	case SizeA5Landscape:
		buf.WriteString("A5_Landscape")
	case SizeA4Portrait:
		buf.WriteString("A4_Portrait")
	case SizeA4Landscape:
		buf.WriteString("A4_Landscape")
	case SizeA3Portrait:
		buf.WriteString("A3_Portrait")
	case SizeA3Landscape:
		buf.WriteString("A3_Landscape")
	case SizeA2Portrait:
		buf.WriteString("A2_Portrait")
	case SizeA2Landscape:
		buf.WriteString("A2_Landscape")
	case SizeA1Portrait:
		buf.WriteString("A1_Portrait")
	case SizeA1Landscape:
		buf.WriteString("A1_Landscape")
	case SizeA0Portrait:
		buf.WriteString("A0_Portrait")
	case SizeA0Landscape:
		buf.WriteString("A0_Landscape")
	case SizeLetterPortrait:
		buf.WriteString("Letter_Portrait")
	case SizeLetterLandscape:
		buf.WriteString("Letter_Landscape")
	case SizeLegalPortrait:
		buf.WriteString("Legal_Portrait")
	case SizeLegalLandscape:
		buf.WriteString("Legal_Landscape")
	case SizeSlide4X3:
		buf.WriteString("Slide_4_3")
	case SizeSlide16X9:
		buf.WriteString("Slide_16_9")
	case SizeSlide16X10:
		buf.WriteString("Slide_16_10")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (p *PaperSizeKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "A6_Portrait":
		*p = SizeA6Portrait
	case "A6_Landscape":
		*p = SizeA6Landscape
	case "A5_Portrait":
		*p = SizeA5Portrait
	case "A5_Landscape":
		*p = SizeA5Landscape
	case "A4_Portrait":
		*p = SizeA4Portrait
	case "A4_Landscape":
		*p = SizeA4Landscape
	case "A3_Portrait":
		*p = SizeA3Portrait
	case "A3_Landscape":
		*p = SizeA3Landscape
	case "A2_Portrait":
		*p = SizeA2Portrait
	case "A2_Landscape":
		*p = SizeA2Landscape
	case "A1_Portrait":
		*p = SizeA1Portrait
	case "A1_Landscape":
		*p = SizeA1Landscape
	case "A0_Portrait":
		*p = SizeA0Portrait
	case "A0_Landscape":
		*p = SizeA0Landscape
	case "Letter_Portrait":
		*p = SizeLetterPortrait
	case "Letter_Landscape":
		*p = SizeLetterLandscape
	case "Legal_Portrait":
		*p = SizeLegalPortrait
	case "Legal_Landscape":
		*p = SizeLegalLandscape
	case "Slide_4_3":
		*p = SizeSlide4X3
	case "Slide_16_9":
		*p = SizeSlide16X9
	case "Slide_16_10":
		*p = SizeSlide16X10
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RoutingKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case RoutingDirect:
		buf.WriteString("Direct")
	case RoutingCurved:
		buf.WriteString("Curved")
	case RoutingOrthogonal:
		buf.WriteString("Orthogonal")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (r *RoutingKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "Direct":
		*r = RoutingDirect
	case "Curved":
		*r = RoutingCurved
	case "Orthogonal":
		*r = RoutingOrthogonal
	}
	return nil
}

// MarshalJSON replaces the constant value with the proper structurizr schema
// string value.
func (r RankDirectionKind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	switch r {
	case RankTopBottom:
		buf.WriteString("TopBottom")
	case RankBottomTop:
		buf.WriteString("BottomTop")
	case RankLeftRight:
		buf.WriteString("LeftRight")
	case RankRightLeft:
		buf.WriteString("RightLeft")
	}
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshalJSON sets the constant from its JSON representation.
func (r *RankDirectionKind) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case "TopBottom":
		*r = RankTopBottom
	case "BottomTop":
		*r = RankBottomTop
	case "LeftRight":
		*r = RankLeftRight
	case "RightLeft":
		*r = RankRightLeft
	}
	return nil
}
