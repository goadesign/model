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
	nodeSpacing: 160,      // Compact horizontal spacing
	layerSpacing: 90,      // Balanced vertical spacing between layers  
	componentSpacing: 80,  // Tighter separation between disconnected components
	padding: 40,           // Less padding around the entire layout
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
		'elk.algorithm': 'layered',
		'elk.direction': direction,
		'elk.spacing.nodeNode': spacing.nodeSpacing.toString(),
		'elk.spacing.componentComponent': spacing.componentSpacing.toString(),
		'elk.padding': `[top=${spacing.padding},left=${spacing.padding},bottom=${spacing.padding},right=${spacing.padding}]`,
		
		// Layer spacing - balanced for clean routing
		'elk.layered.spacing.nodeNodeBetweenLayers': spacing.layerSpacing.toString(),
		'elk.layered.spacing.edgeNodeBetweenLayers': '70', // Sufficient space for edge routing around nodes
		'elk.layered.spacing.edgeEdgeBetweenLayers': '35',  // More space between parallel edges to prevent crowding
		
		// ORTHOGONAL edge routing with automatic vertex creation
		'elk.edgeRouting': 'ORTHOGONAL',
		'elk.layered.unnecessaryBendpoints': 'true', // Clean up unnecessary bend points
		
		// Improved orthogonal routing settings
		'elk.layered.edgeRouting.orthogonal.mode': 'BOX', // Try BOX mode for cleaner routing
		'elk.layered.edgeRouting.orthogonal.spacing': '30', // More spacing for cleaner orthogonal routing
		'elk.layered.edgeRouting.orthogonal.bendPoint': 'AUTO',
		
		// Standard crossing minimization
		'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
		'elk.layered.crossingMinimization.greedySwitch.type': 'TWO_SIDED', 
		'elk.layered.crossingMinimization.greedySwitchCrossingReduction': 'true',
		
		// Node placement optimized for edge routing
		'elk.layered.nodePlacement.strategy': 'NETWORK_SIMPLEX',
		'elk.layered.nodePlacement.favorStraightEdges': 'false',
		
		// Cycle breaking for complex graphs
		'elk.layered.cycleBreaking.strategy': 'GREEDY_MODEL_ORDER',
		
		// Prevent edge merging to maintain individual routing
		'elk.layered.mergeEdges': 'false',
		'elk.layered.mergeHierarchyEdges': 'false',
		
		// Separate components to reduce complexity
		'elk.separateConnectedComponents': 'true',
		
		// Enable hierarchical layout handling for groups
		'elk.hierarchyHandling': 'INCLUDE_CHILDREN',
		
		// Enhanced edge label handling - prevent overlaps
		'elk.edgeLabels.placement': 'CENTER',
		'elk.edgeLabels.inline': 'false',
		'elk.spacing.edgeLabel': '40', // More space around edge labels to prevent overlaps
		'elk.edgeLabels.considerModelOrder': 'true',
		'elk.edgeLabels.useShortLabels': 'false',
		
		// Additional orthogonal routing options for better vertex creation
		'elk.layered.edgeRouting.orthogonal.addUnnecessaryBendpoints': 'false',
		'elk.layered.edgeRouting.orthogonal.nodeOverlapRatio': '0.05' // Less overlap tolerance for cleaner routing
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
		elkGraph.children.push({
			id: node.id,
			// Provide current position as hint to ELK
			x: node.x,
			y: node.y,
			width: node.width || 200,
			height: node.height || 100,
			layoutOptions: {
				// Allow ELK to move nodes but consider current positions
				'elk.position': '',
				// Force ELK to use our exact dimensions
				'elk.nodeSize.constraints': '[FIXED_SIZE]'
			}
		});
	});

	// Handle grouped nodes hierarchically with ELK hierarchyHandling enabled
	const processedGroups = new Set<string>();
	const nodeToGroupMap = new Map<string, string>();
	
	graph.groupsMap.forEach(group => {
		if (!group.id) return; // Skip groups without IDs
		
		if (!processedGroups.has(group.id)) {
			processedGroups.add(group.id);
			
			// Get group-specific spacing with more compaction
			const groupSpacing = getEffectiveSpacing(options, true);
			
			// Create group node
			const groupNode = {
				id: group.id,
				children: [] as any[],
				edges: [] as any[],
				layoutOptions: getELKOptions(groupSpacing, options)
			};

			// Add group members as children and track membership  
			group.nodes.forEach(member => {
				if (!isGroup(member)) {
					nodeToGroupMap.set(member.id, group.id);
					// Move node from root to group
					const nodeIndex = elkGraph.children.findIndex(n => n.id === member.id);
					if (nodeIndex >= 0) {
						const node = elkGraph.children.splice(nodeIndex, 1)[0];
						// Ensure moved node has proper structure and dimensions
						if (!node.layoutOptions) {
							node.layoutOptions = {};
						}
						// Force ELK to use exact node dimensions in groups too
						node.layoutOptions['elk.nodeSize.constraints'] = '[FIXED_SIZE]';
						groupNode.children.push(node);
					}
				}
			});

			// Add group to root if it has children
			if (groupNode.children.length > 0) {
				elkGraph.children.push(groupNode);
			}
		}
	});

	// Add all edges to root level for optimal ELK routing visibility
	let addedEdges = 0;
	graph.edges.forEach(edge => {
		// Skip edges without proper IDs
		if (!edge.id || !edge.from?.id || !edge.to?.id) return;
		addedEdges++;
		
		const elkEdge = {
			id: edge.id,
			sources: [edge.from.id],
			targets: [edge.to.id],
			// Include comprehensive label information for ELK awareness
			labels: edge.label && edge.label.trim() ? [{
				id: `${edge.id}-label`,
				text: edge.label,
				// Estimate label dimensions for ELK layout calculations
				width: Math.max(edge.label.length * 8, 50), // Rough character width estimation
				height: 20, // Standard label height
				layoutOptions: {
					'elk.edgeLabels.placement': 'CENTER',
					'elk.edgeLabels.inline': 'false',
					// Force ELK to consider label dimensions
					'elk.nodeSize.constraints': '[FIXED_SIZE]'
				}
			}] : []
		};

		// All edges at root level for maximum ELK visibility and routing
		elkGraph.edges.push(elkEdge);
	});

	// Basic validation - ensure root graph structure is valid
	if (!elkGraph.id || !elkGraph.children) {
		throw new Error('Invalid ELK graph structure');
	}

	try {
		const layoutedGraph = await elk.layout(elkGraph);
		
		// Extract results
		const nodes: Array<{id: string, x: number, y: number}> = [];
		const edges: Array<{id: string, vertices: Array<{x: number, y: number}>, label?: {x: number, y: number}}> = [];

		// Extract nodes from layout result
		const extractNodes = (container: any, offsetX = 0, offsetY = 0) => {
			container.children?.forEach((child: any) => {
				if (child.children) {
					// This is a group, recurse
					extractNodes(child, offsetX + (child.x || 0), offsetY + (child.y || 0));
				} else {
					// This is a node
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

				// Find the original edge to get label info
				const originalEdge = graph.edges.find(e => e.id === elkEdge.id);
				
				// Calculate label position if edge has a label and vertices
				if (originalEdge?.label && originalEdge.label.trim() && vertices.length >= 2) {
					const midIndex = Math.floor(vertices.length / 2);
					if (vertices.length % 2 === 0) {
						// Even number of vertices, interpolate between middle two
						const v1 = vertices[midIndex - 1];
						const v2 = vertices[midIndex];
						label = {
							x: (v1.x + v2.x) / 2,
							y: (v1.y + v2.y) / 2
						};
					} else {
						// Odd number of vertices, use middle vertex
						label = vertices[midIndex];
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

function isGroup(member: Node | Group): member is Group {
	return 'nodes' in member;
}