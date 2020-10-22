import React, {FC, useEffect, useRef, useState} from "react";
import {buildGraph, GraphData} from "./graph";

interface Props {
	data: GraphData
	zoom: number
	onInit: () => void
	onSelect: (nodeName: string) => void
}

export const Graph: FC<Props> = ({data, zoom, onInit, onSelect}) => {
	const [graphState, setGraphState] = useState(null)
	useEffect(() => {
		if (graphState) return;
		const g = buildGraph(data, n => onSelect(n ? n.id: null))
		ref.current.append(g.svg)
		setGraphState(g)
		onInit()
	}, [])
	if (graphState) {
		graphState.setZoom(zoom)
	}
	const ref = useRef<HTMLDivElement>()
	return <div className="graph" ref={ref}/>
}