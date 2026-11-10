import {applyStyle, NodeStyle, SVG_STYLES} from "./constants";
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
	// The icon shares the content stack so text never draws over it.
	icon?: string;
}

const HORIZONTAL_PADDING = 18;
const VERTICAL_PADDING = 18;
const TITLE_GAP = 6;
const METADATA_GAP = 10;
const ICON_SIZE = 32;
const ICON_GAP = 12;

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
	style: NodeStyle,
): NodeContentLayout {
	const textWidth = Math.max(nodeWidth - HORIZONTAL_PADDING * 2, 80);
	const blocks = [textBlock(title, textWidth, fontSize, true, TITLE_GAP, "name")];
	if (style.metadata !== false && subtitle) {
		blocks.push(textBlock(`[${subtitle}]`, textWidth, fontSize * 0.7, false, METADATA_GAP));
	}
	if (style.description !== false && description) {
		blocks.push(textBlock(description, textWidth, Math.min(fontSize * 0.8, 16), false, 0, "description"));
	}
	blocks[blocks.length - 1].gapAfter = 0;
	const textHeight = blocks.reduce(
		(height, block) => height + block.lines.length * block.lineHeight + block.gapAfter,
		0,
	);

	return {
		blocks,
		textHeight,
		minimumHeight: textHeight + VERTICAL_PADDING * 2 + (style.icon ? ICON_SIZE + ICON_GAP : 0),
		icon: style.icon,
	};
}

export function buildNodeContent(layout: NodeContentLayout, color?: string): SVGGElement {
	const group = create.element("g") as SVGGElement;
	const iconHeight = layout.icon ? ICON_SIZE + ICON_GAP : 0;
	let top = -(layout.textHeight + iconHeight) / 2;
	if (layout.icon) {
		group.append(create.element("image", {
			href: layout.icon,
			x: -ICON_SIZE / 2,
			y: top,
			width: ICON_SIZE,
			height: ICON_SIZE,
			preserveAspectRatio: "xMidYMid meet",
		}, "nodeIcon"));
		top += iconHeight;
	}

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
