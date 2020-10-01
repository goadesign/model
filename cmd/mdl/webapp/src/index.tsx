import ReactDOM from 'react-dom';
import React from 'react';
import {Root} from "./Root";
import './style.css';

Promise.all([
	fetch('data/model.json').then(r => r.json()),
	fetch('data/layout.json').then(r => r.json())])
.then(([model, layout]) => {
	ReactDOM.render(<Root model={model} layout={layout}/>, document.getElementById('root'));
})
