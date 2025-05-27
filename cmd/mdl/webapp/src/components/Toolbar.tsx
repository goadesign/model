import React, { FC, useState } from 'react';
import { getZoomAuto, GraphData, setZoom, getZoom } from '../graph-view/graph';
import { listViews } from '../parseModel';
import { camelToWords } from '../utils';
import { getModifierKeyName } from '../utils/platform';

// Types
interface ToolbarProps {
  model: any;
  currentID: string;
  onViewChange: (id: string) => void;
  graph: GraphData;
  onAutoLayout: () => void;
  onSave: () => void;
  onToggleHelp: () => void;
  saving: boolean;
  layouting: boolean;
  dragMode: 'pan' | 'select';
  setDragMode: (mode: 'pan' | 'select') => void;
}

export const Toolbar: FC<ToolbarProps> = ({
  model, currentID, onViewChange, graph, 
  onAutoLayout, onSave, onToggleHelp, saving, layouting,
  dragMode, setDragMode
}) => {
  const views = listViews(model);
  
  return (
    <div className="toolbar">
      <ViewSelector 
        views={views}
        currentID={currentID}
        onViewChange={onViewChange}
      />
      <ToolbarActions
        graph={graph}
        onAutoLayout={onAutoLayout}
        onSave={onSave}
        onToggleHelp={onToggleHelp}
        saving={saving}
        layouting={layouting}
        dragMode={dragMode}
        setDragMode={setDragMode}
      />
    </div>
  );
};

const ViewSelector: FC<{
  views: any[];
  currentID: string;
  onViewChange: (id: string) => void;
}> = ({ views, currentID, onViewChange }) => (
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
      <span style={{ marginLeft: '8px', fontWeight: 'bold' }}>
        {views[0] ? camelToWords(views[0].section) + ': ' + views[0].title : 'No views available'}
      </span>
    )}
  </div>
);

const ToolbarActions: FC<{
  graph: GraphData;
  onAutoLayout: () => void;
  onSave: () => void;
  onToggleHelp: () => void;
  saving: boolean;
  layouting: boolean;
  dragMode: 'pan' | 'select';
  setDragMode: (mode: 'pan' | 'select') => void;
}> = ({
  graph, onAutoLayout, onSave, onToggleHelp, saving, layouting,
  dragMode, setDragMode
}) => (
  <div style={{ display: 'flex', alignItems: 'center' }}>
    <div className="toolbar-group">
      <DragModeButton dragMode={dragMode} setDragMode={setDragMode} />
    </div>
    <div className="toolbar-group">
      <UndoRedoButtons graph={graph} />
    </div>
    <div className="toolbar-group">
      <AlignmentButtons graph={graph} />
    </div>
    <div className="toolbar-group">
      <LayoutControls onAutoLayout={onAutoLayout} layouting={layouting} />
    </div>
    <div className="toolbar-group">
      <GridControls graph={graph} />
    </div>
    <div className="toolbar-group">
      <ZoomControls graph={graph} />
    </div>
    <div className="toolbar-group">
      <SaveButton onSave={onSave} saving={saving} />
    </div>
    <div className="toolbar-group">
      <HelpButton onToggleHelp={onToggleHelp} />
    </div>
  </div>
);

const DragModeButton: FC<{
  dragMode: 'pan' | 'select';
  setDragMode: (mode: 'pan' | 'select') => void;
}> = ({ dragMode, setDragMode }) => (
  <button 
    className={`mode-toggle ${dragMode === 'select' ? 'select-mode' : 'pan-mode'}`}
    onClick={() => setDragMode(dragMode === 'pan' ? 'select' : 'pan')} 
    data-tooltip={dragMode === 'pan' ? "Pan Mode: Drag to pan the view (T)" : "Select Mode: Drag to select elements, Shift+click to add/remove selection (T)"}
  >
    {dragMode === 'pan' ? <i className="fas fa-hand-paper"></i> : <i className="fas fa-mouse-pointer"></i>}
  </button>
);

const UndoRedoButtons: FC<{ graph: GraphData }> = ({ graph }) => {
  const modKey = getModifierKeyName();
  return (
    <>
      <button onClick={() => graph.undo()} data-tooltip={`Undo the last change made to the diagram (${modKey}+Z)`}>
        <i className="fas fa-undo"></i>
      </button>
      <button onClick={() => graph.redo()} data-tooltip={`Redo the last undone action (${modKey}+Shift+Z / ${modKey}+Y)`}>
        <i className="fas fa-redo"></i>
      </button>
    </>
  );
};

