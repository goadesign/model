import {GraphData, Node, Group} from "./graph";
import ELK from "elkjs/lib/elk.bundled.js";

// Layout algorithm options
export type LayoutAlgorithm = 'layered' | 'stress' | 'mrtree' | 'force' | 'radial' | 'disco';

export interface LayoutOptions {
	algorithm?: LayoutAlgorithm;
	direction?: 'UP' | 'DOWN' | 'LEFT' | 'RIGHT';
	nodeSpacing?: number;
	layerSpacing?: number;
	edgeRouting?: 'POLYLINE' | 'ORTHOGONAL' | 'SPLINES';
	favorStraightEdges?: boolean;
	compactLayout?: boolean;
}

// Systematic spacing configuration
interface SpacingConfig {
	// Base spacing values (these are the source of truth)
	nodeSpacing: number;
	layerSpacing: number;
	componentSpacing: number;
	padding: number;
	
	// Context-specific multipliers
	groupMultiplier: number;
	
	// Algorithm-specific adjustments (minimal, only when absolutely necessary)
	algorithmAdjustments?: {
		nodeSpacing?: number;
		layerSpacing?: number;
		componentSpacing?: number;
	};
}

// Centralized spacing configurations for each algorithm
const SPACING_CONFIGS: Record<LayoutAlgorithm, SpacingConfig> = {
	layered: {
		nodeSpacing: 60,
		layerSpacing: 80,
		componentSpacing: 40,
		padding: 50,
		groupMultiplier: 0.8, // Slightly tighter in groups
	},
	
	stress: {
		nodeSpacing: 80,
		layerSpacing: 120,
		componentSpacing: 60,
		padding: 50,
		groupMultiplier: 0.7, // Tighter in groups to prevent excessive spread
	},
	
	force: {
		nodeSpacing: 20, // Very small node spacing for force layout
		layerSpacing: 25, // Very reduced layer spacing
		componentSpacing: 15, // Very reduced component spacing
		padding: 10, // Very reduced padding
		groupMultiplier: 0.4, // Much tighter spacing in groups
	},
	
	mrtree: {
		nodeSpacing: 70,
		layerSpacing: 90,
		componentSpacing: 50,
		padding: 50,
		groupMultiplier: 0.8,
	},
	
	radial: {
		nodeSpacing: 90,
		layerSpacing: 120,
		componentSpacing: 60,
		padding: 50,
		groupMultiplier: 0.7,
	},
	
	disco: {
		nodeSpacing: 100,
		layerSpacing: 150,
		componentSpacing: 80,
		padding: 50,
		groupMultiplier: 0.8,
	}
};

// Helper function to get effective spacing for a context
function getEffectiveSpacing(
	algorithm: LayoutAlgorithm, 
	userOptions: LayoutOptions = {},
	isGroup: boolean = false
): SpacingConfig {
	const baseConfig = SPACING_CONFIGS[algorithm];
	
	// Apply user overrides to base config
	const effectiveConfig: SpacingConfig = {
		nodeSpacing: userOptions.nodeSpacing ?? baseConfig.nodeSpacing,
		layerSpacing: userOptions.layerSpacing ?? baseConfig.layerSpacing,
		componentSpacing: baseConfig.componentSpacing,
		padding: baseConfig.padding,
		groupMultiplier: baseConfig.groupMultiplier,
		algorithmAdjustments: baseConfig.algorithmAdjustments
	};
	
	// Apply group multiplier if in group context
	if (isGroup) {
		effectiveConfig.nodeSpacing = Math.max(
			effectiveConfig.nodeSpacing * effectiveConfig.groupMultiplier,
			20 // Much tighter minimum spacing in groups
		);
		effectiveConfig.layerSpacing = Math.max(
			effectiveConfig.layerSpacing * effectiveConfig.groupMultiplier,
			25 // Much tighter minimum layer spacing in groups
		);
		// Also reduce component spacing and padding for groups
		effectiveConfig.componentSpacing = Math.max(
			effectiveConfig.componentSpacing * effectiveConfig.groupMultiplier,
			15
		);
		effectiveConfig.padding = Math.max(
			effectiveConfig.padding * effectiveConfig.groupMultiplier,
			15
		);
	}
	
	return effectiveConfig;
}

