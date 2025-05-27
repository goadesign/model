import {parseView} from "../parseModel";
import {buildGraphView} from "../graph-view/graph";

document.querySelectorAll('.c4-diagram').forEach(async el => {
	const dataEl = document.getElementById(el.getAttribute('data-model'))
	const model = JSON.parse(dataEl.textContent)

	const graph = parseView(model as any, {}, el.getAttribute('data-view-key'))
	const svg = buildGraphView(graph)
	el.append(svg)
	
	try {
		await graph.autoLayout()
	} catch (error) {
		console.error('Auto layout failed for static diagram:', error)
	}
})
