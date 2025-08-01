import {defs} from "./defs";
import {create, setPosition} from "./svg-create";
import {cursorInteraction} from "svg-editor-tools/lib/cursor-interaction";
import {shapes} from "./shapes";
import {
	boxesOverlap,
	cabDistance,
	insideBox,
	intersectPolylineBox,
	project,
	scaleBox,
	Segment,
	uncenterBox
} from "./intersect";
import {autoLayout} from "./layout";
import {Undo} from "./undo";
import {
	ADD_LABEL_VERTEX,
	ADD_VERTEX,
	DEL_VERTEX,
	DESELECT,
	findShortcut,
	REDO,
	SELECT_ALL,
	UNDO,
	ZOOM_100,
	ZOOM_FIT,
	ZOOM_IN,
	ZOOM_OUT
} from "../shortcuts";
import {
	Point,
	BBox,
	NodeStyle,
	EdgeStyle,
	DEFAULT_EDGE_STYLE,
	DEFAULT_NODE_STYLE,
	SVG_STYLES,
	SVG_PADDING,
	DEFAULT_GRID_SIZE,
	applyStyle,
	calculateDistance
} from "./constants";
import {
	calculateEdgeVertices,
	calculateLabelPosition,
	createEdgeSegments
} from "./edge-utils";


// Point and BBox interfaces are now imported from constants

export interface Group extends BBox {
	id: string;
	name: string;
	nodes: (Node | Group)[];
	ref?: SVGGElement;
	style: NodeStyle;
}

export interface Node extends BBox {
	id: string;
	title: string;
	sub: string;
	description: string;

	ref?: SVGGElement;
	selected?: boolean;

	intersect: (p: Point) => Point

	style: NodeStyle
}

// NodeStyle interface is now imported from constants

// EdgeStyle interface is now imported from constants

// Default styles are now imported from constants
const defaultEdgeStyle = DEFAULT_EDGE_STYLE;
const defaultNodeStyle = DEFAULT_NODE_STYLE;

// Edge and EdgeVertex interfaces are now defined in edge-utils.ts
// Using local interfaces for compatibility with existing code
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
	labelVertex?: EdgeVertex; // ELK-calculated label position (separate from routing vertices)
}

interface EdgeVertex extends Point {
	id: string
	selected?: boolean
	edge: Edge
	ref?: SVGElement
	label?: boolean
	auto?: boolean
}

interface Layout {
	[k: string]: Point | (Point & { label: boolean })[] | boolean
}

export class GraphData {
	id: string;
	name: string;
	nodesMap: Map<string, Node>;
	edges: Edge[];
	edgeVertices: Map<string, EdgeVertex>
	groupsMap: Map<string, Group>;
	metadata: any;
	private _undo: Undo<Layout>;
	private _gridVisible: boolean = true;
	private _snapToGrid: boolean = true;
	private _gridSize: number = 25;
	private _skipAutoFit: boolean = false;

	constructor(id?: string, name?: string) {
		this.id = id;
		this.name = name;

		this.edges = [];
		this.edgeVertices = new Map;
		this.nodesMap = new Map;
		this.groupsMap = new Map;

		this._undo = new Undo<Layout>(
			this.id,
			() => this.exportLayout(true),
			(lo) => this.importLayout(lo, true)
		)

		// @ts-ignore
		window.graph = this
	}

	// after the graph model is build using addNode, addEdge etc, initialize
	init(layout?: Layout) {
		layout && this.importLayout(layout)
		this._undo = new Undo<Layout>(
			this.id,
			() => this.exportLayout(true),
			(lo) => this.importLayout(lo, true)
		)
		if (this._undo.length()) {
			this.importLayout(this._undo.currentState())
		}
	}

	addNode(id: string, label: string, sub: string, description: string, style: NodeStyle) {
		if (this.nodesMap.has(id)) throw Error('duplicate node: ' + id)
		
		// ALWAYS use fixed dimensions - ignore any width/height from incoming styles
		// This ensures visual consistency regardless of model data
		const shape = style.shape || 'Box';
		const isPersonShape = shape.toLowerCase() === 'person';
		const isCylinderShape = shape.toLowerCase() === 'cylinder';
		
		// Fixed dimensions based only on shape type - completely ignore style.width/height
		let width: number;
		let height: number;
		
		if (isPersonShape) {
			width = 200;
			height = 240;
		} else if (isCylinderShape) {
			// Make cylinders wider to better accommodate text content
			width = 280; // Same width as boxes for consistency
			height = 180; // Same height as boxes for consistency
		} else {
			// Standard rectangular shapes
			width = 280;
			height = 180;
		}
		
		// Remove any width/height from style to prevent conflicts
		const cleanStyle = {...style};
		delete cleanStyle.width;
		delete cleanStyle.height;
		
		const n: Node = {
			id, title: label, sub, description, style: {...defaultNodeStyle, ...cleanStyle},
			x: 0, y: 0, width, height, intersect: null
		}
		this.nodesMap.set(n.id, n)
	}

	nodes() {
		return Array.from(this.nodesMap.values())
	}

	addEdge(id: string, fromNode: string, toNode: string, label: string, vertices: Point[], style: EdgeStyle) {
		vertices && vertices.forEach((p, i) => {
			const v = p as EdgeVertex
			v.id = `v-${id}-${i}`
			this.edgeVertices.set(v.id, v)
		})
		const randomID = () => Math.round(Math.random() * 1e10).toString(36)
		const initVertex = (p: Point) => {
			const v = p as EdgeVertex
			if (!v.id) {
				v.id = `v-${randomID()}`
				this.edgeVertices.set(v.id, v)
			}
			v.edge = edge
			return p as EdgeVertex
		}
		const edge = {
			id,
			from: this.nodesMap.get(fromNode),
			to: this.nodesMap.get(toNode),
			label,
			vertices: null as EdgeVertex[],
			style: {...defaultEdgeStyle, ...style},
			initVertex,
			userDeletedVertices: false
		}
		this.edges.push(edge)
		if (vertices) {
			edge.vertices = vertices.map(p => edge.initVertex(p))
		}
	}

	addGroup(id: string, name: string, nodesOrGroups: string[], style: NodeStyle) {
		if (this.groupsMap.has(id)) {
			console.error(`Group exists: ${id} ${name}`)
			return
		}
		const group: Group = {
			id, name, x: null, y: null, width: null, height: null,
			nodes: nodesOrGroups.map(k => {
				const n = this.nodesMap.get(k) || this.groupsMap.get(k)
				if (!n) console.error(`Node or group ${k} not found for group ${id} "${name}"`)
				return n
			}).filter(Boolean),
			style
		}
		this.groupsMap.set(id, group)
	}

	// private rebuildNode(node: Node) {
	// 	const p = node.ref.parentElement;
	// 	p.removeChild(node.ref)
	// 	node.ref = buildNode(node, this)
	// 	p.appendChild(node.ref)
	// 	this.redrawEdges(node)
	// 	this.redrawGroups(node)
	// }

	setNodeSelected(node: Node, selected: boolean) {
		node.selected = selected
		selected ?
			node.ref.classList.add('selected') :
			node.ref.classList.remove('selected')
		this.updateEdgesSel()
	}

	private updateEdgesSel() {
		this.edges.forEach(e => {
			if (e.to.selected || e.from.selected) {
				e.ref.classList.add('selected')
			} else {
				e.ref.classList.remove('selected')
			}
		})
	}

	moveNode(n: Node, x: number, y: number, disableSnap: boolean = false, skipUndo: boolean = false) {
		if (!n) return
		
		// Apply snap-to-grid if enabled and not explicitly disabled
		if (this._snapToGrid && !disableSnap) {
			const snapped = this.snapToGrid(x, y);
			x = snapped.x;
			y = snapped.y;
		}
		
		if (n.x == x && n.y == y) return
		
		if (!skipUndo) {
			this._undo.beforeChange()
		}
		n.x = x;
		n.y = y;
		setPosition(n.ref, x, y)
		this.redrawEdges(n);
		this.redrawGroups(n)
		if (!skipUndo) {
			this._undo.change()
		}
	}

	moveEdgeVertex(v: EdgeVertex, x: number, y: number, disableSnap: boolean = false, skipUndo: boolean = false) {
		
		if (this._snapToGrid && !disableSnap) {
			const snapped = this.snapToGrid(x, y);
			x = snapped.x;
			y = snapped.y;
		}
		// Use exact coordinates (no rounding needed with modern grid system)
		
		if (v.x == x && v.y == y) return
		if (!skipUndo) {
			this._undo.beforeChange()
		}
		v.x = x;
		v.y = y;
		this.redrawEdge(v.edge)
		if (!skipUndo) {
			this._undo.change()
		}
	}

