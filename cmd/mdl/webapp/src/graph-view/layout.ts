import {GraphData, Node, Group} from "./graph";

export interface LayoutOptions {
	direction?: 'UP' | 'DOWN' | 'LEFT' | 'RIGHT';
	nodeSpacing?: number;
	layerSpacing?: number;
	compactLayout?: boolean;
}

// Simplified spacing configuration
interface SpacingConfig {
	nodeSpacing: number;
	layerSpacing: number;
	componentSpacing: number;
	padding: number;
	groupMultiplier: number;
}

// Balanced spacing configuration - optimized for clean routing
const DEFAULT_SPACING: SpacingConfig = {
	nodeSpacing: 120,      // More compact horizontal spacing
	layerSpacing: 70,      // Tighter vertical spacing between layers  
	componentSpacing: 60,  // Closer separation between disconnected components
	padding: 30,           // Less padding around the entire layout
	groupMultiplier: 0.65, // Moderate compaction within groups
};

// Helper function to get effective spacing for a context
function getEffectiveSpacing(
	userOptions: LayoutOptions = {},
	isGroup: boolean = false
): SpacingConfig {
	// Apply user overrides to base config
	const effectiveConfig: SpacingConfig = {
		nodeSpacing: userOptions.nodeSpacing ?? DEFAULT_SPACING.nodeSpacing,
		layerSpacing: userOptions.layerSpacing ?? DEFAULT_SPACING.layerSpacing,
		componentSpacing: DEFAULT_SPACING.componentSpacing,
		padding: DEFAULT_SPACING.padding,
		groupMultiplier: DEFAULT_SPACING.groupMultiplier,
	};
	
	// Apply group multiplier if in group context
	if (isGroup) {
		effectiveConfig.nodeSpacing = Math.max(
			effectiveConfig.nodeSpacing * effectiveConfig.groupMultiplier,
			30  // Minimum 30px spacing within groups
		);
		effectiveConfig.layerSpacing = Math.max(
			effectiveConfig.layerSpacing * effectiveConfig.groupMultiplier,
			35  // Minimum 35px layer spacing within groups
		);
		effectiveConfig.componentSpacing = Math.max(
			effectiveConfig.componentSpacing * effectiveConfig.groupMultiplier,
			25  // Minimum 25px component spacing within groups
		);
		effectiveConfig.padding = Math.max(
			effectiveConfig.padding * effectiveConfig.groupMultiplier,
			15  // Minimum 15px padding within groups
		);
	}
	
	return effectiveConfig;
}

// Simplified ELK layout options builder
function getELKOptions(
	spacing: SpacingConfig,
	userOptions: LayoutOptions
): Record<string, string> {
	const {
		direction = 'DOWN',
		compactLayout = false
	} = userOptions;
	
	const baseOptions: Record<string, string> = {
		'elk.algorithm': 'layered', // Back to layered for better orthogonal routing
		'elk.direction': direction,
		'elk.spacing.nodeNode': spacing.nodeSpacing.toString(),
		'elk.spacing.componentComponent': spacing.componentSpacing.toString(),
		'elk.padding': `[top=${spacing.padding},left=${spacing.padding},bottom=${spacing.padding},right=${spacing.padding}]`,
		
		// Layer spacing for clean routing
		'elk.layered.spacing.nodeNodeBetweenLayers': spacing.layerSpacing.toString(),
		'elk.layered.spacing.edgeNodeBetweenLayers': '40', // Tighter spacing around nodes
		'elk.layered.spacing.edgeEdgeBetweenLayers': '30',  // Tighter space between edges
		
		// ORTHOGONAL edge routing for cleaner layout
		'elk.edgeRouting': 'ORTHOGONAL',
		'elk.layered.unnecessaryBendpoints': 'false',
		
		// Orthogonal routing configuration - try to minimize detours
		'elk.layered.edgeRouting.orthogonal.mode': 'DIRECTION_BASED',
		'elk.layered.edgeRouting.orthogonal.spacing': '15', // Tighter edge spacing
		'elk.layered.edgeRouting.orthogonal.nodeOverlapRatio': '0.1',
		
		// Compaction options
		'elk.layered.compaction.connectedComponents': 'true',
		'elk.layered.compaction.postCompaction.strategy': 'LEFT_RIGHT',
		
		// Separate components to reduce complexity
		'elk.separateConnectedComponents': 'true',
		
		// Flatten hierarchy for better edge routing
		'elk.hierarchyHandling': 'SEPARATE_CHILDREN',
		'elk.layered.considerModelOrder.strategy': 'NONE', // Ignore model ordering constraints
		
		// Enhanced edge label handling - keep labels close to edges
		'elk.edgeLabels.placement': 'CENTER',
		'elk.edgeLabels.inline': 'false',
		'elk.spacing.edgeLabel': '20', // Reduced spacing - keep labels closer
		'elk.edgeLabels.avoidOverlap': 'false', // Disable aggressive collision avoidance
		'elk.edgeLabels.considerModelOrder': 'false',
		'elk.layered.edgeLabels.sideSelection': 'ALWAYS_DOWN', // Consistent label placement
	};
	
	// Additional compact layout options if requested
	if (compactLayout) {
		baseOptions['elk.spacing.nodeNode'] = Math.max(spacing.nodeSpacing * 0.7, 60).toString();
		baseOptions['elk.layered.spacing.nodeNodeBetweenLayers'] = Math.max(spacing.layerSpacing * 0.7, 80).toString();
	}
	
	return baseOptions;
}

