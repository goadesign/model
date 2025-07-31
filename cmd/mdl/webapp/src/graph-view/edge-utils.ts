import { Point, BBox, calculateDistance } from './constants';
import { intersectPolylineBox, Segment } from './intersect';

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

interface GraphData {
	id: string;
	name: string;
	nodesMap: Map<string, Node>;
	edges: Edge[];
	edgeVertices: Map<string, EdgeVertex>;
	groupsMap: Map<string, any>;
	metadata: any;
}

/**
 * Calculate edge vertices, handling multi-edge scenarios and auto-vertices
 */
export function calculateEdgeVertices(edge: Edge, data: GraphData): Point[] {
	const n1 = edge.from, n2 = edge.to;
	
	// if vertices exists, follow them
	let vertices: Point[] = edge.vertices ? edge.vertices.concat() : [];
	// Only remove label auto vertices (not routing vertices from ELK)
	const tmp = (vertices as EdgeVertex[]);
	tmp.forEach(v => v.auto && v.label && data.edgeVertices.delete(v.id))
	vertices = tmp.filter(v => !(v.auto && v.label))

	if (vertices.length == 0 && !edge.userDeletedVertices) {
		// Only create auto vertices if user hasn't explicitly deleted them
		// for edges with same "from" and "to", we must spread the labels so they don't overlap
		// lookup the other "same" edges
		const sameEdges = data.edges.filter(e => e.from == edge.from && e.to == edge.to)
		let spreadPos = 0
		if (sameEdges.length > 1) {
			const idx = sameEdges.indexOf(edge) // my index in the list of same edges
			spreadPos = idx - (sameEdges.length - 1) / 2

			let spreadX = 0, spreadY = 0;
			if (Math.abs(n1.x - n2.x) > Math.abs(n1.y - n2.y)) {
				spreadY = spreadPos * 70
			} else {
				spreadX = spreadPos * 200
			}
			const v = edge.initVertex({
				x: (n1.x + n2.x) / 2 + spreadX,
				y: (n1.y + n2.y) / 2 + spreadY
			})
			v.label = true
			v.auto = true
			vertices.push(v)
		} else {
			// If there are no user-defined vertices and not a multi-edge scenario,
			// we don't create any auto-vertices here. AutoLayout will provide them.
			// The path will be a straight line between n1 and n2 (after intersection points are calculated).
			// ELK/autoLayout is responsible for providing bend points for non-straight lines.
		}
	}

	vertices.unshift(n1)
	vertices.push(n2)

	vertices[0] = n1.intersect(vertices[1])
	vertices[vertices.length - 1] = n2.intersect(vertices[vertices.length - 2])
	
	return vertices;
}

/**
 * Calculate the position for edge label along the edge path
 */
export function calculateLabelPosition(vertices: Point[], position: number, fallback: Point): Point {
	// where along the edge is the label?
	// position of label
	let pLabel: Point = vertices.find(v => (v as any).label)
	if (!pLabel) {
		let sum = 0 // total length of the edge, sum of segments
		for (let i = 1; i < vertices.length; i++) {
			sum += calculateDistance(vertices[i - 1], vertices[i])
		}
		pLabel = {x: fallback.x, y: fallback.y} // fallback for corner cases
		let acc = 0
		for (let i = 1; i < vertices.length; i++) {
			const d = calculateDistance(vertices[i - 1], vertices[i])
			if (acc + d > sum * position) {
				const pos = (sum * position - acc) / d
				pLabel = {
					x: vertices[i - 1].x + (vertices[i].x - vertices[i - 1].x) * pos,
					y: vertices[i - 1].y + (vertices[i].y - vertices[i - 1].y) * pos
				}
				break
			}
			acc += d
		}
	}
	return pLabel;
}

/**
 * Create edge segments and generate SVG path
 */
export function createEdgeSegments(vertices: Point[], bbox: BBox, n1: Point, n2: Point): { segments: Segment[], path: string } {
	const segments: Segment[] = []
	for (let i = 1; i < vertices.length; i++) {
		segments.push({p: vertices[i - 1], q: vertices[i]})
	}
	// splice edge over label box
	intersectPolylineBox(segments, bbox)

	// Generate path based on routing style - SIMPLIFIED
	let path: string
	
	// Always draw a polyline through the segments.
	// If segments.length is 1 (meaning direct connection or only label vertex), it will be a straight line.
	// If autoLayout provided bend points, those will be in `segments`.
	if (segments.length > 0) {
		path = `M${segments[0].p.x},${segments[0].p.y}`
		for (let i = 0; i < segments.length; i++) {
			const s = segments[i]
			// For polylines, we just draw line segments to each vertex point.
			// The ELK 'POLYLINE' routing should give us the necessary bend points.
			path += ` L${s.q.x},${s.q.y}`
		}
	} else {
		// Fallback for edges with no segments (should ideally not happen if n1 and n2 are defined)
		// Draw a straight line between n1 and n2 directly if no vertices/segments exist.
		// Note: Intersection points are calculated before this, so n1/n2 are already adjusted.
		path = `M${n1.x},${n1.y} L${n2.x},${n2.y}`
	}
	
	return { segments, path };
}