	moveSelected(dx: number, dy: number, disableSnap: boolean = false) {
		this.nodes().forEach(n => n.selected && this.moveNode(n, n.x + dx, n.y + dy, disableSnap, false))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, v.x + dx, v.y + dy, disableSnap, false))
	}

	insertEdgeVertex(edge: Edge, p: Point, pos: number, isLabel: boolean) {
		this._undo.beforeChange()
		const v = edge.initVertex(p)
		v.selected = true
		if (isLabel) { // when shift down, make it label position
			edge.vertices.forEach(v => v.label = false)
			v.label = true
		}
		edge.vertices.splice(pos - 1, 0, v)
		this.redrawEdge(edge)
		this._undo.change()
	}

	deleteEdgeVertex(v: EdgeVertex) {
		this._undo.beforeChange()
		
		const index = v.edge.vertices.indexOf(v)
		if (index >= 0) {
			v.edge.vertices.splice(index, 1)
			this.edgeVertices.delete(v.id)
			
			// Mark that user explicitly deleted vertices from this edge
			v.edge.userDeletedVertices = true
		}
		
		this.redrawEdge(v.edge)
		this._undo.change()
	}

	changed() {
		return this._undo.changed()
	}

	undo() {
		this._undo.undo()
	}

	// moves the entire graph to be aligned top-left of the drawing area
	// used to bring back to visible the nodes that end up at negative coordinates
	alignTopLeft() {
		const contentBounds = this.calculateContentBounds()
		const padding = 100 // Reasonable padding for viewport
		
		const offsetX = -contentBounds.x + padding
		const offsetY = -contentBounds.y + padding
		
		// Set flag to prevent React useEffect from calling fitToView during this operation
		this._skipAutoFit = true
		
		this._undo.beforeChange()
		
		this.nodesMap.forEach(node => {
			this.moveNode(node, node.x + offsetX, node.y + offsetY, true, true) // Disable snap and undo during reset
		})
		
		this.edgeVertices.forEach(vertex => {
			this.moveEdgeVertex(vertex, vertex.x + offsetX, vertex.y + offsetY, true, true) // Disable snap and undo during reset
		})
		
		this._undo.change()
		
		// DON'T clear view state here - let resetPanTransform handle it to avoid React useEffect recursion
	}
	
	// Reset pan transform to (0,0) while preserving zoom
	resetPanTransform() {
		const currentZoom = getZoom()
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (zoomGroup) {
			zoomGroup.setAttribute('transform', `scale(${currentZoom}) translate(0, 0)`)
			updatePanningOptimized(this)
		}
		
		// Clear view state so this reset is not overridden
		clearViewState(this.id)
		
		// Reset the skip auto fit flag after reset is complete
		this._skipAutoFit = false
	}
	
	// Check if auto-fit should be skipped (used by React useEffect)
	shouldSkipAutoFit(): boolean {
		return this._skipAutoFit
	}

	// Reset view to default state: 100% zoom, centered at origin
	resetView() {
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (zoomGroup) {
			// Reset to 100% zoom, centered at origin
			zoomGroup.setAttribute('transform', 'scale(1) translate(0, 0)')
			updatePanning()
		}
		
		// Clear any saved view state so this reset position is not overridden
		clearViewState(this.id)
	}

	//redraw connected edges
	private redrawEdges(n: Node) {
		this.edges.forEach(e => (n == e.from || n == e.to) && this.redrawEdge(e))
		this.updateEdgesSel()
	}

	redrawEdge(e: Edge) {
		const p = e.ref.parentElement;
		p.removeChild(e.ref)
		e.ref = buildEdge(this, e)
		p.append(e.ref)
	}

	private redrawGroups(node: Node) {
		this.groupsMap.forEach(group => {
			//if (group.nodes.indexOf(node) == -1) return
			const p = group.ref.parentElement
			p.removeChild(group.ref)
			buildGroup(group)
			p.append(group.ref)
		})
	}

	exportSVG() {
		// Get the original SVG
		const originalSvg: SVGSVGElement = document.querySelector('svg#graph')
		const elastic = originalSvg.querySelector('rect.elastic')
		
		// Clone the SVG for export (completely separate from the live one)
		const exportSvg = originalSvg.cloneNode(true) as SVGSVGElement
		
		// Remove elastic element from export
		const exportElastic = exportSvg.querySelector('rect.elastic')
		if (exportElastic) {
			exportElastic.remove()
		}
		
		// Calculate actual content bounds including all elements
		const contentBounds = this.calculateContentBounds()
		
		// Add padding around content
		const padding = 50
		
		// Calculate final export dimensions (always positive)
		const exportWidth = contentBounds.width + (padding * 2)
		const exportHeight = contentBounds.height + (padding * 2)
		
		// Calculate offset to move content to start at (padding, padding) within the export area
		const offsetX = -contentBounds.x + padding
		const offsetY = -contentBounds.y + padding
		
		// Apply export positioning to the cloned SVG elements
		const exportZoomGroup = exportSvg.querySelector('g.zoom') as SVGGElement
		if (exportZoomGroup) {
			// Reset zoom to 1 and apply offset transform to position content properly
			exportZoomGroup.setAttribute('transform', `scale(1) translate(${offsetX}, ${offsetY})`)
		}
		
		// Set proper viewBox and dimensions for export - viewBox always starts at (0,0)
		exportSvg.setAttribute('viewBox', `0 0 ${exportWidth} ${exportHeight}`)
		exportSvg.setAttribute('width', String(exportWidth))
		exportSvg.setAttribute('height', String(exportHeight))
		
		// Add required SVG namespace for browser compatibility
		exportSvg.setAttribute('xmlns', 'http://www.w3.org/2000/svg')
		
		// Inject metadata with current layout
		const script = document.createElement('script')
		script.setAttribute('type', 'application/json')
		this.metadata.layout = this.exportLayout()
		script.append('<![CDATA[' + escapeCdata(JSON.stringify(this.metadata, null, 2)) + ']]>')
		exportSvg.insertBefore(script, exportSvg.firstChild)
		
		// Get the export SVG as string
		const src = exportSvg.outerHTML
		
		// No restoration needed since we never touched the original SVG!
		return src
	}

	// Calculate the actual bounds of all content including nodes, edges, and groups
	calculateContentBounds(): BBox {
		let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity
		
		// Process nodes (including their actual dimensions)
		this.nodes().forEach(node => {
			const left = node.x - node.width / 2
			const right = node.x + node.width / 2
			const top = node.y - node.height / 2
			const bottom = node.y + node.height / 2
			
			minX = Math.min(minX, left)
			maxX = Math.max(maxX, right)
			minY = Math.min(minY, top)
			maxY = Math.max(maxY, bottom)
		})
		
		// Process edge vertices (much faster than complex label calculations)
		this.edgeVertices.forEach(vertex => {
			minX = Math.min(minX, vertex.x - 5)
			maxX = Math.max(maxX, vertex.x + 5)
			minY = Math.min(minY, vertex.y - 5)
			maxY = Math.max(maxY, vertex.y + 5)
		})
		
		// Process groups
		this.groupsMap.forEach(group => {
			const left = group.x - group.width / 2
			const right = group.x + group.width / 2
			const top = group.y - group.height / 2
			const bottom = group.y + group.height / 2
			
			minX = Math.min(minX, left)
			maxX = Math.max(maxX, right)
			minY = Math.min(minY, top)
			maxY = Math.max(maxY, bottom)
		})
		
		// Process edges (simplified - just endpoints and vertices, skip complex label calculations)
		this.edges.forEach(edge => {
			// Edge endpoints
			minX = Math.min(minX, edge.from.x - 10, edge.to.x - 10)
			maxX = Math.max(maxX, edge.from.x + 10, edge.to.x + 10)
			minY = Math.min(minY, edge.from.y - 10, edge.to.y - 10)
			maxY = Math.max(maxY, edge.from.y + 10, edge.to.y + 10)
			
			// Edge vertices (if any)
			if (edge.vertices) {
				edge.vertices.forEach(vertex => {
					minX = Math.min(minX, vertex.x - 10)
					maxX = Math.max(maxX, vertex.x + 10)
					minY = Math.min(minY, vertex.y - 10)
					maxY = Math.max(maxY, vertex.y + 10)
				})
			}
			
			// Simplified label bounds (avoid expensive path calculations)
			if (edge.label && edge.label.trim()) {
				// Just use approximate center between from and to nodes
				const centerX = (edge.from.x + edge.to.x) / 2
				const centerY = (edge.from.y + edge.to.y) / 2
				const approxLabelSize = edge.label.length * 10 + 50 // Rough estimate
				
				minX = Math.min(minX, centerX - approxLabelSize)
				maxX = Math.max(maxX, centerX + approxLabelSize)
				minY = Math.min(minY, centerY - 25)
				maxY = Math.max(maxY, centerY + 25)
			}
		})
		
		// Handle empty graph
		if (minX === Infinity) {
			return { x: 0, y: 0, width: 100, height: 100 }
		}
		
		return {
			x: minX,
			y: minY,
			width: maxX - minX,
			height: maxY - minY
		}
	}



	/**
	 * @param full when true, the edges without vertices are saved too, used for undo buffer
	 *        for saving, full is false
	 */
	exportLayout(full = false) {
		const ret: Layout = {}
		this.nodes().forEach(n => ret[n.id] = {x: n.x, y: n.y})
		this.edges.forEach(e => {
			if (!e.vertices) return
			// Save all vertices (both user and auto-generated), preserving their properties
			const lst = e.vertices.map(v => ({
				x: v.x, 
				y: v.y, 
				label: v.label,
				auto: v.auto // Preserve auto flag so we know which are ELK-generated
			}));
			if (lst.length || full) {
				ret[`e-${e.id}`] = lst
			}
			// Also save the userDeletedVertices flag as metadata
			if (e.userDeletedVertices) {
				ret[`e-${e.id}-deleted`] = true
			}
		})
		return ret
	}

	setSaved() {
		this._undo.setSaved()
	}

	importLayout(layout: { [key: string]: any }, rerender = false) {
		// First pass: collect all coordinate values to find bounds
		const coordinates: Array<{x: number, y: number}> = [];
		
		Object.entries(layout).forEach(([k, v]) => {
			if (!k.startsWith('e-') && v.x !== undefined && v.y !== undefined) {
				// Node coordinates
				coordinates.push({x: v.x, y: v.y});
			} else if (k.startsWith('e-') && Array.isArray(v)) {
				// Edge vertex coordinates
				v.forEach((vertex: any) => {
					if (vertex.x !== undefined && vertex.y !== undefined) {
						coordinates.push({x: vertex.x, y: vertex.y});
					}
				});
			}
		});
		
		// Calculate normalization offset if we have coordinates
		let offsetX = 0;
		let offsetY = 0;
		
		if (coordinates.length > 0) {
			const minX = Math.min(...coordinates.map(c => c.x));
			const minY = Math.min(...coordinates.map(c => c.y));
			
			// Only normalize if coordinates are problematic (negative or very large)
			if (minX < -100 || minY < -100 || Math.max(...coordinates.map(c => c.x)) > 3000 || Math.max(...coordinates.map(c => c.y)) > 2000) {
				const padding = 50;
				offsetX = -minX + padding;
				offsetY = -minY + padding;
			}
		}
		
		// Second pass: apply coordinates with normalization
		Object.entries(layout).forEach(([k, v]) => {
			// nodes
			const n = this.nodesMap.get(k)
			if (n) {
				n.x = v.x + offsetX
				n.y = v.y + offsetY
			} else
				// edge vertices
			if (k.startsWith('e-') && !k.endsWith('-deleted')) {
				const edge = this.edges.find(e => e.id == k.slice(2))
				if (!edge) return;
				edge.vertices && edge.vertices.forEach(v => this.edgeVertices.delete(v.id))
				edge.vertices = v.map((p: Point) => {
					const normalizedPoint = { 
						x: p.x + offsetX, 
						y: p.y + offsetY 
					} as Point;
					// Preserve any additional properties like 'label' and 'auto'
					Object.assign(normalizedPoint, p, { x: p.x + offsetX, y: p.y + offsetY });
					const vertex = edge.initVertex(normalizedPoint);
					// Ensure auto flag is preserved after initVertex
					if ((p as any).auto) {
						vertex.auto = true;
					}
					return vertex;
				})
				return;
			}
			if (k.endsWith('-deleted')) {
				const edgeId = k.slice(2, -8); // Remove 'e-' prefix and '-deleted' suffix
				const edge = this.edges.find(e => e.id == edgeId)
				if (edge && v === true) {
					edge.userDeletedVertices = true
				}
				return;
			}
		})
		if (rerender) {
			this.nodes().forEach(n => setPosition(n.ref, n.x, n.y))
			this.edges.forEach(e => this.redrawEdge(e))
			this.updateEdgesSel()
			this.redrawGroups(null)
		}
	}

	async autoLayout(options?: import('./layout').LayoutOptions) {
		try {
			const auto = await autoLayout(this, options)
			
			this._undo.beforeChange()
			
			// Apply node positions
			auto.nodes.forEach(an => {
				const n = this.nodesMap.get(an.id)
				if (n) {
					this.moveNode(n, an.x, an.y, false, true) // Skip undo for individual moves
				}
			})
			
			// Apply edge routing from ELK layout
			auto.edges.forEach(ae => {
				const edge = this.edges.find(e => e.id == ae.id)
				if (edge) {
					// Clear existing vertices for this edge only
					if (edge.vertices) {
						edge.vertices.forEach(v => {
							if (v.id) {
								this.edgeVertices.delete(v.id)
							}
						})
					}
					edge.vertices = []
					edge.userDeletedVertices = false
					
					// Add routing vertices from ELK (these are proper bend points, not nodes)
					if (ae.vertices && ae.vertices.length > 0) {
						edge.vertices = ae.vertices.map(p => {
							const vertex = edge.initVertex(p)
							vertex.auto = true // Mark as auto-generated
							return vertex
						})
					}
					
					// Handle edge label positioning - create proper interactive label vertices
					if (ae.label) {
						
						// Remove any existing label vertices (ELK or user-created)
						if (edge.vertices) {
							edge.vertices.forEach(v => {
								if (v.label) {
									this.edgeVertices.delete(v.id)
								}
							})
							edge.vertices = edge.vertices.filter(v => !v.label)
						}
						
						// Create a proper label vertex that behaves like a user-created vertex
						const labelVertex = edge.initVertex(ae.label)
						labelVertex.label = true
						labelVertex.auto = true // Mark as auto-generated so it can be cleaned up
						
						// Insert label vertex at the optimal position in the routing path
						// Find the best position to insert it and project it onto that line segment
						edge.vertices = edge.vertices || []
						const insertPos = findOptimalLabelPosition(edge.vertices, ae.label, edge.from, edge.to)
						
						// Project the label position onto the line segment where it will be inserted
						const projectedPos = projectLabelOntoSegment(edge.vertices, ae.label, insertPos, edge.from, edge.to)
						labelVertex.x = projectedPos.x
						labelVertex.y = projectedPos.y
						
						edge.vertices.splice(insertPos, 0, labelVertex)
						this.edgeVertices.set(labelVertex.id, labelVertex)
						
					}
					
					// Redraw the edge with new routing
					this.redrawEdge(edge)
				}
			})
			
			// Fit the layout to the viewport with optimal positioning
			this.fitToView()
			
			this._undo.change()
			
		} catch (error) {
			console.error('Auto layout failed:', error)
			// Could show user notification here
		}
	}

	alignSelectionV() {
		const lst: Point[] = this.nodes().filter(n => n.selected)
		lst.push(...Array.from(this.edgeVertices.values()).filter(v => v.selected))
		let minY = Math.min(...lst.map(p => p.y))
		this.nodesMap.forEach(n => n.selected && this.moveNode(n, n.x, minY, false, false))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, v.x, minY, false, false))
	}

	alignSelectionH() {
		const lst: Point[] = this.nodes().filter(n => n.selected)
		lst.push(...Array.from(this.edgeVertices.values()).filter(v => v.selected))
		let minX = Math.min(...lst.map(p => p.x))
		this.nodesMap.forEach(n => n.selected && this.moveNode(n, minX, n.y, false, false))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, minX, v.y, false, false))
	}

	distributeSelectionH() {
		const selectedNodes = this.nodes().filter(n => n.selected)
		const selectedVertices = Array.from(this.edgeVertices.values()).filter(v => v.selected)
		
		if (selectedNodes.length + selectedVertices.length < 3) return // Need at least 3 elements to distribute
		
		this._undo.beforeChange()
		
		// Combine and sort by X coordinate
		const allElements = [...selectedNodes, ...selectedVertices]
		allElements.sort((a, b) => a.x - b.x)
		
		const minX = allElements[0].x
		const maxX = allElements[allElements.length - 1].x
		const spacing = (maxX - minX) / (allElements.length - 1)
		
		// Distribute elements evenly between leftmost and rightmost
		allElements.forEach((element, index) => {
			const newX = minX + (index * spacing)
			if ('title' in element) {
				// It's a Node
				this.moveNode(element as Node, newX, element.y, false, true)
			} else {
				// It's an EdgeVertex
				this.moveEdgeVertex(element as EdgeVertex, newX, element.y, false, true)
			}
		})
		
		this._undo.change()
	}

	distributeSelectionV() {
		const selectedNodes = this.nodes().filter(n => n.selected)
		const selectedVertices = Array.from(this.edgeVertices.values()).filter(v => v.selected)
		
		if (selectedNodes.length + selectedVertices.length < 3) return // Need at least 3 elements to distribute
		
		this._undo.beforeChange()
		
		// Combine and sort by Y coordinate
		const allElements = [...selectedNodes, ...selectedVertices]
		allElements.sort((a, b) => a.y - b.y)
		
		const minY = allElements[0].y
		const maxY = allElements[allElements.length - 1].y
		const spacing = (maxY - minY) / (allElements.length - 1)
		
		// Distribute elements evenly between topmost and bottommost
		allElements.forEach((element, index) => {
			const newY = minY + (index * spacing)
			if ('title' in element) {
				// It's a Node
				this.moveNode(element as Node, element.x, newY, false, true)
			} else {
				// It's an EdgeVertex
				this.moveEdgeVertex(element as EdgeVertex, element.x, newY, false, true)
			}
		})
		
		this._undo.change()
	}

	// Set edge selection state
	setEdgeSelected(edge: Edge, selected: boolean) {
		// Mark the edge as selected by selecting its connected nodes
		if (selected) {
			this.setNodeSelected(edge.from, true)
			this.setNodeSelected(edge.to, true)
		}
		// Update visual selection
		this.updateEdgesSel()
	}

	// Fit the entire graph to the current viewport
	fitToView() {
		const contentBounds = this.calculateContentBounds()
		
		// Handle edge case where there's no content
		if (contentBounds.width === 0 || contentBounds.height === 0) {
			return
		}
		
		// Get viewport dimensions
		const viewportWidth = svg.parentElement?.clientWidth || 800
		const viewportHeight = svg.parentElement?.clientHeight || 600
		
		// Add padding around content
		const padding = 40
		
		// Calculate zoom to fit content with padding
		const zoomX = (viewportWidth - padding * 2) / contentBounds.width
		const zoomY = (viewportHeight - padding * 2) / contentBounds.height
		const optimalZoom = Math.min(zoomX, zoomY)
		
		// Clamp zoom between reasonable bounds
		const finalZoom = Math.max(Math.min(optimalZoom, 2), 0.1)
		
		// Calculate content center in drawing coordinates
		const contentCenterX = contentBounds.x + contentBounds.width / 2
		const contentCenterY = contentBounds.y + contentBounds.height / 2
		
		// Calculate viewport center in screen coordinates
		const viewportCenterX = viewportWidth / 2
		const viewportCenterY = viewportHeight / 2
		
		// Calculate translation needed to center content in viewport
		// With translate(x,y) scale(zoom), translation is in screen coordinates
		const translateX = viewportCenterX - (contentCenterX * finalZoom)
		const translateY = viewportCenterY - (contentCenterY * finalZoom)
		
		// Apply zoom and translation transform
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (zoomGroup) {
			zoomGroup.setAttribute('transform', `translate(${translateX}, ${translateY}) scale(${finalZoom})`)
		}
		
		// Update panning
		updatePanning()
		
		// Save view state so this fit position is preserved after reload
		saveViewState(this.id)
	}

	// Save current layout state for restoration
	private saveLayoutState(): Layout {
		return this.exportLayout(true) // Include all vertices for complete state
	}

	// Restore layout state
	private restoreLayoutState(state: Layout) {
		this.importLayout(state, true) // Rerender after restoring
	}

	// Grid functionality
	isGridVisible(): boolean {
		return this._gridVisible;
	}

	isSnapToGrid(): boolean {
		return this._snapToGrid;
	}

	getGridSize(): number {
		return this._gridSize;
	}

	toggleGrid() {
		this._gridVisible = !this._gridVisible;
		this.updateGridDisplay();
		// Force toolbar update by dispatching a custom event
		window.dispatchEvent(new CustomEvent('gridStateChanged'));
	}

	toggleSnapToGrid() {
		this._snapToGrid = !this._snapToGrid;
		// Force toolbar update by dispatching a custom event
		window.dispatchEvent(new CustomEvent('gridStateChanged'));
	}

	snapAllToGrid() {
		if (!this._snapToGrid) return;
		
		this._undo.beforeChange();
		this.nodes().forEach(node => {
			const snappedX = Math.round(node.x / this._gridSize) * this._gridSize;
			const snappedY = Math.round(node.y / this._gridSize) * this._gridSize;
			this.moveNode(node, snappedX, snappedY, false, true);
		});
		this._undo.change();
	}

	// Helper method to snap a point to grid
	private snapToGrid(x: number, y: number): { x: number, y: number } {
		return {
			x: Math.round(x / this._gridSize) * this._gridSize,
			y: Math.round(y / this._gridSize) * this._gridSize
		};
	}

	updateGridDisplay() {
		if (!svg) return;

		// Remove existing grid pattern and background
		const existingGrid = svg.querySelector('#grid-pattern');
		if (existingGrid) {
			existingGrid.remove();
		}

		const existingGridRect = svg.querySelector('#grid-background');
		if (existingGridRect) {
			existingGridRect.remove();
		}

		if (!this._gridVisible) return;

		// Create grid pattern in defs
		let defs = svg.querySelector('defs');
		if (!defs) {
			defs = document.createElementNS('http://www.w3.org/2000/svg', 'defs');
			svg.insertBefore(defs, svg.firstChild);
		}

		const pattern = document.createElementNS('http://www.w3.org/2000/svg', 'pattern');
		pattern.id = 'grid-pattern';
		pattern.setAttribute('width', this._gridSize.toString());
		pattern.setAttribute('height', this._gridSize.toString());
		pattern.setAttribute('patternUnits', 'userSpaceOnUse');

		const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
		path.setAttribute('d', `M ${this._gridSize} 0 L 0 0 0 ${this._gridSize}`);
		path.setAttribute('fill', 'none');
		path.setAttribute('stroke', '#d0d0d0');
		path.setAttribute('stroke-width', '1');
		path.setAttribute('opacity', '0.8');

		pattern.appendChild(path);
		defs.appendChild(pattern);

		// Create grid background rectangle
		const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
		rect.id = 'grid-background';
		rect.setAttribute('x', '-10000');
		rect.setAttribute('y', '-10000');
		rect.setAttribute('width', '20000');
		rect.setAttribute('height', '20000');
		rect.setAttribute('fill', 'url(#grid-pattern)');
		rect.setAttribute('pointer-events', 'none');

		// Insert grid as first child of zoom group so it transforms with content
		const zoomGroup = svg.querySelector('g.zoom');
		if (zoomGroup) {
			zoomGroup.insertBefore(rect, zoomGroup.firstChild);
		}
	}
}

