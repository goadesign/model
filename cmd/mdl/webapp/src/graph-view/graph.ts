import {defs} from "./defs";
import {create, setPosition} from "./svg-create";
import {cursorInteraction} from "svg-editor-tools/lib/cursor-interaction";
import {shapes} from "./shapes";
import {intersectPolylineBox, Segment} from "./intersect";


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
	expanded?: boolean;
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
	vertices: Point[];
	ref?: SVGGElement;
	style: EdgeStyle;
}

export class GraphData {
	id: string;
	name: string;
	nodesMap: Map<string, Node>;
	edges: Edge[];
	groupsMap: Map<string, Group>;
	metadata: any;

	constructor(id?: string, name?: string) {
		this.id = id;
		this.name = name;

		this.edges = [];
		this.nodesMap = new Map;
		this.groupsMap = new Map;
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
		this.edges.push({
			id,
			from: this.nodesMap.get(fromNode),
			to: this.nodesMap.get(toNode),
			label,
			vertices,
			style: {...defaultEdgeStyle, ...style}
		})
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

	setExpanded(node: Node, ex: boolean) {
		node.expanded = ex;
		this.rebuildNode(node)
		updatePanning()
	}

	private rebuildNode(node: Node) {
		const p = node.ref.parentElement;
		p.removeChild(node.ref)
		node.ref = buildNode(node, this)
		p.appendChild(node.ref)
		this.redrawEdges(node)
		this.redrawGroups(node)
	}

	setSelected(nodes: Node[]) {
		this.nodesMap.forEach(n => {
			n.selected = false
			n.ref.classList.remove('selected')
		})
		nodes.forEach(n => {
			n.selected = true
			n.ref.classList.add('selected')
		});
		this.updateEdgesSel()
		//console.log(nodes.map(n => `'${n.name}'`).join(', '))
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
		n.x = x;
		n.y = y;
		setPosition(n.ref, x, y)
		this.redrawEdges(n);
		this.redrawGroups(n)
	}

	// moves the entire graph to be aligned top-left of the drawing area
	// used to bring back to visible the nodes that end up at negative coordinates
	alignTopLeft() {
		let minX: number = Math.min(...this.nodes().map(n => n.x)) - 250
		let minY: number = Math.min(...this.nodes().map(n => n.y)) - 250
		this.nodesMap.forEach(n => {
			this.moveNode(n, n.x - minX, n.y - minY)
		})
	}

	//redraw connected edges
	private redrawEdges(n: Node) {
		this.edges.forEach(e => {
			if (e.from == n) {
				const p = e.ref.parentElement;
				p.removeChild(e.ref)
				e.ref = buildEdge(this, e)
				p.append(e.ref)
			}
			if (e.to == n) {
				const p = e.ref.parentElement;
				p.removeChild(e.ref)
				e.ref = buildEdge(this, e)
				p.append(e.ref)
			}
		})
		this.updateEdgesSel()
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

	//call this from console: JSON.stringify(gdata.exportLayout())
	exportLayout() {
		return Array.from(this.nodesMap.values())
			.reduce<{ [key: string]: { x: number, y: number } }>(
				(o, n) => {
					o[n.id] = {x: n.x, y: n.y};
					return o
				},
				{}
			)
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

	importLayout(layout: { [key: string]: any }) {
		Object.entries(layout).forEach(([k, v]) => {
			const n = this.nodesMap.get(k)
			if (!n) return
			n.x = v.x
			n.y = v.y
			n.expanded = v.ex
		})
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
		let el = (e.target as any).closest('.node > .expand');
		if (el) {
			const n: Node = el.parentElement.__data
			data.setExpanded(n, !n.expanded)
			return
		}
	}

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

	const elasticEl = create.rect(300, 300, 50, 50, 0, 'elastic')
	svg.append(elasticEl)

	return {
		svg,
		setZoom,
	}
}

function buildEdge(data: GraphData, edge: Edge) {
	const n1 = edge.from, n2 = edge.to;

	const g = create.element('g', {}, 'edge') as SVGGElement
	g.setAttribute('id', edge.id)
	g.setAttribute('data-from', edge.from.id)
	g.setAttribute('data-to', edge.to.id)

	const position = (edge.style.position || 50) / 100

	// if vertices exists, follow them
	const vertices = edge.vertices ? edge.vertices.concat() : [];

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
			vertices.push({
				x: (n1.x + n2.x) / 2 + spreadX,
				y: (n1.y + n2.y) / 2 + spreadY
			})
		// only if no vertices and no splitting, obey routing style Orthogonal
		} else if (edge.style.routing == 'Orthogonal') {
			vertices.push({x: n1.x, y: n2.y})
		}
	}

	vertices.unshift(n1)
	vertices.push(n2)

	vertices[0] = n1.intersect(vertices[1])
	vertices[vertices.length - 1] = n2.intersect(vertices[vertices.length - 2])

	//where along the edge is the label?
	let iLabel: number // the segment index where we place the label
	let pLabel: Point // position of label
	function distance(p1: Point, p2: Point) {
		return Math.sqrt((p2.x - p1.x) * (p2.x - p1.x) + (p2.y - p1.y) * (p2.y - p1.y))
	}

	let sum = 0 // total length of the edge, sum of segments
	for (let i = 1; i < vertices.length; i++) {
		sum += distance(vertices[i - 1], vertices[i])
	}
	let acc = 0
	for (let i = 1; i < vertices.length; i++) {
		const d = distance(vertices[i - 1], vertices[i])
		if (acc + d > sum * position) {
			const pos = (sum * position - acc) / d
			pLabel = {
				x: vertices[i - 1].x + (vertices[i].x - vertices[i - 1].x) * pos,
				y: vertices[i - 1].y + (vertices[i].y - vertices[i - 1].y) * pos
			}
			iLabel = i
			break
		}
		acc += d
	}


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
	g.append(bg, txt)
	txt.setAttribute('data-field', 'label')

	bbox.x += bbox.width / 2
	bbox.y += bbox.height / 2

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
	edge.style.dashed && p.setAttribute('stroke-dasharray', '8')
	g.append(p)

	edge.ref = g
	return g
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

function addCursorInteraction(svg: SVGSVGElement) {

	function getData(el: SVGElement) {
		// @ts-ignore
		return el.__data
	}

	const gd = () => (getData(svg) as GraphData)

	cursorInteraction({
		svg: svg,
		nodeFromEvent(e: MouseEvent): Node {
			e.preventDefault()
			let el = (e.target as SVGElement).closest('g.nodes g.node') as SVGElement
			return el && getData(el)
		},
		setSelection(nodes: Node[]) {
			gd().setSelected(nodes)
			selectListener(nodes[0])
		},
		setDragging(d: boolean) {
			dragging = d
		},
		isSelected(node: Node): boolean {
			return node.selected
		},
		getSelection(): Node[] {
			return gd().nodes().filter(n => n.selected)
		},
		getZoom: getZoom,
		moveNode(n: Node, x: number, y: number) {
			gd().moveNode(n, x, y)
		},
		boxSelection(box: DOMRect, add) {
			gd().setSelected(gd().nodes().filter(n => {
				return (add && n.selected) || svg.checkIntersection(n.ref.firstChild as SVGElement, box)
			}))
			selectListener(gd().nodes().find(n => n.selected))
		},
		updatePanning: updatePanning,
	})
}

function getZoom() {
	const el = svg.querySelector('g.zoom') as SVGGElement
	if (el.transform.baseVal.numberOfItems == 0) return 1
	return el.transform.baseVal.getItem(0).matrix.a
}

const svgPadding = 20

function setZoom(zoom: number) {
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
		"stroke-width": "3px"
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