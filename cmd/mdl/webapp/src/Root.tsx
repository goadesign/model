import React, {FC, useState, useCallback} from "react";
import {getZoom, getZoomAuto, GraphData, setZoom} from "./graph-view/graph";
import {Graph} from "./graph-view/graph-react";
import {BrowserRouter as Router, Routes, Route, useSearchParams} from 'react-router-dom'
import {listViews, parseView} from "./parseModel";
import {findShortcut, HELP, Help, SAVE} from "./shortcuts";
import {LayoutAlgorithm, LayoutOptions} from "./graph-view/layout";

// Global state for graphs to preserve edits
const graphs: {[key: string]: GraphData} = {}

// Global functions for keyboard shortcuts
let globalToggleHelp: () => void
let globalSaveLayout: () => void

// Setup global keyboard event listener
window.addEventListener('keydown', e => {
	const shortcut = findShortcut(e)
	if (globalToggleHelp && shortcut === HELP) {
		globalToggleHelp()
	} else if (globalSaveLayout && shortcut === SAVE) {
		globalSaveLayout()
		e.preventDefault()
	}
})

export const Root: FC<{model: any, layout: any}> = ({model, layout}) => (
	<Router>
		<Routes>
			<Route path="/" element={<ModelPane model={model} layouts={layout}/>}/>
		</Routes>
	</Router>
)

export const refreshGraph = () => {
	const currentID = getCurrentViewID()
	delete graphs[currentID]
}

const getCurrentViewID = () => {
	const params = new URLSearchParams(document.location.search)
	return params.get('id') || ''
}

const ModelPane: FC<{model: any, layouts: any}> = ({model, layouts}) => {
	const [searchParams, setSearchParams] = useSearchParams()
	const currentID = decodeURI(searchParams.get('id') || '')
	
	// UI State
	const [saving, setSaving] = useState(false)
	const [helpVisible, setHelpVisible] = useState(false)
	const [layouting, setLayouting] = useState(false)
	
	// Settings State
	const [layoutAlgorithm, setLayoutAlgorithm] = useState<LayoutAlgorithm>('force')
	const [connectionRouting, setConnectionRouting] = useState<'Orthogonal' | 'Polyline' | 'Splines' | 'Curved'>('Orthogonal')

	// Get or create graph for current view
	const graph = getOrCreateGraph(model, layouts, currentID)
	if (!graph) {
		return <ViewRedirect model={model} />
	}

	// Setup global keyboard handlers
	globalToggleHelp = () => setHelpVisible(!helpVisible)
	globalSaveLayout = useCallback(() => saveCurrentLayout(graph, currentID, setSaving), [graph, currentID])

	const handleAutoLayout = useCallback(async () => {
		setLayouting(true)
		try {
			const options: LayoutOptions = {
				algorithm: layoutAlgorithm,
				direction: 'DOWN',
				edgeRouting: 'ORTHOGONAL',
				favorStraightEdges: true,
				compactLayout: true
			}
			await graph.autoLayout(options)
		} catch (error) {
			console.error('Layout failed:', error)
			alert('Layout failed. See console for details.')
		} finally {
			setLayouting(false)
		}
	}, [graph, layoutAlgorithm])

	const handleConnectionStyleChange = useCallback((global: boolean = false) => {
		const styleChanges = { routing: connectionRouting }
		
		if (global) {
			graph.changeAllEdgeStyle(styleChanges)
		} else {
			graph.changeSelectedEdgeStyle(styleChanges)
		}
	}, [graph, connectionRouting])

	return (
		<>
			<Toolbar
				model={model}
				currentID={currentID}
				onViewChange={(id) => setSearchParams({id: encodeURIComponent(id)})}
				graph={graph}
				layoutAlgorithm={layoutAlgorithm}
				onLayoutAlgorithmChange={setLayoutAlgorithm}
				connectionRouting={connectionRouting}
				onConnectionRoutingChange={setConnectionRouting}
				onAutoLayout={handleAutoLayout}
				onConnectionStyleChange={handleConnectionStyleChange}
				onSave={() => globalSaveLayout()}
				onToggleHelp={() => setHelpVisible(!helpVisible)}
				saving={saving}
				layouting={layouting}
			/>
			<Graph 
				key={currentID}
				data={graph}
				onSelect={id => id && console.log(removeEmptyProps(graph.metadata.elements.find((m: any) => m.id === id)))}
			/>
			{helpVisible && <Help/>}
		</>
	)
}

const ViewRedirect: FC<{model: any}> = ({model}) => {
	const views = listViews(model)
	if (views.length > 0) {
		document.location.href = '?id=' + views[0].key
		return <>Redirecting to {views[0].title}</>
	}
	return <>No views available</>
}

