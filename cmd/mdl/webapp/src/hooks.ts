import { useState, useCallback, useEffect } from 'react';
import { GraphData } from './graph-view/graph';
import { parseView } from './parseModel';
import { LayoutOptions } from './graph-view/layout';
import { 
  findShortcut, 
  HELP, 
  SAVE, 
  TOGGLE_DRAG_MODE,
  ALIGN_HORIZONTAL,
  ALIGN_VERTICAL,
  DISTRIBUTE_HORIZONTAL,
  DISTRIBUTE_VERTICAL,
  AUTO_LAYOUT,
  RESET_POSITION,
  TOGGLE_GRID,
  TOGGLE_SNAP_TO_GRID,
  SNAP_ALL_TO_GRID,
  MOVE_LEFT,
  MOVE_RIGHT,
  MOVE_UP,
  MOVE_DOWN,
  MOVE_LEFT_FINE,
  MOVE_RIGHT_FINE,
  MOVE_UP_FINE,
  MOVE_DOWN_FINE
} from './shortcuts';

// Global state for graphs to preserve edits
const graphs: { [key: string]: GraphData } = {};

// Custom hook for graph management
export const useGraph = (model: any, layouts: any, currentID: string): GraphData | null => {
  if (graphs[currentID]) {
    return graphs[currentID];
  }
  
  const graph = parseView(model, layouts, currentID);
  if (graph) {
    graphs[currentID] = graph;
  }
  
  return graph;
};

// Custom hook for auto layout functionality
export const useAutoLayout = (graph: GraphData) => {
  const [layouting, setLayouting] = useState(false);

  const handleAutoLayout = useCallback(async () => {
    setLayouting(true);
    try {
      const options: LayoutOptions = {
        direction: 'DOWN',
        compactLayout: true
      };
      await graph.autoLayout(options);
    } catch (error) {
      console.error('Layout failed:', error);
      alert('Layout failed. See console for details.');
    } finally {
      setLayouting(false);
    }
  }, [graph]);

  return { layouting, handleAutoLayout };
};

// Custom hook for save functionality
export const useSave = (graph: GraphData, currentID: string) => {
  const [saving, setSaving] = useState(false);

  const handleSave = useCallback(async () => {
    setSaving(true);
    
    try {
      const response = await fetch('data/save?id=' + encodeURIComponent(currentID), {
        method: 'post',
        body: graph.exportSVG()
      });
      
      if (response.status !== 202) {
        alert('Error saving\nSee terminal output.');
      } else {
        graph.setSaved();
      }
    } catch (error) {
      console.error('Save failed:', error);
      alert('Save failed. See console for details.');
    } finally {
      setSaving(false);
    }
  }, [graph, currentID]);

  return { saving, handleSave };
};

// Custom hook for keyboard shortcuts
export const useKeyboardShortcuts = (
  toggleHelp: () => void,
  saveLayout: () => void,
  graph?: GraphData,
  dragMode?: 'pan' | 'select',
  setDragMode?: (mode: 'pan' | 'select') => void,
  onAutoLayout?: () => void
) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const shortcut = findShortcut(e);
      
      // Prevent browser default for all recognized shortcuts
      if (shortcut) {
        e.preventDefault();
      }
      
      if (shortcut === HELP) {
        toggleHelp();
      } else if (shortcut === SAVE) {
        saveLayout();
      } else if (shortcut === TOGGLE_DRAG_MODE && setDragMode && dragMode) {
        setDragMode(dragMode === 'pan' ? 'select' : 'pan');
      } else if (graph) {
        // Graph-dependent shortcuts
        if (shortcut === ALIGN_HORIZONTAL) {
          graph.alignSelectionH();
        } else if (shortcut === ALIGN_VERTICAL) {
          graph.alignSelectionV();
        } else if (shortcut === DISTRIBUTE_HORIZONTAL) {
          graph.distributeSelectionH();
        } else if (shortcut === DISTRIBUTE_VERTICAL) {
          graph.distributeSelectionV();
        } else if (shortcut === AUTO_LAYOUT && onAutoLayout) {
          onAutoLayout();
        } else if (shortcut === RESET_POSITION) {
          graph.resetView();
        } else if (shortcut === TOGGLE_GRID) {
          graph.toggleGrid();
        } else if (shortcut === TOGGLE_SNAP_TO_GRID) {
          graph.toggleSnapToGrid();
        } else if (shortcut === SNAP_ALL_TO_GRID) {
          graph.snapAllToGrid();
        } else if (shortcut === MOVE_LEFT) {
          graph.moveSelected(-graph.getGridSize(), 0);
        } else if (shortcut === MOVE_LEFT_FINE) {
          graph.moveSelected(-1, 0, true); // Disable snap for fine movement
        } else if (shortcut === MOVE_RIGHT) {
          graph.moveSelected(graph.getGridSize(), 0);
        } else if (shortcut === MOVE_RIGHT_FINE) {
          graph.moveSelected(1, 0, true); // Disable snap for fine movement
        } else if (shortcut === MOVE_UP) {
          graph.moveSelected(0, -graph.getGridSize());
        } else if (shortcut === MOVE_UP_FINE) {
          graph.moveSelected(0, -1, true); // Disable snap for fine movement
        } else if (shortcut === MOVE_DOWN) {
          graph.moveSelected(0, graph.getGridSize());
        } else if (shortcut === MOVE_DOWN_FINE) {
          graph.moveSelected(0, 1, true); // Disable snap for fine movement
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [toggleHelp, saveLayout, graph, dragMode, setDragMode, onAutoLayout]);
};

// Utility function to clear graph cache
export const clearGraphCache = (currentID?: string) => {
  if (currentID) {
    delete graphs[currentID];
  } else {
    Object.keys(graphs).forEach(key => delete graphs[key]);
  }
}; 