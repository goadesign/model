import {GraphData} from "./graph-view/graph";

interface Element {
	id: string;
	name: string;
	technology: string;
	description: string;
	parent: Element;
	tags: string;
	location: string;
	containers?: Element[];
	components?: Element[];
	relationships: Relation[];
	properties: { [key: string]: string }
}

interface Relation {
	id: string;
	description: string;
	tags: string;
	sourceId: string;
	destinationId: string;
	technology: string;
	interactionStyle: string;
}

interface View {
	key: string;
	title: string;
	elements: any[];
	relationships: any[];
	softwareSystemId: string;
}

export const layoutDx = 150
export const layoutDy = 100
export const layoutScale = 1

export const parseModel = (model: any, layouts: any) => {

	const elements = new Map<string, Element>();
	const relations = new Map<string, Relation>();
	Object.entries(model.model).forEach(([k, v]) => {
		if (Array.isArray(v)) {
			v.forEach((el: Element) => {
				elements.set(el.id, el)
				if (k == 'softwareSystems' && !el.technology) {
					el.technology = 'Software System';
				}
				if (Array.isArray(el.relationships)) {
					el.relationships.forEach(rel => {
						relations.set(rel.id, rel)
					})
				}

				if (Array.isArray(el.containers)) {
					el.containers.forEach((el1: Element) => {
						el1.parent = el;
						elements.set(el1.id, el1)
						if (Array.isArray(el1.relationships)) {
							el1.relationships.forEach(rel => {
								relations.set(rel.id, rel)
							})
						}
						if (Array.isArray(el1.components)) {
							el1.components.forEach((el2: Element) => {
								el2.parent = el1;
								elements.set(el2.id, el2)
								if (Array.isArray(el2.relationships)) {
									el2.relationships.forEach(rel => {
										relations.set(rel.id, rel)
									})
								}
							})
						}
					})
				}
			})
		}
	})

	const parseView = (view: View, section: string) => {
		const data = new GraphData(view.key, section + ' - ' +(view.title || view.key))
		//nodes
		view.elements.forEach((ref) => {
			const id = ref.id;
			const el = elements.get(ref.id)

			let shape = 'rect'
			let sub = el ? el.technology : ''
			if (el) {
				const tags = el.tags.split(',')
				if (tags.some(t => t.toLowerCase() == 'database'))
					shape = 'cylinder'
				else if (tags.some(t => t.toLowerCase() == 'person')) {
					shape = 'person'
					sub = 'Person'
				} else if (tags.some(t => t.toLowerCase() == 'container')) {
					sub = 'Container'
					if (el.technology)
						sub += ': ' + el.technology
				}
			}
			data.addNode(
				ref.id,
				el ? el.name : ref.id,
				sub,
				(el && el.description) ? el.description : '',
				shape
			)

		})
		//edges
		if (Array.isArray(view.relationships)) {
			view.relationships.forEach(ref => {
				const rel = relations.get(ref.id)
				if (!rel) return;
				if (!data.nodesMap.has(rel.sourceId)) {
					if (elements.has(rel.sourceId)) {
						const el = elements.get(rel.sourceId)
						console.warn('Element not found in this view: ', el.id, el.name)
					} else {
						console.warn('Element not found: ', rel.sourceId)
					}
					return;
				}
				if (!data.nodesMap.has(rel.destinationId)) {
					if (elements.has(rel.destinationId)) {
						const el = elements.get(rel.destinationId)
						console.warn('Element not found in this view: ', el.id, el.name)
					} else {
						console.warn('Element not found: ', rel.destinationId)
					}
					return;
				}
				data.addEdge(rel.sourceId, rel.destinationId, rel.description)
			})
		}
		//groups
		if (view.softwareSystemId && elements.has(view.softwareSystemId)) {
			const systemEl = elements.get(view.softwareSystemId)
			data.addGroup(systemEl.name, view.elements.map(ref => elements.get(ref.id)).filter(el => el.parent == systemEl).map(el => el.id))
		}

		//layout
		if (data.id in layouts && 'elements' in layouts[data.id]) {
			const layout = layouts[data.id].elements.reduce((o: any, item: any) => {
				o[item.id] = {x: item.x * layoutScale + layoutDx, y: item.y * layoutScale + layoutDy};
				return o
			}, {})
			data.importLayout(layout)
		}
		return data
	}

	const graphs: GraphData[] = []
	const sections = Object.keys(model.views).filter(section => section.endsWith('Views'))
	sections.forEach(s => {
		model.views[s].forEach((v: View) => graphs.push(parseView(v, s)))
	})
	return graphs
}
