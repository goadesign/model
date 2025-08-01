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
	// Don't remove label vertices - they should be preserved for rendering
	// (The autoLayout process handles replacing old ones with new ones)
	const tmp = (vertices as EdgeVertex[]);

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

	// Calculate intersection points with node boundaries
	// Find first non-label vertex for start intersection
	let firstRoutingVertex = vertices[vertices.length - 1]; // Default to end node
	for (let i = 1; i < vertices.length - 1; i++) {
		if (!(vertices[i] as any).label) {
			firstRoutingVertex = vertices[i];
			break;
		}
	}
	
	// Find last non-label vertex for end intersection  
	let lastRoutingVertex = vertices[0]; // Default to start node
	for (let i = vertices.length - 2; i > 0; i--) {
		if (!(vertices[i] as any).label) {
			lastRoutingVertex = vertices[i];
			break;
		}
	}
	
	// For connections without routing vertices, ensure we have proper direction
	// The defaults are already correct: firstRoutingVertex = n2, lastRoutingVertex = n1
	
	
	// Calculate proper node boundary intersection
	const calculateNodeIntersection = (node: any, targetPoint: Point): Point => {
		const nodeShape = node.style?.shape?.toLowerCase() || 'box';
		const dx = targetPoint.x - node.x;
		const dy = targetPoint.y - node.y;
		const nodeCenter = { x: node.x, y: node.y };
		
		// If target is at center, default to right edge
		if (Math.abs(dx) < 0.01 && Math.abs(dy) < 0.01) {
			return { x: node.x + node.width / 2, y: node.y };
		}
		
		if (nodeShape === 'cylinder') {
			// Cylinder shape intersection (same as shapes.ts)
			const w = node.width;
			const rx = w / 2;
			const ry = rx / (5.5 + w / 70);
			const halfHeight = node.height / 2;
			
			// First calculate rectangular bounds intersection
			const angle = Math.atan2(dy, dx);
			const cos = Math.cos(angle);
			const sin = Math.sin(angle);
			
			// Check intersection with rectangular bounds
			let t = Infinity;
			if (Math.abs(cos) > 0.01) {
				t = Math.min(t, Math.abs(rx / cos));
			}
			if (Math.abs(sin) > 0.01) {
				t = Math.min(t, Math.abs(halfHeight / sin));
			}
			
			const rectX = node.x + cos * t;
			const rectY = node.y + sin * t;
			
			// Check if we need elliptical intersection for top/bottom curves
			const topCurveY = node.y - halfHeight + ry;
			const bottomCurveY = node.y + halfHeight - ry;
			
			if (rectY < topCurveY || rectY > bottomCurveY) {
				// Use ellipse intersection for curved parts
				const ellipseY = rectY < topCurveY ? node.y - halfHeight + ry : node.y + halfHeight - ry;
				// Solve for ellipse intersection
				const a = 1 / (rx * rx);
				const b = -2 * node.x / (rx * rx);
				const c = (node.x * node.x) / (rx * rx) + ((ellipseY - node.y) * (ellipseY - node.y)) / (ry * ry) - 1;
				
				const discriminant = b * b - 4 * a * c;
				if (discriminant >= 0) {
					const sqrt_d = Math.sqrt(discriminant);
					const x1 = (-b + sqrt_d) / (2 * a);
					const x2 = (-b - sqrt_d) / (2 * a);
					
					// Choose the intersection in the direction of the target
					const intersectX = dx > 0 ? Math.max(x1, x2) : Math.min(x1, x2);
					return { x: intersectX, y: ellipseY };
				}
			}
			
			return { x: rectX, y: rectY };
			
		} else if (nodeShape === 'circle') {
			const radius = node.width / 2;
			const angle = Math.atan2(dy, dx);
			return {
				x: node.x + Math.cos(angle) * radius,
				y: node.y + Math.sin(angle) * radius
			};
			
		} else if (nodeShape === 'ellipse') {
			const rx = node.width * 0.55;
			const ry = node.width * 0.45;
			const angle = Math.atan2(dy, dx);
			const cos = Math.cos(angle);
			const sin = Math.sin(angle);
			
			// Parametric ellipse intersection
			const t = Math.sqrt((rx * rx * sin * sin) + (ry * ry * cos * cos));
			return {
				x: node.x + (rx * cos * ry) / t,
				y: node.y + (ry * sin * rx) / t
			};
			
		} else {
			// Default rectangular intersection
			const halfWidth = node.width / 2;
			const halfHeight = node.height / 2;
			const angle = Math.atan2(dy, dx);
			const cos = Math.cos(angle);
			const sin = Math.sin(angle);
			
			// Calculate which edge we hit first
			let t = Infinity;
			if (Math.abs(cos) > 0.01) {
				t = Math.min(t, Math.abs(halfWidth / cos));
			}
			if (Math.abs(sin) > 0.01) {
				t = Math.min(t, Math.abs(halfHeight / sin));
			}
			
			return {
				x: node.x + cos * t,
				y: node.y + sin * t
			};
		}
	};
	
	// Calculate intersections, but use the direction from center to NEXT vertex in sequence
	// This ensures the line exits the node in the direction it needs to go
	
	// For start intersection: use direction from node center to first routing vertex
	let startIntersection = calculateNodeIntersection(n1, firstRoutingVertex);
	
	// For end intersection: use direction from node center to last routing vertex  
	let endIntersection = calculateNodeIntersection(n2, lastRoutingVertex);
	
	// Start intersection: NO offset needed - intersection already gives perfect boundary point
	// End intersection: NO offset needed - the arrow marker now has refX="0" so the tip is at the endpoint
	
	vertices[0] = startIntersection;
	vertices[vertices.length - 1] = endIntersection;
	
	return vertices;
}

/**
 * Calculate the position for edge label along the edge path
 */
export function calculateLabelPosition(vertices: Point[], position: number, fallback: Point, edge?: any): Point {
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
	
	// Debug the final segment that will have the arrow
	if (segments.length > 0) {
		const lastSegment = segments[segments.length - 1];
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