function escapeCdata(code: string) {
	return code.replace(/]]>/g, ']]]>]><![CDATA[')
}



let svg: SVGSVGElement = document.querySelector('svg#graph')
if (!svg) {
	svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
	svg.setAttribute('id', 'graph')
	svg.addEventListener('click', e => clickListener(e))
	// addCursorInteraction(svg) // Call will be updated in buildGraph
}
svg.setAttribute('width', '100%')
svg.setAttribute('height', '100%')

let clickListener: (e: MouseEvent) => void
let dragging = false;
let selectListener: (n: Node) => void


export const buildGraph = (data: GraphData, onNodeSelect: (n: Node) => void, dragMode: 'pan' | 'select') => {
	// empty svg
	svg.innerHTML = defs
	document.body.append(svg) // make sure svg element is connected, we will measure texts sizes
	// @ts-ignore
	svg.__data = data

	selectListener = onNodeSelect

	//use event delegation
	clickListener = e => {
		if (dragging) {
			return;
		}
		// const el = (e.target as any).closest('.node > .expand');
	}

	_buildGraph(data)
	const elasticEl = create.rect(300, 300, 50, 50, 0, 'elastic')
	svg.append(elasticEl)

	// Initialize grid display now that the zoom group exists
	data.updateGridDisplay()

	// Call addCursorInteraction with dragMode
	addCursorInteraction(svg, dragMode)

	return {
		svg,
		setZoom,
	}
}

