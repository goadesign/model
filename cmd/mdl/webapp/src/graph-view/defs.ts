export const defs = `
<defs>
	<marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" 
		markerWidth="8" markerHeight="8" orient="auto" markerUnits="strokeWidth">
		<path fill="context-stroke" stroke="none" d="M0,0 L10,5 L0,10 Z" class="arrowHead"/>
	</marker>
	<g id="icon-circle" transform="translate(0,-12)" class="icon">
		<circle cx="7" cy="7" r="7"/>
	</g>
	<g id="icon-cube" transform="translate(0, -14)" class="icon">
		<path d="M 5 0 h 10 v 10 l -5 5 h -10 v -10 l 5 -5 M 0 5 h 10 v 10 M 15 0 l -5 5"/>
	</g>
	<filter id="shadow">
		<feDropShadow dx="3" dy="3" stdDeviation="5" flood-color="#00000030" flood-opacity="1"></feDropShadow>
	</filter>
</defs>`
