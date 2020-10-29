import React, {FC, useEffect, useRef, useState} from "react";
import {buildGraph, getZoomAuto, GraphData, setZoom} from "./graph";

interface Props {
	data: GraphData
	onSelect: (nodeName: string) => void
}

export const Graph: FC<Props> = ({data, onSelect}) => {
	const [graphState, setGraphState] = useState(null)
	useEffect(() => {
		if (graphState) return;
		const g = buildGraph(data, n => onSelect(n ? n.id: null))
		ref.current.append(g.svg)
		setGraphState(g)
		setZoom(getZoomAuto())
	}, [])
	const ref = useRef<HTMLDivElement>()
	return <div className="graph" ref={ref}/>
}