export const buildGraphView = (data: GraphData) => {
	svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
	svg.setAttribute('id', 'graph')
	_buildGraph(data)
	return svg
}

const _buildGraph = (data: GraphData) => {
	//toplevel groups
	const zoomG = create.element('g', {}, 'zoom') as SVGGElement
	const nodesG = create.element('g', {}, 'nodes') as SVGGElement
	const edgesG = create.element('g', {}, 'edges') as SVGGElement
	const groupsG = create.element('g', {}, 'groups') as SVGGElement
	zoomG.append(groupsG, edgesG, nodesG)


	data.nodesMap.forEach((n) => {
		buildNode(n, data)
		nodesG.append(n.ref)
	})

	data.edges.forEach(e => {
		buildEdge(data, e)
		edgesG.append(e.ref)
	})

	data.groupsMap.forEach((group) => {
		buildGroup(group)
		groupsG.append(group.ref)
	})

	svg.append(zoomG)
}

function buildEdge(data: GraphData, edge: Edge) {
	const n1 = edge.from, n2 = edge.to;

	const g = create.element('g', {}, 'edge') as SVGGElement
	g.setAttribute('id', edge.id)
	g.setAttribute('data-from', edge.from.id)
	g.setAttribute('data-to', edge.to.id)

	const position = (edge.style.position || 50) / 100

	// Calculate edge vertices using utility function
	const vertices = calculateEdgeVertices(edge, data)

	// Calculate label position using utility function - pass edge for ELK label lookup
	const pLabel = calculateLabelPosition(vertices, position, n1, edge)

	const {bg, txt, bbox} = buildEdgeLabel(pLabel, edge)
	g.append(bg, txt)

	// Create edge segments and path using utility function
	const {segments, path} = createEdgeSegments(vertices, bbox, n1, n2)

	const p = create.path(path, {'marker-end': 'url(#arrow)'}, 'edge')
	p.setAttribute('fill', 'none')
	p.setAttribute('stroke', edge.style.color)
	p.setAttribute('stroke-width', String(edge.style.thickness))
	p.setAttribute('stroke-linecap', 'round')
	edge.style.dashed && p.setAttribute('stroke-dasharray', '8')
	g.append(p)
	
	// Debug visualization removed - arrow issue fixed

	// drag handlers
	edge.vertices = vertices.slice(1, -1).map(p => {
		// Preserve existing EdgeVertex objects to maintain IDs and selection state
		if ('id' in p && 'edge' in p) {
			// This is already an EdgeVertex, preserve it
			const v = p as EdgeVertex;
			v.edge = edge; // Ensure edge reference is correct
			return v;
		} else {
			// This is a new Point, convert to EdgeVertex
			return edge.initVertex(p);
		}
	})
	edge.vertices.forEach((p, i) => {
		const v = p as EdgeVertex
		v.ref = create.element('circle', {id: v.id, cx: p.x, cy: p.y, r: 7, fill: 'none'}, 'v-dot')
		v.selected && v.ref.classList.add('selected')
		v.auto && v.ref.classList.add('auto')
		g.append(v.ref)
	})

	edge.ref = g
	return g
}

