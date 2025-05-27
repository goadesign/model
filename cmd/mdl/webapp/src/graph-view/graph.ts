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


interface Point {
	x: number;
	y: number;
}

interface BBox extends Point {
	width: number;
	height: number;
}

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

export interface NodeStyle {
	// Width of element, in pixels.
	width?: number
	// Height of element, in pixels.
	height?: number
	// Background color of element as HTML RGB hex string (e.g. "#ffffff")
	background?: string
	// Stroke color of element as HTML RGB hex string (e.g. "#000000")
	stroke?: string
	// Foreground (text) color of element as HTML RGB hex string (e.g. "#ffffff")
	color?: string
	// Standard font size used to render text, in pixels.
	fontSize?: number
	// Shape used to render element.
	shape?: string
	// URL of PNG/JPG/GIF file or Base64 data URI representation.
	icon?: string
	// Type of border used to render element.
	border?: string
	// Opacity used to render element; 0-100.
	opacity?: number
	// Whether element metadata should be shown.
	metadata?: boolean
	// Whether element description should be shown.
	description?: boolean
}

// RelationshipStyle defines a relationship style.
export interface EdgeStyle {
	// Thickness of line, in pixels.
	thickness?: number
	// Color of line as HTML RGB hex string (e.g. "#ffffff").
	color?: string
	// Standard font size used to render relationship annotation, in pixels.
	fontSize?: number
	// Width of relationship annotation, in pixels.
	width?: number
	// Whether line is rendered dashed or not.
	dashed?: boolean
	// Position of annotation along the line; 0 (start) to 100 (end).
	position?: number
	// Opacity used to render line; 0-100.
	opacity?: number
	// Arrow style
	arrowStyle?: 'normal' | 'large' | 'small' | 'none'
}

const defaultEdgeStyle: EdgeStyle = {
	thickness: 3,
	color: '#999',
	opacity: 1,
	fontSize: 22,
	dashed: true,
}

const defaultNodeStyle: NodeStyle = {
	width: 280,
	height: 180,
	background: 'rgba(255, 255, 255, .9)',
	color: '#666',
	opacity: .9,
	stroke: '#999',
	fontSize: 22,
	shape: 'Rectangle'
}

interface Edge {
	id: string;
	label: string;
	from: Node;
	to: Node;
	vertices?: EdgeVertex[];
	ref?: SVGGElement;
	style: EdgeStyle;
	initVertex: (p: Point) => EdgeVertex
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
	[k: string]: Point | (Point & { label: boolean })[]
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
	private _gridSize: number = 20;

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
		const shape = style.shape || 'Rectangle';
		const isPersonShape = shape.toLowerCase() === 'person';
		const isCylinderShape = shape.toLowerCase() === 'cylinder';
		
		// Fixed dimensions based only on shape type - completely ignore style.width/height
		let width: number;
		let height: number;
		
