import React, {FC, useState} from "react";
import {getZoom, getZoomAuto, GraphData, setZoom} from "./graph-view/graph";
import {Graph} from "./graph-view/graph-react";
import {BrowserRouter as Router, Routes, Route, useNavigate, useSearchParams} from 'react-router-dom'
import {listViews, parseView, ViewsList} from "./parseModel";
import {findShortcut, HELP, Help, SAVE} from "./shortcuts";


export const Root: FC<{model: any, layout: any}> = ({model, layout}) => <Router><Routes>
	<Route path="/" element={<ModelPane key={getCrtID()} model={model} layouts={layout}/>}/>
</Routes></Router>

const getCrtID = () => {
	const p = new URLSearchParams(document.location.search)
	return p.get('id') || ''
}


// we keep graphs here, in case they are edited but not saved
const graphs: {[key: string]: GraphData} = {}
export const refreshGraph = () => {
	delete graphs[getCrtID()]
}

let toggHelp: () => void
let saveLayout: () => void
window.addEventListener('keydown', e => {
	const shortcut = findShortcut(e)
	if (toggHelp && shortcut == HELP) {
		toggHelp()
	} else if (saveLayout && shortcut == SAVE) {
		saveLayout()
		e.preventDefault()
	}
})

const ModelPane: FC<{model: any, layouts: any}> = ({model, layouts}) => {

	const [saving, setSaving] = useState(false)
	const [helpOn, setHelpOn] = useState(false)

	const [searchParams, setSearchParams] = useSearchParams()
	const crtID = searchParams.get('id') || ''

	const graph = graphs[crtID] || parseView(model, layouts, crtID)
	if (!graph) {
		const lst = listViews(model)
		document.location.href = '?id=' + lst[0].key
		return <>Redirecting to {lst[0].title}</>
	}
	graphs[crtID] = graph

	saveLayout = () => {
		setSaving(true)

		fetch('data/save?id=' + encodeURIComponent(crtID), {
			method: 'post',
			body: graph.exportSVG()
		}).then(ret => {
			if (ret.status != 202) {
				alert('Error saving\nSee terminal output.')
			}
			setSaving(false)
			graph.setSaved()
		})
	}

	toggHelp = () => setHelpOn(!helpOn)

	return <>
		<div className="toolbar">
			<div>
				View:
				<select
					onChange={e => setSearchParams({id: encodeURIComponent(e.target.value)})} value={crtID}>
					<option disabled value="" hidden>...</option>
					{listViews(model).map(m => <option key={m.key}
											value={m.key}>{camelToWords(m.section) + ': ' + m.title}</option>)}
				</select>
			</div>
			<div>
				<button onClick={() => graph.undo()} title="Undo last change">Undo</button>
				<button className="grp" onClick={() => graph.redo()} title="Redo undone actions">Redo</button>

				<button onClick={() => graph.alignSelectionH()} title="Align selected objects horizontally">H Align</button>
				<button className="grp" onClick={() => graph.alignSelectionV()} title="Align selected objects vertically">V Align</button>
				<button className="grp" onClick={() => graph.autoLayout()} title="Automatic layout using DagreJS">Auto Layout</button>
				<button onClick={() => setZoom(getZoom() - .05)} title="Zoom out">Zoom -</button>
				<button onClick={() => setZoom(getZoom() + .05)} title="Zoom in">Zoom +</button>
				<button onClick={() => {
					graph.alignTopLeft()
					setZoom(getZoomAuto())
				}} title="Zoom/Move to make all graph visible">Fit</button>
				<button onClick={() => setZoom(1)}>Zoom 100%</button>
				<button className="action" disabled={saving} onClick={() => saveLayout()}>Save View</button>
				<button onClick={() => setHelpOn(!helpOn)}>Help</button>
			</div>
		</div>
		<Graph key={crtID}
			   data={graph}
			   // print metadata in console
			   onSelect={id => id && console.log(removeEmptyProps(graph.metadata.elements.find((m: any) => m.id == id)))}
			/>
		{helpOn && <Help/>}
	</>
}

function removeEmptyProps(o: any) {
	return JSON.parse(JSON.stringify(o))
}

function camelToWords(camel: string) {
	let split = camel.replace( /([A-Z])/g, " $1" );
	return split.charAt(0).toUpperCase() + split.slice(1);
}