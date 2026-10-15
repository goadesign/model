import {afterEach, beforeEach, describe, expect, it, vi} from "vitest";

import {parseView} from "./parseModel";

vi.hoisted(() => {
	vi.stubGlobal("document", {
		querySelector: () => ({setAttribute: () => {}}),
	});
});

// Parser tests use fixed text sizes; browser tests verify the rendered geometry.
vi.mock("./graph-view/svg-text", () => ({
	measureSVGText: (text: string) => ({width: text.length * 8, height: 16}),
	svgTextWrap: (text: string) => ({lines: [text], maxW: text.length * 8}),
}));

beforeEach(() => {
	vi.stubGlobal("window", {});
	vi.useFakeTimers();
});

afterEach(() => {
	vi.clearAllTimers();
	vi.useRealTimers();
	vi.unstubAllGlobals();
});

describe("component view boundaries", () => {
	it.each([undefined, false, true])("keeps the owner boundary with external boundaries %s", external => {
		const model = componentModel();
		model.views.componentViews[0].externalContainersBoundariesVisible = external;
		const graph = parseView(model, {}, "components");
		expect([...graph.groupsMap.keys()]).toEqual(external ? ["owner", "neighbor"] : ["owner"]);
		expect(graph.groupsMap.get("owner").nodes.map(node => node.id)).toEqual(["handler", "store"]);
		expect(graph.groupsMap.get("neighbor")?.nodes.map(node => node.id)).toEqual(external ? ["worker"] : undefined);
		expect([...graph.nodesMap.keys()]).toEqual(["handler", "store", "worker", "database"]);
		expect(graph.edges.map(edge => [edge.id, edge.from.id, edge.to.id])).toEqual([
			["save", "handler", "store"],
			["submit", "handler", "worker"],
			["read", "worker", "database"],
		]);
	});

	it("keeps an explicitly selected container and its relationship as nodes", () => {
		const model = componentModel();
		const view = model.views.componentViews[0];
		view.externalContainersBoundariesVisible = true;
		view.elements.push({id: "neighbor"});
		view.relationships.push({id: "notify", vertices: [], routing: ""});
		const graph = parseView(model, {}, "components");
		expect([...graph.groupsMap.keys()]).toEqual(["owner"]);
		expect(graph.nodesMap.has("neighbor")).toBe(true);
		expect(graph.edges[graph.edges.length - 1]?.to.id).toBe("neighbor");
	});

	it("does not create an empty owner boundary when only external elements are selected", () => {
		const model = componentModel();
		model.views.componentViews[0].elements = [{id: "database"}];
		model.views.componentViews[0].relationships = [];
		const graph = parseView(model, {}, "components");
		expect(graph.groupsMap.size).toBe(0);
		expect([...graph.nodesMap.keys()]).toEqual(["database"]);
	});
});

// componentModel includes two components in one container, an external component,
// and a standalone container so grouping cannot mistake a system for a container.
function componentModel(): Parameters<typeof parseView>[0] {
	return {
		name: "Component ownership", description: "", version: "1",
		model: {
			enterprise: {name: "Example"}, people: [], deploymentNodes: [],
			softwareSystems: [{
				id: "system", name: "Product", tags: "Software System", containers: [
					{id: "owner", name: "Application", tags: "Container", components: [
						{id: "handler", name: "Handler", tags: "Component", relationships: [
							relation("save", "handler", "store"), relation("submit", "handler", "worker"),
							relation("notify", "handler", "neighbor"),
						]},
						{id: "store", name: "Store", tags: "Component"},
					]},
					{id: "neighbor", name: "Jobs", tags: "Container", components: [
						{id: "worker", name: "Worker", tags: "Component", relationships: [relation("read", "worker", "database")]},
					]},
					{id: "database", name: "Database", tags: "Container"},
				],
			}],
		},
		views: {
			systemLandscapeViews: [], containerViews: [], dynamicViews: [], deploymentViews: [],
			styles: {elements: [], relationships: []},
			componentViews: [{
				key: "components", title: "Components", description: "", containerId: "owner",
				elements: ["handler", "store", "worker", "database"].map(id => ({id})),
				relationships: ["save", "submit", "read"].map(id => ({id, vertices: [] as {x: number; y: number}[], routing: ""})),
			}],
		},
	};
}

function relation(id: string, sourceId: string, destinationId: string) {
	return {id, sourceId, destinationId, description: id, tags: "Relationship", technology: "", interactionStyle: "Synchronous"};
}
