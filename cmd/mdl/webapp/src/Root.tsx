import React, { FC, useState, useCallback, useEffect } from "react";
import { GraphData } from "./graph-view/graph";
import { Graph } from "./graph-view/graph-react";
import { BrowserRouter as Router, Routes, Route, useSearchParams } from 'react-router-dom';
import { listViews } from "./parseModel";
import { Help } from "./shortcuts";
import { useGraph, useAutoLayout, useSave, useKeyboardShortcuts, clearGraphCache } from "./hooks";
import { Toolbar } from "./components/Toolbar";
import { removeEmptyProps, getCurrentViewID } from "./utils";

// Types
interface ModelData {
  model: any;
  layout: any;
}

export const Root: FC<ModelData> = ({ model, layout }) => (
  <Router>
    <Routes>
      <Route path="/" element={<ModelPane model={model} layouts={layout} />} />
    </Routes>
  </Router>
);

export const refreshGraph = () => {
  const currentID = getCurrentViewID();
  clearGraphCache(currentID);
};

const ModelPane: FC<{ model: any; layouts: any }> = ({ model, layouts }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const currentID = decodeURI(searchParams.get('id') || '');
  
  // UI State
  const [helpVisible, setHelpVisible] = useState(false);
  const [dragMode, setDragMode] = useState<'pan' | 'select'>('pan');
  
  // Get or create graph for current view
  const graph = useGraph(model, layouts, currentID);
  
  // Custom hooks for functionality
  const { layouting, handleAutoLayout } = useAutoLayout(graph || ({} as GraphData));
  const { saving, handleSave } = useSave(graph || ({} as GraphData), currentID);
  
  if (!graph) {
    return <ViewRedirect model={model} />;
  }

  const handleToggleHelp = useCallback(() => {
    setHelpVisible(!helpVisible);
  }, [helpVisible]);

  // Update document title when view changes
  useEffect(() => {
    if (graph && graph.name) {
      document.title = `${graph.name} - Model`;
    }
  }, [graph]);

  // Setup keyboard shortcuts
  useKeyboardShortcuts(handleToggleHelp, handleSave, graph, dragMode, setDragMode, handleAutoLayout);

  const handleViewChange = useCallback((id: string) => {
    setSearchParams({ id: encodeURIComponent(id) });
  }, [setSearchParams]);

  const handleSelect = useCallback((id: string | null) => {
    if (id) {
      const element = graph.metadata.elements.find((m: any) => m.id === id);
      console.log(removeEmptyProps(element));
    }
  }, [graph]);

	return (
		<>
			<Toolbar
				model={model}
				currentID={currentID}
				onViewChange={handleViewChange}
				graph={graph}
				onAutoLayout={handleAutoLayout}
				onSave={handleSave}
				onToggleHelp={handleToggleHelp}
				saving={saving}
				layouting={layouting}
				dragMode={dragMode}
				setDragMode={setDragMode}
			/>
			<Graph 
				key={currentID}
				data={graph}
				onSelect={handleSelect}
				dragMode={dragMode}
			/>
			{helpVisible && <Help />}
		</>
	);
};

const ViewRedirect: FC<{ model: any }> = ({ model }) => {
  const views = listViews(model);
  
  React.useEffect(() => {
    // Set default title when no view is selected
    document.title = 'Model - Architecture Diagrams as Code';
    
    if (views.length > 0) {
      document.location.href = '?id=' + views[0].key;
    }
  }, [views]);

  if (views.length > 0) {
    return <>Redirecting to {views[0].title}</>;
  }
  return <>No views available</>;
};

