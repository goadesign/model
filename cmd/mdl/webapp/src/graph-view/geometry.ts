// This file defines the diagram sizes and positions shared by layout and SVG
// drawing. Automatic and manual layout both produce ResolvedLayout. The SVG
// drawing code accepts a layout only after this file has checked it.

/** Point is one position in the diagram. */
export interface Point {
	x: number;
	y: number;
}

/** Size is a width and height. Both values must be zero or greater. */
export interface Size {
	width: number;
	height: number;
}

/** Rect is a box that does not rotate. */
export interface Rect extends Point, Size {}

/** MeasuredNode holds the exact box and text sizes needed before layout. */
export interface MeasuredNode {
	id: string;
	parentId?: string;
	size: Size;
	contentSize: Size;
	centerOffset: Point;
	shape: string;
}

/** MeasuredGroup describes a boundary that must contain its children and title. */
export interface MeasuredGroup {
	id: string;
	parentId?: string;
	titleSize: Size;
}

/** MeasuredEdge names the two connected nodes and the exact label size. */
export interface MeasuredEdge {
	id: string;
	sourceId: string;
	targetId: string;
	labelSize?: Size;
}

/**
 * MeasuredDiagram is the complete input to automatic or manual layout.
 * It contains sizes, but no positions.
 */
export interface MeasuredDiagram {
	id: string;
	nodes: MeasuredNode[];
	groups: MeasuredGroup[];
	edges: MeasuredEdge[];
}

/** ResolvedNode is one measured node placed in the diagram. */
export interface ResolvedNode {
	id: string;
	parentId?: string;
	bounds: Rect;
}

/** ResolvedGroup is one placed boundary and its title box. */
export interface ResolvedGroup {
	id: string;
	parentId?: string;
	bounds: Rect;
	titleBounds: Rect;
}

/**
 * ResolvedEdge contains a complete route. Each section is an ordered list of
 * points. labelBounds is the final label box chosen by layout.
 */
export interface ResolvedEdge {
	id: string;
	sourceId: string;
	targetId: string;
	sections: Point[][];
	labelBounds?: Rect;
}

/**
 * ResolvedLayout is the unchecked result of automatic or manual layout.
 * Drawing code must not use it directly.
 */
export interface ResolvedLayout {
	diagramId: string;
	bounds: Rect;
	nodes: ResolvedNode[];
	groups: ResolvedGroup[];
	edges: ResolvedEdge[];
}

/** GeometryIssueCode names one rule that the layout broke. */
export type GeometryIssueCode =
	| "duplicate_id"
	| "missing_element"
	| "unexpected_element"
	| "non_finite"
	| "size_mismatch"
	| "content_overflow"
	| "node_overlap"
	| "group_overlap"
	| "false_containment"
	| "edge_through_node"
	| "edge_through_group"
	| "label_node_overlap"
	| "label_label_overlap"
	| "label_edge_overlap"
	| "label_title_overlap"
	| "out_of_bounds"
	| "disconnected_edge";

/** GeometryIssue reports the stable IDs of the affected diagram elements. */
export interface GeometryIssue {
	code: GeometryIssueCode;
	elementIds: string[];
	message: string;
}

const validatedLayoutBrand: unique symbol = Symbol("ValidatedLayout");

/**
 * ValidatedLayout is the only layout the SVG drawing code accepts. The private
 * marker prevents other code from creating one without running the checks.
 */
export type ValidatedLayout = ResolvedLayout & {
	readonly [validatedLayoutBrand]: true;
};

/** GeometryValidationError reports every layout problem found in one check. */
export class GeometryValidationError extends Error {
	readonly issues: GeometryIssue[];

	constructor(issues: GeometryIssue[]) {
		super(issues.map(issue => issue.message).join("\n"));
		this.name = "GeometryValidationError";
		this.issues = issues;
	}
}

/**
 * validateResolvedLayout checks that every element is present, every number is
 * valid, boundaries show true ownership, elements do not collide, and all
 * content stays inside the diagram.
 */