const Toolbar: FC<{
	model: any
	currentID: string
	onViewChange: (id: string) => void
	graph: GraphData
	layoutAlgorithm: LayoutAlgorithm
	onLayoutAlgorithmChange: (algorithm: LayoutAlgorithm) => void
	connectionRouting: string
	onConnectionRoutingChange: (routing: any) => void
	onAutoLayout: () => void
	onConnectionStyleChange: (global: boolean) => void
	onSave: () => void
	onToggleHelp: () => void
	saving: boolean
	layouting: boolean
}> = ({
	model, currentID, onViewChange, graph, layoutAlgorithm, onLayoutAlgorithmChange,
	connectionRouting, onConnectionRoutingChange, onAutoLayout, onConnectionStyleChange,
	onSave, onToggleHelp, saving, layouting
}) => {
	const views = listViews(model)
	
	return (
		<div className="toolbar">
			<ViewSelector 
				views={views}
				currentID={currentID}
				onViewChange={onViewChange}
			/>
			<ToolbarActions
				graph={graph}
				layoutAlgorithm={layoutAlgorithm}
				onLayoutAlgorithmChange={onLayoutAlgorithmChange}
				connectionRouting={connectionRouting}
				onConnectionRoutingChange={onConnectionRoutingChange}
				onAutoLayout={onAutoLayout}
				onConnectionStyleChange={onConnectionStyleChange}
				onSave={onSave}
				onToggleHelp={onToggleHelp}
				saving={saving}
				layouting={layouting}
			/>
		</div>
	)
}

const ViewSelector: FC<{
	views: any[]
	currentID: string
	onViewChange: (id: string) => void
}> = ({views, currentID, onViewChange}) => (
	<div>
		View:
		{views.length > 1 ? (
			<select onChange={e => onViewChange(e.target.value)} value={currentID}>
				<option disabled value="" hidden>...</option>
				{views.map(view => (
					<option key={view.key} value={view.key}>
						{camelToWords(view.section) + ': ' + view.title}
					</option>
				))}
			</select>
		) : (
			<span style={{marginLeft: '8px', fontWeight: 'bold'}}>
				{views[0] ? camelToWords(views[0].section) + ': ' + views[0].title : 'No views available'}
			</span>
		)}
	</div>
)

const ToolbarActions: FC<{
	graph: GraphData
	layoutAlgorithm: LayoutAlgorithm
	onLayoutAlgorithmChange: (algorithm: LayoutAlgorithm) => void
	connectionRouting: string
	onConnectionRoutingChange: (routing: any) => void
	onAutoLayout: () => void
	onConnectionStyleChange: (global: boolean) => void
	onSave: () => void
	onToggleHelp: () => void
	saving: boolean
	layouting: boolean
}> = ({
	graph, layoutAlgorithm, onLayoutAlgorithmChange, connectionRouting, onConnectionRoutingChange,
	onAutoLayout, onConnectionStyleChange, onSave, onToggleHelp, saving, layouting
}) => (
	<div>
		<UndoRedoButtons graph={graph} />
		<AlignmentButtons graph={graph} />
		<LayoutControls
			layoutAlgorithm={layoutAlgorithm}
			onLayoutAlgorithmChange={onLayoutAlgorithmChange}
			onAutoLayout={onAutoLayout}
			layouting={layouting}
		/>
		<ConnectionControls
			connectionRouting={connectionRouting}
			onConnectionRoutingChange={onConnectionRoutingChange}
			onConnectionStyleChange={onConnectionStyleChange}
		/>
		<ZoomControls graph={graph} />
		<ActionButtons
			onSave={onSave}
			onToggleHelp={onToggleHelp}
			saving={saving}
		/>
	</div>
)

const UndoRedoButtons: FC<{graph: GraphData}> = ({graph}) => (
	<>
		<button onClick={() => graph.undo()} data-tooltip="Undo the last change made to the diagram">
			<i className="fas fa-undo"></i>
		</button>
		<button className="grp" onClick={() => graph.redo()} data-tooltip="Redo the last undone action">
			<i className="fas fa-redo"></i>
		</button>
	</>
)

const AlignmentButtons: FC<{graph: GraphData}> = ({graph}) => (
	<>
		<button onClick={() => graph.alignSelectionH()} data-tooltip="Align all selected elements horizontally (left edges)">
			<i className="fas fa-arrows-alt-v"></i>
		</button>
		<button onClick={() => graph.alignSelectionV()} data-tooltip="Align all selected elements vertically (top edges)">
			<i className="fas fa-arrows-alt-h"></i>
		</button>
		<button onClick={() => graph.distributeSelectionH()} data-tooltip="Distribute selected elements evenly horizontally (equal spacing)">
			<i className="fas fa-ellipsis-h"></i>
		</button>
		<button className="grp" onClick={() => graph.distributeSelectionV()} data-tooltip="Distribute selected elements evenly vertically (equal spacing)">
			<i className="fas fa-ellipsis-v"></i>
		</button>
	</>
)

