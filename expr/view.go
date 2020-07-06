package expr

import (
	"bytes"
	"encoding/json"
	"strings"
)

type (
	// ViewProps contains common properties of a view.
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
		// AnimationSteps describes the animation steps if any.
		AnimationSteps []*AnimationStep `json:"animationSteps,omitempty"`
	}

	// ElementView describes an instance of a model element (Person,
	// Software System, Container or Component) in a View.
	ElementView struct {
		// ID of element.
		ID string `json:"id"`
		// Horizontal position of element when rendered.
		X int `json:"x"`
		// Vertical position of element when rendered.
		Y int `json:"y"`
		// Correpsonding model element.
		Element *Element `json:"-"`
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
		Vertices []*Vertex `json:"vertices"`
		// Routing algorithm used to render relationship.
		Routing RoutingKind `json:"routing"`
		// Position of annotation along line; 0 (start) to 100 (end).
		Position int `json:"position"`
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

	// AnimationStep represents an animation step.
	AnimationStep struct {
		// Order of animation step.
		Order string `json:"order"`
		// Set of element IDs that should be included.
		Elements []string `json:"elements,omitempty"`
		// Set of relationship IDs tat should be included.
		Relationships []string `json:"relationships,omitempty"`
	}

	// Layout describes an automatic layout.
	Layout struct {
		// Algorithm rank direction.
		RankDirection RankDirectionKind `json:"rankDirection"`
		// RankSep defines the separation between ranks in pixels.
		RankSep int `json:"rankSeparation"`
		// NodeSep defines the separation between nodes in pixels.
		NodeSep int `json:"nodeSeparation"`
		// EdgeSep defines the separation between edges in pixels.
		EdgeSep int `json:"edgeSeparation"`
		// Render vertices if true.
		Vertices bool
	}
)

type (
	// View is the common interface for all views.
	View interface {
		GetElementView(string) *ElementView
		GetRelationshipView(string) *RelationshipView
		AddElements(...ElementHolder)
		AddRelationships(...*Relationship)
		Remove(id string)
		RemoveTagged(tag string)
		RemoveUnreachable(*Element)
		RemoveUnrelated()
	}

	// Processor function that processes elements. The argument must be
	// either a slice of structs that implement ElementHolder.
	Processor func(...ElementHolder)
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
	SizeA6Portrait PaperSizeKind = iota + 1
	SizeA6Landscape
	SizeA5Portrait
	SizeA5Landscape
	SizeA4Portrait
	SizeA4Landscape
	SizeA3Portrait
	SizeA3Landscape
	SizeA2Portrait
	SizeA2Landscape
	SizeA1Portrait
	SizeA1Landscape
	SizeA0Portrait
	SizeA0Landscape
	SizeLetterPortrait
	SizeLetterLandscape
	SizeLegalPortrait
	SizeLegalLandscape
	SizeSlide4X3
	SizeSlide16X9
	SizeSlide16X10
)

const (
	RoutingDirect RoutingKind = iota + 1
	RoutingCurved
	RoutingOrthogonal
)

const (
	RankTopBottom RankDirectionKind = iota + 1
	RankBottomTop
	RankLeftRight
	RankRightLeft
)

// AddElements adds the given elements to the view if not already present.
func (v *ViewProps) AddElements(ehs ...ElementHolder) {
loop:
	for _, eh := range ehs {
		id := eh.GetElement().ID
		for _, e := range v.ElementViews {
			if e.ID == id {
				continue loop
			}
		}
		v.ElementViews = append(v.ElementViews, &ElementView{ID: id, Element: eh.GetElement()})
		v.completeRelationships(id)
	}
}

// AddRelationships adds the given relationships to the view if not already
// present. It does nothing if the relationship source and destination are not
// already in the view.
func (v *ViewProps) AddRelationships(rels ...*Relationship) {
loop:
	for _, r := range rels {
		var src, dest bool
		for _, ev := range v.ElementViews {
			if ev.ID == r.SourceID {
				src = true
				if dest {
					break
				}
			}
			if ev.ID == r.DestinationID {
				dest = true
				if src {
					break
				}
			}
		}
		if !src || !dest {
			continue loop
		}
		for _, rv := range v.RelationshipViews {
			if rv.ID == r.ID {
				continue loop
			}
		}
		v.RelationshipViews = append(v.RelationshipViews, &RelationshipView{ID: r.ID, Relationship: r})
	}
}

// Remove removes the element with the given ID from the view if present.
func (v *ViewProps) Remove(id string) {
	idx := v.index(id)
	if idx == -1 {
		return
	}
	v.ElementViews = append(v.ElementViews[:idx], v.ElementViews[idx+1:]...)

	// Remove corresponding relationships.
	var ids []string
	for _, r := range v.RelationshipViews {
		if r.Relationship.SourceID == id {
			ids = append(ids, id)
		} else if r.Relationship.DestinationID == id {
			ids = append(ids, id)
		}
	}
	rvs := v.RelationshipViews
	tmp := rvs[:0]
	for _, r := range rvs {
		remove := false
		for _, id := range ids {
			if r.ID == id {
				remove = true
				break
			}
		}
		if !remove {
			tmp = append(tmp, r)
		}
	}
	v.RelationshipViews = tmp
}

// RemoveTagged removes all elements with the given tag from the view.
func (v *ViewProps) RemoveTagged(tag string) {
	var rm []string
	for _, ev := range v.ElementViews {
		vals := strings.Split(ev.Element.Tags, ",")
		for _, val := range vals {
			if strings.Trim(val, " ") == tag {
				rm = append(rm, ev.ID)
				break
			}
		}
	}
	for _, id := range rm {
		v.Remove(id)
	}
}

// RemoveUnreachable removes all elements that are not related - directly or not
// - to the element.
func (v *ViewProps) RemoveUnreachable(elt *Element) {
	if v.index(elt.ID) == -1 {
		return
	}
	var rm []string
	ids := elt.Reachable()
loop:
	for _, e := range v.ElementViews {
		for _, id := range ids {
			if id == e.ID {
				continue loop
			}
		}
		rm = append(rm, e.ID)
	}
	for _, id := range rm {
		v.Remove(id)
	}
}

// RemoveUnrelated removes all elements that have no relationship to other
// elements in the view.
func (v *ViewProps) RemoveUnrelated() {
	for _, ev := range v.ElementViews {
		related := false
		for _, r := range v.RelationshipViews {
			if r.Relationship.SourceID == ev.ID {
				related = true
				break
			}
			if r.Relationship.DestinationID == ev.ID {
				related = true
				break
			}
		}
		if !related {
			v.Remove(ev.ID)
		}
	}
}

// GetElementView returns the element view with the given ID if any.
func (v *ViewProps) GetElementView(id string) *ElementView {
	for _, e := range v.ElementViews {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// GetRelationshipView returns the relationship view with the given ID if any.
func (v *ViewProps) GetRelationshipView(id string) *RelationshipView {
	for _, r := range v.RelationshipViews {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// completeRelationships adds the relationships for which the element with the
// given id is either a source or a destination and the other end of the
// relationship is already in the view.
func (v *ViewProps) completeRelationships(id string) {
	var rels []*Relationship
	for _, r := range AllRelationships() {
		if r.SourceID == id {
			if v.GetElementView(r.DestinationID) != nil {
				rels = append(rels, r)
			}
		} else if r.DestinationID == id {
			if v.GetElementView(r.SourceID) != nil {
				rels = append(rels, r)
			}
		}
	}
	v.AddRelationships(rels...)
}

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