export function validateResolvedLayout(
	measured: MeasuredDiagram,
	resolved: ResolvedLayout,
): ValidatedLayout {
	const issues: GeometryIssue[] = [];
	if (measured.id !== resolved.diagramId) {
		issues.push(issue(
			"missing_element",
			[measured.id, resolved.diagramId],
			`layout ${resolved.diagramId} does not resolve measured diagram ${measured.id}`,
		));
	}

	checkUniqueIDs("measured node", measured.nodes.map(node => node.id), issues);
	checkUniqueIDs("measured group", measured.groups.map(group => group.id), issues);
	checkUniqueIDs("measured edge", measured.edges.map(edge => edge.id), issues);
	checkUniqueIDs("resolved node", resolved.nodes.map(node => node.id), issues);
	checkUniqueIDs("resolved group", resolved.groups.map(group => group.id), issues);
	checkUniqueIDs("resolved edge", resolved.edges.map(edge => edge.id), issues);

	checkCoverage(
		"node",
		measured.nodes.map(node => node.id),
		resolved.nodes.map(node => node.id),
		issues,
	);
	checkCoverage(
		"group",
		measured.groups.map(group => group.id),
		resolved.groups.map(group => group.id),
		issues,
	);
	checkCoverage(
		"edge",
		measured.edges.map(edge => edge.id),
		resolved.edges.map(edge => edge.id),
		issues,
	);

	const measuredNodes = new Map(measured.nodes.map(node => [node.id, node]));
	const measuredGroups = new Map(measured.groups.map(group => [group.id, group]));
	const resolvedNodes = new Map(resolved.nodes.map(node => [node.id, node]));
	const resolvedGroups = new Map(resolved.groups.map(group => [group.id, group]));

	checkFiniteRect("diagram", resolved.diagramId, resolved.bounds, issues);
	for (const node of resolved.nodes) {
		checkFiniteRect("node", node.id, node.bounds, issues);
		const expected = measuredNodes.get(node.id);
		if (expected && (
			!nearlyEqual(expected.size.width, node.bounds.width) ||
			!nearlyEqual(expected.size.height, node.bounds.height)
		)) {
			issues.push(issue(
				"size_mismatch",
				[node.id],
				`node ${node.id} size changed after measurement`,
			));
		}
		if (expected && (
			expected.contentSize.width > node.bounds.width ||
			expected.contentSize.height > node.bounds.height
		)) {
			issues.push(issue(
				"content_overflow",
				[node.id],
				`node ${node.id} content exceeds its resolved bounds`,
			));
		}
		checkRectInBounds("node", node.id, node.bounds, resolved.bounds, issues);
	}
	for (const group of resolved.groups) {
		checkFiniteRect("group", group.id, group.bounds, issues);
		checkFiniteRect("group title", group.id, group.titleBounds, issues);
		checkRectInBounds("group", group.id, group.bounds, resolved.bounds, issues);
		if (!contains(group.bounds, group.titleBounds)) {
			issues.push(issue(
				"false_containment",
				[group.id],
				`group ${group.id} does not contain its title`,
			));
		}
		const expected = measuredGroups.get(group.id);
		if (expected && (
			!nearlyEqual(expected.titleSize.width, group.titleBounds.width) ||
			!nearlyEqual(expected.titleSize.height, group.titleBounds.height)
		)) {
			issues.push(issue(
				"size_mismatch",
				[group.id],
				`group ${group.id} title size changed after measurement`,
			));
		}
	}

	checkNodeOverlaps(resolved.nodes, issues);
	checkGroupOverlaps(resolved.groups, issues);
	checkOwnership(
		measured.nodes,
		measured.groups,
		resolvedNodes,
		resolvedGroups,
		issues,
	);
	checkEdges(
		measured,
		resolved,
		resolvedNodes,
		resolvedGroups,
		measuredNodes,
		measuredGroups,
		issues,
	);

	if (issues.length > 0) {
		throw new GeometryValidationError(issues);
	}
	return resolved as ValidatedLayout;
}