// Helper function to build ELK layout options from spacing config
function buildELKLayoutOptions(
	algorithm: LayoutAlgorithm,
	spacing: SpacingConfig,
	userOptions: LayoutOptions
): Record<string, string> {
	const {
		direction = 'DOWN',
		edgeRouting = 'ORTHOGONAL',
		favorStraightEdges = true,
		compactLayout = true
	} = userOptions;
	
	const baseOptions: Record<string, string> = {
		'elk.algorithm': algorithm,
		'elk.direction': direction,
		'elk.spacing.nodeNode': spacing.nodeSpacing.toString(),
		'elk.spacing.componentComponent': spacing.componentSpacing.toString(),
		'elk.padding': `[top=${spacing.padding},left=${spacing.padding},bottom=${spacing.padding},right=${spacing.padding}]`,
	};
	
	// Algorithm-specific options (only essential ones that can't be generalized)
	switch (algorithm) {
		case 'layered':
			return {
				...baseOptions,
				'elk.layered.spacing.nodeNodeBetweenLayers': spacing.layerSpacing.toString(),
				'elk.layered.spacing.edgeNodeBetweenLayers': '20',
				'elk.edgeRouting': edgeRouting,
				'elk.layered.edgeRouting.selfLoopDistribution': 'EQUALLY',
				'elk.layered.edgeRouting.selfLoopOrdering': 'STACKED',
				'elk.layered.nodePlacement.favorStraightEdges': favorStraightEdges.toString(),
				'elk.layered.nodePlacement.linearSegmentsDeflectionDampening': '0.3',
				'elk.layered.nodePlacement.strategy': 'BRANDES_KOEPF',
				'elk.layered.cycleBreaking.strategy': 'DEPTH_FIRST',
				...(compactLayout && {
					'elk.layered.compaction.postCompaction.strategy': 'EDGE_LENGTH',
					'elk.layered.compaction.connectedComponents': 'true',
				}),
				'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
				'elk.layered.crossingMinimization.greedySwitch.type': 'TWO_SIDED',
			};
			
		case 'stress':
			return {
				...baseOptions,
				'elk.stress.desiredEdgeLength': spacing.layerSpacing.toString(),
				'elk.stress.dimension': 'XY',
				'elk.stress.epsilon': '0.01',
				'elk.stress.iterationLimit': '1000',
			};
			
		case 'force':
			return {
				...baseOptions,
				'elk.force.model': 'FRUCHTERMAN_REINGOLD', // Use Fruchterman-Reingold model
				'elk.force.iterations': '2000', // More iterations for better convergence
				'elk.force.repulsivePower': '0.1', // Much lower repulsive power to allow closer positioning
				'elk.force.temperature': '0.01', // Higher temperature for more movement
				'elk.separateConnectedComponents': 'false', // Don't separate components - keep everything together
			};
			
		case 'mrtree':
			return {
				...baseOptions,
				'elk.mrtree.searchOrder': 'DFS',
				'elk.mrtree.weighting': 'DESCENDANTS',
			};
			
		case 'radial':
			return {
				...baseOptions,
				'elk.radial.radius': '1000',
				'elk.radial.compactor': 'IMPROVE_STRAIGHTNESS',
				'elk.radial.wedgeCriteria': 'RAYLIKE',
				'elk.radial.optimizationCriteria': 'COMPACTNESS',
			};
			
		case 'disco':
			return {
				...baseOptions,
				'elk.disco.componentCompaction.strategy': 'IMPROVE_STRAIGHTNESS',
			};
			
		default:
			return baseOptions;
	}
}

// Create ELK instance
const elk = new ELK();