function buildEdgeLabel(pLabel: Point, edge: Edge) {
	// label
	const fontSize = edge.style.fontSize
	let {txt, dy, maxW} = create.textArea(edge.label, 200, fontSize, false, pLabel.x, pLabel.y, 'middle')
	//move text up to center relative to the edge
	dy -= fontSize / 2
	txt.setAttribute('y', String(pLabel.y - dy / 2))

	maxW += fontSize
	applyStyle(txt, styles.edgeText)
	txt.setAttribute('stroke', 'none')
	txt.setAttribute('font-size', String(edge.style.fontSize))
	txt.setAttribute('fill', edge.style.color)

	const bbox = {x: pLabel.x - maxW / 2, y: pLabel.y - dy / 2, width: maxW, height: dy}
	const bg = create.rect(bbox.width, bbox.height, bbox.x, bbox.y)
	applyStyle(bg, styles.edgeRect)
	txt.setAttribute('data-field', 'label')

	bbox.x += bbox.width / 2
	bbox.y += bbox.height / 2
	return {bg, txt, bbox}
}


function buildNode(n: Node, data: GraphData) {
	// @ts-ignore
	window.gdata = data

	// Use the pre-calculated dimensions from addNode
	const w = n.width;
	const h = n.height;

	const g = create.element('g', {}, 'node') as SVGGElement
	g.setAttribute('id', n.id)
	n.selected && g.classList.add('selected')
	setPosition(g, n.x, n.y)

	// Ensure we use the correct shape from style, defaulting to Box
	const shapeType = n.style.shape || 'Box';
	const shapeFn = shapes[shapeType.toLowerCase()] || shapes.box
	const shape: SVGElement = shapeFn(g, n);

	shape.classList.add('nodeBorder')

	// Apply generic styles first
	applyStyle(shape, styles.nodeBorder)
	// Then apply tag-specific styles to override generic ones
	shape.setAttribute('fill', n.style.background)
	shape.setAttribute('stroke', n.style.stroke)
	// Consistent border width for all elements
	shape.setAttribute('stroke-width', '3')
	shape.setAttribute('opacity', String(n.style.opacity))
	setBorderStyle(shape, n.style.border)

	const tg = create.element('g') as SVGGElement
	let cy = Number(g.getAttribute('label-offset-y')) || 0
	
	// Consistent text width calculation for all shapes
	const textPadding = 12; // Standard padding
	const maxTextWidth = Math.max(w - (textPadding * 2), 80); // Ensure reasonable minimum width
	
	{
		const fontSize = n.style.fontSize
		// Use 95% of available width for titles - much more aggressive
		const {txt, dy} = create.textArea(n.title, maxTextWidth * 0.95, fontSize, true, 0, cy, 'middle')
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('data-field', 'name')

		tg.append(txt)
		cy += dy + 6 // Reduced spacing after title
	}
	{
		const txt = create.text(`[${n.sub}]`, {x: 0, y: cy, 'text-anchor': 'middle'})
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('font-size', String(0.75 * n.style.fontSize))
		tg.append(txt)
		cy += 12 // Reduced spacing
	}
	{
		cy += 6 // Reduced spacing before description
		const fontSize = Math.min(n.style.fontSize * 0.8, 16) // Keep smaller description text
		// Use 95% of available width for descriptions too - much more aggressive
		const {txt, dy} = create.textArea(n.description, maxTextWidth * 0.95, fontSize, false, 0, cy, 'middle')
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('data-field', 'description')
		tg.append(txt)
		cy += dy
	}

	// Better vertical centering with improved padding
	setPosition(tg, 0, -cy / 2)
	g.append(tg)

	// @ts-ignore
	g.__data = n;
	n.ref = g;

	return g
}


function buildGroup(group: Group) {
	if (group.nodes.length == 0) {
		return
	}
	const g = create.element('g', {}, 'group') as SVGGElement

	let p0: Point = {x: 1e100, y: 1e100}, p1: Point = {x: 0, y: 0}
	group.nodes.forEach(n => {
		const b = {x: n.x - n.width / 2, y: n.y - n.height / 2, width: n.width, height: n.height}
		p0.x = Math.min(p0.x, b.x)
		p0.y = Math.min(p0.y, b.y)
		p1.x = Math.max(p1.x, b.x + b.width)
		p1.y = Math.max(p1.y, b.y + b.height)
	})
	const pad = 25
	const w = Math.max(p1.x - p0.x, 200)
	const h = p1.y - p0.y + pad * 1.5
	const bb = {
		x: p0.x - pad,
		y: p0.y - pad,
		width: w + pad * 2,
		height: h + pad * 2,
	}
	const r = create.rect(bb.width, bb.height, bb.x, bb.y)
	group.x = bb.x + bb.width / 2
	group.y = bb.y + bb.height / 2
	group.width = bb.width
	group.height = bb.height
	applyStyle(r, styles.groupRect)
	group.style.stroke && r.setAttribute('stroke', group.style.stroke)
	group.style.background && r.setAttribute('fill', group.style.background)

	const txt = create.text(group.name, {x: p0.x, y: bb.y + bb.height - styles.groupText["font-size"]})
	applyStyle(txt, styles.groupText)
	group.style.color && txt.setAttribute('fill', group.style.color)

	g.append(r, txt)
	group.ref = g
}

function findClosestSegment(graph: GraphData, p: Point) {
	// find the closest point on a segment
	let fnd = {dst: Number.POSITIVE_INFINITY, pos: -1, edge: null as Edge, prj: null as Point}
	graph.edges.forEach(edge => {
		const vertices = edge.vertices || []
		const pts = [edge.from, ...vertices, edge.to]
		for (let i = 1; i < pts.length; i++) {
			const prj = project(p, pts[i - 1], pts[i])
			const dst = cabDistance(p, prj)
			if (dst > 50) continue
			if (dst < fnd.dst) {
				fnd = {dst, pos: i, prj, edge}
			}
		}
	})
	return fnd.edge ? fnd : null
}

function mouseToDrawing(e: MouseEvent): Point {
	// transform event coords to drawing coords
	const b = svg.getBoundingClientRect()
	const z = getZoom()
	
	// Get current pan transform
	const currentTransform = getCurrentTransform()
	
	// Convert screen coordinates to drawing coordinates accounting for zoom and pan
	// The transform values are in screen coordinates, so we need to divide by zoom and subtract
	return {
		x: (e.clientX - b.x - currentTransform.x) / z,
		y: (e.clientY - b.y - currentTransform.y) / z
	}
}

interface Handle extends Point {
	id: string
	selected?: boolean
	ref?: SVGElement
}

