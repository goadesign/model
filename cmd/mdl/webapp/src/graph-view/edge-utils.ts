import { Point, calculateDistance } from './constants';
import { Segment } from './intersect';

// Define interfaces locally since they're not exported from graph.ts
interface NodeStyle {
	background?: string;
	stroke?: string;
	opacity?: number;
	fontSize?: number;
	shape?: string;
	border?: string;
}

interface Node extends Point {
	id: string;
	title: string;
	sub: string;
	description: string;
	width: number;
	height: number;
	ref?: SVGGElement;
	selected?: boolean;
	intersect: (p: Point) => Point;
	style: NodeStyle;
}

interface EdgeVertex extends Point {
	id: string;
	selected?: boolean;
	ref?: SVGElement;
	label?: boolean;
	auto?: boolean;
}

interface EdgeStyle {
	color?: string;
	thickness?: number;
	fontSize?: number;
	position?: number;
	dashed?: boolean;
}

interface Edge {
	id: string;
	label: string;
	from: Node;
	to: Node;
	vertices?: EdgeVertex[];
	ref?: SVGGElement;
	style: EdgeStyle;
	initVertex: (p: Point) => EdgeVertex;
	userDeletedVertices?: boolean; // Track if user explicitly deleted vertices
}

export interface EdgeLabelPlacement extends Point {
	orientation: 'horizontal' | 'vertical';
	segment?: Segment;
	movable: boolean;
}

/**
 * Calculate edge vertices, handling multi-edge scenarios and auto-vertices
 */
export function calculateEdgeVertices(edge: Edge): Point[] {
	const n1 = edge.from, n2 = edge.to;
	
	
	// if vertices exists, follow them
	let vertices: Point[] = edge.vertices ? edge.vertices.concat() : [];
	// Don't remove label vertices - they should be preserved for rendering
	// (The autoLayout process handles replacing old ones with new ones)

	vertices.unshift(n1)
	vertices.push(n2)

	const interior = vertices.slice(1, -1) as EdgeVertex[];
	const firstRoutingVertex = interior.find(vertex => !vertex.label) ?? n2;
	let lastRoutingVertex: Point = n1;
	for (let index = interior.length - 1; index >= 0; index--) {
		if (!interior[index].label) {
			lastRoutingVertex = interior[index];
			break;
		}
	}
	if (!n1.intersect || !n2.intersect) {
		throw new Error(`edge ${edge.id} nodes must be drawn before routing`);
	}
	vertices[0] = n1.intersect(firstRoutingVertex);
	vertices[vertices.length - 1] = n2.intersect(lastRoutingVertex);
	
	return vertices;
}

/**
 * Calculate the label anchor and the orientation of the path segment that owns
 * it. Renderers use the orientation to place text beside the relationship line
 * instead of centering text over vertical segments.
 */
export function calculateLabelPlacement(
	vertices: Point[],
	position: number,
	defaultPoint: Point,
): EdgeLabelPlacement {
	let point = {x: defaultPoint.x, y: defaultPoint.y};
	let segment: Segment | undefined;
	const labelIndex = vertices.findIndex(vertex => (vertex as EdgeVertex).label);
	let movable = true;

	if (labelIndex >= 0) {
		const labelVertex = vertices[labelIndex] as EdgeVertex;
		point = labelVertex;
		movable = labelVertex.auto === true;
		const adjacentSegments: Segment[] = [];
		if (labelIndex > 0) {
			adjacentSegments.push({p: vertices[labelIndex - 1], q: point});
		}
		if (labelIndex < vertices.length - 1) {
			adjacentSegments.push({p: point, q: vertices[labelIndex + 1]});
		}
		segment = adjacentSegments.reduce<Segment | undefined>((longest, candidate) => {
			if (!longest) {
				return candidate;
			}
			return calculateDistance(candidate.p, candidate.q) >
				calculateDistance(longest.p, longest.q)
				? candidate
				: longest;
		}, undefined);
	} else {
		const totalLength = vertices.slice(1).reduce(
			(sum, vertex, index) => sum + calculateDistance(vertices[index], vertex),
			0,
		);
		const targetLength = totalLength * position;
		let traversed = 0;
		for (let index = 1; index < vertices.length; index++) {
			const candidate = {p: vertices[index - 1], q: vertices[index]};
			const length = calculateDistance(candidate.p, candidate.q);
			if (length > 0 && traversed + length >= targetLength) {
				const segmentPosition = (targetLength - traversed) / length;
				point = {
					x: candidate.p.x + (candidate.q.x - candidate.p.x) * segmentPosition,
					y: candidate.p.y + (candidate.q.y - candidate.p.y) * segmentPosition,
				};
				segment = candidate;
				break;
			}
			traversed += length;
		}
	}

	const horizontalDistance = segment ? Math.abs(segment.q.x - segment.p.x) : 0;
	const verticalDistance = segment ? Math.abs(segment.q.y - segment.p.y) : 0;
	return {
		...point,
		orientation: verticalDistance > horizontalDistance ? 'vertical' : 'horizontal',
		segment,
		movable,
	};
}

/** createEdgePath draws the given points without cutting or moving the route. */
export function createEdgePath(vertices: Point[]): string {
	if (vertices.length < 2) {
		throw new Error("edge route must contain at least two points");
	}
	return vertices.slice(1).reduce(
		(path, point) => `${path} L${point.x},${point.y}`,
		`M${vertices[0].x},${vertices[0].y}`,
	);
}