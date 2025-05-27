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

// Helper function to break long words that exceed width
const breakLongWord = (word: string, maxWidth: number, attrs: { [key: string]: string }, mt: any): string[] => {
	const parts: string[] = [];
	let currentPart = '';
	
	for (let i = 0; i < word.length; i++) {
		const testPart = currentPart + word[i];
		const size = mt.measure(testPart, attrs);
		
		if (size.width > maxWidth && currentPart.length > 0) {
			parts.push(currentPart);
			currentPart = word[i];
		} else {
			currentPart = testPart;
		}
	}
	
	if (currentPart.length > 0) {
		parts.push(currentPart);
	}
	
	return parts;
}

// split a text in lines wrapped at a certain width
export const svgTextWrap = (text: string, width: number, attrs: { [key: string]: string }) => {
	const mt = textMeasure()
	let maxW = 0;
	
	const ret = text.trim().split('\n').map(text => { //split paragraphs
		//do one paragraph
		const words = text.trim().split(/\s+/);
		let lines: string[] = [];
		let currentLine: string[] = [];
		
		words.forEach(word => {
			// First check if the single word exceeds the width
			const wordSize = mt.measure(word, attrs);
			if (wordSize.width > width) {
				// If we have content in current line, finish it first
				if (currentLine.length > 0) {
					lines.push(currentLine.join(' '));
					currentLine = [];
				}
				// Break the long word into smaller parts
				const brokenParts = breakLongWord(word, width, attrs, mt);
				// Add all but the last part as complete lines
				for (let i = 0; i < brokenParts.length - 1; i++) {
					lines.push(brokenParts[i]);
					maxW = Math.max(maxW, mt.measure(brokenParts[i], attrs).width);
				}
				// Start new line with the last part
				if (brokenParts.length > 0) {
					currentLine = [brokenParts[brokenParts.length - 1]];
				}
			} else {
				// Normal word processing
				const newLine = [...currentLine, word];
				const size = mt.measure(newLine.join(' '), attrs);
				if (size.width > width && currentLine.length > 0) {
					lines.push(currentLine.join(' '));
					currentLine = [word];
				} else {
					maxW = Math.max(maxW, size.width)
					currentLine = newLine;
				}
			}
		});

		if (currentLine.length > 0) {
			lines.push(currentLine.join(' '));
		}
		return lines;
	}).reduce((a, v) => a.concat(v), []) //flatten

	mt.clean()
	return {lines: ret, maxW};
};


