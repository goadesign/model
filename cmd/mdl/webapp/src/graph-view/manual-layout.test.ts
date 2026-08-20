import {describe, expect, it} from "vitest";

import {MeasuredDiagram, ResolvedLayout} from "./geometry";
import {resolveManual} from "./manual-layout";

describe("resolveManual", () => {
	it("keeps complete authored positions unchanged", () => {
		const measured = measuredDiagram();
		const snapshot = completeSnapshot();
		const validated = resolveManual(measured, snapshot);
		expect(validated).toBe(snapshot);
	});

	it("rejects a snapshot that omits a node", () => {
		const measured = measuredDiagram();
		const snapshot = completeSnapshot();
		snapshot.nodes = [];
		expect(() => resolveManual(measured, snapshot)).toThrow("node node is missing");
	});
});

/** measuredDiagram creates the smallest complete manual-layout input. */
function measuredDiagram(): MeasuredDiagram {
	return {
		id: "manual",
		nodes: [{
			id: "node",
			size: {width: 100, height: 60},
			contentSize: {width: 90, height: 50},
			centerOffset: {x: 50, y: 30},
			shape: "box",
		}],
		groups: [],
		edges: [],
	};
}

/** completeSnapshot places every measured element inside one view. */
function completeSnapshot(): ResolvedLayout {
	return {
		diagramId: "manual",
		bounds: {x: 0, y: 0, width: 200, height: 120},
		nodes: [{
			id: "node",
			bounds: {x: 50, y: 30, width: 100, height: 60},
		}],
		groups: [],
		edges: [],
	};
}
