// This file is the only owner of automatic positions. It gives ELK the exact
// measured sizes, reads ELK's complete result, and does not repair that result.

import type {
	ElkExtendedEdge,
	ElkLabel,
	ElkNode,
	ElkPoint,
} from "elkjs/lib/elk-api";

import {
	MeasuredDiagram,
	MeasuredEdge,
	MeasuredGroup,
	MeasuredNode,
	Point,
	ResolvedEdge,
	ResolvedGroup,
	ResolvedLayout,
	ResolvedNode,
} from "./geometry";

/** AutomaticLayoutOptions are the choices a caller may make before layout. */
export interface AutomaticLayoutOptions {
	direction: "UP" | "DOWN" | "LEFT" | "RIGHT";
	nodeSpacing: number;
	layerSpacing: number;
}

const defaultOptions: AutomaticLayoutOptions = {
	direction: "DOWN",
	nodeSpacing: 80,
	layerSpacing: 80,
};

/**
 * resolveAutomatic runs one compound ELK layout with exact measured sizes.
 * ELK errors and incomplete results are returned to the caller unchanged.
 */
export async function resolveAutomatic(
	diagram: MeasuredDiagram,
	options: Partial<AutomaticLayoutOptions> = {},
): Promise<ResolvedLayout> {
	const settings = {...defaultOptions, ...options};
	const graph = buildLayoutTree(diagram, settings);
	const ELK = await import("elkjs/lib/elk.bundled.js").then(module => module.default);
	const elk = new ELK();
	const result = await elk.layout(graph);
	return readResolvedLayout(diagram, result);
}

/**
 * buildLayoutTree creates one ELK tree. Nodes belong to their declared group,
 * groups belong to their declared parent, and each edge belongs to the nearest
 * group shared by both connected nodes.
 */
function buildLayoutTree(
	diagram: MeasuredDiagram,
	options: AutomaticLayoutOptions,
): ElkNode {
	const nodes = new Map(diagram.nodes.map(node => [node.id, node]));
	const groups = new Map(diagram.groups.map(group => [group.id, group]));
	checkMeasuredLinks(diagram, nodes, groups);

	const childrenByParent = new Map<string | undefined, MeasuredNode[]>();
	for (const node of diagram.nodes) {
		append(childrenByParent, node.parentId, node);
	}
	const groupsByParent = new Map<string | undefined, MeasuredGroup[]>();
	for (const group of diagram.groups) {
		append(groupsByParent, group.parentId, group);
	}

	const elkGroups = new Map<string, ElkNode>();
	const buildGroup = (group: MeasuredGroup): ElkNode => {
		const child: ElkNode = {
			id: group.id,
			children: [
				...(groupsByParent.get(group.id) ?? []).map(buildGroup),
				...(childrenByParent.get(group.id) ?? []).map(buildNode),
			],
			edges: [],
			labels: [{
				text: group.id,
				width: group.titleSize.width,
				height: group.titleSize.height,
			}],
			layoutOptions: groupLayoutOptions(options, group.titleSize.height),
		};
		elkGroups.set(group.id, child);
		return child;
	};

	const root: ElkNode = {
		id: diagram.id,
		children: [
			...(groupsByParent.get(undefined) ?? []).map(buildGroup),
			...(childrenByParent.get(undefined) ?? []).map(buildNode),
		],
		edges: [],
		layoutOptions: rootLayoutOptions(options),
	};

	for (const edge of diagram.edges) {
		const owner = lowestCommonGroup(edge, nodes, groups);
		const container = owner ? elkGroups.get(owner) : root;
		if (!container) {
			throw new Error(`edge ${edge.id} has missing layout group ${owner}`);
		}
		container.edges?.push(buildEdge(edge));
	}
	return root;
}

/** buildNode passes the measured node size to ELK without padding or guesses. */
function buildNode(node: MeasuredNode): ElkNode {
	return {
		id: node.id,
		width: node.size.width,
		height: node.size.height,
		layoutOptions: {
			"elk.nodeSize.constraints": "FIXED_SIZE",
			"elk.portConstraints": "FREE",
		},
	};
}

/** buildEdge passes the exact label size and stable endpoint IDs to ELK. */
function buildEdge(edge: MeasuredEdge): ElkExtendedEdge {
	return {
		id: edge.id,
		sources: [edge.sourceId],
		targets: [edge.targetId],
		labels: edge.labelSize ? [{
			text: edge.id,
			width: edge.labelSize.width,
			height: edge.labelSize.height,
			layoutOptions: {
				"elk.edgeLabels.inline": "false",
				"elk.edgeLabels.placement": "CENTER",
			},
		}] : [],
	};
}

