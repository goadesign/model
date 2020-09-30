interface Point {
	x: number;
	y: number;
}

interface BBox extends Point{
	width: number;
	height: number;
}

interface D3Node extends BBox {
	intersect: (p: Point) => Point
}

class D3Element {
	private readonly _el: SVGElement

	constructor(el: SVGElement) {
		this._el = el;
	}

	node() {
		return this._el;
	}

	attr(name: string, value: string | number) {
		this._el.setAttribute(name, String(value))
		return this;
	}

	insert(type: string, pos: string) {
		const el = document.createElementNS('http://www.w3.org/2000/svg', type)
		const el2 = this._el.insertBefore(el, this._el.querySelector(pos))
		return new D3Element(el2)
	}
}

export function intersectRect(node: BBox, point: Point) {
	const x = node.x;
	const y = node.y;

	// Rectangle intersection algorithm from:
	// http://math.stackexchange.com/questions/108113/find-edge-between-two-boxes
	const dx = point.x - x;
	const dy = point.y - y;
	let w = node.width / 2;
	let h = node.height / 2;

	let sx, sy;
	if (Math.abs(dy) * w > Math.abs(dx) * h) {
		// Intersection is top or bottom of rect.
		if (dy < 0) {
			h = -h;
		}
		sx = dy === 0 ? 0 : h * dx / dy;
		sy = h;
	} else {
		// Intersection is left or right of rect.
		if (dx < 0) {
			w = -w;
		}
		sx = w;
		sy = dx === 0 ? 0 : w * dy / dx;
	}

	return {x: x + sx, y: y + sy};
}

function rect(parent: D3Element, bbox: BBox, node: D3Node) {
	const shapeSvg = parent.insert("rect", ":first-child")
		.attr("rx", 3)
		.attr("ry", 3)
		.attr("x", -bbox.width / 2)
		.attr("y", -bbox.height / 2)
		.attr("width", bbox.width)
		.attr("height", bbox.height);

	node.intersect = function (point) {
		return intersectRect(node, point);
	};

	return shapeSvg;
}

function intersectEllipse(ellCenter: Point, rx: number, ry: number, nodeCenter: Point, point: Point) {

	//translate all to center ellipse
	const p1 = {x: point.x - ellCenter.x, y: point.y - ellCenter.y}
	const p2 = {x: nodeCenter.x - ellCenter.x, y: nodeCenter.y - ellCenter.y}

	if (p2.x == p1.x) { //hack to avoid singularity
		p1.x += .0000001
	}

	const s = (p2.y - p1.y) / (p2.x - p1.x);
	const si = p2.y - (s * p2.x);
	const a = (ry * ry) + (rx * rx * s * s);
	const b = 2 * rx * rx * si * s;
	const c = rx * rx * si * si - rx * rx * ry * ry;

	const radicand_sqrt = Math.sqrt((b * b) - (4 * a * c));
	const x = p1.x > p2.x ?
		(-b + radicand_sqrt) / (2 * a) :
		(-b - radicand_sqrt) / (2 * a)
	const pos = {
		x: x,
		y: s * x + si
	}
	//translate back
	pos.x += ellCenter.x;
	pos.y += ellCenter.y

	return pos;
}


function cylinder(parent: D3Element, bbox: BBox, node: D3Node) {
	const w = bbox.width;
	const rx = w / 2;
	const ry = rx / (2.5 + w / 70);
	const h = bbox.height;

	const shape =
		`M 0,${ry} a${rx},${ry} 0,0,0 ${w} 0 a ${rx},${ry} 0,0,0 ${-w} 0 l 0,${h - 2 * ry} a ${rx},${ry} 0,0,0 ${w} 0 l 0,${-h + 2 * ry}`;

	const shapeSvg = parent
		.attr('label-offset-y', 2 * ry)
		.insert('path', ':first-child')
		.attr('d', shape)
		.attr('transform', 'translate(' + -w / 2 + ',' + -(h / 2) + ')');

	node.intersect = function (point: Point) {
		const pos = intersectRect(node, point)
		let cy = node.y + node.height / 2 - ry
		if (pos.y > cy)
			return intersectEllipse({x: node.x, y: cy}, rx, ry, node, point)

		cy = node.y - node.height / 2 + ry
		if (pos.y < cy)
			return intersectEllipse({x: node.x, y: cy}, rx, ry, node, point)

		return pos;
	};

	return shapeSvg;
}

function person(parent: D3Element, bbox: BBox, node: D3Node) {
	const w = bbox.width;
	const h = bbox.height;

	const shape =
		`M ${.38 * w},${h / 3} A${w / 2},${h / 2} 0,0,0 0 ${h / 2}
		L${w / 11},${h} L${w - w / 11},${h} L${w},${h / 2}
		A${w / 2},${h / 2} 0,0,0 ${w - .38 * w} ${h / 3} 
		A${w / 6},${w / 6} 0,1,0 ${.38 * w} ${h / 3}`;

	const shapeSvg = parent
		.attr('label-offset-y', h * .4)
		.insert('path', ':first-child')
		.attr('d', shape)
		.attr('transform', 'translate(' + -w / 2 + ',' + -(h / 2) + ')');

	node.intersect = function (point: Point) {
		const pos = intersectRect(node, point)
		return pos;
	};

	return shapeSvg;
}


export const shapes: { [key: string]: (parent: SVGElement, node: D3Node) => SVGElement } = {
	rect: (parent: SVGElement, node: D3Node) => rect(new D3Element(parent), node, node).node(),
	cylinder: (parent: SVGElement, node: D3Node) => cylinder(new D3Element(parent), node, node).node(),
	person: (parent: SVGElement, node: D3Node) => person(new D3Element(parent), node, node).node(),
}