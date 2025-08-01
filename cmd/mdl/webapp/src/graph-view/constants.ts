// Constants and default configurations for the graph view

export interface Point {
	x: number;
	y: number;
}

export interface BBox extends Point {
	width: number;
	height: number;
}

export interface NodeStyle {
	// Width of element, in pixels.
	width?: number
	// Height of element, in pixels.
	height?: number
	// Background color of element as HTML RGB hex string (e.g. "#ffffff")
	background?: string
	// Stroke color of element as HTML RGB hex string (e.g. "#000000")
	stroke?: string
	// Foreground (text) color of element as HTML RGB hex string (e.g. "#ffffff")
	color?: string
	// Standard font size used to render text, in pixels.
	fontSize?: number
	// Shape used to render element.
	shape?: string
	// URL of PNG/JPG/GIF file or Base64 data URI representation.
	icon?: string
	// Type of border used to render element.
	border?: string
	// Opacity used to render element; 0-100.
	opacity?: number
	// Whether element metadata should be shown.
	metadata?: boolean
	// Whether element description should be shown.
	description?: boolean
}

export interface EdgeStyle {
	// Thickness of line, in pixels.
	thickness?: number
	// Color of line as HTML RGB hex string (e.g. "#ffffff").
	color?: string
	// Standard font size used to render relationship annotation, in pixels.
	fontSize?: number
	// Width of relationship annotation, in pixels.
	width?: number
	// Whether line is dashed.
	dashed?: boolean
	// Position of label along edge (0-100).
	position?: number
	// Opacity used to render relationship; 0-100.
	opacity?: number
	// Arrow style for the edge.
	arrowStyle?: 'normal' | 'large' | 'small' | 'none'
}

// Default styles
export const DEFAULT_EDGE_STYLE: EdgeStyle = {
	thickness: 3,
	color: '#999',
	opacity: 1,
	fontSize: 22,
	dashed: true,
};

export const DEFAULT_NODE_STYLE: NodeStyle = {
	width: 280,
	height: 180,
	background: 'rgba(255, 255, 255, .9)',
	color: '#666',
	opacity: .9,
	stroke: '#999',
	fontSize: 22,
	shape: 'Box'
};

// SVG styles
export const SVG_STYLES = {
	nodeBorder: {
		fill: "rgba(255, 255, 255, 0.86)",
		stroke: "#aaa",
		filter: 'url(#shadow)',
	},
	nodeText: {
		'font-family': 'Arial, sans-serif',
		stroke: "none"
	},
	edgeText: {
		'font-family': 'Arial, sans-serif',
		stroke: "none"
	},
	edgeRect: {
		fill: "none",
		stroke: "none",
	},
	groupRect: {
		fill: "rgba(0, 0, 0, 0.02)",
		stroke: "#666",
		'stroke-width': 3,
		"stroke-dasharray": 4,
	},
	groupText: {
		fill: "#666",
		"font-size": 22,
		"font-weight": "bold",
		cursor: "default"
	}
};

// Configuration constants
export const SVG_PADDING = 20;
export const DEFAULT_GRID_SIZE = 25;
export const EDGE_SPREAD_DISTANCE = 70;
export const EDGE_SPREAD_DISTANCE_X = 200;

// Utility function to apply styles to SVG elements
export const applyStyle = (el: SVGElement, style: { [key: string]: string | number }) => {
	Object.keys(style).forEach(key => {
		const value = style[key];
		if (typeof value === 'number') {
			el.style.setProperty(key, value.toString());
		} else {
			el.style.setProperty(key, value);
		}
	});
};

// Utility function to calculate distance between two points
export const calculateDistance = (p1: Point, p2: Point): number => {
	return Math.sqrt((p2.x - p1.x) * (p2.x - p1.x) + (p2.y - p1.y) * (p2.y - p1.y));
};