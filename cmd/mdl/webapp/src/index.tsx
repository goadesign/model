import { createRoot } from 'react-dom/client';
import React, { Suspense, lazy, useEffect, useState } from 'react';
import { refreshGraph } from "./Root";
import './style.css';
import '@fortawesome/fontawesome-free/css/all.css';
import { RefreshConnector } from "./websocket";

const Root = lazy(() => import('./Root').then(module => ({ default: module.Root })));

interface ModelData {
	model: any;
	layout: any;
}

interface AppState {
	data: ModelData | null;
	error: string | null;
	loading: boolean;
}

const App: React.FC = () => {
	const [state, setState] = useState<AppState>({
		data: null,
		error: null,
		loading: true
	});

	const loadData = async () => {
		setState(prev => ({ ...prev, loading: true, error: null }));
		
		try {
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

			setState({
				data: { model, layout },
				error: null,
				loading: false
			});
		} catch (error) {
			console.error('Failed to load data:', error);
			setState({
				data: null,
				error: error instanceof Error ? error.message : 'Unknown error occurred',
				loading: false
			});
		}
	};

	const handleFileChange = (path: string) => {
		if (path.endsWith('.svg')) {
			return; // Ignore SVG changes to avoid infinite loops
		}
		
		console.log('File changed:', path);
		refreshGraph();
		loadData();
	};

	useEffect(() => {
		// Setup refresh connector
		const refreshConnector = new RefreshConnector(handleFileChange);
		refreshConnector.connect();

		// Initial data load
		loadData();

		// Cleanup function
		return () => {
			// RefreshConnector cleanup would go here if it had a disconnect method
		};
	}, []);

	if (state.loading) {
		return <LoadingScreen />;
	}

	if (state.error) {
		return <ErrorScreen error={state.error} onRetry={loadData} />;
	}

	if (!state.data) {
		return <ErrorScreen error="No data available" onRetry={loadData} />;
	}

	return (
		<Suspense fallback={<LoadingScreen />}>
			<Root model={state.data.model} layout={state.data.layout} />
		</Suspense>
	);
};

const LoadingScreen: React.FC = () => (
	<div style={{
		display: 'flex',
		justifyContent: 'center',
		alignItems: 'center',
		height: '100vh',
		fontFamily: 'Arial, sans-serif'
	}}>
		<div>Loading...</div>
	</div>
);

const ErrorScreen: React.FC<{ error: string; onRetry: () => void }> = ({ error, onRetry }) => (
	<div style={{
		padding: '20px',
		color: 'red',
		fontFamily: 'monospace',
		whiteSpace: 'pre-wrap',
		display: 'flex',
		flexDirection: 'column',
		alignItems: 'center',
		justifyContent: 'center',
		height: '100vh'
	}}>
		<h2>Error loading application</h2>
		<p>{error}</p>
		<button 
			onClick={onRetry}
			style={{
				padding: '10px 20px',
				fontSize: '16px',
				cursor: 'pointer',
				backgroundColor: '#007bff',
				color: 'white',
				border: 'none',
				borderRadius: '4px'
			}}
		>
			Retry
		</button>
	</div>
);

// Initialize the application
const container = document.getElementById('root');
if (!container) {
	throw new Error('Root container not found');
}

const root = createRoot(container);
root.render(<App />);