/** checkUniqueIDs requires each model ID to appear only once. */
function checkUniqueIDs(label: string, ids: string[], issues: GeometryIssue[]): void {
	const seen = new Set<string>();
	for (const id of ids) {
		if (seen.has(id)) {
			issues.push(issue("duplicate_id", [id], `${label} ${id} is duplicated`));
		}
		seen.add(id);
	}
}

/** checkCoverage rejects missing elements and old elements left in a layout. */
function checkCoverage(
	label: string,
	expected: string[],
	actual: string[],
	issues: GeometryIssue[],
): void {
	const expectedIDs = new Set(expected);
	const actualIDs = new Set(actual);
	for (const id of expectedIDs) {
		if (!actualIDs.has(id)) {
			issues.push(issue("missing_element", [id], `${label} ${id} is missing from resolved layout`));
		}
	}
	for (const id of actualIDs) {
		if (!expectedIDs.has(id)) {
			issues.push(issue("unexpected_element", [id], `${label} ${id} is not in the measured diagram`));
		}
	}
}

/** checkFiniteRect rejects invalid numbers and negative sizes. */
function checkFiniteRect(
	label: string,
	id: string,
	rect: Rect,
	issues: GeometryIssue[],
): void {
	if (!isFiniteRect(rect) || rect.width < 0 || rect.height < 0) {
		issues.push(issue("non_finite", [id], `${label} ${id} has invalid geometry`));
	}
}

/** checkRectInBounds requires the whole box to fit inside the diagram. */
function checkRectInBounds(
	label: string,
	id: string,
	rect: Rect,
	bounds: Rect,
	issues: GeometryIssue[],
): void {
	if (!contains(bounds, rect)) {
		issues.push(issue("out_of_bounds", [id], `${label} ${id} lies outside the resolved view bounds`));
	}
}

/** checkNodeOverlaps prevents two nodes from sharing the same area. */
function checkNodeOverlaps(nodes: ResolvedNode[], issues: GeometryIssue[]): void {
	forEachPair(nodes, (first, second) => {
		if (overlaps(first.bounds, second.bounds)) {
			issues.push(issue(
				"node_overlap",
				[first.id, second.id],
				`nodes ${first.id} and ${second.id} overlap`,
			));
		}
	});
}

/**
 * checkGroupOverlaps allows one real boundary inside another, but rejects
 * boundaries that partly cover each other and suggest false ownership.
 */
function checkGroupOverlaps(groups: ResolvedGroup[], issues: GeometryIssue[]): void {
	forEachPair(groups, (first, second) => {
		if (
			overlaps(first.bounds, second.bounds) &&
			!contains(first.bounds, second.bounds) &&
			!contains(second.bounds, first.bounds)
		) {
			issues.push(issue(
				"group_overlap",
				[first.id, second.id],
				`groups ${first.id} and ${second.id} partially overlap`,
			));
		}
	});
}

/**
 * checkOwnership requires each child to be inside the parent named by the
 * model. Visual overlap never creates ownership.
 */
function checkOwnership(
	nodes: MeasuredNode[],
	groups: MeasuredGroup[],
	resolvedNodes: Map<string, ResolvedNode>,
	resolvedGroups: Map<string, ResolvedGroup>,
	issues: GeometryIssue[],
): void {
	for (const node of nodes) {
		const child = resolvedNodes.get(node.id);
		if (child && child.parentId !== node.parentId) {
			issues.push(issue(
				"false_containment",
				[node.id],
				`node ${node.id} has the wrong parent ID`,
			));
		}
		if (!node.parentId) {
			continue;
		}
		const parent = resolvedGroups.get(node.parentId);
		if (child && parent && !contains(parent.bounds, child.bounds)) {
			issues.push(issue(
				"false_containment",
				[node.id, node.parentId],
				`node ${node.id} lies outside owning group ${node.parentId}`,
			));
		}
	}
	for (const group of groups) {
		const child = resolvedGroups.get(group.id);
		if (child && child.parentId !== group.parentId) {
			issues.push(issue(
				"false_containment",
				[group.id],
				`group ${group.id} has the wrong parent ID`,
			));
		}
		if (!group.parentId) {
			continue;
		}
		const parent = resolvedGroups.get(group.parentId);
		if (child && parent && !contains(parent.bounds, child.bounds)) {
			issues.push(issue(
				"false_containment",
				[group.id, group.parentId],
				`group ${group.id} lies outside owning group ${group.parentId}`,
			));
		}
	}
}