// Custom cursor interaction that prioritizes panning over selection
function addCustomCursorInteraction(svg: SVGSVGElement, conn: {
	nodeFromEvent(e: MouseEvent): Handle | null;
	setSelection(handles: Handle[]): void;
	setDragging(dragging: boolean): void;
	isSelected(handle: Handle): boolean;
	getSelection(): Handle[];
	getZoom(): number;
	moveNode(h: Handle, x: number, y: number): void;
	boxSelection(box: DOMRect, add: boolean): void;
	updatePanning(): void;
}, dragMode: 'pan' | 'select') {
	let ini: { x: number; y: number; n: Handle }[] = []
	let elastic: any = null
	let isPanning = false
	let panStartX = 0
	let panStartY = 0
	let initialTransform = { x: 0, y: 0 }
	let pendingSelectionChange: { node: Handle; shiftKey: boolean } | null = null
	let hasDragged = false
	
	// Store event listeners for cleanup
	const eventListeners: Array<{ element: Element | Window, event: string, handler: EventListener }> = []

	// Simple elastic selection box implementation - use mouseToDrawing for proper coordinate conversion
	function createElastic() {
		let startDrawingX = 0, startDrawingY = 0, rect: SVGRectElement | null = null
		
		return {
			ini(e: MouseEvent) {
				// Use mouseToDrawing to get proper drawing coordinates (accounts for zoom and pan)
				const pt = mouseToDrawing(e)
				startDrawingX = pt.x
				startDrawingY = pt.y
				
				rect = document.createElementNS("http://www.w3.org/2000/svg", "rect")
				rect.setAttribute('fill', 'rgba(0, 100, 255, 0.1)')
				rect.setAttribute('stroke', 'rgba(0, 100, 255, 0.5)')
				rect.setAttribute('stroke-width', '1')
				rect.setAttribute('stroke-dasharray', '3,3')
				rect.setAttribute('x', String(startDrawingX))
				rect.setAttribute('y', String(startDrawingY))
				rect.setAttribute('width', '0')
				rect.setAttribute('height', '0')
				
				// Add to the zoom group so it transforms with the content
				const zoomGroup = svg.querySelector('g.zoom')
				if (zoomGroup) {
					zoomGroup.appendChild(rect)
				} else {
					svg.appendChild(rect)
				}
			},
			update(e: MouseEvent) {
				if (!rect) return
				
				// Convert current mouse position to drawing coordinates (accounts for zoom and pan)
				const currentPt = mouseToDrawing(e)
				const currentDrawingX = currentPt.x
				const currentDrawingY = currentPt.y
				
				// Calculate rectangle bounds in drawing coordinates
				const x = Math.min(startDrawingX, currentDrawingX)
				const y = Math.min(startDrawingY, currentDrawingY)
				const width = Math.abs(currentDrawingX - startDrawingX)
				const height = Math.abs(currentDrawingY - startDrawingY)
				
				rect.setAttribute('x', String(x))
				rect.setAttribute('y', String(y))
				rect.setAttribute('width', String(width))
				rect.setAttribute('height', String(height))
			},
			end(): DOMRect | null {
				if (!rect) return null
				
				// Get final rectangle bounds in drawing coordinates
				const x = parseFloat(rect.getAttribute('x') || '0')
				const y = parseFloat(rect.getAttribute('y') || '0')
				const width = parseFloat(rect.getAttribute('width') || '0')
				const height = parseFloat(rect.getAttribute('height') || '0')
				
				rect.remove()
				rect = null
				
				// Return drawing coordinates directly for boxSelection since we're now working in the same coordinate system
				if (width > 5 && height > 5) {
					return {
						x: x, y: y, width: width, height: height,
						left: x, top: y, right: x + width, bottom: y + height
					} as DOMRect
				}
				return null
			}
		}
	}

	function getCurrentTransformLocal() {
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (!zoomGroup) return { x: 0, y: 0 }
		
		const transform = zoomGroup.getAttribute('transform') || ''
		const translateMatch = transform.match(/translate\(([^,]+),([^)]+)\)/)
		if (translateMatch) {
			return {
				x: parseFloat(translateMatch[1]) || 0,
				y: parseFloat(translateMatch[2]) || 0
			}
		}
		return { x: 0, y: 0 }
	}

	function setTransform(x: number, y: number) {
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (!zoomGroup) return
		
		const zoom = getZoom()
		zoomGroup.setAttribute('transform', `translate(${x}, ${y}) scale(${zoom})`)
	}

	function onMouseDown(e: MouseEvent) {
		e.preventDefault();
		hasDragged = false
		pendingSelectionChange = null

		const node = conn.nodeFromEvent(e)
		
		// Determine effective mode: invert if shift is held
		const effectiveMode = e.shiftKey ? (dragMode === 'pan' ? 'select' : 'pan') : dragMode
		
		if (!node) { // Clicked on empty space
			if (effectiveMode === 'pan') {
				// Pan mode: pan and deselect
				isPanning = true;
				elastic = null;
				panStartX = e.clientX;
				panStartY = e.clientY;
				initialTransform = getCurrentTransformLocal();
				ini = [];
				// Deselect all elements when clicking empty space in pan mode
				conn.setSelection([]);
			} else {
				// Select mode: ONLY select, no panning
				isPanning = false;
				elastic = createElastic();
				if (elastic) elastic.ini(e);
				ini = [];
			}
			return;
		}

		// Clicked on a node/vertex - behavior depends on effective mode
		if (effectiveMode === 'pan') {
			// Pan mode: select the element, don't pan
			isPanning = false;
			elastic = null;
			
			const nodes = conn.getSelection()
			if (conn.isSelected(node)) {
				// Clicking on a selected node - prepare to drag all selected elements
				ini = nodes.map(n => ({ x: n.x, y: n.y, n }))
				// No selection change needed since we're clicking on already selected element
			} else {
				// Clicking on an unselected node - select only this element
				conn.setSelection([node]);
				ini = [{ x: node.x, y: node.y, n: node }];
			}
		} else {
			// Select mode: selection/drag logic
			isPanning = false; // Ensure no panning if a node is clicked
			elastic = null; // Ensure no selection box if a node is clicked
			const nodes = conn.getSelection()
			
			if (e.shiftKey && dragMode === 'select') {
				// Shift+click in select mode: immediately change selection (no dragging expected)
				if (conn.isSelected(node)) {
					const index = nodes.findIndex(n => n.id === node.id)
					if (index >= 0) nodes.splice(index, 1)
				} else {
					nodes.push(node)
				}
				conn.setSelection(nodes)
				ini = nodes.map(n => ({ x: n.x, y: n.y, n }))
			} else {
				// Regular click: defer selection change until we know if it's a drag or click
				if (conn.isSelected(node)) {
					// Clicking on a selected node - prepare to drag all selected elements
					ini = nodes.map(n => ({ x: n.x, y: n.y, n }))
					// No pending selection change needed since we're clicking on already selected element
					pendingSelectionChange = null
				} else {
					// Clicking on an unselected node - defer selection change until we determine if it's a click or drag
					pendingSelectionChange = { node, shiftKey: e.shiftKey }
					// For now, prepare to drag just the clicked node (we'll update selection when drag starts)
					ini = [{ x: node.x, y: node.y, n: node }]
				}
			}
		}
	}

	function onMouseMove(e: MouseEvent, dx: number, dy: number) {
		// Check if we've moved enough to consider this a drag
		const dragThreshold = 3 // pixels
		if (!hasDragged && (Math.abs(dx) > dragThreshold || Math.abs(dy) > dragThreshold)) {
			hasDragged = true
			
			// If we have a pending selection change and we're now dragging, apply it
			if (pendingSelectionChange) {
				const nodes = conn.getSelection()
				nodes.length = 0
				nodes.push(pendingSelectionChange.node)
				conn.setSelection(nodes)
				// Update ini to drag only the newly selected node
				ini = [{ x: pendingSelectionChange.node.x, y: pendingSelectionChange.node.y, n: pendingSelectionChange.node }]
				pendingSelectionChange = null
			}
		}
		
		if (isPanning) {
			// Pan the view - apply mouse delta directly (no zoom division needed)
			// setTransform expects screen coordinates, dx/dy are already screen pixel deltas
			const newX = initialTransform.x + dx
			const newY = initialTransform.y + dy
			setTransform(newX, newY)
		} else if (ini.length > 0 && hasDragged) {
			// Move selected nodes/vertices (only if we've actually started dragging)
			// dx, dy are screen pixel deltas, convert to drawing coordinate deltas
			const zoom = conn.getZoom()
			const drawingDx = dx / zoom
			const drawingDy = dy / zoom
			ini.forEach(item => {
				// item.x, item.y are initial drawing coordinates
				// Add the drawing coordinate delta to get new position
				conn.moveNode(item.n, item.x + drawingDx, item.y + drawingDy)
			})
			conn.setDragging(true)
		} else if (elastic) {
			// Update selection box
			elastic.update(e)
			conn.setDragging(true)
		}
	}

	function onMouseUp(e: MouseEvent) {
		conn.setDragging(false)
		
		// If we have a pending selection change and didn't drag, apply it now (it was just a click)
		if (pendingSelectionChange && !hasDragged) {
			const nodes = conn.getSelection()
			nodes.length = 0
			nodes.push(pendingSelectionChange.node)
			conn.setSelection(nodes)
		}
		
		if (elastic) {
			const box = elastic.end()
			if (box) {
				conn.boxSelection(box, e.shiftKey)
			} else if (!ini.length) {
				// Deselect if no box was drawn
				conn.setSelection([])
			}
			elastic = null
		}
		
		// Save view state if user was panning (user-initiated view change)
		if (isPanning && hasDragged) {
			const graphData = (svg as any).__data as GraphData
			if (graphData && graphData.id) {
				saveViewState(graphData.id)
			}
		}
		
		// Reset state
		pendingSelectionChange = null
		hasDragged = false
		isPanning = false
		conn.updatePanning()
	}

	// Add drag and drop functionality
	function addDnd(element: SVGSVGElement) {
		let md: { ex: number; ey: number } | null = null

		function convertEvent(e: MouseEvent | TouchEvent): MouseEvent {
			if ('changedTouches' in e && e.changedTouches) {
				return e.changedTouches[0] as any
			}
			return e as MouseEvent
		}

		function onMouseMoveHandler(e: MouseEvent | TouchEvent) {
			if (!md) return
			e = convertEvent(e)
			onMouseMove(e, e.clientX - md.ex, e.clientY - md.ey)
		}

		function removeListeners() {
			document.removeEventListener('touchmove', onMouseMoveHandler as any)
			document.removeEventListener('mousemove', onMouseMoveHandler as any)
			document.removeEventListener('mouseup', onMouseUpHandler)
			document.removeEventListener('touchend', onMouseUpHandler)
		}

		function onMouseUpHandler(e: MouseEvent | TouchEvent) {
			removeListeners()
			onMouseUp(convertEvent(e))
			md = null
		}

		function onMouseDownHandler(e: MouseEvent | TouchEvent) {
			e = convertEvent(e)
			md = { ex: e.clientX, ey: e.clientY }
			onMouseDown(e)
			document.addEventListener('touchmove', onMouseMoveHandler as any)
			document.addEventListener('mousemove', onMouseMoveHandler as any)
			document.addEventListener('mouseup', onMouseUpHandler)
			document.addEventListener('touchend', onMouseUpHandler)
		}

		element.addEventListener('mousedown', onMouseDownHandler as any)
		element.addEventListener('touchstart', onMouseDownHandler as any)
		
		// Track these listeners for cleanup
		eventListeners.push(
			{ element, event: 'mousedown', handler: onMouseDownHandler as any },
			{ element, event: 'touchstart', handler: onMouseDownHandler as any }
		)
	}

	addDnd(svg)
	
	// Return cleanup function
	return () => {
		eventListeners.forEach(({ element, event, handler }) => {
			element.removeEventListener(event, handler)
		})
	}
}