export async function autoLayout(graph: GraphData, options: LayoutOptions = {}): Promise<{
	nodes: Array<{id: string, x: number, y: number}>,
	edges: Array<{id: string, vertices: Array<{x: number, y: number}>, label?: {x: number, y: number}}>
}> {
	// Dynamically import ELK only when auto-layout is used
	const ELK = await import('elkjs/lib/elk.bundled.js').then(module => module.default);
	const elk = new ELK();
	// Get systematic spacing configuration
	const rootSpacing = getEffectiveSpacing(options, false);
	
	// Build ELK graph structure
	const elkGraph = {
		id: "root",
		layoutOptions: getELKOptions(rootSpacing, options),
		children: [] as any[],
		edges: [] as any[]
	};

	// Add ONLY actual nodes to ELK graph (exclude edge vertices)
	const nodeMap = new Map<string, Node>();
	graph.nodesMap.forEach(node => {
		if (!node.id) return; // Skip nodes without IDs
		
		nodeMap.set(node.id, node);
		
		// Ensure minimum dimensions and validate node size data
		const nodeWidth = Math.max(node.width || 200, 150); // Min width 150px
		const nodeHeight = Math.max(node.height || 100, 80); // Min height 80px
		
		// Add padding to node dimensions for ELK to account for arrow size
		// This makes ELK route edges to a slightly larger boundary
		const arrowPadding = 40; // About 8 arrow lengths (arrow is 5px in markerWidth)
		
		elkGraph.children.push({
			id: node.id,
			// Provide current position as hint to ELK
			x: node.x,
			y: node.y,
			width: nodeWidth + (arrowPadding * 2),
			height: nodeHeight + (arrowPadding * 2),
			layoutOptions: {
				// Allow ELK to move nodes but consider current positions
				'elk.position': '',
				// Force ELK to use our exact dimensions
				'elk.nodeSize.constraints': '[FIXED_SIZE]'
			}
		});
	});

	// EXPERIMENTAL: Skip group processing - treat all nodes as flat
	// This should eliminate group-based routing detours

	// Add all edges to root level for optimal ELK routing visibility
	let addedEdges = 0;
	graph.edges.forEach(edge => {
		// Skip edges without proper IDs
		if (!edge.id || !edge.from?.id || !edge.to?.id) return;
		
		// Verify source and target nodes exist in our node map
		if (!nodeMap.has(edge.from.id) || !nodeMap.has(edge.to.id)) {
			console.warn(`Skipping edge ${edge.id}: source ${edge.from.id} or target ${edge.to.id} not found in nodes`);
			return;
		}
		
		addedEdges++;
		
		// Calculate more accurate label dimensions for ELK
		const labelWidth = edge.label && edge.label.trim() ? 
			Math.min(edge.label.length * 7, 200) : 0; // More realistic width estimate
		
		
		const elkEdge = {
			id: edge.id,
			sources: [edge.from.id],
			targets: [edge.to.id],
			// Include label information with much smaller dimensions
			labels: edge.label && edge.label.trim() ? [{
				id: `${edge.id}-label`,
				text: edge.label,
				// Much more conservative label size estimates
				width: labelWidth,
				height: 20, // More realistic label height
				layoutOptions: {
					'elk.edgeLabels.placement': 'CENTER',
					'elk.edgeLabels.inline': 'false'
					// Remove the FIXED_SIZE constraint that might be forcing detours
				}
			}] : []
		};

		// All edges at root level for maximum ELK visibility and routing
		elkGraph.edges.push(elkEdge);
	});

	// Enhanced validation - ensure ELK gets complete data
	if (!elkGraph.id || !elkGraph.children) {
		throw new Error('Invalid ELK graph structure');
	}
	

	try {
		const layoutedGraph = await elk.layout(elkGraph);
		
		// Extract results
		const nodes: Array<{id: string, x: number, y: number}> = [];
		const edges: Array<{id: string, vertices: Array<{x: number, y: number}>, label?: {x: number, y: number}}> = [];

		// Same padding we added to nodes for ELK
		const arrowPadding = 40;

		// Extract nodes from layout result
		const extractNodes = (container: any, offsetX = 0, offsetY = 0) => {
			container.children?.forEach((child: any) => {
				if (child.children) {
					// This is a group, recurse
					extractNodes(child, offsetX + (child.x || 0), offsetY + (child.y || 0));
				} else {
					// This is a node
					// Adjust for the padding we added - ELK positioned based on padded size
					// So we need to shift by the padding amount to get the real center
					nodes.push({
						id: child.id,
						x: offsetX + (child.x || 0) + (child.width || 0) / 2,
						y: offsetY + (child.y || 0) + (child.height || 0) / 2
					});
				}
			});
		};

		// Process edges from ELK layout result to get routing information
		const processEdgesFromELK = (container: any, offsetX = 0, offsetY = 0) => {
			container.edges?.forEach((elkEdge: any) => {
				const vertices: Array<{x: number, y: number}> = [];
				let label: {x: number, y: number} | undefined;

				// Process edge sections to get bend points
				if (elkEdge.sections && elkEdge.sections.length > 0) {
					elkEdge.sections.forEach((section: any, sectionIndex: number) => {
						// Add start point if it exists
						if (section.startPoint) {
							vertices.push({
								x: offsetX + section.startPoint.x, 
								y: offsetY + section.startPoint.y
							});
						}
						
						// Add bend points (this is where ELK puts the routing vertices!)
						if (section.bendPoints && section.bendPoints.length > 0) {
							section.bendPoints.forEach((bp: any, bpIndex: number) => {
								vertices.push({
									x: offsetX + bp.x, 
									y: offsetY + bp.y
								});
							});
						}
						
						// Add end point if it exists
						if (section.endPoint) {
							vertices.push({
								x: offsetX + section.endPoint.x, 
								y: offsetY + section.endPoint.y
							});
						}
					});
				}

				// Extract ELK's calculated label positions (respect collision avoidance!)
				const originalEdge = graph.edges.find(e => e.id === elkEdge.id);
				if (originalEdge?.label && originalEdge.label.trim()) {
					if (elkEdge.labels && elkEdge.labels.length > 0) {
						const elkLabel = elkEdge.labels[0]; // Get first label
						if (elkLabel.x !== undefined && elkLabel.y !== undefined) {
							label = {
								x: offsetX + elkLabel.x + (elkLabel.width || 0) / 2, // Center of label
								y: offsetY + elkLabel.y + (elkLabel.height || 0) / 2
							};
						}
					} else if (vertices.length >= 2) {
						// Fallback: use middle of edge if ELK didn't provide label position
						const midIndex = Math.floor(vertices.length / 2);
						if (vertices.length % 2 === 0) {
							const v1 = vertices[midIndex - 1];
							const v2 = vertices[midIndex];
							label = { x: (v1.x + v2.x) / 2, y: (v1.y + v2.y) / 2 };
						} else {
							label = vertices[midIndex];
						}
					}
				}


				edges.push({
					id: elkEdge.id,
					vertices,
					label
				});
			});
			
			// Also process edges in child containers (groups)
			container.children?.forEach((child: any) => {
				if (child.edges && child.edges.length > 0) {
					processEdgesFromELK(child, offsetX + (child.x || 0), offsetY + (child.y || 0));
				}
			});
		};


		extractNodes(layoutedGraph);
		processEdgesFromELK(layoutedGraph);

		// Normalize coordinates to start near (0,0) to prevent huge canvas sizes
		// while preserving relative positioning between elements
		if (nodes.length > 0) {
			// Find the minimum coordinates across all elements
			const minX = Math.min(...nodes.map(n => n.x));
			const minY = Math.min(...nodes.map(n => n.y));
			
			// Add some padding so content doesn't start at exact (0,0)
			const padding = 50;
			const offsetX = -minX + padding;
			const offsetY = -minY + padding;
			
			// Normalize all node positions
			nodes.forEach(node => {
				node.x += offsetX;
				node.y += offsetY;
			});
			
			// Normalize all edge positions
			edges.forEach(edge => {
				edge.vertices.forEach(vertex => {
					vertex.x += offsetX;
					vertex.y += offsetY;
				});
				if (edge.label) {
					edge.label.x += offsetX;
					edge.label.y += offsetY;
				}
			});
		}

		return { nodes, edges };

	} catch (error) {
		console.warn('ELK layout failed, using fallback layout. Error:', error);
		return createFallbackLayout(graph);
	}
}

// Simplified fallback layout
function createFallbackLayout(graph: GraphData): {
	nodes: Array<{id: string, x: number, y: number}>,
	edges: Array<{id: string, vertices: Array<{x: number, y: number}>}>
} {
	const nodes: Array<{id: string, x: number, y: number}> = [];
	const edges: Array<{id: string, vertices: Array<{x: number, y: number}>}> = [];

	// Simple grid layout for nodes
	let x = 0, y = 0;
	const spacing = 300;
	const maxCols = Math.ceil(Math.sqrt(graph.nodesMap.size));

	let col = 0;
	graph.nodesMap.forEach(node => {
		nodes.push({
			id: node.id,
			x: x,
			y: y
		});

		col++;
		if (col >= maxCols) {
			col = 0;
			x = 0;
			y += spacing;
		} else {
			x += spacing;
		}
	});

	// Simple straight line edges
	graph.edges.forEach(edge => {
		edges.push({
			id: edge.id,
			vertices: []
		});
	});

	return { nodes, edges };
}

// Removed unused isGroup function since we skip group processing