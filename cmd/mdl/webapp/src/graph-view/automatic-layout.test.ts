import {describe, expect, it} from "vitest";

import {resolveAutomatic} from "./automatic-layout";
import {MeasuredDiagram, validateResolvedLayout} from "./geometry";

describe("resolveAutomatic", () => {
	it("uses exact sizes and returns the same compound layout each time", async () => {
		const diagram = compoundDiagram();
		const layouts = await Promise.all([
			resolveAutomatic(diagram),
			resolveAutomatic(diagram),
			resolveAutomatic(diagram),
		]);
		for (const layout of layouts) {
			expect(() => validateResolvedLayout(diagram, layout)).not.toThrow();
			for (const node of layout.nodes) {
				const measured = diagram.nodes.find(candidate => candidate.id === node.id);
				expect(node.bounds.width).toBe(measured?.size.width);
				expect(node.bounds.height).toBe(measured?.size.height);
			}
			for (const edge of layout.edges) {
				const measured = diagram.edges.find(candidate => candidate.id === edge.id);
				expect(edge.labelBounds?.width).toBe(measured?.labelSize?.width);
				expect(edge.labelBounds?.height).toBe(measured?.labelSize?.height);
			}
		}
		expect(JSON.stringify(layouts[1])).toBe(JSON.stringify(layouts[0]));
		expect(JSON.stringify(layouts[2])).toBe(JSON.stringify(layouts[0]));
	});

	it("fails when a node names a missing parent group", async () => {
		const diagram = compoundDiagram();
		diagram.nodes[0].parentId = "missing";
		await expect(resolveAutomatic(diagram)).rejects.toThrow(
			"node source has missing parent group missing",
		);
	});

	it("lays out several labeled relationships between the same groups", async () => {
		const diagram = compoundDiagram();
		diagram.edges.push(
			edge("evidence", "reader", "trigger", 95),
			edge("status", "reader", "trigger", 65),
		);
		const layout = await resolveAutomatic(diagram);
		expect(() => validateResolvedLayout(diagram, layout)).not.toThrow();
	});
});

/** compoundDiagram covers groups, cross-group edges, labels, and exact node sizes. */
function compoundDiagram(): MeasuredDiagram {
	return {
		id: "compound",
		nodes: [
			node("source", "edge", 180, 100),
			node("reader", "edge", 220, 120),
			node("trigger", "cloud", 190, 110),
			node("action", "cloud", 200, 100),
		],
		groups: [
			{id: "edge", titleSize: {width: 100, height: 24}},
			{id: "cloud", titleSize: {width: 110, height: 24}},
		],
		edges: [
			edge("readings", "source", "reader", 80),
			edge("events", "reader", "trigger", 70),
			edge("actions", "trigger", "action", 75),
		],
	};
}

/** node creates one measured fixture node with content smaller than its box. */
function node(
	id: string,
	parentId: string,
	width: number,
	height: number,
) {
	return {
		id,
		parentId,
		size: {width, height},
		contentSize: {width: width - 20, height: height - 20},
		centerOffset: {x: width / 2, y: height / 2},
		shape: "box",
	};
}

/** edge creates one measured fixture edge with a fixed-height label. */
function edge(
	id: string,
	sourceId: string,
	targetId: string,
	labelWidth: number,
) {
	return {
		id,
		sourceId,
		targetId,
		labelSize: {width: labelWidth, height: 24},
	};
}
