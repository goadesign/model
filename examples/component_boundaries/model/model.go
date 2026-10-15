// Package model defines component views with a shared owner and an optional
// external container boundary so editor and exported diagrams can be compared.
package model

import (
	. "goa.design/model/dsl"
	"goa.design/model/expr"
)

var _ = Design("Component Boundaries", "Component ownership and external dependencies", func() {
	var handler, store, worker *expr.Component
	var database, application, jobs, client *expr.Container
	SoftwareSystem("Job Processing", "Runs background jobs", func() {
		jobs = Container("Jobs", "Accepts submitted work", "Go", func() {
			worker = Component("Worker", "Processes accepted jobs", "Go")
		})
	})
	SoftwareSystem("Product", "Handles user requests", func() {
		database = Container("Database", "Stores completed work", "SQL")
		client = Container("Web Client", "Submits user requests", "JavaScript")
		application = Container("Application", "Handles requests and stores work", "Go", func() {
			handler = Component("Handler", "Accepts requests", "Go", func() {
				Uses(worker, "Submits jobs")
				Uses(jobs, "Checks queue")
			})
			store = Component("Store", "Saves completed work", "Go", func() {
				Uses(database, "Writes results")
			})
		})
	})
	Views(func() {
		ComponentView(application, "Owner Boundary", func() {
			Add(handler)
			Add(store)
			Add(worker)
			Add(database)
			Add(client)
			AutoLayout(RankLeftRight)
		})
		ComponentView(application, "External Boundaries", func() {
			Add(handler)
			Add(store)
			Add(worker)
			Add(database)
			Add(client)
			ContainerBoundariesVisible()
			AutoLayout(RankLeftRight)
		})
		ComponentView(application, "Container Node", func() {
			Add(handler)
			Add(store)
			Add(worker)
			Add(database)
			Add(jobs)
			ContainerBoundariesVisible()
			AutoLayout(RankLeftRight)
		})
	})
})
