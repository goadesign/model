interface Point {
	x: number;
	y: number;
}

interface BBox extends Point {
	width: number;
	height: number;
}


export function insideBox(p: Point, b: BBox): boolean {
	return p.x > b.x - b.width / 2 && p.x < b.x + b.width / 2 && p.y > b.y - b.height / 2 && p.y < b.y + b.height / 2
}

// intersect 2 segments (p1->q1) with (p2, q2)
// if the lines intersect, the result contains the x and y of the intersection (treating the lines as infinite)
// and booleans for whether line segment 1 or line segment 2 contain the point
function segmentIntersection(p1: Point, q1: Point, p2: Point, q2: Point) {
	let denominator, a, b, numerator1, numerator2,
		result: { x: number, y: number, onLine1: boolean, onLine2: boolean } = {
			x: null,
			y: null,
			onLine1: false,
			onLine2: false
		};
	denominator = ((q2.y - p2.y) * (q1.x - p1.x)) - ((q2.x - p2.x) * (q1.y - p1.y));
	if (denominator == 0) {
		return result;
	}
	a = p1.y - p2.y;
	b = p1.x - p2.x;
	numerator1 = ((q2.x - p2.x) * a) - ((q2.y - p2.y) * b);
	numerator2 = ((q1.x - p1.x) * a) - ((q1.y - p1.y) * b);
	a = numerator1 / denominator;
	b = numerator2 / denominator;

	// if we cast these lines infinitely in both directions, they intersect here:
	result.x = p1.x + (a * (q1.x - p1.x));
	result.y = p1.y + (a * (q1.y - p1.y));

	// if line1 is a segment and line2 is infinite, they intersect if:
	if (a > 0 && a < 1) {
		result.onLine1 = true;
	}
	// if line2 is a segment and line1 is infinite, they intersect if:
	if (b > 0 && b < 1) {
		result.onLine2 = true;
	}
	// if line1 and line2 are segments, they intersect if both of the above are true
	return result;
}

// intersects a segment (p1->p2) with a box
export function intersectRectFull(p1: Point, p2: Point, box: BBox): Point[] {
	const w = box.width / 2
	const h = box.height / 2
	const segs: { p: Point; q: Point }[] = [
		{p: {x: box.x - w, y: box.y - h}, q: {x: box.x - w, y: box.y + h}},
		{p: {x: box.x - w, y: box.y - h}, q: {x: box.x + w, y: box.y - h}},
		{p: {x: box.x + w, y: box.y - h}, q: {x: box.x + w, y: box.y + h}},
		{p: {x: box.x - w, y: box.y + h}, q: {x: box.x + w, y: box.y + h}},
	]
	return segs.map(s => segmentIntersection(p1, p2, s.p, s.q)).filter(ret => ret.onLine1 && ret.onLine2)
}

// intersects a line that goes from p to the center of the box
export function intersectRect(box: BBox, p: Point): Point {
	if (insideBox(p, box)) return {x: box.x, y: box.y}
	return intersectRectFull(box, p, box)[0]
}

export function intersectEllipse(ellCenter: Point, rx: number, ry: number, nodeCenter: Point, point: Point) {

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

export interface Segment {
	p: Point;
	q: Point;
}

// given a polyline as a list of segments, interrupt it over the box so no line is inside the box
export function intersectPolylineBox(segments: Segment[], box: BBox) {
	for (let i = 0; i < segments.length; i++) {
		const s = segments[i]
		if (insideBox(s.p, box)) {
			if (insideBox(s.q, box)) { // segment both ends inside box
				segments.splice(i, 1)
				i -= 1
			} else { // segment start inside box
				s.p = intersectRectFull(s.p, s.q, box)[0]
			}
		} else {
			if (insideBox(s.q, box)) { // segment end inside box
				s.q = intersectRectFull(s.p, s.q, box)[0]
			} else { // both ends outside
				const ret = intersectRectFull(s.p, s.q, box)
				if (ret.length == 2) {  // intersects the box, splice segment
					// order the intersection points, closest first
					const dst1 = Math.abs(ret[0].x - s.p.x) + Math.abs(ret[0].y - s.p.y)
					const dst2 = Math.abs(ret[1].x - s.p.x) + Math.abs(ret[1].y - s.p.y)
					if (dst1 > dst2) ret.reverse()
					// split the segment in 2
					const s2 = {p: ret[1], q: s.q}
					s.q = ret[0]
					segments.splice(i + 1, 0, s2)
					i += 1
				}
			}
		}
	}
}