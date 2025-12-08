import {GraphData} from "./graph-view/graph";


interface Model {
	name: string
	description: string
	version: string
	model: {
		enterprise: {
			name: string
		}
		people: Element[]
		softwareSystems: Element[]
		deploymentNodes: Element[]
	}
	views: {
		systemLandscapeViews: View[]
		containerViews: View[]
		componentViews: View[]
		dynamicViews: View[]
		deploymentViews: View[]
		styles: {
			elements: {
				[key: string]: string
			}[];
			relationships: {
				[key: string]: string
			}[]
		}
	}
}

interface Layouts {
	[key: string]: { // keyed by view key
		[key: string]: { x: number; y: number } // keyed by element id
	}
}

interface Element {
	id: string;
	name: string;
	technology?: string;
	description?: string;
	parent?: Element;
	tags?: string;
	location?: string;
	containers?: Element[];
	components?: Element[];
	relationships?: Relation[];
	properties?: { [key: string]: string }
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
	description: string
	elements: {
		id: string
	}[];
	relationships: {
		id: string;
		vertices: { x: number; y: number }[];
		routing: string; // takes priority over style
	}[];
	softwareSystemId: string;
}

interface Metadata {
	name: string
	description: string
	version: string
	elements: {
		id: string
		tags?: string;
		location?: string;
		properties?: { [key: string]: string };
		elementViewKey?: string;
		technology?: string;
	}[]
}

export type ViewsList = {
	key: string;
	title: string;
	section: string;
}[]