export async function autoLayout(graph: GraphData, options: LayoutOptions = {}): Promise<{
	nodes: Array<{id: string, x: number, y: number}>,
	edges: Array<{id: string, vertices: Array<{x: number, y: number}>, label?: {x: number, y: number}}>
}> {
	const algorithm = options.algorithm || 'force';
	
	// Get systematic spacing configuration
	const rootSpacing = getEffectiveSpacing(algorithm, options, false);
	
	// Build ELK graph structure
	const elkGraph = {
		id: "root",
		layoutOptions: buildELKLayoutOptions(algorithm, rootSpacing, options),
		children: [] as any[],
		edges: [] as any[]
	};

	// Add ONLY actual nodes to ELK graph (exclude edge vertices)
	const nodeMap = new Map<string, Node>();
	graph.nodesMap.forEach(node => {
		nodeMap.set(node.id, node);
		elkGraph.children.push({
			id: node.id,
			width: node.width,
			height: node.height,
			// Don't provide existing positions to allow ELK to optimize layout
		});
	});

	// Handle grouped nodes (hierarchical layout) - MOVED BEFORE EDGE PROCESSING
	const processedGroups = new Set<string>();
	const nodeToGroupMap = new Map<string, string>(); // Track which group each node belongs to
	
	graph.groupsMap.forEach(group => {
		if (!processedGroups.has(group.id)) {
			processedGroups.add(group.id);
			
			// Get group-specific spacing (applies group multiplier)
			const groupSpacing = getEffectiveSpacing(algorithm, options, true);
			
			// Create group node
			const groupNode = {
				id: group.id,
				children: [] as any[],
				edges: [] as any[],
				layoutOptions: buildELKLayoutOptions(algorithm, groupSpacing, options)
			};

			// Add group members as children and track membership
			group.nodes.forEach(member => {
				if (!isGroup(member)) {
					nodeToGroupMap.set(member.id, group.id);
					// Move node from root to group
					const nodeIndex = elkGraph.children.findIndex(n => n.id === member.id);
					if (nodeIndex >= 0) {
						const node = elkGraph.children.splice(nodeIndex, 1)[0];
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

	// Add edges to ELK graph with proper hierarchical placement
	graph.edges.forEach(edge => {
		const fromGroup = nodeToGroupMap.get(edge.from.id);
		const toGroup = nodeToGroupMap.get(edge.to.id);
		
		const edgeData = {
			id: edge.id,
			sources: [edge.from.id],
			targets: [edge.to.id],
			// Add edge labels if they exist
			labels: edge.label ? [{
				id: `${edge.id}-label`,
				text: edge.label,
				width: Math.max(edge.label.length * 8, 60), // Better label width estimation
				height: 25   // Estimated label height
			}] : []
		};

		// Determine where to place the edge based on node group membership
		if (fromGroup && toGroup && fromGroup === toGroup) {
			// Both nodes are in the same group - add edge to that group
			const groupNode = elkGraph.children.find(child => child.id === fromGroup);
			if (groupNode && groupNode.edges) {
				groupNode.edges.push(edgeData);
			}
		} else {
			// Nodes are in different groups or at least one is external - add to root level
			elkGraph.edges.push(edgeData);
		}
	});

	try {
		// Perform layout calculation
		const layoutedGraph = await elk.layout(elkGraph);

		// Extract results
		const nodes: Array<{id: string, x: number, y: number}> = [];
		const edges: Array<{id: string, vertices: Array<{x: number, y: number}>, label?: {x: number, y: number}}> = [];

		// Process nodes (including nested group nodes)
		const extractNodes = (container: any, offsetX = 0, offsetY = 0) => {
			if (container.children) {
				container.children.forEach((child: any) => {
					if (child.children) {
						// This is a group node, recurse into it
						extractNodes(child, (child.x || 0) + offsetX, (child.y || 0) + offsetY);
					} else {
						// This is a regular node
						nodes.push({
							id: child.id,
							x: (child.x || 0) + offsetX,
							y: (child.y || 0) + offsetY
						});
					}
				});
			}
		};

		extractNodes(layoutedGraph);

		// Process edges with proper vertex handling from both root and group levels
		const processEdges = (container: any) => {
			if (container.edges) {
				container.edges.forEach((edge: any) => {
					const edgeResult: any = {
						id: edge.id,
						vertices: []
					};

					// Extract edge routing points from ELK sections
					if (edge.sections && edge.sections.length > 0) {
						const section = edge.sections[0];
						
						// Build vertex list from ELK routing
						const routingPoints: Array<{x: number, y: number}> = [];
						
						// Add start point if available
						if (section.startPoint) {
							routingPoints.push({
								x: section.startPoint.x,
								y: section.startPoint.y
							});
						}

						// Add bend points (these are the actual routing vertices)
						if (section.bendPoints && section.bendPoints.length > 0) {
							section.bendPoints.forEach((point: any) => {
								routingPoints.push({
									x: point.x,
									y: point.y
								});
							});
						}

						// Add end point if available
						if (section.endPoint) {
							routingPoints.push({
								x: section.endPoint.x,
								y: section.endPoint.y
							});
						}

						// Only keep intermediate points as vertices (exclude start/end which are node connection points)
						if (routingPoints.length > 2) {
							edgeResult.vertices = routingPoints.slice(1, -1);
						}
					}

					// Extract label position
					if (edge.labels && edge.labels.length > 0) {
						const label = edge.labels[0];
						edgeResult.label = {
							x: label.x || 0,
							y: label.y || 0
						};
					}

					edges.push(edgeResult);
				});
			}
			
			// Recursively process edges in child groups
			if (container.children) {
				container.children.forEach((child: any) => {
					if (child.edges) {
						processEdges(child);
					}
				});
			}
		};

		// Process edges from root level and all group levels
		processEdges(layoutedGraph);

		return { nodes, edges };

	} catch (error) {
		console.error('ELK layout failed:', error);
		// Fallback to a simple grid layout
		return createFallbackLayout(graph);
	}
}

// Fallback layout in case ELK fails
function createFallbackLayout(graph: GraphData): {
	nodes: Array<{id: string, x: number, y: number}>,
	edges: Array<{id: string, vertices: Array<{x: number, y: number}>}>
} {
	const nodes: Array<{id: string, x: number, y: number}> = [];
	const edges: Array<{id: string, vertices: Array<{x: number, y: number}>}> = [];

	// Simple grid layout as fallback
	let x = 100, y = 100;
	const spacing = 150;
	let nodesInRow = 0;
	const maxNodesPerRow = Math.ceil(Math.sqrt(graph.nodesMap.size));

	graph.nodesMap.forEach(node => {
		nodes.push({ id: node.id, x, y });
		
		nodesInRow++;
		if (nodesInRow >= maxNodesPerRow) {
			x = 100;
			y += spacing;
			nodesInRow = 0;
		} else {
			x += spacing;
		}
	});

	// Simple straight-line edges
	graph.edges.forEach(edge => {
		const fromNode = nodes.find(n => n.id === edge.from.id);
		const toNode = nodes.find(n => n.id === edge.to.id);
		
		if (fromNode && toNode) {
			edges.push({
				id: edge.id,
				vertices: [
					{ x: fromNode.x, y: fromNode.y },
					{ x: toNode.x, y: toNode.y }
				]
			});
		}
	});

	return { nodes, edges };
}

function isGroup(member: Node | Group): member is Group {
	return 'nodes' in member;
}