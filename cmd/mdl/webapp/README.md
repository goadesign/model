# Model Diagram Layout (MDL) Webapp

This webapp provides an interactive interface for viewing and editing model diagrams with automatic layout capabilities.

## Features

### Interactive Diagram Editing

The webapp provides intuitive mouse interactions for navigating and editing diagrams:

#### Mouse Interactions

- **Pan View**: Drag on empty space to pan around the diagram
- **Zoom**: Use scroll wheel to zoom in/out at cursor position
- **Select Elements**: Click on nodes or edges to select them
- **Multi-Selection**: Shift+click to add/remove elements from selection
- **Box Selection**: Shift+drag on empty space to select multiple elements
- **Move Elements**: Drag selected elements to reposition them
- **Undo/Redo**: All changes are tracked for easy reversal

#### Toolbar Features

- **Alignment Tools**: Align selected elements horizontally or vertically
- **Distribution Tools**: Evenly distribute selected elements with equal spacing
- **Auto Layout**: Apply automatic layout algorithms to organize the entire diagram
- **Connection Routing**: Change how connections are drawn (orthogonal, straight, curved)
- **Zoom Controls**: Zoom in/out, fit to screen, or reset to 100%
- **Save**: Save the current layout and positions

### Auto Layout Engine

The webapp now uses **ELK.js** (Eclipse Layout Kernel) for automatic diagram layout, providing significant improvements over the previous Dagre.js implementation:

#### Available Layout Algorithms

1. **Layered (Hierarchical)** - Default algorithm, ideal for directed graphs with clear hierarchy
   - Optimized for node-link diagrams with inherent direction
   - Excellent edge routing with minimal crossings
   - Supports grouping and nested structures

2. **Stress (Force-directed)** - Physics-based layout for general graphs
   - Good for understanding overall graph structure
   - Minimizes edge lengths and overlaps

3. **Tree** - Specialized for tree structures
   - Optimized spacing for hierarchical data
   - Clean, readable layouts for tree-like diagrams

4. **Force** - Alternative force-directed algorithm
   - Different physics simulation approach
   - Good for dense graphs

5. **Radial** - Circular layout with central focus
   - Places important nodes at the center
   - Good for showing relationships radiating from key elements

6. **Disco** - Disconnected components layout
   - Handles graphs with multiple disconnected parts
   - Optimizes space usage for complex diagrams

#### Layout Features

- **Improved Space Optimization**: Better node spacing and edge routing
- **Orthogonal Edge Routing**: Cleaner, more readable edge paths
- **Group Support**: Proper handling of nested diagram elements
- **Configurable Options**: Adjustable spacing, direction, and algorithm-specific parameters
- **Fallback Handling**: Graceful degradation if layout fails

#### Usage

1. Select your preferred layout algorithm from the dropdown menu
2. Click "Auto Layout" to apply the selected algorithm
3. The system will automatically optimize node positions and edge routing
4. Use "Fit" to zoom and center the resulting layout

### Keyboard Shortcuts

#### Navigation & Zoom
- **Scroll Wheel**: Zoom in/out
- **Ctrl + =**: Zoom in
- **Ctrl + -**: Zoom out
- **Ctrl + 9**: Fit diagram to screen
- **Ctrl + 0**: Reset zoom to 100%

#### Selection & Editing
- **Ctrl + A**: Select all elements
- **Esc**: Deselect all elements
- **Delete/Backspace**: Remove selected edge vertices
- **Arrow Keys**: Move selected elements (1px)
- **Shift + Arrow Keys**: Move selected elements (10px)

#### File Operations
- **Ctrl + S**: Save current layout
- **Ctrl + Z**: Undo last change
- **Ctrl + Shift + Z** or **Ctrl + Y**: Redo last undone change

#### Help
- **Shift + ?** or **Shift + F1**: Show/hide keyboard shortcuts help

#### Advanced Editing
- **Alt + Click**: Add vertex to relationship edge (Option + Click on Mac)
- **Alt + Shift + Click**: Add label anchor to relationship edge (Option + Shift + Click on Mac)

## Technical Implementation

- **ELK.js v0.10.0**: Modern layout algorithms with TypeScript support
- **Async Layout**: Non-blocking layout calculation with progress indication
- **Error Handling**: Robust fallback mechanisms for layout failures
- **Performance**: Optimized for both small and large diagrams

## Development

```bash
npm install
npm start    # Development server
npm run build # Production build
```

## Migration from Dagre.js

The migration from Dagre.js to ELK.js provides:
- **Better algorithms**: More sophisticated layout techniques
- **Improved edge routing**: Cleaner, more readable connections
- **Enhanced grouping**: Better support for nested elements
- **Modern codebase**: Active development and TypeScript support
- **Multiple algorithms**: Choose the best layout for your diagram type

## Editor webapp

The mdl editor is a web application that runs in the browser.
It is embedded in the go executable, so it can be served
as static files. `mdl` runs an HTTP server that serves the
editor as static files. `mdl` also serves the model and layout data
dynamically, for the editor to load. The editor then renders the
selected view as an SVG, allowing the user to edit positions of elements
and shapes of relationships.

### Development setup

To develop the mdl and the editor, start the TypeScript compiler in watch mode and run mdl go program in devmode.

`mdl` can be instructed to serve the editor files from disk instead
of the embedded copies, to allow for easy development.
```
DEVMODE=1 go run ./cmd/mdl ... mdl params
``` 

Compile and run the TypeScript application in watch mode
```
yarn install
yarn watch
```

`yarn watch` will watch for changes in the webapp files and recompile. 
Simply refresh the browser to see the changes.
