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
	children?: Element[];
	infrastructureNodes?: Element[];
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

export const parseModel = (model: any, layouts: any) => {

	const elements = new Map<string, Element>();
	const relations = new Map<string, Relation>();

	const collectRels = (el: Element) => {
		if (Array.isArray(el.relationships)) {
			el.relationships.forEach(rel => {
				relations.set(rel.id, rel)
			})
		}
	}

	// People
	model.model.people.forEach((el: Element) => {
		elements.set(el.id, el)
		if (Array.isArray(el.relationships)) {
			el.relationships.forEach(rel => {
				relations.set(rel.id, rel)
			})
		}
	})
	// Software Systems
	model.model.softwareSystems.forEach((el: Element) => {
		elements.set(el.id, el)
		el.technology || (el.technology = 'Software System');
		collectRels(el)

		if (Array.isArray(el.containers)) {
			el.containers.forEach((el1: Element) => {
				el1.parent = el;
				elements.set(el1.id, el1)
				collectRels(el1)
				if (Array.isArray(el1.components)) {
					el1.components.forEach((el2: Element) => {
						el2.parent = el1;
						elements.set(el2.id, el2)
						collectRels(el2)
					})
				}
			})
		}
	})

	// Deployment Nodes
	const containerInstances = (el: any) => {
		el.containerInstances && el.containerInstances.forEach((item: any) => {
			const el1 = {...elements.get(item.containerId), id: item.id}
			elements.set(el1.id, el1)
			el1.parent = el
			collectRels(item)
		})
	}

	const recAddNodes = (el:Element, parent: Element) => {
		el.parent = parent;
		elements.set(el.id, el)
		collectRels(el)
		containerInstances(el)
		el.children && el.children.forEach((el1: Element) => recAddNodes(el1, el))
		el.infrastructureNodes && el.infrastructureNodes.forEach((el1: Element) => recAddNodes(el1, el))
	}

	model.model.deploymentNodes.forEach((el: Element) => recAddNodes(el, null))


	// Views
	const parseView = (view: View, section: string) => {
		const data = new GraphData(view.key, section + ' - ' + (view.title || view.key))

		//grouping rules - elements that are groups will not be nodes
		const groupingIDs: {[key: string]: boolean} = {}
		if (section == 'deploymentViews') {
			view.elements.forEach((ref) => {
				const el = elements.get(ref.id)
				if (el && el.parent) {
					groupingIDs[el.parent.id] = true
				}
			})
		} else if (view.softwareSystemId) {
			groupingIDs[view.softwareSystemId] = true
		}
		console.log(view.key, 'grouping:', Object.keys(groupingIDs).map(id => elements.get(id)))

		//nodes
		view.elements.forEach((ref) => {
			// except grouping elements
			if (groupingIDs[ref.id]) return

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
				el ? (el.name || ref.id) : ref.id,
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
		//sort by depth to solve dependency
		const level = (el: Element) => {
			let i = 0
			for (let p = el.parent; p; p = p.parent) i++;
			return i
		}
		const gElements = Object.keys(groupingIDs)
			.map(id => elements.get(id))
			.sort((a, b) => level(a) > level(b) ? -1 : 1)

		gElements.forEach(parent => {
			data.addGroup(parent.id, parent.name,
				view.elements
					.map(ref => elements.get(ref.id))
					.filter(el => el && el.parent == parent)
					.map(el => el.id)
			)
		})

		//layout
		if (data.id in layouts) {
			data.importLayout(layouts[data.id])
		}
		return data
	}

	// Graph
	const graphs: GraphData[] = []
	const sections = Object.keys(model.views).filter(section => section.endsWith('Views'))
	sections.forEach(s => {
		model.views[s].forEach((v: View) => graphs.push(parseView(v, s)))
	})
	return graphs
}