/** rootLayoutOptions set the rules shared by the full compound layout. */
function rootLayoutOptions(options: AutomaticLayoutOptions): Record<string, string> {
	return {
		"elk.algorithm": "layered",
		"elk.direction": options.direction,
		"elk.edgeRouting": "ORTHOGONAL",
		"elk.hierarchyHandling": "INCLUDE_CHILDREN",
		"elk.padding": "[top=40,left=40,bottom=40,right=40]",
		"elk.spacing.nodeNode": String(options.nodeSpacing),
		"elk.spacing.componentComponent": String(options.nodeSpacing),
		"elk.layered.spacing.nodeNodeBetweenLayers": String(options.layerSpacing),
		"elk.layered.spacing.edgeNodeBetweenLayers": "30",
		"elk.layered.spacing.edgeEdgeBetweenLayers": "20",
		"elk.layered.nodePlacement.strategy": "NETWORK_SIMPLEX",
		"elk.layered.nodePlacement.favorStraightEdges": "true",
		"elk.layered.crossingMinimization.strategy": "LAYER_SWEEP",
		"elk.layered.unnecessaryBendpoints": "true",
		"elk.edgeLabels.inline": "false",
		"elk.edgeLabels.avoidOverlap": "true",
		"elk.layered.edgeLabels.sideSelection": "SMART_DOWN",
		"elk.spacing.edgeLabel": "12",
	};
}

/** groupLayoutOptions reserve title space and keep each group's children together. */
function groupLayoutOptions(
	options: AutomaticLayoutOptions,
	titleHeight: number,
): Record<string, string> {
	return {
		...rootLayoutOptions(options),
		"elk.padding": `[top=${titleHeight + 30},left=30,bottom=30,right=30]`,
		"elk.nodeLabels.placement": "[H_LEFT,V_TOP,INSIDE]",
		"elk.spacing.labelNode": "10",
	};
}

/**
 * checkMeasuredLinks rejects broken ownership and endpoints before calling ELK,
 * so ELK never has to guess what the model meant.
 */
function checkMeasuredLinks(
	diagram: MeasuredDiagram,
	nodes: Map<string, MeasuredNode>,
	groups: Map<string, MeasuredGroup>,
): void {
	for (const node of diagram.nodes) {
		if (node.parentId && !groups.has(node.parentId)) {
			throw new Error(`node ${node.id} has missing parent group ${node.parentId}`);
		}
	}
	for (const group of diagram.groups) {
		if (group.parentId && !groups.has(group.parentId)) {
			throw new Error(`group ${group.id} has missing parent group ${group.parentId}`);
		}
	}
	for (const edge of diagram.edges) {
		if (!nodes.has(edge.sourceId) || !nodes.has(edge.targetId)) {
			throw new Error(
				`edge ${edge.id} connects missing node ${edge.sourceId} or ${edge.targetId}`,
			);
		}
	}
}

/** lowestCommonGroup finds the nearest boundary that contains both endpoints. */
function lowestCommonGroup(
	edge: MeasuredEdge,
	nodes: Map<string, MeasuredNode>,
	groups: Map<string, MeasuredGroup>,
): string | undefined {
	const sourceGroups = parentGroups(nodes.get(edge.sourceId)?.parentId, groups);
	const targetGroups = new Set(parentGroups(nodes.get(edge.targetId)?.parentId, groups));
	return sourceGroups.find(group => targetGroups.has(group));
}

/** parentGroups returns parent IDs from nearest to farthest. */
function parentGroups(
	parentId: string | undefined,
	groups: Map<string, MeasuredGroup>,
): string[] {
	const result: string[] = [];
	let current = parentId;
	while (current) {
		result.push(current);
		current = groups.get(current)?.parentId;
	}
	return result;
}

/** readResolvedLayout changes ELK's parent-relative values into diagram positions. */
function readResolvedLayout(
	diagram: MeasuredDiagram,
	root: ElkNode,
): ResolvedLayout {
	if (root.width === undefined || root.height === undefined) {
		throw new Error(`ELK did not size diagram ${diagram.id}`);
	}
	const nodes: ResolvedNode[] = [];
	const groups: ResolvedGroup[] = [];
	const edges: ResolvedEdge[] = [];
	const measuredNodes = new Map(diagram.nodes.map(node => [node.id, node]));
	const measuredGroups = new Map(diagram.groups.map(group => [group.id, group]));

	readContainer(
		root,
		0,
		0,
		measuredNodes,
		measuredGroups,
		nodes,
		groups,
		edges,
	);
	return {
		diagramId: diagram.id,
		bounds: {x: 0, y: 0, width: root.width, height: root.height},
		nodes,
		groups,
		edges,
	};
}

