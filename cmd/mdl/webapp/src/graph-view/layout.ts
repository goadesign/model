// This file connects the editable graph to the layout pipeline. It measures
// the current SVG once, asks ELK for one complete answer, validates that answer,
// and returns positions without changing them.

import {resolveAutomatic} from "./automatic-layout";
import {
	MeasuredDiagram,
	MeasuredEdge,
	MeasuredGroup,
	MeasuredNode,
	Rect,
	ResolvedLayout,
	ValidatedLayout,
	validateResolvedLayout,
} from "./geometry";
import {Edge, GraphData, Group, LayoutDirection, Node} from "./graph";
import {calculateEdgeVertices} from "./edge-utils";
import {resolveManual} from "./manual-layout";

/** LayoutOptions are the explicit choices available to editor and CLI callers. */
export interface LayoutOptions {
	direction?: LayoutDirection;
	nodeSpacing?: number;
	layerSpacing?: number;
	compactLayout?: boolean;
}

/** GraphLayout carries checked positions back to the editable graph. */
export interface GraphLayout {
	validated: ValidatedLayout;
	nodes: Array<{id: string; x: number; y: number}>;
	groups: Array<{id: string; bounds: Rect; titleBounds: Rect}>;
	edges: Array<{
		id: string;
		sections: Array<Array<{x: number; y: number}>>;
		labelBounds?: Rect;
	}>;
}

/**
 * autoLayout measures, lays out, and validates a graph. It does not catch
 * errors because callers must not save a diagram after any failed step.
 */
export async function autoLayout(
	graph: GraphData,
	options: LayoutOptions = {},
): Promise<GraphLayout> {
	await document.fonts.ready;
	const measured = measureGraph(graph);
	const spacing = options.compactLayout ? 56 : 80;
	const resolved = await resolveAutomatic(measured, {
		direction: options.direction ?? graph.layoutDirection ?? "DOWN",
		nodeSpacing: options.nodeSpacing ?? spacing,
		layerSpacing: options.layerSpacing ?? spacing,
	});
	const validated = validateResolvedLayout(measured, resolved);
	const measuredNodes = new Map(measured.nodes.map(node => [node.id, node]));
	return {
		validated,
		nodes: validated.nodes.map(node => {
			const measuredNode = measuredNodes.get(node.id);
			if (!measuredNode) {
				throw new Error(`validated layout contains unknown node ${node.id}`);
			}
			return {
				id: node.id,
				x: node.bounds.x + measuredNode.centerOffset.x,
				y: node.bounds.y + measuredNode.centerOffset.y,
			};
		}),
		groups: validated.groups.map(group => ({
			id: group.id,
			bounds: group.bounds,
			titleBounds: group.titleBounds,
		})),
		edges: validated.edges.map(edge => ({
			id: edge.id,
			sections: edge.sections,
			labelBounds: edge.labelBounds,
		})),
	};
}

/**
 * validateCurrentLayout adapts the fully drawn editor state into the manual
 * layout contract. It never reads automatic positions to fill missing data.
 */
export function validateCurrentLayout(graph: GraphData): ValidatedLayout {
	const measured = measureGraph(graph);
	const measuredNodes = new Map(measured.nodes.map(node => [node.id, node]));
	const nodes = graph.nodes().map(node => {
		const measuredNode = measuredNodes.get(node.id);
		if (!measuredNode) {
			throw new Error(`measured diagram omitted node ${node.id}`);
		}
		return {
			id: node.id,
			parentId: measuredNode.parentId,
			bounds: {
				x: node.x - measuredNode.centerOffset.x,
				y: node.y - measuredNode.centerOffset.y,
				...measuredNode.size,
			},
		};
	});
	const groups = Array.from(graph.groupsMap.values()).map(group => {
		if (!group.ref?.isConnected) {
			throw new Error(`group ${group.id} must be drawn before saving`);
		}
		const boundary = group.ref.querySelector("rect");
		const title = group.ref.querySelector("text");
		if (!(boundary instanceof SVGGraphicsElement) || !(title instanceof SVGGraphicsElement)) {
			throw new Error(`group ${group.id} is missing its boundary or title`);
		}
		return {
			id: group.id,
			parentId: measured.groups.find(candidate => candidate.id === group.id)?.parentId,
			bounds: copyBox(boundary.getBBox()),
			titleBounds: copyBox(title.getBBox()),
		};
	});
	const edges = graph.edges.map(edge => {
		const sections = edge.layoutSections ?? [calculateEdgeVertices(edge)];
		const labelBounds = edge.layoutLabelBounds ?? edge.labelBounds;
		if (edge.label.trim().length > 0 && !labelBounds) {
			throw new Error(`edge ${edge.id} is missing its label position`);
		}
		return {
			id: edge.id,
			sourceId: edge.from.id,
			targetId: edge.to.id,
			sections,
			labelBounds: labelBounds ? {...labelBounds} : undefined,
		};
	});
	const snapshot: ResolvedLayout = {
		diagramId: graph.id,
		bounds: contentBounds(nodes.map(node => node.bounds), groups, edges),
		nodes,
		groups,
		edges,
	};
	return resolveManual(measured, snapshot);
}

