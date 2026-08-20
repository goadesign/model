import {describe, expect, it} from "vitest";

import {
	GeometryIssueCode,
	GeometryValidationError,
	MeasuredDiagram,
	ResolvedLayout,
	validateResolvedLayout,
} from "./geometry";

describe("validateResolvedLayout", () => {
	it("accepts complete finite geometry", () => {
		expect(() => validateResolvedLayout(measuredDiagram(), resolvedLayout())).not.toThrow();
	});

	it.each([
		{
			name: "missing service",
			code: "missing_element",
			change: (layout: ResolvedLayout) => {
				layout.nodes.pop();
			},
		},
		{
			name: "node title overflow",
			code: "content_overflow",
			change: (_layout: ResolvedLayout, measured: MeasuredDiagram) => {
				measured.nodes[0].contentSize.width = 100;
			},
		},
		{
			name: "overlapping services",
			code: "node_overlap",
			change: (layout: ResolvedLayout) => {
				layout.nodes[1].bounds.x = 30;
			},
		},
		{
			name: "false system containment",
			code: "false_containment",
			change: (layout: ResolvedLayout) => {
				layout.nodes[0].bounds.x = 290;
			},
		},
		{
			name: "Flows relationship through service",
			code: "edge_through_node",
			change: (layout: ResolvedLayout, measured: MeasuredDiagram) => {
				measured.nodes.push({
					id: "runs-service",
					parentId: "platform",
					size: {width: 40, height: 30},
					contentSize: {width: 30, height: 20},
					centerOffset: {x: 20, y: 15},
					shape: "box",
				});
				layout.nodes.push({
					id: "runs-service",
					parentId: "platform",
					bounds: {x: 110, y: 60, width: 40, height: 30},
				});
			},
		},
		{
			name: "AURA relationship label over service",
			code: "label_node_overlap",
			change: (layout: ResolvedLayout) => {
				layout.edges[0].labelBounds = {x: 195, y: 65, width: 50, height: 20};
			},
		},
		{
			name: "relationship label over boundary title",
			code: "label_title_overlap",
			change: (layout: ResolvedLayout) => {
				layout.edges[0].labelBounds = {x: 20, y: 10, width: 50, height: 20};
			},
		},
		{
			name: "route outside view",
			code: "out_of_bounds",
			change: (layout: ResolvedLayout) => {
				layout.edges[0].sections[0][1].x = 400;
			},
		},
		{
			name: "non-finite position",
			code: "non_finite",
			change: (layout: ResolvedLayout) => {
				layout.nodes[0].bounds.x = Number.NaN;
			},
		},
		{
			name: "disconnected route",
			code: "disconnected_edge",
			change: (layout: ResolvedLayout) => {
				layout.edges[0].sections = [];
			},
		},
	] satisfies Array<{
		name: string;
		code: GeometryIssueCode;
		change: (layout: ResolvedLayout, measured: MeasuredDiagram) => void;
	}>)("rejects $name", ({code, change}) => {
		const measured = measuredDiagram();
		const resolved = resolvedLayout();
		change(resolved, measured);
		expectIssue(measured, resolved, code);
	});

	it("rejects Production Workspaces relationship-label overlap", () => {
		const measured = measuredDiagram();
		const resolved = resolvedLayout();
		measured.edges.push({
			id: "control",
			sourceId: "source",
			targetId: "target",
			labelSize: {width: 40, height: 20},
		});
		resolved.edges.push({
			id: "control",
			sourceId: "source",
			targetId: "target",
			sections: [[{x: 50, y: 105}, {x: 200, y: 105}]],
			labelBounds: {x: 130, y: 90, width: 40, height: 20},
		});
		expectIssue(measured, resolved, "label_label_overlap");
	});

	it("rejects a route through an unrelated system boundary", () => {
		const measured = measuredDiagram();
		const resolved = resolvedLayout();
		measured.groups.push({
			id: "unrelated",
			titleSize: {width: 30, height: 10},
		});
		resolved.groups.push({
			id: "unrelated",
			bounds: {x: 100, y: 50, width: 60, height: 100},
			titleBounds: {x: 105, y: 55, width: 30, height: 10},
		});
		expectIssue(measured, resolved, "edge_through_group");
	});

	it("rejects partially overlapping system boundaries", () => {
		const measured = measuredDiagram();
		const resolved = resolvedLayout();
		resolved.bounds = {x: 0, y: 0, width: 400, height: 300};
		measured.groups.push({
			id: "partial",
			titleSize: {width: 30, height: 10},
		});
		resolved.groups.push({
			id: "partial",
			bounds: {x: 250, y: 150, width: 100, height: 100},
			titleBounds: {x: 255, y: 155, width: 30, height: 10},
		});
		expectIssue(measured, resolved, "group_overlap");
	});

	it("rejects duplicate resolved IDs", () => {
		const measured = measuredDiagram();
		const resolved = resolvedLayout();
		resolved.nodes.push({...resolved.nodes[0]});
		expectIssue(measured, resolved, "duplicate_id");
	});
});

function expectIssue(
	measured: MeasuredDiagram,
	resolved: ResolvedLayout,
	code: GeometryIssueCode,
): void {
	try {
		validateResolvedLayout(measured, resolved);
		expect.fail(`expected ${code}`);
	} catch (error) {
		expect(error).toBeInstanceOf(GeometryValidationError);
		const validationError = error as GeometryValidationError;
		expect(validationError.issues.map(issue => issue.code)).toContain(code);
	}
}

function measuredDiagram(): MeasuredDiagram {
	return {
		id: "regression",
		nodes: [
			{
				id: "source",
				parentId: "platform",
				size: {width: 40, height: 30},
				contentSize: {width: 30, height: 20},
				centerOffset: {x: 20, y: 15},
				shape: "box",
			},
			{
				id: "target",
				parentId: "platform",
				size: {width: 40, height: 30},
				contentSize: {width: 30, height: 20},
				centerOffset: {x: 20, y: 15},
				shape: "box",
			},
		],
		groups: [{
			id: "platform",
			titleSize: {width: 80, height: 20},
		}],
		edges: [{
			id: "flow",
			sourceId: "source",
			targetId: "target",
			labelSize: {width: 40, height: 20},
		}],
	};
}

function resolvedLayout(): ResolvedLayout {
	return {
		diagramId: "regression",
		bounds: {x: 0, y: 0, width: 300, height: 200},
		nodes: [
			{
				id: "source",
				parentId: "platform",
				bounds: {x: 10, y: 60, width: 40, height: 30},
			},
			{
				id: "target",
				parentId: "platform",
				bounds: {x: 200, y: 60, width: 40, height: 30},
			},
		],
		groups: [{
			id: "platform",
			bounds: {x: 0, y: 0, width: 300, height: 200},
			titleBounds: {x: 10, y: 10, width: 80, height: 20},
		}],
		edges: [{
			id: "flow",
			sourceId: "source",
			targetId: "target",
			sections: [[{x: 50, y: 75}, {x: 200, y: 75}]],
			labelBounds: {x: 120, y: 90, width: 40, height: 20},
		}],
	};
}