		if (isPersonShape) {
			width = 200;
			height = 240;
		} else if (isCylinderShape) {
			width = 200;
			height = 160;
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
		// console.log(label, id, style, {...defaultNodeStyle, ...style})
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
			initVertex
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

	moveNode(n: Node, x: number, y: number, disableSnap: boolean = false) {
		if (!n) return
		
		// Apply snap-to-grid if enabled and not explicitly disabled
		if (this._snapToGrid && !disableSnap) {
			const snapped = this.snapToGrid(x, y);
			x = snapped.x;
			y = snapped.y;
		}
		
		if (n.x == x && n.y == y) return
		
		this._undo.beforeChange()
		n.x = x;
		n.y = y;
		setPosition(n.ref, x, y)
		this.redrawEdges(n);
		this.redrawGroups(n)
		this._undo.change()
	}

	moveEdgeVertex(v: EdgeVertex, x: number, y: number, disableSnap: boolean = false) {
		// Apply snap-to-grid if enabled and not explicitly disabled
		if (this._snapToGrid && !disableSnap) {
			const snapped = this.snapToGrid(x, y);
			x = snapped.x;
			y = snapped.y;
		}
		// Use exact coordinates (no rounding needed with modern grid system)
		
		if (v.x == x && v.y == y) return
		this._undo.beforeChange()
		v.x = x;
		v.y = y;
		this.redrawEdge(v.edge)
		this._undo.change()
	}

	moveSelected(dx: number, dy: number, disableSnap: boolean = false) {
		this.nodes().forEach(n => n.selected && this.moveNode(n, n.x + dx, n.y + dy, disableSnap))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, v.x + dx, v.y + dy, disableSnap))
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
		if (v.auto) {
			v.auto = true
		} else {
			v.edge.vertices.splice(v.edge.vertices.indexOf(v), 1)
			this.edgeVertices.delete(v.id)
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

	redo() {
		this._undo.redo()
	}

	// moves the entire graph to be aligned top-left of the drawing area
	// used to bring back to visible the nodes that end up at negative coordinates
	alignTopLeft() {
		const contentBounds = this.calculateContentBounds()
		const padding = 100 // Reasonable padding for viewport
		
		const offsetX = -contentBounds.x + padding
		const offsetY = -contentBounds.y + padding
		
		this._undo.beforeChange()
		
		this.nodesMap.forEach(node => {
			this.moveNode(node, node.x + offsetX, node.y + offsetY)
		})
		
		this.edgeVertices.forEach(vertex => {
			this.moveEdgeVertex(vertex, vertex.x + offsetX, vertex.y + offsetY)
		})
		
		this._undo.change()
		
		// Clear view state so this reset position is not overridden
		clearViewState(this.id)
	}
	
	// Reset pan transform to (0,0) while preserving zoom
	resetPanTransform() {
		const currentZoom = getZoom()
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (zoomGroup) {
			zoomGroup.setAttribute('transform', `scale(${currentZoom}) translate(0, 0)`)
			updatePanning()
		}
		
		// Clear view state so this reset is not overridden
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
		const exportBounds = {
			x: contentBounds.x - padding,
			y: contentBounds.y - padding,
			width: contentBounds.width + (padding * 2),
			height: contentBounds.height + (padding * 2)
		}
		
		// Calculate offset for export positioning
		const offsetX = -contentBounds.x + padding
		const offsetY = -contentBounds.y + padding
		
		// Apply export positioning to the cloned SVG elements
		const exportZoomGroup = exportSvg.querySelector('g.zoom') as SVGGElement
		if (exportZoomGroup) {
			// Reset zoom to 1 and apply offset transform
			exportZoomGroup.setAttribute('transform', `scale(1) translate(${offsetX}, ${offsetY})`)
		}
		
		// Set proper viewBox and dimensions for export
		exportSvg.setAttribute('viewBox', `${exportBounds.x} ${exportBounds.y} ${exportBounds.width} ${exportBounds.height}`)
		exportSvg.setAttribute('width', String(exportBounds.width))
		exportSvg.setAttribute('height', String(exportBounds.height))
		
		// Inject metadata with current layout
		const script = document.createElement('script')
		script.setAttribute('type', 'application/json')
		this.metadata.layout = this.exportLayout()
		script.append('<![CDATA[' + escapeCdata(JSON.stringify(this.metadata, null, 2)) + ']]>')
		exportSvg.insertBefore(script, exportSvg.firstChild)
		
		// Get the export SVG as string
		const src = exportSvg.outerHTML.replace(/^<svg/, '<svg xmlns="http://www.w3.org/2000/svg"')
		
		// No restoration needed since we never touched the original SVG!
		return src
	}

	// Calculate the actual bounds of all content including nodes, edges, and groups
	calculateContentBounds(): BBox {
		const elements: BBox[] = []
		
		// Add node bounds (including their actual dimensions)
		this.nodes().forEach(node => {
			elements.push({
				x: node.x - node.width / 2,
				y: node.y - node.height / 2,
				width: node.width,
				height: node.height
			})
		})
		
		// Add edge vertex bounds
		this.edgeVertices.forEach(vertex => {
			elements.push({
				x: vertex.x - 5, // Small padding for vertex dots
				y: vertex.y - 5,
				width: 10,
				height: 10
			})
		})
		
		// Add group bounds
		this.groupsMap.forEach(group => {
			elements.push({
				x: group.x - group.width / 2,
				y: group.y - group.height / 2,
				width: group.width,
				height: group.height
			})
		})
		
		// Add edge path bounds (approximate)
		this.edges.forEach(edge => {
			// Add edge endpoints
			elements.push({
				x: edge.from.x - 10,
				y: edge.from.y - 10,
				width: 20,
				height: 20
			})
			elements.push({
				x: edge.to.x - 10,
				y: edge.to.y - 10,
				width: 20,
				height: 20
			})
			
			// Add edge vertices
			if (edge.vertices) {
				edge.vertices.forEach(vertex => {
					elements.push({
						x: vertex.x - 10,
						y: vertex.y - 10,
						width: 20,
						height: 20
					})
				})
			}
		})
		
		if (elements.length === 0) {
			return { x: 0, y: 0, width: 100, height: 100 }
		}
		
		// Calculate overall bounds
		const minX = Math.min(...elements.map(e => e.x))
		const minY = Math.min(...elements.map(e => e.y))
		const maxX = Math.max(...elements.map(e => e.x + e.width))
		const maxY = Math.max(...elements.map(e => e.y + e.height))
		
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
			const lst = e.vertices.filter(v => !v.auto).map(v => ({x: v.x, y: v.y, label: v.label}));
			(lst.length || full) && (ret[`e-${e.id}`] = lst)
		})
		return ret
	}

	setSaved() {
		this._undo.setSaved()
	}

	importLayout(layout: { [key: string]: any }, rerender = false) {
		Object.entries(layout).forEach(([k, v]) => {
			// nodes
			const n = this.nodesMap.get(k)
			if (n) {
				n.x = v.x
				n.y = v.y
			} else
				// edge vertices
			if (k.startsWith('e-')) {
				const edge = this.edges.find(e => e.id == k.slice(2))
				if (!edge) return;
				edge.vertices && edge.vertices.forEach(v => this.edgeVertices.delete(v.id))
				edge.vertices = v.map((p: Point) => edge.initVertex(p))
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
			
			// Apply node positions
			auto.nodes.forEach(an => {
				const n = this.nodesMap.get(an.id)
				if (n) {
					this.moveNode(n, an.x, an.y)
				}
			})
			
			// Clear existing edge vertices before applying new layout
			this.edgeVertices.clear()
			
			// Apply edge routing from ELK layout
			auto.edges.forEach(ae => {
				const edge = this.edges.find(e => e.id == ae.id)
				if (edge) {
					// Clear existing vertices
					edge.vertices = []
					
					// Add routing vertices from ELK (these are proper bend points, not nodes)
					if (ae.vertices && ae.vertices.length > 0) {
						edge.vertices = ae.vertices.map(p => {
							const vertex = edge.initVertex(p)
							vertex.auto = true // Mark as auto-generated
							return vertex
						})
					}
					
					// Handle edge label positioning
					if (ae.label) {
						const labelVertex = edge.initVertex(ae.label)
						labelVertex.label = true
						labelVertex.auto = true
						edge.vertices.push(labelVertex)
					}
					
					// Redraw the edge with new routing
					this.redrawEdge(edge)
				}
			})
			
			// Fit the layout to the viewport with optimal positioning
			this.fitToView()
			
		} catch (error) {
			console.error('Auto layout failed:', error)
			// Could show user notification here
		}
	}

	alignSelectionV() {
		const lst: Point[] = this.nodes().filter(n => n.selected)
		lst.push(...Array.from(this.edgeVertices.values()).filter(v => v.selected))
		let minY = Math.min(...lst.map(p => p.y))
		this.nodesMap.forEach(n => n.selected && this.moveNode(n, n.x, minY))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, v.x, minY))
	}

	alignSelectionH() {
		const lst: Point[] = this.nodes().filter(n => n.selected)
		lst.push(...Array.from(this.edgeVertices.values()).filter(v => v.selected))
		let minX = Math.min(...lst.map(p => p.x))
		this.nodesMap.forEach(n => n.selected && this.moveNode(n, minX, n.y))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, minX, v.y))
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
				this.moveNode(element as Node, newX, element.y)
			} else {
				// It's an EdgeVertex
				this.moveEdgeVertex(element as EdgeVertex, newX, element.y)
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
				this.moveNode(element as Node, element.x, newY)
			} else {
				// It's an EdgeVertex
				this.moveEdgeVertex(element as EdgeVertex, element.x, newY)
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
		
		// Calculate viewport center in drawing coordinates (after zoom)
		const viewportCenterX = viewportWidth / (2 * finalZoom)
		const viewportCenterY = viewportHeight / (2 * finalZoom)
		
		// Calculate translation needed to center content in viewport
		const translateX = viewportCenterX - contentCenterX
		const translateY = viewportCenterY - contentCenterY
		
		// Apply zoom and translation transform
		const zoomGroup = svg.querySelector('g.zoom') as SVGGElement
		if (zoomGroup) {
			zoomGroup.setAttribute('transform', `scale(${finalZoom}) translate(${translateX}, ${translateY})`)
		}
		
		// Update panning
		updatePanning()
		
		// Clear view state so this fit is not overridden
		clearViewState(this.id)
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
			this.moveNode(node, snappedX, snappedY);
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
		path.setAttribute('stroke', '#e0e0e0');
		path.setAttribute('stroke-width', '1');
		path.setAttribute('opacity', '0.5');

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

	// if vertices exists, follow them
	let vertices: Point[] = edge.vertices ? edge.vertices.concat() : [];
	// remove auto vertices, they will be regenerated
	const tmp = (vertices as EdgeVertex[]);
	tmp.forEach(v => v.auto && data.edgeVertices.delete(v.id))
	vertices = tmp.filter(v => !v.auto)

	if (vertices.length == 0) {
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

	// where along the edge is the label?
	// position of label
	let pLabel: Point = vertices.find(v => (v as EdgeVertex).label)
	if (!pLabel) {
		const distance = (p1: Point, p2: Point) =>
			Math.sqrt((p2.x - p1.x) * (p2.x - p1.x) + (p2.y - p1.y) * (p2.y - p1.y))

		let sum = 0 // total length of the edge, sum of segments
		for (let i = 1; i < vertices.length; i++) {
			sum += distance(vertices[i - 1], vertices[i])
		}
		pLabel = {x: n1.x, y: n1.y} // fallback for corner cases
		let acc = 0
		for (let i = 1; i < vertices.length; i++) {
			const d = distance(vertices[i - 1], vertices[i])
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

	const {bg, txt, bbox} = buildEdgeLabel(pLabel, edge)
	g.append(bg, txt)

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

	const p = create.path(path, {'marker-end': 'url(#arrow)'}, 'edge')
	p.setAttribute('fill', 'none')
	p.setAttribute('stroke', edge.style.color)
	p.setAttribute('stroke-width', String(edge.style.thickness))
	p.setAttribute('stroke-linecap', 'round')
	edge.style.dashed && p.setAttribute('stroke-dasharray', '8')
	g.append(p)

	// drag handlers
	edge.vertices = vertices.slice(1, -1).map(p => edge.initVertex(p))
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

	// Ensure we use the correct shape from style, defaulting to Rectangle
	const shapeType = n.style.shape || 'Rectangle';
	const shapeFn = shapes[shapeType.toLowerCase()] || shapes.rectangle
	const shape: SVGElement = shapeFn(g, n);

	shape.classList.add('nodeBorder')

	applyStyle(shape, styles.nodeBorder)
	shape.setAttribute('fill', n.style.background)
	shape.setAttribute('stroke', n.style.stroke)
	// Consistent border width for all elements
	shape.setAttribute('stroke-width', '3')
	shape.setAttribute('opacity', String(n.style.opacity))
	setBorderStyle(shape, n.style.border)

	const tg = create.element('g') as SVGGElement
	let cy = Number(g.getAttribute('label-offset-y')) || 0
	
	// Optimized padding strategy - use minimal padding but ensure text fits
	const textPadding = 12; // Much smaller padding for better space utilization
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
	return {
		x: (e.clientX - b.x) / z - currentTransform.x,
		y: (e.clientY - b.y) / z - currentTransform.y
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
			update(dx: number, dy: number) {
				if (!rect) return
				const zoom = conn.getZoom()
				
				// Calculate current drawing position by adding scaled delta to start position
				const currentDrawingX = startDrawingX + dx / zoom
				const currentDrawingY = startDrawingY + dy / zoom
				
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
		zoomGroup.setAttribute('transform', `scale(${zoom}) translate(${x}, ${y})`)
	}

	function onMouseDown(e: MouseEvent) {
		e.preventDefault(); 
		hasDragged = false
		pendingSelectionChange = null

		const node = conn.nodeFromEvent(e)
		
		if (!node) { // Clicked on empty space
			if (dragMode === 'pan') {
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

		// Clicked on a node/vertex - selection/drag logic (works in both modes)
		isPanning = false; // Ensure no panning if a node is clicked
		elastic = null; // Ensure no selection box if a node is clicked
		const nodes = conn.getSelection()
		
		if (e.shiftKey) {
			// Shift+click: immediately change selection (no dragging expected)
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
			if (conn.isSelected(node) && nodes.length > 1) {
				// Clicking on a selected node in a multi-selection - prepare to drag all
				ini = nodes.map(n => ({ x: n.x, y: n.y, n }))
				// Don't change selection yet - wait to see if user drags
			} else if (!conn.isSelected(node)) {
				// Clicking on an unselected node - defer selection change
				pendingSelectionChange = { node, shiftKey: e.shiftKey }
				// For now, prepare to drag just this node
				ini = [{ x: node.x, y: node.y, n: node }]
			} else {
				// Clicking on the only selected node - prepare to drag it
				ini = [{ x: node.x, y: node.y, n: node }]
			}
		}
	}

	function onMouseMove(dx: number, dy: number) {
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
				ini = [{ x: pendingSelectionChange.node.x, y: pendingSelectionChange.node.y, n: pendingSelectionChange.node }]
				pendingSelectionChange = null
			}
		}
		
		if (isPanning) {
			// Pan the view - scale mouse delta by zoom factor
			const zoom = conn.getZoom()
			const newX = initialTransform.x + dx / zoom
			const newY = initialTransform.y + dy / zoom
			setTransform(newX, newY)
		} else if (ini.length > 0 && hasDragged) {
			// Move selected nodes/vertices (only if we've actually started dragging)
			ini.forEach(item => {
				const sc = conn.getZoom()
				conn.moveNode(item.n, item.x + dx / sc, item.y + dy / sc)
			})
			conn.setDragging(true)
		} else if (elastic) {
			// Update selection box
			elastic.update(dx, dy)
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
			onMouseMove(e.clientX - md.ex, e.clientY - md.ey)
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
			setZoomCentered(newZoom, e.clientX, e.clientY)
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
				Array.from(gd().edgeVertices.values()).filter(v => v.selected).forEach(v => {
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
				const newZoomIn = Math.min(5, getZoom() + .05)
				setZoomCentered(newZoomIn)
				saveViewState(gd().id) // Save after user keyboard zoom
				break
			case ZOOM_OUT:
				const newZoomOut = Math.max(0.1, getZoom() - .05)
				setZoomCentered(newZoomOut)
				saveViewState(gd().id) // Save after user keyboard zoom
				break
			case ZOOM_100:
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
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (el.transform.baseVal.numberOfItems == 0) return 1
	return el.transform.baseVal.getItem(0).matrix.a
}

const svgPadding = 20

export function setZoom(zoom: number) {
	const el = svg.querySelector('g.zoom') as SVGGElement
	
	// Preserve existing translation when setting zoom
	const currentTransform = getCurrentTransform()
	el.setAttribute('transform', `scale(${zoom}) translate(${currentTransform.x}, ${currentTransform.y})`)
	
	// also set panning size
	updatePanning()
}

export function setZoomCentered(newZoom: number, centerX?: number, centerY?: number) {
	const el = svg.querySelector('g.zoom') as SVGGElement
	const oldZoom = getZoom()
	
	// If no center point provided, use viewport center
	if (centerX === undefined || centerY === undefined) {
		const rect = svg.getBoundingClientRect()
		centerX = rect.width / 2
		centerY = rect.height / 2
	}
	
	// Get current transform
	const currentTransform = getCurrentTransform()
	
	// Calculate the point in the drawing coordinate system
	const drawingX = (centerX / oldZoom) - currentTransform.x
	const drawingY = (centerY / oldZoom) - currentTransform.y
	
	// Calculate new translation to keep the same point under the cursor
	const newTranslateX = (centerX / newZoom) - drawingX
	const newTranslateY = (centerY / newZoom) - drawingY
	
	// Apply the new transform
	el.setAttribute('transform', `scale(${newZoom}) translate(${newTranslateX}, ${newTranslateY})`)
	
	// Update panning
	updatePanning()
}

function getCurrentTransform() {
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
	const el = svg.querySelector('g.zoom') as SVGGElement
	const bb = el.getBBox()
	const zoom = getZoom()
	const w = Math.max(svg.parentElement.clientWidth / zoom, bb.x + bb.width + svgPadding)
	const h = Math.max(svg.parentElement.clientHeight / zoom, bb.y + bb.height + svgPadding)
	svg.setAttribute('width', String(w * zoom))
	svg.setAttribute('height', String(h * zoom))
	
	// Note: View state saving removed from here to prevent interference with reset/fit functions
	// View state is now only saved on user interactions and page unload
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
		fill: "rgba(255, 255, 255, 0.86)",
		stroke: "#aaa",
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
		cursor: "default"
	}
}

const applyStyle = (el: SVGElement, style: { [key: string]: string | number }) => {
	Object.keys(style).forEach(name => {
		if (name == 'font-size') {
			if (typeof (style[name]) != 'number') {
				console.error(`All font-sizes in styles have to be numbers representing px! Found:`, style)
			}
			el.setAttribute(name, style[name] + 'px')
		} else {
			el.setAttribute(name, String(style[name]))
		}
	})
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