/**
 * measureGraph reads the exact node and label boxes already drawn by the
 * browser. Parent IDs come only from the model's declared groups.
 */
function measureGraph(graph: GraphData): MeasuredDiagram {
	const {nodeParents, groupParents} = declaredParents(graph);
	return {
		id: graph.id,
		nodes: graph.nodes().map(node => measureNode(node, nodeParents.get(node.id))),
		groups: Array.from(graph.groupsMap.values()).map(group =>
			measureGroup(group, groupParents.get(group.id)),
		),
		edges: graph.edges.map(measureEdge),
	};
}

/** measureNode returns the complete box drawn for one node and its center. */
function measureNode(node: Node, parentId: string | undefined): MeasuredNode {
	if (!node.ref?.isConnected) {
		throw new Error(`node ${node.id} must be drawn before layout`);
	}
	const bounds = node.ref.getBBox();
	checkMeasuredBox(node.id, bounds);
	return {
		id: node.id,
		parentId,
		size: {width: bounds.width, height: bounds.height},
		contentSize: {width: bounds.width, height: bounds.height},
		centerOffset: {x: -bounds.x, y: -bounds.y},
		shape: node.style.shape ?? "Box",
	};
}

/** measureGroup reads the title box already drawn by the browser. */
function measureGroup(group: Group, parentId: string | undefined): MeasuredGroup {
	if (!group.ref?.isConnected) {
		throw new Error(`group ${group.id} must be drawn before layout`);
	}
	const title = group.ref.querySelector("text");
	if (!(title instanceof SVGGraphicsElement)) {
		throw new Error(`group ${group.id} is missing its drawn title`);
	}
	const bounds = title.getBBox();
	checkMeasuredBox(group.id, bounds);
	return {
		id: group.id,
		parentId,
		titleSize: {width: bounds.width, height: bounds.height},
	};
}

/** measureEdge uses the label background that the current SVG already drew. */
function measureEdge(edge: Edge): MeasuredEdge {
	const hasLabel = edge.label.trim().length > 0;
	if (hasLabel && !edge.labelBounds) {
		throw new Error(`edge ${edge.id} must be drawn before layout`);
	}
	return {
		id: edge.id,
		sourceId: edge.from.id,
		targetId: edge.to.id,
		labelSize: hasLabel ? {
			width: edge.labelBounds?.width ?? 0,
			height: edge.labelBounds?.height ?? 0,
		} : undefined,
	};
}

/**
 * declaredParents records model ownership and rejects elements that appear in
 * more than one group. Visual position never decides ownership.
 */
function declaredParents(graph: GraphData): {
	nodeParents: Map<string, string>;
	groupParents: Map<string, string>;
} {
	const nodeParents = new Map<string, string>();
	const groupParents = new Map<string, string>();
	for (const group of graph.groupsMap.values()) {
		for (const member of group.nodes) {
			const parents = isGroup(member) ? groupParents : nodeParents;
			const existing = parents.get(member.id);
			if (existing && existing !== group.id) {
				throw new Error(
					`${isGroup(member) ? "group" : "node"} ${member.id} belongs to both ` +
					`${existing} and ${group.id}`,
				);
			}
			parents.set(member.id, group.id);
		}
	}
	return {nodeParents, groupParents};
}

/** checkMeasuredBox rejects missing or invalid browser measurements. */
function checkMeasuredBox(id: string, box: Rect): void {
	if (
		!Number.isFinite(box.x) ||
		!Number.isFinite(box.y) ||
		!Number.isFinite(box.width) ||
		!Number.isFinite(box.height) ||
		box.width <= 0 ||
		box.height <= 0
	) {
		throw new Error(`browser returned invalid measured box for ${id}`);
	}
}

/** contentBounds returns the smallest box that contains all drawn geometry. */
function contentBounds(
	nodes: Rect[],
	groups: Array<{bounds: Rect; titleBounds: Rect}>,
	edges: Array<{sections: Array<Array<{x: number; y: number}>>; labelBounds?: Rect}>,
): Rect {
	const boxes = [
		...nodes,
		...groups.flatMap(group => [group.bounds, group.titleBounds]),
		...edges.flatMap(edge => edge.labelBounds ? [edge.labelBounds] : []),
	];
	const points = edges.flatMap(edge => edge.sections.flat());
	if (boxes.length === 0 && points.length === 0) {
		return {x: 0, y: 0, width: 0, height: 0};
	}
	const left = Math.min(
		...boxes.map(box => box.x),
		...points.map(point => point.x),
	);
	const top = Math.min(
		...boxes.map(box => box.y),
		...points.map(point => point.y),
	);
	const right = Math.max(
		...boxes.map(box => box.x + box.width),
		...points.map(point => point.x),
	);
	const bottom = Math.max(
		...boxes.map(box => box.y + box.height),
		...points.map(point => point.y),
	);
	return {x: left, y: top, width: right - left, height: bottom - top};
}

/** copyBox copies the browser's read-only SVG box into the layout contract. */
function copyBox(box: DOMRect): Rect {
	return {x: box.x, y: box.y, width: box.width, height: box.height};
}

/** isGroup tells group boundaries from nodes without using display names. */
function isGroup(member: Node | Group): member is Group {
	return "nodes" in member;
}
