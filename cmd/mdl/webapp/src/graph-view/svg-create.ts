import {svgTextWrap} from "./svg-text";

export const create = {
	element(type: string, attrs: Record<string, string | number> = {}, className?: string) {
		const el = document.createElementNS('http://www.w3.org/2000/svg', type);
		Object.entries(attrs).forEach(([k, v]) => el.setAttribute(k, String(v)));
		if (className) el.classList.add(className);
		return el;
	},

	use(id: string, attrs: Record<string, string | number> = {}) {
		const el = this.element('use', attrs);
		el.setAttributeNS('http://www.w3.org/1999/xlink', 'xlink:href', '#' + id);
		return el;
	},

	path(path: string, attrs: Record<string, string | number> = {}, className?: string) {
		const p = this.element("path", {...attrs, d: path}, className);
		return p;
	},

	text(text: string, attrs: Record<string, string | number> = {}) {
		const t = this.element('text', attrs) as SVGTextElement;
		if (text) t.textContent = text;
		return t;
	},

	textArea(text: string, width: number, fontSize: number, bold: boolean, x = 0, y = 0, anchor = '') {
		const attrs: Record<string, string> = {
			'font-size': `${fontSize}px`,
			'font-weight': bold ? 'bold' : 'normal'
		};
		const {lines, maxW} = svgTextWrap(text, width, attrs);
		const txt = this.text('', {x: 0, y, 'text-anchor': anchor || undefined});
		
		lines.forEach((line, i) => {
			const span = this.element('tspan', {x, dy: `${fontSize + 2}px`, ...attrs});
			span.textContent = line;
			txt.append(span);
		});
		
		return {txt, dy: (lines.length + 1) * (fontSize + 2), maxW};
	},

	rect(width: number, height: number, x = 0, y = 0, r = 0, className?: string) {
		return this.element('rect', {x, y, rx: r, ry: r, width, height}, className) as SVGRectElement;
	},

	icon(icon: string, x = 0, y = 0) {
		return this.use(icon, {x, y});
	},

	expand(x: number, y: number, expanded: boolean) {
		const g = this.element('g', {transform: `translate(${x},${y})`}, 'expand') as SVGGElement;
		g.append(
			this.rect(19, 19, 0, 0, 1),
			this.text(expanded ? '-' : '+', {x: 10, y: 14, 'text-anchor': 'middle'})
		);
		return g;
	}
};

export function setPosition(g: SVGGElement, x: number, y: number) {
	g.setAttribute('transform', `translate(${x},${y})`);
}