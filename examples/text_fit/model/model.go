package design

import (
	. "goa.design/model/dsl"
	"goa.design/model/expr"
)

var _ = Design("Text Fit", "Exercises content-aware element sizing.", func() {
	var auditAPI *expr.Container
	SoftwareSystem("Audit System", "Owns the external audit trail.", func() {
		auditAPI = Container(
			"Audit API",
			"Accepts immutable lifecycle records from trusted platform services.",
			"Go and Goa",
		)
	})

	system := SoftwareSystem("Text Fit System", "Owns the text-fit regression fixture.", func() {
		Container(
			"Long-Lived Workflow Coordination Service",
			"Owns execution records, validates compiled definitions, persists idempotent lifecycle state, and serves complete run history without shortening the architecture contract.",
			"Go, Goa, Temporal Worker, and OpenTelemetry",
			func() {
				Uses(auditAPI, "Publishes immutable lifecycle records to", "gRPC", Synchronous)
			},
		)
	})

	Views(func() {
		ContainerView(system, "Text Fit", "Text must remain inside its element boundary.", func() {
			AddAll()
			Add(auditAPI)
			SystemBoundariesVisible()
			AutoLayout(RankTopBottom)
		})
	})
})
