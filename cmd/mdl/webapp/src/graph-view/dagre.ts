import {GraphData, Node, Group} from "./graph";
import {graphlib, layout} from "dagre";

export function autoLayout(graph: GraphData) {
	const g = new graphlib.Graph({multigraph: true, compound: true})
	g.setGraph({marginx: 250, marginy: 250, nodesep: 150, ranksep: 150})
	graph.nodesMap.forEach(n => {
		g.setNode(n.id, {width: n.width+1, height: n.height+1})
	})
	graph.edges.forEach(e => g.setEdge(e.from.id, e.to.id, {labelpos: 'c', width: 200, height: 30}, e.id))

	// use grouping info for layout
	graph.groupsMap.forEach(group => {
		g.setNode(group.id, {})
		group.nodes.forEach(member => {
			if (!isGroup(member)) {
				g.setParent(member.id, group.id)
			}
		})
	})

	layout(g)

	return {
		nodes: g.nodes().map(id => {return {id, x: g.node(id).x, y: g.node(id).y}}),
		edges: g.edges().map(e => {
			const edge = g.edge(e)
			return {id: e.name, vertices: edge.points, label: {x: edge.x, y: edge.y}}
		})
	}
}

function isGroup(member: Node | Group): member is Group {
	return 'nodes' in member
}