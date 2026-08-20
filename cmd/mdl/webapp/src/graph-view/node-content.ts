import {applyStyle, SVG_STYLES} from "./constants";
import {create} from "./svg-create";
import {svgTextWrap} from "./svg-text";

interface TextBlockLayout {
	lines: string[];
	fontSize: number;
	lineHeight: number;
	bold: boolean;
	field?: string;
	gapAfter: number;
}

export interface NodeContentLayout {
	blocks: TextBlockLayout[];
	textHeight: number;
	minimumHeight: number;
}

const HORIZONTAL_PADDING = 18;
const VERTICAL_PADDING = 18;
const TITLE_GAP = 6;
const METADATA_GAP = 10;

function textBlock(
	text: string,
	width: number,
	fontSize: number,
	bold: boolean,
	gapAfter: number,
	field?: string,
): TextBlockLayout {
	const attrs = {
		"font-family": String(SVG_STYLES.nodeText["font-family"]),
		"font-size": `${fontSize}px`,
		"font-weight": bold ? "bold" : "normal",
	};
	const wrapped = svgTextWrap(text, width, attrs);
	const lines = wrapped.lines.length > 0 ? wrapped.lines : [""];

	return {
		lines,
		fontSize,
		lineHeight: fontSize + 2,
		bold,
		field,
		gapAfter,
	};
}

export function layoutNodeContent(
	title: string,
	subtitle: string,
	description: string,
	nodeWidth: number,
	fontSize: number,
): NodeContentLayout {
	const textWidth = Math.max(nodeWidth - HORIZONTAL_PADDING * 2, 80);
	const blocks = [
		textBlock(title, textWidth, fontSize, true, TITLE_GAP, "name"),
		textBlock(`[${subtitle}]`, textWidth, fontSize * 0.75, false, METADATA_GAP),
		textBlock(description, textWidth, Math.min(fontSize * 0.8, 16), false, 0, "description"),
	];
	const textHeight = blocks.reduce(
		(height, block) => height + block.lines.length * block.lineHeight + block.gapAfter,
		0,
	);

	return {
		blocks,
		textHeight,
		minimumHeight: textHeight + VERTICAL_PADDING * 2,
	};
}

export function buildNodeContent(layout: NodeContentLayout, color?: string): SVGGElement {
	const group = create.element("g") as SVGGElement;
	let top = -layout.textHeight / 2;

	layout.blocks.forEach((block) => {
		const text = create.text("", {"text-anchor": "middle"});
		applyStyle(text, SVG_STYLES.nodeText);
		if (color) {
			text.setAttribute("fill", color);
		}
		if (block.field) {
			text.setAttribute("data-field", block.field);
		}

		block.lines.forEach((line, index) => {
			const span = create.element("tspan", {
				x: 0,
				y: top + block.fontSize + index * block.lineHeight,
				"font-size": `${block.fontSize}px`,
				"font-weight": block.bold ? "bold" : "normal",
			});
			span.textContent = line;
			text.append(span);
		});

		group.append(text);
		top += block.lines.length * block.lineHeight + block.gapAfter;
	});

	return group;
}
