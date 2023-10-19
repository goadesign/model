const textMeasure = () => {
	const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
	document.body.appendChild(svg);

	return {
		measure: (text: string, attrs: { [key: string]: string }) => {
			const node = document.createElementNS('http://www.w3.org/2000/svg', 'text')
			node.setAttribute('x', '0');
			node.setAttribute('y', '0');
			for (let attr in attrs) {
				node.setAttribute(attr, attrs[attr]);
			}
			node.appendChild(document.createTextNode(text));

			svg.appendChild(node);
			const {width, height} = node.getBBox();
			svg.removeChild(node);
			return {width, height};
		},
		clean: () => {
			document.body.removeChild(svg);
		}
	}
}

// split a text in lines wrapped at a certain width
export const svgTextWrap = (text: string, width: number, attrs: { [key: string]: string }) => {
	const mt = textMeasure()
	let maxW = 0;
	const ret = text.trim().split('\n').map(text => { //split paragraphs
		//do one paragraph
		const words = text.trim().split(/\s+/);
		let lines = [];
		let currentLine: string[] = [];
		words.forEach(word => {
			const newLine = [...currentLine, word];
			const size = mt.measure(newLine.join(' '), attrs);
			if (size.width > width) {
				lines.push(currentLine.join(' '));
				currentLine = [word];
			} else {
				maxW = Math.max(maxW, size.width)
				currentLine.push(word);
			}
		});

		lines.push(currentLine.join(' '));
		return lines;
	}).reduce((a, v) => a.concat(v), []) //flatten

	mt.clean()
	return {lines: ret, maxW};
};


