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

// Default spacing configuration
const DEFAULT_SPACING: SpacingConfig = {
	nodeSpacing: 80,
	layerSpacing: 100,
	componentSpacing: 40,
	padding: 50,
	groupMultiplier: 0.6,
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
			20
		);
		effectiveConfig.layerSpacing = Math.max(
			effectiveConfig.layerSpacing * effectiveConfig.groupMultiplier,
			25
		);
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

// Simplified ELK layout options builder
function getELKOptions(
	spacing: SpacingConfig,
	userOptions: LayoutOptions
): Record<string, string> {
	const {
		direction = 'DOWN',
		compactLayout = true
	} = userOptions;
	
	const baseOptions: Record<string, string> = {
		'elk.algorithm': 'layered',
		'elk.direction': direction,
		'elk.spacing.nodeNode': spacing.nodeSpacing.toString(),
		'elk.spacing.componentComponent': spacing.componentSpacing.toString(),
		'elk.padding': `[top=${spacing.padding},left=${spacing.padding},bottom=${spacing.padding},right=${spacing.padding}]`,
		'elk.layered.spacing.nodeNodeBetweenLayers': spacing.layerSpacing.toString(),
		'elk.layered.spacing.edgeNodeBetweenLayers': '40',
		'elk.edgeRouting': 'POLYLINE',
		'elk.layered.edgeRouting.selfLoopDistribution': 'EQUALLY',
		'elk.layered.edgeRouting.selfLoopOrdering': 'STACKED',
		'elk.layered.nodePlacement.favorStraightEdges': 'true',
		'elk.layered.nodePlacement.linearSegmentsDeflectionDampening': '0.3',
		'elk.layered.nodePlacement.strategy': 'NETWORK_SIMPLEX',
		'elk.layered.cycleBreaking.strategy': 'GREEDY',
		'elk.layered.crossingMinimization.strategy': 'LAYER_SWEEP',
	};
	
	if (compactLayout) {
		baseOptions['elk.layered.compaction.postCompaction.strategy'] = 'EDGE_LENGTH';
		baseOptions['elk.layered.compaction.connectedComponents'] = 'true';
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
		nodeMap.set(node.id, node);
		elkGraph.children.push({
			id: node.id,
			width: node.width,
			height: node.height,
		});
	});

	// Handle grouped nodes (hierarchical layout)
	const processedGroups = new Set<string>();
	const nodeToGroupMap = new Map<string, string>();
	
	graph.groupsMap.forEach(group => {
		if (!processedGroups.has(group.id)) {
			processedGroups.add(group.id);
			
			// Get group-specific spacing
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
		
		const elkEdge = {
			id: edge.id,
			sources: [edge.from.id],
			targets: [edge.to.id],
			sections: edge.vertices?.length ? [{
				startPoint: {x: edge.vertices[0].x, y: edge.vertices[0].y},
				endPoint: {x: edge.vertices[edge.vertices.length - 1].x, y: edge.vertices[edge.vertices.length - 1].y},
				bendPoints: edge.vertices.slice(1, -1).map(v => ({x: v.x, y: v.y}))
			}] : undefined
		};

		// Place edge in appropriate container
		if (fromGroup && toGroup && fromGroup === toGroup) {
			// Both nodes in same group
			const groupContainer = elkGraph.children.find(c => c.id === fromGroup);
			if (groupContainer) {
				groupContainer.edges.push(elkEdge);
			}
		} else {
			// Cross-group or root-level edge
			elkGraph.edges.push(elkEdge);
		}
	});

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

		// Extract edges from layout result
		const processEdges = (container: any) => {
			container.edges?.forEach((edge: any) => {
				const vertices: Array<{x: number, y: number}> = [];
				let label: {x: number, y: number} | undefined;

				if (edge.sections && edge.sections.length > 0) {
					const section = edge.sections[0];
					
					// Add start point
					if (section.startPoint) {
						vertices.push({x: section.startPoint.x, y: section.startPoint.y});
					}
					
					// Add bend points
					if (section.bendPoints) {
						section.bendPoints.forEach((bp: any) => {
							vertices.push({x: bp.x, y: bp.y});
						});
					}
					
					// Add end point
					if (section.endPoint) {
						vertices.push({x: section.endPoint.x, y: section.endPoint.y});
					}

					// Calculate label position (middle of edge)
					if (vertices.length >= 2) {
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
				}

				edges.push({
					id: edge.id,
					vertices,
					label
				});
			});

			// Process child containers
			container.children?.forEach((child: any) => {
				if (child.edges) {
					processEdges(child);
				}
			});
		};

		extractNodes(layoutedGraph);
		processEdges(layoutedGraph);

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
		console.warn('ELK layout failed, using fallback:', error);
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