export function addCursorInteraction(svg: SVGSVGElement, dragMode: 'pan' | 'select') {
	// Clean up any existing event listeners to prevent conflicts
	const existingCleanup = (svg as any).__cursorInteractionCleanup
	if (existingCleanup) {
		existingCleanup()
	}

	function getData(el: SVGElement) {
		// @ts-ignore
		return el.__data
	}

	const gd = () => (getData(svg) as GraphData)
	
	// Store event listeners for cleanup
	const eventListeners: Array<{ element: Element | Window, event: string, handler: EventListener }> = []

	const beforeUnloadHandler = (e: BeforeUnloadEvent) => {
		if (!gd().changed()) return
		e.preventDefault()
		e.returnValue = ''
	}
	window.addEventListener("beforeunload", beforeUnloadHandler)
	eventListeners.push({ element: window, event: 'beforeunload', handler: beforeUnloadHandler })

	function setDotSelected(d: Handle, selected: boolean) {
		d.selected = selected
		const dotEl = svg.querySelector('#' + d.id)
		d.selected ? dotEl.classList.add('selected') : dotEl.classList.remove('selected')
	}

	// show moving dot along edge when ALT is pressed
	const mouseMoveHandler = (e: MouseEvent) => {
		if (!e.altKey) return
		const fnd = findClosestSegment(gd(), mouseToDrawing(e))
		if (fnd) {
			const {prj} = fnd
			const parent = svg.querySelector('g.edges')
			let dot = parent.querySelector('#prj')
			if (!dot) {
				dot = create.element('circle', {id: 'prj', cx: prj.x, cy: prj.y, r: 7})
				parent.append(dot)
			}
			dot.setAttribute('cx', String(prj.x))
			dot.setAttribute('cy', String(prj.y))
		} else {
			removePrjDot()
		}
	}
	svg.addEventListener('mousemove', mouseMoveHandler)
	eventListeners.push({ element: svg, event: 'mousemove', handler: mouseMoveHandler })

	const keyUpHandler = (e: KeyboardEvent) => {
		const key = findShortcut(e, true)
		if (key == ADD_VERTEX || key == ADD_LABEL_VERTEX) return
		removePrjDot()
	}
	window.addEventListener('keyup', keyUpHandler)
	eventListeners.push({ element: window, event: 'keyup', handler: keyUpHandler })

	function removePrjDot() {
		const el = svg.querySelector('g.edges #prj')
		el && el.parentElement.removeChild(el)
	}

	const clickHandler = (e: MouseEvent) => {
		const key = findShortcut(e, true)
		if (key != ADD_LABEL_VERTEX && key != ADD_VERTEX) return
		const fnd = findClosestSegment(gd(), mouseToDrawing(e))
		if (fnd) {
			const {edge, pos, prj} = fnd
			// depending on keyboard modifier, make it label position
			gd().insertEdgeVertex(edge, prj, pos, key == ADD_LABEL_VERTEX)
			removePrjDot()
		}
	}
	svg.addEventListener('click', clickHandler)
	eventListeners.push({ element: svg, event: 'click', handler: clickHandler })

	const wheelHandler = (e: WheelEvent) => {
		// Handle wheel zoom directly without relying on shortcuts
		// deltaY > 0 means scrolling down (zoom out), deltaY < 0 means scrolling up (zoom in)
		const delta = Math.sign(e.deltaY) * 0.1 // Normalize to 0.1 zoom steps
		const currentZoom = getZoom()
		const newZoom = Math.max(0.1, Math.min(5, currentZoom - delta)) // Clamp zoom between 0.1 and 5
		
		if (newZoom !== currentZoom) {
			// Convert absolute screen coordinates to SVG-relative coordinates
			const rect = svg.getBoundingClientRect()
			const svgX = e.clientX - rect.left
			const svgY = e.clientY - rect.top
			setZoomCentered(newZoom, svgX, svgY)
			e.preventDefault()
			
			// Save view state after user wheel zoom
			const graphData = (svg as any).__data as GraphData
			if (graphData && graphData.id) {
				saveViewState(graphData.id)
			}
		}
	}
	svg.addEventListener('wheel', wheelHandler)
	eventListeners.push({ element: svg, event: 'wheel', handler: wheelHandler })

	const keyDownHandler = (e: KeyboardEvent) => {
		const shortcut = findShortcut(e)
		
		
		if (shortcut) {
			e.preventDefault() // Prevent browser default for all recognized shortcuts
		}
		
		switch (shortcut) {
			case DEL_VERTEX:
				const selectedVertices = Array.from(gd().edgeVertices.values()).filter(v => v.selected);
				selectedVertices.forEach(v => {
					gd().deleteEdgeVertex(v)
				})
				break
			case UNDO:
				gd().undo()
				break
			case REDO:
				gd().redo()
				break
			case ZOOM_IN:
				const newZoomIn = Math.min(5, getZoom() * 1.2)
				// Center zoom on viewport center like mouse wheel
				setZoomCentered(newZoomIn)
				saveViewState(gd().id) // Save after user keyboard zoom
				break
			case ZOOM_OUT:
				const newZoomOut = Math.max(0.1, getZoom() / 1.2)
				// Center zoom on viewport center like mouse wheel
				setZoomCentered(newZoomOut)
				saveViewState(gd().id) // Save after user keyboard zoom
				break
			case ZOOM_100:
				// Center zoom on viewport center like mouse wheel
				setZoomCentered(1)
				saveViewState(gd().id) // Save after user keyboard zoom
				break
			case ZOOM_FIT:
				gd().fitToView()
				// Don't save view state here - fitToView should not be persisted
				break
			case SELECT_ALL:
				gd().nodes().forEach(n => gd().setNodeSelected(n, true))
				gd().edgeVertices.forEach(v => setDotSelected(v, true))
				break
			case DESELECT:
				gd().nodes().forEach(n => gd().setNodeSelected(n, false))
				gd().edgeVertices.forEach(v => setDotSelected(v, false))
				break
		}
	}
	window.addEventListener('keydown', keyDownHandler)
	eventListeners.push({ element: window, event: 'keydown', handler: keyDownHandler })

	// Custom cursor interaction with pan-first behavior
	const customInteractionCleanup = addCustomCursorInteraction(svg, {
		nodeFromEvent(e: MouseEvent): Handle {
			e.preventDefault()
			// node clicked
			let el = (e.target as SVGElement).closest('g.nodes g.node') as SVGElement
			if (el) return getData(el)
			// vertex dot clicked
			el = (e.target as SVGElement).closest('g.edges g.edge .v-dot') as SVGElement
			if (el) {
				return gd().edgeVertices.get(el.id)
			}
			return null
		},
		setSelection(handles: Handle[]) {
			// nodes
			gd().nodes().forEach(n => gd().setNodeSelected(n, handles.some(h => h.id == n.id)))
			// dots
			gd().edgeVertices.forEach(d => setDotSelected(d, handles.some(h => h.id == d.id)))
			selectListener(gd().nodes().find(n => n.selected))
		},
		setDragging(d: boolean) {
			dragging = d
		},
		isSelected(handle: Handle): boolean {
			return handle.selected
		},
		getSelection(): Handle[] {
			const ret: Handle[] = gd().nodes().filter(n => n.selected)
			gd().edgeVertices.forEach(d => d.selected && ret.push(d))
			return ret
		},
		getZoom: getZoom,
		moveNode(h: Handle, x: number, y: number) {
			if (gd().nodesMap.has(h.id))
				gd().moveNode(h as Node, x, y)
			else {
				(h as EdgeVertex).auto = false
				gd().moveEdgeVertex(h as EdgeVertex, x, y)
			}
		},
		boxSelection(box: DOMRect, add) {
			// Box is now already in drawing coordinates, no need to scale
			// nodes
			gd().nodesMap.forEach(n => {
				const inBox = boxesOverlap(uncenterBox(n), box)
				if (inBox) {
					// Toggle selection for elements in the box
					gd().setNodeSelected(n, !n.selected)
				} else if (!add) {
					// If not holding shift and element is outside box, deselect it
					gd().setNodeSelected(n, false)
				}
			})
			// dots
			gd().edgeVertices.forEach(d => {
				const inBox = insideBox(d, box, false)
				if (inBox) {
					// Toggle selection for elements in the box
					setDotSelected(d, !d.selected)
				} else if (!add) {
					// If not holding shift and element is outside box, deselect it
					setDotSelected(d, false)
				}
			})

			selectListener(gd().nodes().find(n => n.selected))
		},
		updatePanning: updatePanning,
	}, dragMode)
	
	// Store cleanup function on the SVG element for later use
	const cleanup = () => {
		eventListeners.forEach(({ element, event, handler }) => {
			element.removeEventListener(event, handler)
		})
		if (customInteractionCleanup) {
			customInteractionCleanup()
		}
	}
	
	// Store cleanup function on SVG element
	;(svg as any).__cursorInteractionCleanup = cleanup
}

export function getZoom() {
	if (!svg) return 1
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (!el) return 1
	
	// Parse zoom from transform attribute to match how we set it
	const transform = el.getAttribute('transform') || ''
	const scaleMatch = transform.match(/scale\(([^)]+)\)/)
	if (scaleMatch) {
		return parseFloat(scaleMatch[1]) || 1
	}
	return 1
}

// svgPadding is now imported as SVG_PADDING from constants.ts

export function setZoom(zoom: number) {
	if (!svg) return
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (!el) return
	
	// Preserve existing translation when setting zoom
	const currentTransform = getCurrentTransform()
	el.setAttribute('transform', `translate(${currentTransform.x}, ${currentTransform.y}) scale(${zoom})`)
	
	// also set panning size
	updatePanning()
}

