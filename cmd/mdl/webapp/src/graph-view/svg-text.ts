type TextAttributes = Record<string, string>;

interface TextSize {
	width: number;
	height: number;
}

/** WrappedText contains the exact lines and widest rendered line. */
export interface WrappedText {
	lines: string[];
	maxW: number;
}

const measuredText = new Map<string, TextSize>();
let measurementSVG: SVGSVGElement | undefined;

/**
 * measureSVGText returns the browser's rendered text size. Results are cached
 * because node and edge layout often measures the same words many times.
 */
export function measureSVGText(text: string, attrs: TextAttributes): TextSize {
	const key = measurementKey(text, attrs);
	const cached = measuredText.get(key);
	if (cached) {
		return cached;
	}

	const node = document.createElementNS("http://www.w3.org/2000/svg", "text");
	node.setAttribute("x", "0");
	node.setAttribute("y", "0");
	for (const [name, value] of Object.entries(attrs)) {
		node.setAttribute(name, value);
	}
	node.textContent = text;

	const svg = getMeasurementSVG();
	svg.appendChild(node);
	const bounds = node.getBBox();
	node.remove();
	const size = {width: bounds.width, height: bounds.height};
	measuredText.set(key, size);
	return size;
}

/**
 * svgTextWrap splits text into lines that fit the given width. It uses the
 * same browser text sizes that SVG drawing uses.
 */
export function svgTextWrap(
	text: string,
	width: number,
	attrs: TextAttributes,
): WrappedText {
	let maxW = 0;
	const lines = text.trim().split("\n").flatMap(paragraph => {
		const paragraphLines: string[] = [];
		let currentLine: string[] = [];
		for (const word of paragraph.trim().split(/\s+/)) {
			if (measureSVGText(word, attrs).width > width) {
				if (currentLine.length > 0) {
					paragraphLines.push(currentLine.join(" "));
					currentLine = [];
				}
				const parts = breakLongWord(word, width, attrs);
				for (const part of parts.slice(0, -1)) {
					paragraphLines.push(part);
					maxW = Math.max(maxW, measureSVGText(part, attrs).width);
				}
				currentLine = parts.slice(-1);
				continue;
			}

			const nextLine = [...currentLine, word];
			const size = measureSVGText(nextLine.join(" "), attrs);
			if (size.width > width && currentLine.length > 0) {
				paragraphLines.push(currentLine.join(" "));
				currentLine = [word];
			} else {
				maxW = Math.max(maxW, size.width);
				currentLine = nextLine;
			}
		}
		if (currentLine.length > 0) {
			const line = currentLine.join(" ");
			paragraphLines.push(line);
			maxW = Math.max(maxW, measureSVGText(line, attrs).width);
		}
		return paragraphLines;
	});
	return {lines, maxW};
}

/** getMeasurementSVG creates the single hidden SVG used for text measurement. */
function getMeasurementSVG(): SVGSVGElement {
	if (measurementSVG?.isConnected) {
		return measurementSVG;
	}
	measurementSVG = document.createElementNS("http://www.w3.org/2000/svg", "svg");
	measurementSVG.setAttribute("aria-hidden", "true");
	measurementSVG.style.position = "absolute";
	measurementSVG.style.visibility = "hidden";
	measurementSVG.style.pointerEvents = "none";
	document.body.appendChild(measurementSVG);
	return measurementSVG;
}

/** measurementKey gives equal text and font settings the same cache entry. */
function measurementKey(text: string, attrs: TextAttributes): string {
	const attributes = Object.entries(attrs).sort(([first], [second]) =>
		first.localeCompare(second),
	);
	return JSON.stringify([text, attributes]);
}

/** breakLongWord splits one word only when it cannot fit on an empty line. */
function breakLongWord(
	word: string,
	maxWidth: number,
	attrs: TextAttributes,
): string[] {
	const parts: string[] = [];
	let currentPart = "";
	for (const character of word) {
		const nextPart = currentPart + character;
		if (measureSVGText(nextPart, attrs).width > maxWidth && currentPart.length > 0) {
			parts.push(currentPart);
			currentPart = character;
		} else {
			currentPart = nextPart;
		}
	}
	if (currentPart.length > 0) {
		parts.push(currentPart);
	}
	return parts;
}


