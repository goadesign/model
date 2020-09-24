import React, {FC, useState} from "react";
import {getZoomAuto, GraphData} from "./graph-view/graph";
import {Graph} from "./graph-view/graph-react";
import {BrowserRouter as Router, Route} from 'react-router-dom'
import {useHistory} from "react-router";
import {layoutDx, layoutDy, layoutScale} from "./models";


export const Root: FC<{models: GraphData[]}> = ({models}) => <Router>
	<Route path="/" component={() => <ModelPane models={models}/>}/>
</Router>

const getCrtID = () => {
	const p = new URLSearchParams(document.location.search)
	return p.get('id') || ''
}

const DomainSelect: FC<{ models: GraphData[]; crtID: string}> = ({models, crtID}) => {
	const history = useHistory();
	return <select
		onChange={e => history.push('?id=' + encodeURIComponent(e.target.value))} value={crtID}>
		<option disabled value="" hidden>...</option>
		{models.map(m => <option key={m.id} value={m.id}>{m.name}</option>)}
	</select>
}

const ModelPane: FC<{models:GraphData[]}> = ({models}) => {
	const crtID = getCrtID()
	const [zoom, setZoom] = useState(1)
	const [saving, setSaving] = useState(false)

	const graph = models.find(o => o.id == crtID)
	if (!graph) {
		return <div style={{padding:30}}><DomainSelect models={models} crtID=""/></div>
	}

	function saveLayout() {
		setSaving(true)
		fetch('data/save?id=' + encodeURIComponent(crtID), {
			method: 'post',
			body: JSON.stringify(graph.exportLayout().map(({id, x, y}) => {
				return {id, x: (x - layoutDx)/layoutScale, y: (y-layoutDy)/layoutScale}
			}))
		}).then(ret => {
			if (ret.status != 202) {
				alert('Error saving')
			}
			setSaving(false)
		})
	}
	return <>
		<div className="toolbar">
			<div>
				View: <DomainSelect models={models} crtID={crtID}/>
				{' '}
			</div>
			<div>
				<button onClick={() => setZoom(zoom - .05)}>Zoom -</button>
				{' '}
				<button onClick={() => setZoom(zoom + .05)}>Zoom +</button>
				{' '}
				<button onClick={() => setZoom(getZoomAuto())}>Fit</button>
				{' '}
				<button onClick={() => setZoom(1)}>Zoom 100%</button>
				{' '}
				<button className="action" disabled={saving} onClick={() => saveLayout()}>Save Layout</button>
			</div>
		</div>
		<Graph key={crtID}
			   data={graph}
			   zoom={zoom}
			   onSelect={name => null}
			   onInit={() => setZoom(getZoomAuto())}/>
	</>
}

