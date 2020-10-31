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
import {autoLayout} from "./dagre";
import {Undo} from "./undo";
import {
	ADD_LABEL_VERTEX,
	ADD_VERTEX,
	DEL_VERTEX,
	DESELECT,
	findShortcut,
	MOVE_DOWN,
	MOVE_DOWN_FAST,
	MOVE_LEFT,
	MOVE_LEFT_FAST,
	MOVE_RIGHT,
	MOVE_RIGHT_FAST,
	MOVE_UP,
	MOVE_UP_FAST,
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

interface Group extends BBox {
	id: string;
	name: string;
	nodes: (Node | Group)[];
	ref?: SVGGElement;
}

interface Node extends BBox {
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
	// Routing algorithm used to render lines.
	routing?: string
	// Position of annotation along the line; 0 (start) to 100 (end).
	position?: number
	// Opacity used to render line; 0-100.
	opacity?: number
}

const defaultEdgeStyle: EdgeStyle = {
	thickness: 3,
	color: '#999',
	opacity: 1,
	fontSize: 22,
	dashed: true,
}

const defaultNodeStyle: NodeStyle = {
	width: 300,
	height: 300,
	background: 'rgba(255, 255, 255, .9)',
	color: '#666',
	opacity: .9,
	stroke: '#999',
	fontSize: 22,
	shape: 'Rect'
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
		const n: Node = {
			id, title: label, sub, description, style: {...defaultNodeStyle, ...style},
			x: 0, y: 0, width: style.width, height: style.height, intersect: null
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

	addGroup(id: string, name: string, nodesOrGroups: string[]) {
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
			}).filter(Boolean)
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

	moveNode(n: Node, x: number, y: number) {
		if (n.x == x && n.y == y) return
		this._undo.beforeChange()
		n.x = x;
		n.y = y;
		setPosition(n.ref, x, y)
		this.redrawEdges(n);
		this.redrawGroups(n)
		this._undo.change()
	}

	moveEdgeVertex(v: EdgeVertex, x: number, y: number) {
		if (v.x == x && v.y == y) return
		this._undo.beforeChange()
		v.x = x;
		v.y = y;
		this.redrawEdge(v.edge)
		this._undo.change()
	}

	moveSelected(dx: number, dy: number) {
		this.nodes().forEach(n => n.selected && this.moveNode(n, n.x + dx, n.y + dy))
		this.edgeVertices.forEach(v => v.selected && this.moveEdgeVertex(v, v.x + dx, v.y + dy))
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
		const nodes: Point[] = this.nodes()
		const vertices = Array.from(this.edgeVertices.values())
		const all: Point[] = nodes.concat(vertices)
		let minX: number = Math.min(...all.map(n => n.x)) - 250
		let minY: number = Math.min(...all.map(n => n.y)) - 250
		this.nodesMap.forEach(n => this.moveNode(n, n.x - minX, n.y - minY))
		vertices.forEach(v => this.moveEdgeVertex(v, v.x - minX, v.y - minY))
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
		//save svg html
		let svg: SVGSVGElement = document.querySelector('svg#graph')
		const elastic = svg.querySelector('rect.elastic')
		const p = elastic.parentElement
		p.removeChild(elastic)
		const zoom = getZoom()
		setZoom(1)
		// inject metadata
		const script = document.createElement('script')
		script.setAttribute('type', 'application/json')
		this.metadata.layout = this.exportLayout()
		script.append('<![CDATA[' + escapeCdata(JSON.stringify(this.metadata, null, 2)) + ']]>')
		svg.insertBefore(script, svg.firstChild)
		// read the SVG
		let src = svg.outerHTML
		// restore all
		svg.removeChild(script)
		p.append(elastic)
		setZoom(zoom)
		return src.replace(/^<svg/, '<svg xmlns="http://www.w3.org/2000/svg"')
	}


	/**
	 * @param full when true, the edges without vertices are saved too, used for undo buffer
	 *        for saving, full is false
	 */
	exportLayout(full = false) {
		const ret: Layout = {}
		this.nodes().forEach(n => ret[n.id] = {x: n.x, y: n.y})
		this.edges.forEach(e => {
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

	autoLayout() {
		const auto = autoLayout(this)
		auto.nodes.forEach(an => {
			const n = this.nodesMap.get(an.id)
			this.moveNode(n, an.x, an.y)
		})
		this.edgeVertices.clear()
		auto.edges.forEach(ae => {
			const edge = this.edges.find(e => e.id == ae.id)
			const labelVertex = ae.vertices.find(v => ae.label.x == v.x && ae.label.y == v.y) as EdgeVertex
			labelVertex && (labelVertex.label = true)
			edge.vertices = ae.vertices.slice(1, -1).map(p => edge.initVertex(p))
			this.redrawEdge(edge)
		})
		updatePanning()
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
}

function escapeCdata(code: string) {
	return code.replace(/]]>/g, ']]]>]><![CDATA[')
}


let svg: SVGSVGElement = document.querySelector('svg#graph')
if (!svg) {
	svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
	svg.setAttribute('id', 'graph')
	svg.addEventListener('click', e => clickListener(e))
	addCursorInteraction(svg)
}
svg.setAttribute('width', '100%')
svg.setAttribute('height', '100%')

let clickListener: (e: MouseEvent) => void
let dragging = false;
let selectListener: (n: Node) => void


export const buildGraph = (data: GraphData, onNodeSelect: (n: Node) => void) => {
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
		// the expand button was clicked
		// let el = (e.target as any).closest('.node > .expand');
	}

	_buildGraph(data)
	const elasticEl = create.rect(300, 300, 50, 50, 0, 'elastic')
	svg.append(elasticEl)

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
			// only if no vertices and no splitting, obey routing style Orthogonal
		} else if (edge.style.routing == 'Orthogonal') {
			Math.abs(n2.x - n1.x) > Math.abs(n2.y - n1.y) ?
				vertices.push({x: n1.x, y: n2.y, auto: true} as Point) :
				vertices.push({x: n2.x, y: n1.y, auto: true} as Point)
			// vertices.push({x: (n1.x + n2.x)/2, y: n1.y, auto: true} as Point)
			// vertices.push({x: (n1.x + n2.x)/2, y: n2.y, auto: true} as Point)
			// vertices.push({x: n1.x, y: (n1.y + n2.y) /2, auto: true} as Point)
			// vertices.push({x: n2.x, y: (n1.y + n2.y) /2, auto: true} as Point)
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

	const path = segments.map(s => `M${s.p.x},${s.p.y} L${s.q.x},${s.q.y}`).join(' ')

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

	const w = n.style.width;//Math.max(60, textWidth(n.id), ...n.fields.map(f => textWidth(f.name))) + 70
	const h = n.style.height;
	n.width = w;
	n.height = h;

	const g = create.element('g', {}, 'node') as SVGGElement
	g.setAttribute('id', n.id)
	n.selected && g.classList.add('selected')
	setPosition(g, n.x, n.y)

	const shapeFn = shapes[n.style.shape.toLowerCase()] || shapes.box
	const shape: SVGElement = shapeFn(g, n);

	shape.classList.add('nodeBorder')

	applyStyle(shape, styles.nodeBorder)
	shape.setAttribute('fill', n.style.background)
	shape.setAttribute('stroke', n.style.stroke)
	shape.setAttribute('stroke-width', (n.width / 70).toFixed(1))
	shape.setAttribute('opacity', String(n.style.opacity))
	setBorderStyle(shape, n.style.border)

	const tg = create.element('g') as SVGGElement
	let cy = Number(g.getAttribute('label-offset-y')) || 0
	{
		const fontSize = n.style.fontSize
		const {txt, dy} = create.textArea(n.title, w - 40, fontSize, true, 0, cy, 'middle')
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('data-field', 'name')

		tg.append(txt)
		cy += dy
	}
	{
		const txt = create.text(`[${n.sub}]`, 0, cy, 'middle')
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('font-size', String(0.75 * n.style.fontSize))
		tg.append(txt)
		cy += 10
	}
	{
		cy += 10
		const fontSize = n.style.fontSize
		const {txt, dy} = create.textArea(n.description, w - 40, fontSize, false, 0, cy, 'middle')
		applyStyle(txt, styles.nodeText)
		txt.setAttribute('fill', n.style.color)
		txt.setAttribute('data-field', 'description')
		tg.append(txt)
		cy += dy
	}

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

	const txt = create.text(group.name, p0.x, bb.y + bb.height - styles.groupText["font-size"])
	applyStyle(txt, styles.groupText)

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
	return {x: (e.clientX - b.x) / z, y: (e.clientY - b.y) / z}
}

function addCursorInteraction(svg: SVGSVGElement) {

	function getData(el: SVGElement) {
		// @ts-ignore
		return el.__data
	}

	const gd = () => (getData(svg) as GraphData)

	interface Handle extends Point {
		id: string
		selected?: boolean
		ref?: SVGElement
	}

	window.addEventListener("beforeunload", e => {
		if (!gd().changed()) return
		e.preventDefault()
		e.returnValue = ''
	})

	function setDotSelected(d: Handle, selected: boolean) {
		d.selected = selected
		const dotEl = svg.querySelector('#' + d.id)
		d.selected ? dotEl.classList.add('selected') : dotEl.classList.remove('selected')
	}

	// show moving dot along edge when ALT is pressed
	svg.addEventListener('mousemove', e => {
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
	})

	window.addEventListener('keyup', e => {
		const key = findShortcut(e, true)
		if (key == ADD_VERTEX || key == ADD_LABEL_VERTEX) return
		removePrjDot()
	})

	function removePrjDot() {
		const el = svg.querySelector('g.edges #prj')
		el && el.parentElement.removeChild(el)
	}

	svg.addEventListener('click', e => {
		const key = findShortcut(e, true)
		if (key != ADD_LABEL_VERTEX && key != ADD_VERTEX) return
		const fnd = findClosestSegment(gd(), mouseToDrawing(e))
		if (fnd) {
			const {edge, pos, prj} = fnd
			// depending on keyboard modifier, make it label position
			gd().insertEdgeVertex(edge, prj, pos, key == ADD_LABEL_VERTEX)
			removePrjDot()
		}
	})

	svg.addEventListener('wheel', e => {
		const key = findShortcut(e, false, true)
		if (key == ZOOM_OUT || key == ZOOM_IN) {
			const delta = Math.round(e.deltaY / 10) / 100
			setZoom(getZoom() - delta)
			e.preventDefault()
		}
	})

	window.addEventListener('keydown', e => {
		switch (findShortcut(e)) {
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
				setZoom(getZoom() + .05)
				e.preventDefault()
				break
			case ZOOM_OUT:
				setZoom(getZoom() - .05)
				e.preventDefault()
				break
			case ZOOM_100:
				setZoom(1)
				e.preventDefault()
				break
			case ZOOM_FIT:
				gd().alignTopLeft()
				setZoom(getZoomAuto())
				e.preventDefault()
				break
			case SELECT_ALL:
				gd().nodes().forEach(n => gd().setNodeSelected(n, true))
				gd().edgeVertices.forEach(v => setDotSelected(v, true))
				e.preventDefault()
				break
			case DESELECT:
				gd().nodes().forEach(n => gd().setNodeSelected(n, false))
				gd().edgeVertices.forEach(v => setDotSelected(v, false))
				e.preventDefault()
				break
			case MOVE_LEFT:
				gd().moveSelected(-1, 0)
				break
			case MOVE_LEFT_FAST:
				gd().moveSelected(-10, 0)
				break
			case MOVE_RIGHT:
				gd().moveSelected(1, 0)
				break
			case MOVE_RIGHT_FAST:
				gd().moveSelected(10, 0)
				break
			case MOVE_UP:
				gd().moveSelected(0, -1)
				break
			case MOVE_UP_FAST:
				gd().moveSelected(0, -10)
				break
			case MOVE_DOWN:
				gd().moveSelected(0, 1)
				break
			case MOVE_DOWN_FAST:
				gd().moveSelected(0, 10)
				break
		}
	})

	cursorInteraction({
		svg: svg,
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
			const b = scaleBox(box, 1 / getZoom())
			// nodes
			gd().nodesMap.forEach(n => gd().setNodeSelected(n, (add && n.selected) || boxesOverlap(uncenterBox(n), b)))
			// dots
			gd().edgeVertices.forEach(d => setDotSelected(d, (add && d.selected) || insideBox(d, b, false)))

			selectListener(gd().nodes().find(n => n.selected))
		},
		updatePanning: updatePanning,
	})
}

export function getZoom() {
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (el.transform.baseVal.numberOfItems == 0) return 1
	return el.transform.baseVal.getItem(0).matrix.a
}

const svgPadding = 20

export function setZoom(zoom: number) {
	const el = svg.querySelector('g.zoom') as SVGGElement
	el.setAttribute('transform', `scale(${zoom})`)
	// also set panning size
	updatePanning()
}

function updatePanning() {
	const el = svg.querySelector('g.zoom') as SVGGElement
	const bb = el.getBBox()
	const zoom = getZoom()
	const w = Math.max(svg.parentElement.clientWidth / zoom, bb.x + bb.width + svgPadding)
	const h = Math.max(svg.parentElement.clientHeight / zoom, bb.y + bb.height + svgPadding)
	svg.setAttribute('width', String(w * zoom))
	svg.setAttribute('height', String(h * zoom))
}

export const getZoomAuto = () => {
	const el = svg.querySelector('g.zoom') as SVGGElement
	const bb = el.getBBox()
	const zoom = Math.min(
		(svg.parentElement.clientWidth - 20) / (bb.width + bb.x + svgPadding),
		(svg.parentElement.clientHeight - 20) / (bb.height + bb.y + svgPadding),
	)
	return Math.max(Math.min(zoom, 1), .2)
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

	edgeRect: {
		fill: "none",
		stroke: "none",
	},

	//group styles
	groupRect: {
		//fill: "none",
		fill: "rgba(0, 0, 0, 0.02)",
		stroke: "#666",
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