/**
 * checkEdges checks connected node IDs, complete routes, crossings, labels,
 * and boundary crossings for every relationship.
 */
function checkEdges(
	measured: MeasuredDiagram,
	resolved: ResolvedLayout,
	nodes: Map<string, ResolvedNode>,
	groups: Map<string, ResolvedGroup>,
	measuredNodes: Map<string, MeasuredNode>,
	measuredGroups: Map<string, MeasuredGroup>,
	issues: GeometryIssue[],
): void {
	const measuredEdges = new Map(measured.edges.map(edge => [edge.id, edge]));
	const labels: Array<{id: string; bounds: Rect}> = [];
	for (const edge of resolved.edges) {
		const expected = measuredEdges.get(edge.id);
		if (!expected) {
			continue;
		}
		if (
			edge.sourceId !== expected.sourceId ||
			edge.targetId !== expected.targetId
		) {
			issues.push(issue(
				"disconnected_edge",
				[edge.id, edge.sourceId, edge.targetId],
				`edge ${edge.id} endpoints do not match its measured contract`,
			));
		}
		if (edge.sections.length === 0 || edge.sections.some(section => section.length < 2)) {
			issues.push(issue(
				"disconnected_edge",
				[edge.id],
				`edge ${edge.id} has no complete route section`,
			));
			continue;
		}
		checkSectionConnections(edge, nodes, issues);
		for (const section of edge.sections) {
			for (const point of section) {
				if (!isFinitePoint(point) || !pointInRect(point, resolved.bounds, true)) {
					issues.push(issue(
						"out_of_bounds",
						[edge.id],
						`edge ${edge.id} has a route point outside the resolved view`,
					));
					break;
				}
			}
		}
		for (const node of nodes.values()) {
			if (node.id === edge.sourceId || node.id === edge.targetId) {
				continue;
			}
			if (routeIntersectsRect(edge.sections, inset(node.bounds, 0.5))) {
				issues.push(issue(
					"edge_through_node",
					[edge.id, node.id],
					`edge ${edge.id} crosses unrelated node ${node.id}`,
				));
			}
		}
		for (const group of groups.values()) {
			const sourceInside = belongsToGroup(
				edge.sourceId,
				group.id,
				measuredNodes,
				measuredGroups,
			);
			const targetInside = belongsToGroup(
				edge.targetId,
				group.id,
				measuredNodes,
				measuredGroups,
			);
			if (
				sourceInside === targetInside &&
				routeViolatesGroup(edge.sections, group.bounds, sourceInside)
			) {
				issues.push(issue(
					"edge_through_group",
					[edge.id, group.id],
					`edge ${edge.id} violates group ${group.id} boundary`,
				));
			}
		}
		if (edge.labelBounds) {
			checkFiniteRect("edge label", edge.id, edge.labelBounds, issues);
			checkRectInBounds("edge label", edge.id, edge.labelBounds, resolved.bounds, issues);
			if (expected.labelSize && (
				!nearlyEqual(expected.labelSize.width, edge.labelBounds.width) ||
				!nearlyEqual(expected.labelSize.height, edge.labelBounds.height)
			)) {
				issues.push(issue(
					"size_mismatch",
					[edge.id],
					`edge ${edge.id} label size changed after measurement`,
				));
			}
			for (const node of nodes.values()) {
				if (overlaps(edge.labelBounds, node.bounds)) {
					issues.push(issue(
						"label_node_overlap",
						[edge.id, node.id],
						`edge ${edge.id} label overlaps node ${node.id}`,
					));
				}
			}
			for (const group of groups.values()) {
				if (overlaps(edge.labelBounds, group.titleBounds)) {
					issues.push(issue(
						"label_title_overlap",
						[edge.id, group.id],
						`edge ${edge.id} label overlaps group ${group.id} title`,
					));
				}
			}
			labels.push({id: edge.id, bounds: edge.labelBounds});
		} else if (expected.labelSize) {
			issues.push(issue(
				"missing_element",
				[edge.id],
				`edge ${edge.id} is missing its measured label`,
			));
		}
	}
	forEachPair(labels, (first, second) => {
		if (overlaps(first.bounds, second.bounds)) {
			issues.push(issue(
				"label_label_overlap",
				[first.id, second.id],
				`edge labels ${first.id} and ${second.id} overlap`,
			));
		}
	});
	for (const label of labels) {
		for (const edge of resolved.edges) {
			if (routeIntersectsRect(edge.sections, inset(label.bounds, 0.5))) {
				issues.push(issue(
					"label_edge_overlap",
					[label.id, edge.id],
					`edge ${edge.id} crosses label ${label.id}`,
				));
			}
		}
	}
}

