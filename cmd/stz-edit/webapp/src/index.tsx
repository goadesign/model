import ReactDOM from 'react-dom';
import React from 'react';
import {Root} from "./Root";
import './style.css';
import {parseModel} from "./models";

Promise.all([
	fetch('data/model.json').then(r => r.json()),
	fetch('data/model.layout.json').then(r => r.json())])
.then(([modelJSON, layoutJSON]) => {
	const models = parseModel(modelJSON, layoutJSON)
	ReactDOM.render(<Root models={models}/>, document.getElementById('root'));
})
