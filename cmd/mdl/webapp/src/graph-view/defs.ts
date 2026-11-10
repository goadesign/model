import fontFaces from "../fonts.css?inline";
import fontLicense from "@fontsource/inter/LICENSE?raw";

/**
 * defs returns one graph's shared SVG definitions. Every id starts with prefix,
 * the graph's idPrefix, so definitions of graphs on the same page stay distinct.
 */
export const defs = (prefix: string) => `
<defs>
	<style>${fontFaces}</style>
	<metadata id="${prefix}inter-font-license">${fontLicense.replace(/&/g, "&amp;").replace(/</g, "&lt;")}</metadata>
	<marker id="${prefix}arrow" viewBox="0 0 10 10" refX="9" refY="5" 
		markerWidth="8" markerHeight="8" orient="auto" markerUnits="strokeWidth">
		<path fill="context-stroke" stroke="none" d="M0,0 L10,5 L0,10 Z" class="arrowHead"/>
	</marker>
	<g id="${prefix}icon-circle" transform="translate(0,-12)" class="icon">
		<circle cx="7" cy="7" r="7"/>
	</g>
	<g id="${prefix}icon-cube" transform="translate(0, -14)" class="icon">
		<path d="M 5 0 h 10 v 10 l -5 5 h -10 v -10 l 5 -5 M 0 5 h 10 v 10 M 15 0 l -5 5"/>
	</g>
</defs>`