/**
 * checkSectionConnections requires one ordered route from the source border to
 * the target border. Separate or reversed sections are not a complete route.
 */
function checkSectionConnections(
	edge: ResolvedEdge,
	nodes: Map<string, ResolvedNode>,
	issues: GeometryIssue[],
): void {
	for (let index = 1; index < edge.sections.length; index++) {
		const previous = edge.sections[index - 1].slice(-1)[0];
		const current = edge.sections[index][0];
		if (!samePoint(previous, current)) {
			issues.push(issue(
				"disconnected_edge",
				[edge.id],
				`edge ${edge.id} has disconnected route sections`,
			));
			return;
		}
	}
	const source = nodes.get(edge.sourceId);
	const target = nodes.get(edge.targetId);
	const start = edge.sections[0][0];
	const end = edge.sections.slice(-1)[0].slice(-1)[0];
	if (
		!source ||
		!target ||
		!pointOnBorder(start, source.bounds) ||
		!pointOnBorder(end, target.bounds)
	) {
		issues.push(issue(
			"disconnected_edge",
			[edge.id, edge.sourceId, edge.targetId],
			`edge ${edge.id} does not connect both node borders`,
		));
	}
}

/** belongsToGroup follows parent IDs to decide whether a node belongs to a group. */
function belongsToGroup(
	nodeId: string,
	groupId: string,
	nodes: Map<string, MeasuredNode>,
	groups: Map<string, MeasuredGroup>,
): boolean {
	let parent = nodes.get(nodeId)?.parentId;
	while (parent) {
		if (parent === groupId) {
			return true;
		}
		parent = groups.get(parent)?.parentId;
	}
	return false;
}

/**
 * routeViolatesGroup checks the middle of each line. When both connected nodes
 * are inside a group, the line must stay inside it. When both are outside, the
 * line must stay outside it.
 */
function routeViolatesGroup(
	sections: Point[][],
	group: Rect,
	shouldBeInside: boolean,
): boolean {
	for (const section of sections) {
		for (let index = 1; index < section.length; index++) {
			const midpoint = {
				x: (section[index - 1].x + section[index].x) / 2,
				y: (section[index - 1].y + section[index].y) / 2,
			};
			if (pointInRect(midpoint, group, false) !== shouldBeInside) {
				return true;
			}
		}
	}
	return false;
}

/** routeIntersectsRect reports whether any part of a route enters a box. */
function routeIntersectsRect(sections: Point[][], rect: Rect): boolean {
	return sections.some(section => section.some((point, index) =>
		index > 0 && segmentIntersectsRect(section[index - 1], point, rect),
	));
}

/**
 * segmentIntersectsRect checks both straight and diagonal lines. A line may
 * touch the box border, but it may not enter the box.
 */
