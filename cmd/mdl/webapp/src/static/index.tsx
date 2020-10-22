import {parseView} from "../parseModel";
import {buildGraphView} from "../graph-view/graph";

document.querySelectorAll('.c4-diagram').forEach(el => {
	const dataEl = document.getElementById(el.getAttribute('data-model'))
	const model = JSON.parse(dataEl.textContent)

	const graph = parseView(model as any, {}, el.getAttribute('data-view-key'))
	const svg = buildGraphView(graph)
	el.append(svg)
	graph.autoLayout()
})
