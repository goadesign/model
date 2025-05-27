import React, {FC, useEffect, useRef, useState} from "react";
import {buildGraph, getZoomAuto, GraphData, Node, setZoom, addCursorInteraction, restoreViewState, saveViewState} from "./graph";

interface Props {
	data: GraphData;
	onSelect: (nodeName: string | null) => void;
	dragMode: 'pan' | 'select';
}

export const Graph: FC<Props> = ({data, onSelect, dragMode}) => {
	const [graphState, setGraphState] = useState<any>(null);
	const ref = useRef<HTMLDivElement>(null);

	// Single effect for building the graph and handling all setup/cleanup
	useEffect(() => {
		if (!ref.current) return;

		// Clear previous content
		ref.current.innerHTML = '';

		// Build graph with current props
		const g = buildGraph(data, (n: Node | null) => onSelect(n ? n.id : null), dragMode);
		ref.current.append(g.svg);
		setGraphState(g);
		
		// Try to restore previous view state, otherwise use auto zoom
		if (!restoreViewState(data.id)) {
			setZoom(getZoomAuto());
		}

		// Save view state before page unload
		const handleBeforeUnload = () => {
			if (data?.id) {
				saveViewState(data.id);
			}
		};
		
		window.addEventListener('beforeunload', handleBeforeUnload);

		return () => {
			// Save view state before cleanup
			if (data?.id) {
				saveViewState(data.id);
			}
			
			window.removeEventListener('beforeunload', handleBeforeUnload);
			
			if (ref.current) {
				ref.current.innerHTML = '';
			}
		};
	}, [data, onSelect, dragMode]);

	// Effect for updating drag mode on existing graph
	useEffect(() => {
		if (graphState?.svg) {
			// Clean up existing cursor interaction
			const svg = graphState.svg;
			const existingCleanup = (svg as any).__cursorInteractionCleanup;
			if (existingCleanup) {
				existingCleanup();
			}
			
			// Set up cursor interaction with current drag mode
			addCursorInteraction(svg, dragMode);
		}
	}, [dragMode, graphState]);

	return <div className="graph" ref={ref}/>;
}