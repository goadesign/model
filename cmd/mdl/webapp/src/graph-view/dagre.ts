import {GraphData} from "./graph";
import {graphlib, layout} from "dagre";

export function autoLayout(graph: GraphData) {
	const g = new graphlib.Graph({multigraph: true})
	g.setGraph({marginx: 250, marginy: 250, nodesep: 150, ranksep: 150})
	graph.nodesMap.forEach(n => {
		g.setNode(n.id, {width: n.width+1, height: n.height+1})
	})
	graph.edges.forEach(e => g.setEdge(e.from.id, e.to.id, {labelpos: 'c', width: 200, height: 30}, e.id))

	layout(g)

	return {
		nodes: g.nodes().map(id => {return {id, x: g.node(id).x, y: g.node(id).y}}),
		edges: g.edges().map(e => {
			const edge = g.edge(e)
			return {id: e.name, vertices: edge.points, label: {x: edge.x, y: edge.y}}
		})
	}

}