export const parseView = (model: Model, layouts: Layouts, viewKey: string) => {

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
	model.model.people && model.model.people.forEach((el: Element) => {
		elements.set(el.id, el)
		if (Array.isArray(el.relationships)) {
			el.relationships.forEach(rel => {
				relations.set(rel.id, rel)
			})
		}
	})
	// Software Systems
	model.model.softwareSystems && model.model.softwareSystems.forEach((el: Element) => {
		elements.set(el.id, el)
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
	if (model.model.deploymentNodes) {
		const containerInstances = (el: any) => {
			el.containerInstances && el.containerInstances.forEach((item: any) => {
				const el1 = {...elements.get(item.containerId), id: item.id}
				elements.set(el1.id, el1)
				el1.parent = el
				collectRels(item)
			})
		}

		const recAddNodes = (el: Element, parent: Element) => {
			el.parent = parent;
			elements.set(el.id, el)
			collectRels(el)
			containerInstances(el)
			el.children && el.children.forEach((el1: Element) => recAddNodes(el1, el))
			el.infrastructureNodes && el.infrastructureNodes.forEach((el1: Element) => recAddNodes(el1, el))
		}

		model.model.deploymentNodes.forEach((el: Element) => recAddNodes(el, null))
	}

	// Create graph from selected view
	const {view, section} = getView(model, viewKey)

	if (!view) return null

	const graph = new GraphData(view.key, view.title || view.key)
	const metadata: Metadata = {name: graph.name, description: view.description, version: model.version, elements: []}
	graph.metadata = metadata

	if (!view.elements) return graph

	//grouping rules - elements that are groups will not be nodes
	const groupingIDs: { [key: string]: boolean } = {}
	if (section == 'deploymentViews' || section == 'containerViews') {
		view.elements.forEach(ref => {
			const el = elements.get(ref.id)
			if (el?.parent) {
				groupingIDs[el.parent.id] = true
			}
		})
	} else if (view.softwareSystemId) {
		//don't show grouping if the element is listed in the view
		if (!view.elements.find(ref => ref.id == view.softwareSystemId))
			groupingIDs[view.softwareSystemId] = true
	} else if (section == 'systemLandscapeViews') {
		// create a virtual parent element from enterprise
		const p: Element = {id: '__enterprise__', ...model.model.enterprise}
		elements.set(p.id, p)
		if (model.model.people) model.model.people.filter(el => el.location != 'External').forEach(el => el.parent = p)
		if (model.model.softwareSystems) model.model.softwareSystems.filter(el => el.location != 'External').forEach(el => el.parent = p)
		groupingIDs[p.id] = true
	}

	const styles = model.views.styles

	// Build color-to-variable mapping for CSS custom properties theming
	const cssClassName = (tag: string) => tag.toLowerCase().replace(/[^a-z0-9-]/g, '-')
	
	if (styles?.elements) {
		styles.elements.forEach(s => {
			if (s.tag) {
				const varPrefix = `--mdl-${cssClassName(s.tag)}`
				if (s.background) graph.colorToVarMap.set(s.background as string, `${varPrefix}-bg`)
				if (s.color) graph.colorToVarMap.set(s.color as string, `${varPrefix}-color`)
				if (s.stroke) graph.colorToVarMap.set(s.stroke as string, `${varPrefix}-stroke`)
			}
		})
	}

	if (styles?.relationships) {
		styles.relationships.forEach(s => {
			if (s.tag) {
				const varPrefix = `--mdl-rel-${cssClassName(s.tag)}`
				if (s.color) graph.colorToVarMap.set(s.color as string, `${varPrefix}-color`)
			}
		})
	}

	//nodes
	view.elements.forEach((ref) => {
		// except grouping elements
		if (groupingIDs[ref.id]) return

		const el = elements.get(ref.id)

		let sub = ''
		let style = {}
		if (el) {
			const tags = el.tags.split(',')
			sub = tags[tags.length - 1] // subtitle is [<last tag>]
			if (el.technology)
				sub += ': ' + el.technology // or [<technology>: <last tag>]

			tags.forEach(tag => {
				const s = styles && styles.elements && styles.elements.find(s => s.tag == tag)
				s && (style = {...style, ...s})
			})
		}

		graph.addNode(
			ref.id,
			el ? (el.name || ref.id) : ref.id,
			sub,
			(el && el.description) ? el.description : '',
			style
		)
		el && metadata.elements.push({
			id: el.id,
			tags: el.tags,
			location: el.location,
			properties: el.properties,
			elementViewKey: lookupElementKeyView(model, el.id),
			technology: el.technology
		})
	})
	//edges
	if (Array.isArray(view.relationships)) {
		view.relationships.forEach(ref => {
			const rel = relations.get(ref.id)
			if (!rel) return;

			if (!graph.nodesMap.has(rel.sourceId)) {
				if (elements.has(rel.sourceId)) {
					const el = elements.get(rel.sourceId)
					console.warn('Element not found in this view: ', el.id, el.name)
				} else {
					console.warn('Element not found: ', rel.sourceId)
				}
				return;
			}
			if (!graph.nodesMap.has(rel.destinationId)) {
				if (elements.has(rel.destinationId)) {
					const el = elements.get(rel.destinationId)
					console.warn('Element not found in this view: ', el.id, el.name)
				} else {
					console.warn('Element not found: ', rel.destinationId)
				}
				return;
			}
			let style: any = {}
			rel.tags.split(',').forEach(tag => {
				const s = styles && styles.relationships && styles.relationships.find(s => s.tag == tag)
				s && (style = {...style, ...s})
			})
			if (ref.routing) style.routing = ref.routing

			graph.addEdge(rel.id, rel.sourceId, rel.destinationId, rel.description, ref.vertices, style)
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
		let style = {}
		if (section == 'deploymentViews') {
			const el = elements.get(parent.id)
			const tags = el.tags.split(',')
			tags.forEach(tag => {
				const s = styles && styles.elements && styles.elements.find(s => s.tag == tag)
				s && (style = {...style, ...s})
			})
		}
		
		// Filter group members more carefully to respect boundaries
		const groupMembers = view.elements
			.map(ref => elements.get(ref.id))
			.filter(el => {
				if (!el || el.parent !== parent) return false;
				
				// For system landscape views, respect the location-based grouping
				if (section === 'systemLandscapeViews' && parent.id === '__enterprise__') {
					// Only include elements that are explicitly non-external
					return el.location !== 'External';
				}
				
				// For other view types, include if parent matches
				return true;
			})
			.map(el => el.id);
		
		// Only create group if it has members
		if (groupMembers.length > 0) {
			graph.addGroup(
				parent.id,
				parent.name,
				groupMembers,
				style
			)
		}
	})

	//layout if any and init graph
	graph.init(layouts[graph.id])
	return graph
}

// lookup the view in all Views sections in the model. return the view and the section
function getView(model: Model, viewKey: string) {
	let view: View = null, section: string = ''
	Object.keys(model.views).filter(s => s.endsWith('Views')).some((s: string) => {
		return ((model.views as any)[s]).some((v: View) => {
			if (v.key == viewKey) {
				view = v
				section = s
				return true
			}
		})
	})
	return {view, section}
}


function lookupElementKeyView(model: any, softwareSystemId: string) {
	let key: string = undefined
	Object.keys(model.views).filter(s => s.endsWith('Views')).some((s: string) => {
		return ((model.views as any)[s]).some((v: View) => {
			if (v.softwareSystemId == softwareSystemId) {
				key = v.key
				return true
			}
		})
	})
	return key
}

export const listViews = (model: any) => {
	const viewsList: ViewsList = []
	const sections = Object.keys(model.views).filter(section => section.endsWith('Views'))
	sections.forEach(s => {
		model.views[s].forEach((v: View) => {
			viewsList.push({key: v.key, title: v.title || v.key, section: s})
		})
	})
	return viewsList;
}
