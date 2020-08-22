package design

// Design represents a C4 system architecture design complete with models and
// views. This is the top level struct.
type Design struct {
	// Name of design
	Name string `json:"name"`
	// Description for overall architecture
	Description string `json:"description"`
	// Version of design
	Version string `json:"version"`
	// Model defines the C4 software architecture elements that make up the
	// architecture.
	Model *Model `json:"model"`
	// Views defines the C4 views used to visualize the architecture.
	Views *Views `json:"views"`
}
