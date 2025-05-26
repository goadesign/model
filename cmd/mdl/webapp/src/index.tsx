import {createRoot} from 'react-dom/client';
import React from 'react';
import {refreshGraph, Root} from "./Root";
import './style.css';
import '@fortawesome/fontawesome-free/css/all.css';
import {RefreshConnector} from "./websocket";

interface ModelData {
	model: any;
	layout: any;
}

class App {
	private root: any;
	private refreshConnector: RefreshConnector;

	constructor() {
		this.initializeRoot();
		this.setupRefreshConnector();
		this.loadAndRender();
	}

	private initializeRoot() {
		const container = document.getElementById('root');
		if (!container) {
			throw new Error('Root container not found');
		}
		this.root = createRoot(container);
	}

	private setupRefreshConnector() {
		this.refreshConnector = new RefreshConnector(this.handleFileChange.bind(this));
		this.refreshConnector.connect();
	}

	private handleFileChange(path: string) {
		if (path.endsWith('.svg')) {
			return; // Ignore SVG changes to avoid infinite loops
		}
		
		console.log('File changed:', path);
		refreshGraph();
		this.loadAndRender();
	}

	private async loadAndRender() {
		try {
			const data = await this.fetchData();
			this.render(data);
		} catch (error) {
			console.error('Failed to load data:', error);
			this.renderError(error);
		}
	}

	private async fetchData(): Promise<ModelData> {
		const [modelResponse, layoutResponse] = await Promise.all([
			fetch('data/model.json'),
			fetch('data/layout.json')
		]);

		if (!modelResponse.ok) {
			throw new Error(`Failed to fetch model: ${modelResponse.statusText}`);
		}
		
		if (!layoutResponse.ok) {
			throw new Error(`Failed to fetch layout: ${layoutResponse.statusText}`);
		}

		const [model, layout] = await Promise.all([
			modelResponse.json(),
			layoutResponse.json()
		]);

		return { model, layout };
	}

	private render(data: ModelData) {
		this.root.render(<Root model={data.model} layout={data.layout}/>);
	}

	private renderError(error: any) {
		this.root.render(
			<div style={{
				padding: '20px',
				color: 'red',
				fontFamily: 'monospace',
				whiteSpace: 'pre-wrap'
			}}>
				<h2>Error loading application</h2>
				<p>{error.message || 'Unknown error occurred'}</p>
				<button onClick={() => this.loadAndRender()}>
					Retry
				</button>
			</div>
		);
	}
}

// Initialize the application
new App();