function segmentIntersectsRect(start: Point, end: Point, rect: Rect): boolean {
	const dx = end.x - start.x;
	const dy = end.y - start.y;
	let minimum = 0;
	let maximum = 1;
	for (const [origin, delta, lower, upper] of [
		[start.x, dx, rect.x, rect.x + rect.width],
		[start.y, dy, rect.y, rect.y + rect.height],
	] as const) {
		if (Math.abs(delta) < Number.EPSILON) {
			if (origin <= lower || origin >= upper) {
				return false;
			}
			continue;
		}
		const first = (lower - origin) / delta;
		const second = (upper - origin) / delta;
		minimum = Math.max(minimum, Math.min(first, second));
		maximum = Math.min(maximum, Math.max(first, second));
		if (minimum >= maximum) {
			return false;
		}
	}
	return maximum > 0 && minimum < 1;
}

/** overlaps treats touching borders as contact, not shared area. */
function overlaps(first: Rect, second: Rect): boolean {
	return (
		first.x < second.x + second.width &&
		second.x < first.x + first.width &&
		first.y < second.y + second.height &&
		second.y < first.y + first.height
	);
}

/** contains requires the whole inner box to be inside the outer box. */
function contains(outer: Rect, inner: Rect): boolean {
	return (
		outer.x <= inner.x &&
		outer.y <= inner.y &&
		outer.x + outer.width >= inner.x + inner.width &&
		outer.y + outer.height >= inner.y + inner.height
	);
}

/** pointInRect can include or exclude points that only touch the border. */
function pointInRect(point: Point, rect: Rect, inclusive: boolean): boolean {
	if (inclusive) {
		return (
			point.x >= rect.x &&
			point.y >= rect.y &&
			point.x <= rect.x + rect.width &&
			point.y <= rect.y + rect.height
		);
	}
	return (
		point.x > rect.x &&
		point.y > rect.y &&
		point.x < rect.x + rect.width &&
		point.y < rect.y + rect.height
	);
}

/** inset makes a box smaller so touching its border does not count as crossing it. */
function inset(rect: Rect, amount: number): Rect {
	return {
		x: rect.x + amount,
		y: rect.y + amount,
		width: Math.max(0, rect.width - amount * 2),
		height: Math.max(0, rect.height - amount * 2),
	};
}

/** isFinitePoint accepts only real numbers that SVG can draw. */
function isFinitePoint(point: Point): boolean {
	return Number.isFinite(point.x) && Number.isFinite(point.y);
}

/** isFiniteRect checks all four numbers in a box. */
function isFiniteRect(rect: Rect): boolean {
	return isFinitePoint(rect) && Number.isFinite(rect.width) && Number.isFinite(rect.height);
}

/** samePoint compares route points while allowing tiny browser number changes. */
function samePoint(first: Point, second: Point): boolean {
	return nearlyEqual(first.x, second.x) && nearlyEqual(first.y, second.y);
}

/** pointOnBorder checks that a route endpoint touches one side of a node box. */
function pointOnBorder(point: Point, rect: Rect): boolean {
	const withinX = point.x >= rect.x - 0.1 && point.x <= rect.x + rect.width + 0.1;
	const withinY = point.y >= rect.y - 0.1 && point.y <= rect.y + rect.height + 0.1;
	const onVertical = nearlyEqual(point.x, rect.x) ||
		nearlyEqual(point.x, rect.x + rect.width);
	const onHorizontal = nearlyEqual(point.y, rect.y) ||
		nearlyEqual(point.y, rect.y + rect.height);
	return (withinY && onVertical) || (withinX && onHorizontal);
}

/** nearlyEqual hides only browser rounding smaller than one tenth of a unit. */
function nearlyEqual(first: number, second: number): boolean {
	return Math.abs(first - second) <= 0.1;
}

/** forEachPair compares each pair once. */
function forEachPair<T>(items: T[], visit: (first: T, second: T) => void): void {
	for (let first = 0; first < items.length; first++) {
		for (let second = first + 1; second < items.length; second++) {
			visit(items[first], items[second]);
		}
	}
}

/** issue creates one layout problem report. */
function issue(
	code: GeometryIssueCode,
	elementIds: string[],
	message: string,
): GeometryIssue {
	return {code, elementIds, message};
}