export function setZoomCentered(newZoom: number, centerX?: number, centerY?: number) {
	const el = svg.querySelector('g.zoom') as SVGGElement
	const oldZoom = getZoom()
	
	// If no center point provided, use viewport center
	if (centerX === undefined || centerY === undefined) {
		// Use the parent container's dimensions for the visible viewport
		// The SVG might be larger than the visible area due to overflow
		const container = svg.parentElement
		if (container) {
			centerX = container.clientWidth / 2
			centerY = container.clientHeight / 2
		} else {
			// Fallback to SVG dimensions if no parent
			centerX = svg.clientWidth / 2
			centerY = svg.clientHeight / 2
		}
	}
	
	// Get current transform
	const currentTransform = getCurrentTransform()
	
	// Convert screen coordinates to drawing coordinates
	// For transform order translate(tx, ty) scale(s):
	// screen_point = (drawing_point * scale) + translation
	// So: drawing_point = (screen_point - translation) / scale
	const drawingX = (centerX - currentTransform.x) / oldZoom
	const drawingY = (centerY - currentTransform.y) / oldZoom
	
	// Calculate new translation to keep the same drawing point at the same screen position
	// screen_point = (drawing_point * new_scale) + new_translation
	// So: new_translation = screen_point - (drawing_point * new_scale)
	const newTranslateX = centerX - (drawingX * newZoom)
	const newTranslateY = centerY - (drawingY * newZoom)
	
	// Apply the new transform
	el.setAttribute('transform', `translate(${newTranslateX}, ${newTranslateY}) scale(${newZoom})`)
	
	// Update panning
	updatePanning()
}

function getCurrentTransform() {
	if (!svg) return { x: 0, y: 0 }
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (!el) return { x: 0, y: 0 }
	
	const transform = el.getAttribute('transform') || ''
	const translateMatch = transform.match(/translate\(([^,]+),([^)]+)\)/)
	if (translateMatch) {
		return {
			x: parseFloat(translateMatch[1]) || 0,
			y: parseFloat(translateMatch[2]) || 0
		}
	}
	return { x: 0, y: 0 }
}

function updatePanning() {
	if (!svg) return
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (!el) return
	const bb = el.getBBox()
	const zoom = getZoom()
	if (!svg.parentElement) return
	const w = Math.max(svg.parentElement.clientWidth / zoom, bb.x + bb.width + SVG_PADDING)
	const h = Math.max(svg.parentElement.clientHeight / zoom, bb.y + bb.height + SVG_PADDING)
	svg.setAttribute('width', String(w * zoom))
	svg.setAttribute('height', String(h * zoom))
	
	// Note: View state saving removed from here to prevent interference with reset/fit functions
	// View state is now only saved on user interactions and page unload
}

// Optimized version that uses pre-calculated content bounds instead of expensive getBBox()
function updatePanningOptimized(graphData: GraphData) {
	const bb = graphData.calculateContentBounds() // Use already calculated bounds
	const zoom = getZoom()
	const w = Math.max(svg.parentElement.clientWidth / zoom, bb.x + bb.width + SVG_PADDING)
	const h = Math.max(svg.parentElement.clientHeight / zoom, bb.y + bb.height + SVG_PADDING)
	svg.setAttribute('width', String(w * zoom))
	svg.setAttribute('height', String(h * zoom))
}

export const getZoomAuto = () => {
	// Get the graph data to calculate proper content bounds
	const graphData = (svg as any).__data as GraphData
	if (!graphData) {
		return 1 // Default zoom if no graph data
	}
	
	// Use proper content bounds calculation
	const contentBounds = graphData.calculateContentBounds()
	const viewportWidth = svg.parentElement?.clientWidth || 800
	const viewportHeight = svg.parentElement?.clientHeight || 600
	
	// Add padding around content
	const padding = 40
	
	// Calculate zoom to fit content with padding
	const zoomX = (viewportWidth - padding * 2) / contentBounds.width
	const zoomY = (viewportHeight - padding * 2) / contentBounds.height
	const zoom = Math.min(zoomX, zoomY)
	
	// Clamp zoom between reasonable bounds
	return Math.max(Math.min(zoom, 2), 0.1)
}

const setBorderStyle = (el: SVGElement, style: string) => {
	if (style == 'Dashed') el.setAttribute('stroke-dasharray', '4')
	else if (style == 'Dotted') el.setAttribute('stroke-dasharray', '2')
}

const styles = {
	//node styles
	nodeBorder: {
		// Don't set fill and stroke here - let tag-specific styles handle colors
		filter: 'url(#shadow)',
	},
	nodeText: {
		'font-family': 'Arial, sans-serif',
		stroke: "none"
	},

	//edge styles
	edgeText: {
		'font-family': 'Arial, sans-serif',
		stroke: "none"
	},

	edgeRect: {
		fill: "none",
		stroke: "none",
	},

	//group styles
	groupRect: {
		//fill: "none",
		fill: "rgba(0, 0, 0, 0.02)",
		stroke: "#666",
		'stroke-width': 3,
		"stroke-dasharray": 4,
	},
	groupText: {
		fill: "#666",
		"font-size": 22,
		"font-weight": "bold",
		cursor: "default"
	}
}

// View state persistence - similar to undo cache but for zoom/pan
const viewStateCache = new Map<string, { zoom: number; transform: { x: number; y: number } }>();

// Save current view state (zoom and pan)
export function saveViewState(graphId: string) {
	if (!svg) return;
	
	const zoom = getZoom();
	const transform = getCurrentTransform();
	
	const state = {
		zoom,
		transform: { x: transform.x, y: transform.y }
	};
	
	viewStateCache.set(graphId, state);
}

// Restore view state if it exists
export function restoreViewState(graphId: string): boolean {
	if (!svg || !viewStateCache.has(graphId)) {
		return false;
	}
	
	const state = viewStateCache.get(graphId);
	if (!state) {
		return false;
	}
	
	// Restore zoom and transform
	const zoomGroup = svg.querySelector('g.zoom') as SVGGElement;
	if (zoomGroup) {
		zoomGroup.setAttribute('transform', `scale(${state.zoom}) translate(${state.transform.x}, ${state.transform.y})`);
		updatePanning();
	}
	
	return true;
}

// Clear view state for a graph
export function clearViewState(graphId: string) {
	viewStateCache.delete(graphId);
}

/**
 * Find the optimal position to insert a label vertex into the routing path
 * to minimize disruption to the existing route
 */
function findOptimalLabelPosition(vertices: Point[], labelPos: Point, fromNode: Point, toNode: Point): number {
	// If no existing vertices, insert at the beginning
	if (vertices.length === 0) {
		return 0;
	}
	
	// Build the full routing path including start/end nodes
	const fullPath = [fromNode, ...vertices, toNode];
	
	// Find the closest point on the path to the label position
	let minDistance = Infinity;
	let bestSegmentIndex = 0;
	
	for (let i = 0; i < fullPath.length - 1; i++) {
		const segmentStart = fullPath[i];
		const segmentEnd = fullPath[i + 1];
		
		// Calculate distance from label position to this segment
		const distance = distanceToSegment(labelPos, segmentStart, segmentEnd);
		
		if (distance < minDistance) {
			minDistance = distance;
			bestSegmentIndex = i;
		}
	}
	
	// Convert full path index to vertices array index
	// bestSegmentIndex 0 means between fromNode and vertices[0] -> insert at 0
	// bestSegmentIndex 1 means between vertices[0] and vertices[1] -> insert at 1
	// etc.
	return bestSegmentIndex;
}

/**
 * Calculate distance from a point to a line segment
 */
function distanceToSegment(point: Point, segmentStart: Point, segmentEnd: Point): number {
	const A = point.x - segmentStart.x;
	const B = point.y - segmentStart.y;
	const C = segmentEnd.x - segmentStart.x;
	const D = segmentEnd.y - segmentStart.y;
	
	const dot = A * C + B * D;
	const lenSq = C * C + D * D;
	
	if (lenSq === 0) {
		// Segment is actually a point
		return Math.sqrt(A * A + B * B);
	}
	
	let param = dot / lenSq;
	
	let xx, yy;
	
	if (param < 0) {
		xx = segmentStart.x;
		yy = segmentStart.y;
	} else if (param > 1) {
		xx = segmentEnd.x;
		yy = segmentEnd.y;
	} else {
		xx = segmentStart.x + param * C;
		yy = segmentStart.y + param * D;
	}
	
	const dx = point.x - xx;
	const dy = point.y - yy;
	return Math.sqrt(dx * dx + dy * dy);
}

/**
 * Project a label position onto the line segment where it will be inserted
 */
function projectLabelOntoSegment(vertices: Point[], labelPos: Point, insertPos: number, fromNode: Point, toNode: Point): Point {
	// Build the full routing path including start/end nodes
	const fullPath = [fromNode, ...vertices, toNode];
	
	// The segment where we're inserting is between fullPath[insertPos] and fullPath[insertPos + 1]
	const segmentStart = fullPath[insertPos];
	const segmentEnd = fullPath[insertPos + 1];
	
	// Project the label position onto this line segment
	return projectPointOntoSegment(labelPos, segmentStart, segmentEnd);
}

/**
 * Project a point onto a line segment (closest point on the segment)
 */
function projectPointOntoSegment(point: Point, segmentStart: Point, segmentEnd: Point): Point {
	const A = point.x - segmentStart.x;
	const B = point.y - segmentStart.y;
	const C = segmentEnd.x - segmentStart.x;
	const D = segmentEnd.y - segmentStart.y;
	
	const dot = A * C + B * D;
	const lenSq = C * C + D * D;
	
	if (lenSq === 0) {
		// Segment is actually a point, return that point
		return { x: segmentStart.x, y: segmentStart.y };
	}
	
	let param = dot / lenSq;
	
	// Clamp to segment (don't extend beyond endpoints)
	param = Math.max(0, Math.min(1, param));
	
	return {
		x: segmentStart.x + param * C,
		y: segmentStart.y + param * D
	};
}