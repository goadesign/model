// Package model defines a dense service view used to verify that automatic
// relationship labels do not obscure nodes or one another.
package model

import (
	. "goa.design/model/dsl"
	"goa.design/model/expr"
)

var _ = Design("Label Collision", "Exercises collision-aware relationship labels.", func() {
	var pulse *expr.SoftwareSystem
	pulse = SoftwareSystem("Event Bus", "Platform event transport.", func() {
		Tag("External")
	})

	var provider *expr.Container
	SoftwareSystem("Automation", "Agentic automation platform.", func() {
		provider = Container("Workflow Provider", "Bridges automation routines to workflows.", "Go and Goa", func() {
			Uses("Workflows/Definitions Service", "Manages routine-backed workflow definitions", "gRPC", Synchronous)
			Uses("Workflows/Runs Service", "Starts runs and reads run status", "gRPC", Synchronous)
		})
	})

	workflows := SoftwareSystem("Workflows", "Turns platform signals into typed, durable workflows.", func() {
		Container(
			"Signals Service",
			"Stores signal contracts, accepts events, publishes them to the event bus, and serves subscriptions.",
			"Go and Goa",
			func() {
				Uses(pulse, "Publishes accepted signal events to", "gRPC", Asynchronous)
			},
		)
		Container(
			"Definitions Service",
			"Owns workflow templates, provider contracts, compilation, persistence, and runtime projections.",
			"Go and Goa",
		)
		Container(
			"Designer Agent",
			"Builds workflow drafts from deployment context and the signal catalog.",
			"Go and Goa",
			func() {
				Uses("Workflows/Definitions Service", "Compiles draft workflows through", "gRPC", Synchronous)
				Uses("Workflows/Inference Engine", "Generates workflow drafts using", "gRPC", Synchronous)
			},
		)
		Container(
			"Inference Engine",
			"Provides model-independent inference for workflow drafting.",
			"Go and Goa",
		)
		Container(
			"Runs Service",
			"Creates runs, validates compiled definitions, persists lifecycle state, and serves run history.",
			"Go and Goa",
			func() {
				Uses(
					"Workflows/Definitions Service",
					"Validates runs against compiled definitions and consumes activations from",
					"gRPC",
					Synchronous,
				)
				Uses(
					"Workflows/Signals Service",
					"Subscribes to signal events that trigger runs",
					"gRPC",
					Asynchronous,
				)
			},
		)
	})

	Views(func() {
		ContainerView(workflows, "Label Collision", "Labels must remain clear of nodes and peer labels.", func() {
			AddAll()
			Add(provider)
			Add(pulse)
			SystemBoundariesVisible()
			AutoLayout(RankTopBottom)
		})
	})
})
