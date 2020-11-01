import ReactDOM from 'react-dom';
import React from 'react';
import {refreshGraph, Root} from "./Root";
import './style.css';
import {RefreshConnector} from "./websocket";

function reload() {
	Promise.all([
		fetch('data/model.json').then(r => r.json()),
		fetch('data/layout.json').then(r => r.json())])
		.then(([model, layout]) => {
			ReactDOM.render(<Root model={model} layout={layout}/>, document.getElementById('root'));
		})
}

let c = new RefreshConnector(path => {
	console.log('file changed:', path)
	refreshGraph()
	reload()
})
c.connect()

reload()