/**
 * readContainer walks one ELK group, adding its parent's position exactly once
 * to each child, edge point, and label box.
 */
function readContainer(
	container: ElkNode,
	parentX: number,
	parentY: number,
	measuredNodes: Map<string, MeasuredNode>,
	measuredGroups: Map<string, MeasuredGroup>,
	nodes: ResolvedNode[],
	groups: ResolvedGroup[],
	edges: ResolvedEdge[],
): void {
	const containerX = parentX + (container.x ?? 0);
	const containerY = parentY + (container.y ?? 0);
	for (const child of container.children ?? []) {
		const childX = containerX + requiredNumber(child.x, child.id, "x");
		const childY = containerY + requiredNumber(child.y, child.id, "y");
		const width = requiredNumber(child.width, child.id, "width");
		const height = requiredNumber(child.height, child.id, "height");
		const measuredNode = measuredNodes.get(child.id);
		if (measuredNode) {
			nodes.push({
				id: child.id,
				parentId: measuredNode.parentId,
				bounds: {x: childX, y: childY, width, height},
			});
			continue;
		}
		const measuredGroup = measuredGroups.get(child.id);
		if (!measuredGroup) {
			throw new Error(`ELK returned unknown node ${child.id}`);
		}
		const label = requiredLabel(child.labels?.[0], child.id);
		groups.push({
			id: child.id,
			parentId: measuredGroup.parentId,
			bounds: {x: childX, y: childY, width, height},
			titleBounds: {
				x: childX + requiredNumber(label.x, child.id, "title x"),
				y: childY + requiredNumber(label.y, child.id, "title y"),
				width: requiredNumber(label.width, child.id, "title width"),
				height: requiredNumber(label.height, child.id, "title height"),
			},
		});
		readContainer(
			child,
			containerX,
			containerY,
			measuredNodes,
			measuredGroups,
			nodes,
			groups,
			edges,
		);
	}
	for (const edge of container.edges ?? []) {
		edges.push(readEdge(edge, containerX, containerY));
	}
}

/** readEdge copies every ELK section and label without choosing new positions. */
function readEdge(edge: ElkExtendedEdge, offsetX: number, offsetY: number): ResolvedEdge {
	const sections = (edge.sections ?? []).map(section => [
		addOffset(section.startPoint, offsetX, offsetY),
		...(section.bendPoints ?? []).map(point => addOffset(point, offsetX, offsetY)),
		addOffset(section.endPoint, offsetX, offsetY),
	]);
	const label = edge.labels?.[0];
	return {
		id: edge.id,
		sourceId: requiredEndpoint(edge.sources, edge.id, "source"),
		targetId: requiredEndpoint(edge.targets, edge.id, "target"),
		sections,
		labelBounds: label ? {
			x: offsetX + requiredNumber(label.x, edge.id, "label x"),
			y: offsetY + requiredNumber(label.y, edge.id, "label y"),
			width: requiredNumber(label.width, edge.id, "label width"),
			height: requiredNumber(label.height, edge.id, "label height"),
		} : undefined,
	};
}

/** addOffset changes one parent-relative ELK point into a diagram point. */
function addOffset(point: ElkPoint, offsetX: number, offsetY: number): Point {
	return {x: offsetX + point.x, y: offsetY + point.y};
}

/** requiredNumber rejects incomplete ELK output instead of filling in zero. */
function requiredNumber(
	value: number | undefined,
	elementId: string,
	field: string,
): number {
	if (value === undefined) {
		throw new Error(`ELK omitted ${field} for ${elementId}`);
	}
	return value;
}

/** requiredEndpoint requires each MDL relationship to have exactly one endpoint. */
function requiredEndpoint(
	endpoints: string[],
	edgeId: string,
	field: string,
): string {
	if (endpoints.length !== 1) {
		throw new Error(`ELK returned ${endpoints.length} ${field}s for edge ${edgeId}`);
	}
	return endpoints[0];
}

/** requiredLabel rejects a group whose measured title was not placed by ELK. */
function requiredLabel(label: ElkLabel | undefined, groupId: string): ElkLabel {
	if (!label) {
		throw new Error(`ELK omitted title for group ${groupId}`);
	}
	return label;
}

/** append adds a child to one parent entry while keeping model order. */
function append<T>(
	items: Map<string | undefined, T[]>,
	parentId: string | undefined,
	item: T,
): void {
	const siblings = items.get(parentId) ?? [];
	siblings.push(item);
	items.set(parentId, siblings);
}
