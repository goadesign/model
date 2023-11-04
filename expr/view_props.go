package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
)

type (
	// ViewProps contains common properties of a view as well as helper
	// methods to fetch them.
	ViewProps struct {
		Key               string
		Description       string
		Title             string
		AutoLayout        *AutoLayout
		PaperSize         PaperSizeKind
		ElementViews      []*ElementView
		RelationshipViews []*RelationshipView
		AnimationSteps    []*AnimationStep

		// The following fields are used to compute the elements and
		// relationships that should be added to the view.
		AddAll              bool
		AddDefault          bool
		AddNeighbors        []*Element
		RemoveElements      []*Element
		RemoveTags          []string
		RemoveRelationships []*Relationship
		RemoveUnreachable   []*Element
		RemoveUnrelated     bool
	}

	// ElementView describes an instance of a model element (Person,
	// Software System, Container or Component) in a View.
	ElementView struct {
		Element        *Element
		NoRelationship bool
		X              *int
		Y              *int
	}

	// RelationshipView describes an instance of a model relationship in a
	// view.
	RelationshipView struct {
		Source      *Element
		Destination *Element
		Description string
		Order       string
		Vertices    []*Vertex
		Routing     RoutingKind
		Position    *int

		// RelationshipID is computed in finalize.
		RelationshipID string
	}

	// AutoLayout describes an automatic layout.
	AutoLayout struct {
		Implementation ImplementationKind
		RankDirection  RankDirectionKind
		RankSep        *int
		NodeSep        *int
		EdgeSep        *int
		Vertices       *bool
	}

	// Vertex describes the x and y coordinate of a bend in a line.
	Vertex struct {
		X int
		Y int
	}

	// AnimationStep represents an animation step.
	AnimationStep struct {
		Elements        []ElementHolder
		RelationshipIDs []string
		Order           int
		View            View
	}
	// PaperSizeKind is the enum for possible paper kinds.
	PaperSizeKind int

	// RoutingKind is the enum for possible routing algorithms.
	RoutingKind int

	// ImplementationKind is the enum for possible automatic layout implementations
	ImplementationKind int

	// RankDirectionKind is the enum for possible automatic layout rank
	// directions.
	RankDirectionKind int
)

const (
	SizeUndefined PaperSizeKind = iota
	SizeA0Landscape
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
	RoutingUndefined RoutingKind = iota
	RoutingDirect
	RoutingOrthogonal
	RoutingCurved
)

const (
	ImplementationUndefined ImplementationKind = iota
	ImplementationGraphviz
	ImplementationDagre
)

const (
	RankUndefined RankDirectionKind = iota
	RankTopBottom
	RankBottomTop
	RankLeftRight
	RankRightLeft
)

// ElementView returns the element view for the element with the given ID if
// any.
func (v *ViewProps) ElementView(id string) *ElementView {
	for _, ev := range v.ElementViews {
		if ev.Element.ID == id {
			return ev
		}
	}
	return nil
}

// Props returns the underlying properties object.
func (v *ViewProps) Props() *ViewProps { return v }

// EvalName returns the generic expression name used in error messages.
func (v *ViewProps) EvalName() string {
	var suf string
	switch {
	case v.Title != "":
		suf = fmt.Sprintf(" with title %q and key %q", v.Title, v.Key)
	default:
		suf = fmt.Sprintf(" %q", v.Key)
	}
	return fmt.Sprintf("view%s", suf)
}

// EvalName returns the generic expression name used in error messages.
func (l *AnimationStep) EvalName() string { return "animation step" }

// Add adds the given elements to the animation step.
func (l *AnimationStep) Add(eh ElementHolder) {
	l.Elements = append(l.Elements, eh)
}

// EvalName returns the generic expression name used in error messages.
func (l *AutoLayout) EvalName() string { return "automatic layout" }

// EvalName returns the generic expression name used in error messages.
func (v *ElementView) EvalName() string { return "element view" }

// EvalName returns the generic expression name used in error messages.
func (v *RelationshipView) EvalName() string { return "relationship view" }

// Validate makes sure there is a corresponding relationship (and exactly one).
func (v *RelationshipView) Validate() error {
	verr := new(eval.ValidationErrors)
	var rel *Relationship
	found := false
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == v.Source.ID && r.Destination.ID == v.Destination.ID {
			if v.Description != "" {
				if r.Description == v.Description {
					rel = r
				}
			} else {
				rel = r
				if found {
					verr.Add(v, "Link: there exists multiple relationships between %q and %q, specify the relationship description.", v.Source.Name, v.Destination.Name)
				}
				found = true
			}
		}
	})
	if rel == nil {
		var suffix string
		if v.Description != "" {
			suffix = fmt.Sprintf(" with description %q", v.Description)
		}
		verr.Add(v, "Link: no relationship between %q and %q%s", v.Source.Name, v.Destination.Name, suffix)
	}
	return verr
}

// reachable returns the IDs of all elements that can be reached by traversing
// the relationships from the given root.
func reachable(e *Element) (res []string) {
	seen := make(map[string]struct{})
	traverse(e, seen)
	res = make([]string, len(seen))
	for k := range seen {
		res = append(res, k)
	}
	return
}

func traverse(e *Element, seen map[string]struct{}) {
	add := func(nid string) bool {
		for id := range seen {
			if id == nid {
				return false
			}
		}
		seen[nid] = struct{}{}
		return true
	}
	IterateRelationships(func(r *Relationship) {
		if r.Source.ID == e.ID {
			if add(r.Destination.ID) {
				traverse(r.Destination, seen)
			}
		}
		if r.Destination.ID == e.ID {
			if add(r.Source.ID) {
				traverse(r.Source, seen)
			}
		}
	})
}

// Name returns the name of the paper size kind.
func (k PaperSizeKind) Name() string {
	switch k {
	case SizeA0Landscape:
		return "A0_Landscape"
	case SizeA0Portrait:
		return "A0_Portrait"
	case SizeA1Landscape:
		return "A1_Landscape"
	case SizeA1Portrait:
		return "A1_Portrait"
	case SizeA2Landscape:
		return "A2_Landscape"
	case SizeA2Portrait:
		return "A2_Portrait"
	case SizeA3Landscape:
		return "A3_Landscape"
	case SizeA3Portrait:
		return "A3_Portrait"
	case SizeA4Landscape:
		return "A4_Landscape"
	case SizeA4Portrait:
		return "A4_Portrait"
	case SizeA5Landscape:
		return "A5_Landscape"
	case SizeA5Portrait:
		return "A5_Portrait"
	case SizeA6Landscape:
		return "A6_Landscape"
	case SizeA6Portrait:
		return "A6_Portrait"
	case SizeLegalLandscape:
		return "Legal_Landscape"
	case SizeLegalPortrait:
		return "Legal_Portrait"
	case SizeLetterLandscape:
		return "Letter_Landscape"
	case SizeLetterPortrait:
		return "Letter_Portrait"
	case SizeSlide16X10:
		return "Slide_16x10"
	case SizeSlide16X9:
		return "Slide_16x9"
	case SizeSlide4X3:
		return "Slide_4x3"
	default:
		return ""
	}
}

// Name of routing kind.
func (k RoutingKind) Name() string {
	switch k {
	case RoutingDirect:
		return "Direct"
	case RoutingOrthogonal:
		return "Orthogonal"
	case RoutingCurved:
		return "Curved"
	default:
		return ""
	}
}
