import {svgTextWrap} from "./svg-text";

export const create = {
	element(type: string, attrs?: { [key: string]: string | number }, className?:string) {
		const el = document.createElementNS('http://www.w3.org/2000/svg', type)
		attrs && Object.entries(attrs).forEach(([k, v]) => el.setAttribute(k, String(v)))
		className && el.classList.add(className)
		return el
	},

	use(id:string, attrs: {[key: string]: string|number}) {
		const el = create.element('use', attrs)
		el.setAttributeNS('http://www.w3.org/1999/xlink', 'xlink:href', '#'+id)
		return el
	},

	text(text: string, x: number = 0, y: number = 0, anchor = '') {
		const t = create.element('text', {x, y}) as SVGTextElement
		anchor && t.setAttribute('text-anchor', anchor)
		if (text != '') t.append(text)
		return t
	},

	textArea(text: string, width: number, fontSize:number, bold: boolean, x=0, y=0, anchor='') {
		const attrs: {[key: string]: string} = {
			'font-size': `${fontSize}px`,
			'font-weight': bold ? 'bold' : 'normal'
		}
		let {lines, maxW} = svgTextWrap(text, width, attrs)
		const txt = create.text('', 0, y, anchor)
		lines.forEach((line, i) => {
			const span = create.element('tspan', {x, dy: `${fontSize+2}px`})
			for (let attr in attrs) {
				span.setAttribute(attr, attrs[attr]);
			}
			span.append(line)
			txt.append(span)
		})
		return {txt, dy: (lines.length+1) * (fontSize+2), maxW}
	},

	rect(width: number, height: number, x = 0, y = 0, r = 0, className = '') {
		return create.element('rect',
			{x, y, rx: r, ry: r, width, height}, className) as SVGRectElement
	},

	icon(icon: string, x = 0, y = 0) {
		return create.use(icon, {x, y})
	},

	expand(x: number, y: number, ex: boolean) {
		const g = create.element('g', {}, 'expand') as SVGGElement
		setPosition(g, x, y)
		g.append(
			create.rect(19, 19, 0, 0, 1),
			create.text(ex ? '-' : '+', 10, 14, 'middle')
		)
		g.setAttribute('transform', `translate(${x},${y})`)
		return g
	}
}

export function setPosition(g: SVGGElement, x: number, y: number) {
	g.setAttribute('transform', `translate(${x},${y})`)
}