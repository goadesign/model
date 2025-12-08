import {intersectEllipse, intersectRect} from "./intersect";

interface Point {
	x: number;
	y: number;
}

interface BBox extends Point {
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


function rect(parent: D3Element, bbox: BBox, node: D3Node, rounded = false) {
	const shapeSvg = parent.insert("rect", ":first-child")
		.attr("rx", rounded ? node.width / 8 : 3)
		.attr("ry", rounded ? node.width / 8 : 3)
		.attr("x", -bbox.width / 2)
		.attr("y", -bbox.height / 2)
		.attr("width", bbox.width)
		.attr("height", bbox.height);

	node.intersect = function (point) {
		return intersectRect(node, point);
	};

	return shapeSvg;
}


function cylinder(parent: D3Element, bbox: BBox, node: D3Node) {
	const w = bbox.width;
	const rx = w / 2;
	const ry = rx / (5.5 + w / 70); // Make ellipse even flatter - increased to 5.5
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

function _ellipse(parent: D3Element, bbox: BBox, node: D3Node, rx: number, ry: number) {
	const shapeSvg = parent.insert("ellipse", ":first-child")
		.attr("cx", 0)
		.attr("cy", 0)
		.attr('rx', rx)
		.attr('ry', ry)
		.attr("width", node.width)
		.attr("height", node.height);

	node.intersect = function (point) {
		return intersectEllipse(node, rx, ry, node, point)
	};
	return shapeSvg;
}

function circle(parent: D3Element, bbox: BBox, node: D3Node) {
	return _ellipse(parent, bbox, node, node.width / 2, node.width / 2)
}

function ellipse(parent: D3Element, bbox: BBox, node: D3Node) {
	return _ellipse(parent, bbox, node, node.width * .55, node.width * .45)
}

function hexagon(parent: D3Element, bbox: BBox, node: D3Node) {
	const sz = node.width / 2
	// drawing a hexagon from polar coords
	// [0,1,2,3,4,5,6].map(i=>`${Math.sin(Math.PI/3*i+Math.PI/6).toFixed(4)},${Math.cos(Math.PI/3*i+Math.PI/6).toFixed(4)}`).join(',')
	const shapeSvg = parent.insert("polygon", ":first-child")
		.attr("points",
			[0.5000, 0.8660, 1.0000, 0.0000, 0.5000, -0.8660, -0.5000, -0.8660, -1.0000, -0.0000, -0.5000, 0.8660, 0.5000, 0.8660].map(n => n * sz).join(','))
		.attr("width", node.width)
		.attr("height", node.height);

	node.intersect = function (point) {
		return intersectEllipse(node, node.width / 2, node.width / 2, node, point)
	};
	return shapeSvg;
}

function component(parent: D3Element, bbox: BBox, node: D3Node) {
	const dx = node.width / 10
	const shapeSvg = parent.insert('g', ':first-child')
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", 3).attr("ry", 3)
		.attr("x", -node.width / 2 - dx)
		.attr("y", -node.height / 2 + dx)
		.attr("width", dx * 2)
		.attr("height", dx);
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", 3).attr("ry", 3)
		.attr("x", -node.width / 2 - dx)
		.attr("y", -node.height / 2 + dx * 2.5)
		.attr("width", dx * 2)
		.attr("height", dx);
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", 3).attr("ry", 3)
		.attr("x", -node.width / 2)
		.attr("y", -node.height / 2)
		.attr("width", node.width)
		.attr("height", node.height);

	node.intersect = function (point) {
		return intersectRect({x: node.x - dx / 2, y: node.y, width: node.width + dx, height: node.height}, point);
	};

	return shapeSvg;
}

function folder(parent: D3Element, bbox: BBox, node: D3Node) {
	const dy = node.width / 20
	const shapeSvg = parent
		.attr('label-offset-y', dy * 2)
		.insert('g', ':first-child')
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", 3).attr("ry", 3)
		.attr("x", -node.width / 2)
		.attr("y", -node.height / 2 + dy * 2)
		.attr("width", node.width)
		.attr("height", node.height - dy * 2);
	shapeSvg.insert("path", ":first-child")
		.attr('d', `M0,${-node.height / 2 + 2 * dy} l${dy},${-2 * dy} h${node.width / 2 - dy * 2} v${dy * 2}`)

	node.intersect = function (point) {
		return intersectRect({x: node.x, y: node.y + dy / 2, width: node.width, height: node.height + dy}, point);
	};

	return shapeSvg;
}

function mobiledevicelandscape(parent: D3Element, bbox: BBox, node: D3Node, rounded = false) {
	const dx = node.width / 8
	const r = node.width / 14
	const shapeSvg = parent.insert('g', ':first-child')
	shapeSvg.insert('path', ':first-child')
		.attr('d', `M${-node.width / 2},${-node.height / 2} l0,${node.height} M${node.width / 2},${-node.height / 2} l0,${node.height}`)
	shapeSvg.insert('circle', ':first-child')
		.attr('cx', -node.width / 2 - dx / 2)
		.attr('cy', 0)
		.attr('r', r * .4)
	shapeSvg.insert('rect', ':first-child')
		.attr('x', node.width / 2 + dx / 2 - r * .2)
		.attr('y', -r)
		.attr('width', r * .4)
		.attr('height', r * 2)
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", r)
		.attr("ry", r)
		.attr("x", -bbox.width / 2 - dx)
		.attr("y", -bbox.height / 2)
		.attr("width", bbox.width + 2 * dx)
		.attr("height", bbox.height);

	node.intersect = function (point) {
		return intersectRect({x: node.x, y: node.y, width: node.width + 2 * dx, height: node.height}, point);
	};

	return shapeSvg;
}

function mobiledeviceportrait(parent: D3Element, bbox: BBox, node: D3Node) {
	const dy = node.width / 8
	const r = node.width / 14
	const shapeSvg = parent.insert('g', ':first-child')
	shapeSvg.insert('path', ':first-child')
		.attr('d', `M${-node.width / 2},${-node.height / 2} l${node.width},0 M${-node.width / 2},${node.height / 2} l${node.width},0`)
	shapeSvg.insert('circle', ':first-child')
		.attr('cx', 0)
		.attr('cy', node.height / 2 + dy / 2)
		.attr('r', r * .4)
	shapeSvg.insert('rect', ':first-child')
		.attr('x', -r)
		.attr('y', -node.height / 2 - dy / 2 - r * .2)
		.attr('width', r * 2)
		.attr('height', r * .4)
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", r)
		.attr("ry", r)
		.attr("x", -bbox.width / 2)
		.attr("y", -bbox.height / 2 - dy)
		.attr("width", bbox.width)
		.attr("height", bbox.height + 2 * dy);

	node.intersect = function (point) {
		return intersectRect({x: node.x, y: node.y, width: node.width, height: node.height + 2 * dy}, point);
	};

	return shapeSvg;
}

function pipe(parent: D3Element, bbox: BBox, node: D3Node) {
	const w = node.width;
	const h = node.height;
	const ry = h / 2;
	const rx = ry / (2.5 + w / 70);

	const shape =
		`M${-rx},0
		a${rx},${ry} 0,0,1 0,${h}
		a${rx},${ry} 0,0,1 0,${-h}
		l${w},0
		a${rx},${ry} 0,0,1 0,${h}
		l${-w},0`;

	const shapeSvg = parent
		.insert('path', ':first-child')
		.attr('d', shape)
		.attr('transform', 'translate(' + -w / 2 + ',' + -(h / 2) + ')');

	node.intersect = function (point: Point) {
		return intersectRect({x: node.x - rx, y: node.y, width: node.width + 2 * rx, height: node.height}, point)
	};

	return shapeSvg;
}

function robot(parent: D3Element, bbox: BBox, node: D3Node) {
	const w = node.width
	const h = node.height
	
	// Small head at top (like person shape but robot-styled)
	const headSize = Math.min(w * 0.28, h * 0.25)
	const headR = headSize * 0.2
	const antennaH = headSize * 0.25
	const antennaR = headSize * 0.08
	
	// Eye dimensions
	const eyeR = headSize * 0.12
	const eyeSpacing = headSize * 0.22
	
	// Ear dimensions  
	const earW = headSize * 0.12
	const earH = headSize * 0.3
	
	// Body fills most of the space for text
	const bodyW = w
	const bodyTop = -h / 2 + antennaH + headSize
	const bodyH = h - antennaH - headSize
	const bodyR = 3
	
	// Label offset - similar to person shape (h * 0.4)
	const labelOffsetY = h * 0.35
	
	const shapeSvg = parent
		.attr('label-offset-y', labelOffsetY)
		.insert('g', ':first-child')
	
	// Body - main rectangle for text (draw first so it's behind)
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", bodyR)
		.attr("ry", bodyR)
		.attr('x', -bodyW / 2)
		.attr('y', bodyTop)
		.attr('width', bodyW)
		.attr('height', bodyH)
	
	// Head
	const headTop = -h / 2 + antennaH
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", headR)
		.attr("ry", headR)
		.attr('x', -headSize / 2)
		.attr('y', headTop)
		.attr('width', headSize)
		.attr('height', headSize)
	
	// Antenna
	shapeSvg.insert("line", ":first-child")
		.attr('class', 'robot-antenna')
		.attr('x1', 0)
		.attr('y1', headTop)
		.attr('x2', 0)
		.attr('y2', -h / 2 + antennaR * 2)
		.attr('stroke-width', antennaR * 0.6)
		.attr('stroke-linecap', 'round')
	
	// Antenna ball
	shapeSvg.insert("circle", ":first-child")
		.attr('class', 'robot-antenna-ball')
		.attr('cx', 0)
		.attr('cy', -h / 2 + antennaR * 2)
		.attr('r', antennaR * 1.2)
	
	// Eyes
	const eyeY = headTop + headSize * 0.4
	shapeSvg.insert("circle", ":first-child")
		.attr('class', 'robot-eye')
		.attr('cx', -eyeSpacing)
		.attr('cy', eyeY)
		.attr('r', eyeR)
	shapeSvg.insert("circle", ":first-child")
		.attr('class', 'robot-eye')
		.attr('cx', eyeSpacing)
		.attr('cy', eyeY)
		.attr('r', eyeR)
	
	// Mouth - simple smile
	const mouthY = headTop + headSize * 0.7
	const mouthW = headSize * 0.28
	shapeSvg.insert("path", ":first-child")
		.attr('class', 'robot-mouth')
		.attr('d', `M${-mouthW / 2},${mouthY} Q0,${mouthY + mouthW * 0.3} ${mouthW / 2},${mouthY}`)
		.attr('fill', 'none')
		.attr('stroke-width', antennaR * 0.5)
		.attr('stroke-linecap', 'round')
	
	// Ears
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", earW * 0.25)
		.attr("ry", earW * 0.25)
		.attr('x', -headSize / 2 - earW - 1)
		.attr('y', eyeY - earH / 2)
		.attr('width', earW)
		.attr('height', earH)
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", earW * 0.25)
		.attr("ry", earW * 0.25)
		.attr('x', headSize / 2 + 1)
		.attr('y', eyeY - earH / 2)
		.attr('width', earW)
		.attr('height', earH)

	node.intersect = function (point) {
		return intersectRect(node, point);
	};

	return shapeSvg;
}

function webbrowser(parent: D3Element, bbox: BBox, node: D3Node) {
	const dy = node.height / 8
	const shapeSvg = parent
		.attr('label-offset-y', dy)
		.insert('g', ':first-child')
	shapeSvg.insert("path", ":first-child")
		.attr('d', `
			M${-node.width / 2},${-node.height / 2 + dy} h${node.width}
			M${-node.width / 2 + dy / 4},${-node.height / 2 + dy / 4} h${dy / 2} v${dy / 2} h${-dy / 2} z
			M${-node.width / 2 + dy},${-node.height / 2 + dy / 4} h${node.width - dy - dy / 4} v${dy / 2} h${-node.width + dy + dy / 4} z
		`)
	shapeSvg.insert("rect", ":first-child")
		.attr("rx", 3).attr("ry", 3)
		.attr("x", -node.width / 2)
		.attr("y", -node.height / 2)
		.attr("width", node.width)
		.attr("height", node.height);

	node.intersect = function (point) {
		return intersectRect(node, point);
	};

	return shapeSvg;
}

export const shapes: { [key: string]: (parent: SVGElement, node: D3Node) => SVGElement } = {
	box: (parent: SVGElement, node: D3Node) => rect(new D3Element(parent), node, node).node(),
	roundedbox: (parent: SVGElement, node: D3Node) => rect(new D3Element(parent), node, node, true).node(),
	component: (parent: SVGElement, node: D3Node) => component(new D3Element(parent), node, node).node(),
	cylinder: (parent: SVGElement, node: D3Node) => cylinder(new D3Element(parent), node, node).node(),
	person: (parent: SVGElement, node: D3Node) => person(new D3Element(parent), node, node).node(),
	circle: (parent: SVGElement, node: D3Node) => circle(new D3Element(parent), node, node).node(),
	ellipse: (parent: SVGElement, node: D3Node) => ellipse(new D3Element(parent), node, node).node(),
	hexagon: (parent: SVGElement, node: D3Node) => hexagon(new D3Element(parent), node, node).node(),
	folder: (parent: SVGElement, node: D3Node) => folder(new D3Element(parent), node, node).node(),
	mobiledevicelandscape: (parent: SVGElement, node: D3Node) => mobiledevicelandscape(new D3Element(parent), node, node).node(),
	mobiledeviceportrait: (parent: SVGElement, node: D3Node) => mobiledeviceportrait(new D3Element(parent), node, node).node(),
	mobiledevice: (parent: SVGElement, node: D3Node) => mobiledeviceportrait(new D3Element(parent), node, node).node(),
	pipe: (parent: SVGElement, node: D3Node) => pipe(new D3Element(parent), node, node).node(),
	robot: (parent: SVGElement, node: D3Node) => robot(new D3Element(parent), node, node).node(),
	webbrowser: (parent: SVGElement, node: D3Node) => webbrowser(new D3Element(parent), node, node).node(),
}