const AlignmentButtons: FC<{ graph: GraphData }> = ({ graph }) => {
  const modKey = getModifierKeyName();
  return (
    <>
      <button onClick={() => graph.alignSelectionH()} data-tooltip={`Align all selected elements horizontally (left edges) (${modKey}+Shift+H)`}>
        <i className="fas fa-arrows-alt-v"></i>
      </button>
      <button onClick={() => graph.alignSelectionV()} data-tooltip={`Align all selected elements vertically (top edges) (${modKey}+Shift+A)`}>
        <i className="fas fa-arrows-alt-h"></i>
      </button>
      <button onClick={() => graph.distributeSelectionH()} data-tooltip={`Distribute selected elements evenly horizontally (equal spacing) (${modKey}+Alt+H)`}>
        <i className="fas fa-ellipsis-h"></i>
      </button>
      <button onClick={() => graph.distributeSelectionV()} data-tooltip={`Distribute selected elements evenly vertically (equal spacing) (${modKey}+Alt+V)`}>
        <i className="fas fa-ellipsis-v"></i>
      </button>
    </>
  );
};

const LayoutControls: FC<{
  onAutoLayout: () => void;
  layouting: boolean;
}> = ({ onAutoLayout, layouting }) => {
  const modKey = getModifierKeyName();
  return (
    <button 
      onClick={onAutoLayout} 
      disabled={layouting} 
      data-tooltip={`Automatically arrange all elements using the Layered algorithm (${modKey}+L)`}
    >
      {layouting ? <i className="fas fa-spinner fa-spin"></i> : <i className="fas fa-magic"></i>}
    </button>
  );
};

const GridControls: FC<{ graph: GraphData }> = ({ graph }) => {
  const [gridVisible, setGridVisible] = useState(graph.isGridVisible());
  const [snapToGrid, setSnapToGrid] = useState(graph.isSnapToGrid());
  const modKey = getModifierKeyName();
  
  // Update state when graph changes or when grid state changes via shortcuts
  React.useEffect(() => {
    const updateGridState = () => {
      setGridVisible(graph.isGridVisible());
      setSnapToGrid(graph.isSnapToGrid());
    };
    
    // Initial update
    updateGridState();
    
    // Listen for grid state changes from keyboard shortcuts
    window.addEventListener('gridStateChanged', updateGridState);
    
    return () => {
      window.removeEventListener('gridStateChanged', updateGridState);
    };
  }, [graph]);
  
  const handleToggleGrid = () => {
    graph.toggleGrid();
    setGridVisible(graph.isGridVisible());
  };
  
  const handleToggleSnap = () => {
    graph.toggleSnapToGrid();
    setSnapToGrid(graph.isSnapToGrid());
  };
  
  const handleSnapAll = () => {
    graph.snapAllToGrid();
  };
  
  return (
    <>
      <button 
        className={gridVisible ? 'active-toggle' : ''}
        onClick={handleToggleGrid} 
        data-tooltip={`Toggle grid visibility (${modKey}+G)`}
      >
        <i className="fas fa-th"></i>
      </button>
      <button 
        className={snapToGrid ? 'active-toggle' : ''}
        onClick={handleToggleSnap} 
        data-tooltip={`Toggle snap to grid (${modKey}+Shift+G)`}
      >
        <i className="fas fa-magnet"></i>
      </button>
      <button 
        onClick={handleSnapAll} 
        disabled={!snapToGrid}
        data-tooltip={`Snap all elements to grid (${modKey}+Alt+G)`}
      >
        <i className="fas fa-border-all"></i>
      </button>
    </>
  );
};

const ZoomControls: FC<{ graph: GraphData }> = ({ graph }) => {
  const modKey = getModifierKeyName();
  return (
    <>
      <button onClick={() => setZoom(Math.max(0.1, getZoom() - 0.05))} data-tooltip={`Zoom out to see more of the diagram (${modKey}+-)`}>
        <i className="fas fa-search-minus"></i>
      </button>
      <button onClick={() => setZoom(Math.min(5, getZoom() + 0.05))} data-tooltip={`Zoom in to see details more clearly (${modKey}+=)`}>
        <i className="fas fa-search-plus"></i>
      </button>
      <button onClick={() => graph.fitToView()} data-tooltip={`Automatically fit the entire diagram in the visible area (${modKey}+9)`}>
        <i className="fas fa-expand-arrows-alt"></i>
      </button>
      <button onClick={() => setZoom(1)} data-tooltip={`Reset zoom to 100% (actual size) (${modKey}+0)`}>
        <i className="fas fa-search"></i>
      </button>
      <button onClick={() => { graph.alignTopLeft(); graph.resetPanTransform(); }} data-tooltip={`Reset position and view (${modKey}+Home)`}>
        <i className="fas fa-home"></i>
      </button>
    </>
  );
};

const SaveButton: FC<{
  onSave: () => void;
  saving: boolean;
}> = ({ onSave, saving }) => {
  const modKey = getModifierKeyName();
  return (
    <button className="action" disabled={saving} onClick={onSave} data-tooltip={`Save the current diagram layout (${modKey}+S)`}>
      {saving ? <i className="fas fa-spinner fa-spin"></i> : <i className="fas fa-save"></i>}
    </button>
  );
};

const HelpButton: FC<{
  onToggleHelp: () => void;
}> = ({ onToggleHelp }) => {
  return (
    <button onClick={onToggleHelp} data-tooltip="Show keyboard shortcuts and help information (Shift+? / Shift+F1)">
      <i className="fas fa-question-circle"></i>
    </button>
  );
}; 