const LayoutControls: FC<{
	layoutAlgorithm: LayoutAlgorithm
	onLayoutAlgorithmChange: (algorithm: LayoutAlgorithm) => void
	onAutoLayout: () => void
	layouting: boolean
}> = ({layoutAlgorithm, onLayoutAlgorithmChange, onAutoLayout, layouting}) => (
	<>
		<select 
			value={layoutAlgorithm} 
			onChange={e => onLayoutAlgorithmChange(e.target.value as LayoutAlgorithm)}
			data-tooltip="Select which automatic layout algorithm to use"
			disabled={layouting}
		>
			<option value="force">Force</option>
			<option value="layered">Layered (Hierarchical)</option>
			<option value="stress">Stress (Force-directed)</option>
			<option value="mrtree">Tree</option>
			<option value="radial">Radial</option>
			<option value="disco">Disco</option>
		</select>
		<button 
			className="grp" 
			onClick={onAutoLayout} 
			disabled={layouting}
			data-tooltip={`Automatically arrange all elements using ${layoutAlgorithm} algorithm`}
		>
			{layouting ? <i className="fas fa-spinner fa-spin"></i> : <i className="fas fa-magic"></i>}
		</button>
	</>
)

const ConnectionControls: FC<{
	connectionRouting: string
	onConnectionRoutingChange: (routing: any) => void
	onConnectionStyleChange: (global: boolean) => void
}> = ({connectionRouting, onConnectionRoutingChange, onConnectionStyleChange}) => (
	<>
		<select 
			value={connectionRouting} 
			onChange={e => onConnectionRoutingChange(e.target.value)}
			data-tooltip="Choose how connections are routed between elements"
		>
			<option value="Orthogonal">Orthogonal</option>
			<option value="Polyline">Straight</option>
			<option value="Splines">Curved</option>
		</select>
		<button 
			onClick={() => onConnectionStyleChange(false)}
			data-tooltip="Apply connection routing style to currently selected connections only"
		>
			<i className="fas fa-dot-circle"></i>
		</button>
		<button 
			className="grp"
			onClick={() => onConnectionStyleChange(true)}
			data-tooltip="Apply connection routing style to all connections in the diagram"
		>
			<i className="fas fa-globe"></i>
		</button>
	</>
)

const ZoomControls: FC<{graph: GraphData}> = ({graph}) => (
	<>
		<button onClick={() => setZoom(getZoom() - .05)} data-tooltip="Zoom out to see more of the diagram">
			<i className="fas fa-search-minus"></i>
		</button>
		<button onClick={() => setZoom(getZoom() + .05)} data-tooltip="Zoom in to see details more clearly">
			<i className="fas fa-search-plus"></i>
		</button>
		<button onClick={() => {
			graph.alignTopLeft()
			setZoom(getZoomAuto())
		}} data-tooltip="Automatically fit the entire diagram in the visible area">
			<i className="fas fa-expand-arrows-alt"></i>
		</button>
		<button onClick={() => setZoom(1)} data-tooltip="Reset zoom to 100% (actual size)">
			<i className="fas fa-search"></i>
		</button>
	</>
)

const ActionButtons: FC<{
	onSave: () => void
	onToggleHelp: () => void
	saving: boolean
}> = ({onSave, onToggleHelp, saving}) => (
	<>
		<button className="action" disabled={saving} onClick={onSave} data-tooltip="Save the current diagram layout">
			<i className="fas fa-save"></i>
		</button>
		<button onClick={onToggleHelp} data-tooltip="Show keyboard shortcuts and help information">
			<i className="fas fa-question-circle"></i>
		</button>
	</>
)

// Helper functions
function getOrCreateGraph(model: any, layouts: any, currentID: string): GraphData | null {
	if (graphs[currentID]) {
		return graphs[currentID]
	}
	
	const graph = parseView(model, layouts, currentID)
	if (graph) {
		graphs[currentID] = graph
	}
	
	return graph
}

async function saveCurrentLayout(graph: GraphData, currentID: string, setSaving: (saving: boolean) => void) {
	setSaving(true)
	
	try {
		const response = await fetch('data/save?id=' + encodeURIComponent(currentID), {
			method: 'post',
			body: graph.exportSVG()
		})
		
		if (response.status !== 202) {
			alert('Error saving\nSee terminal output.')
		} else {
			graph.setSaved()
		}
	} catch (error) {
		console.error('Save failed:', error)
		alert('Save failed. See console for details.')
	} finally {
		setSaving(false)
	}
}

function removeEmptyProps(obj: any) {
	return JSON.parse(JSON.stringify(obj))
}

function camelToWords(camel: string) {
	const split = camel.replace(/([A-Z])/g, " $1")
	return split.charAt(0).toUpperCase() + split.slice(1)
}