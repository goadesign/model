import styled from "styled-components";
import React, {FC, ReactNode} from "react";

const Separator = styled.div`
	cursor: col-resize;
	background-color: #f0f0f0;
	background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='10' height='30'><path d='M2 0 v30 M5 0 v30 M8 0 v30' fill='none' stroke='black'/></svg>");
	background-repeat: no-repeat;
	background-position: center;
	width: 10px;
	height: 100%;
	user-select: none;
`

interface SplitPaneProps {
	left: ReactNode;
	right: ReactNode;
	leftPercent: number;
}
//dom based minimalistic split pane
function splitter(element: HTMLElement) {
	if (!element || element.dataset.splitterInited) return;
	element.dataset.splitterInited = 'true';

	const [first, , second] = Array.from(
		element.parentElement.children
	) as HTMLDivElement[];

	let md: {
		e: MouseEvent;
		offsetLeft: number;
		offsetTop: number;
		firstWidth: number;
		secondWidth: number;
	};

	function convertEvent(e: any) {
		if (e.changedTouches === undefined) return e;
		return e.changedTouches[0];
	}

	function onMouseMove(e: MouseEvent) {
		e = convertEvent(e);

		let delta = e.clientX - md.e.clientX;
		// prevent negative-sized elements
		delta = Math.min(Math.max(delta, -md.firstWidth), md.secondWidth);

		element.style.left = md.offsetLeft + delta + 'px';
		first.style.width = md.firstWidth + delta + 'px';
		second.style.width = md.secondWidth - delta + 'px';
	}

	function onMouseDown(e: MouseEvent) {
		e = convertEvent(e);

		md = {
			e,
			offsetLeft: element.offsetLeft,
			offsetTop: element.offsetTop,
			firstWidth: first.offsetWidth,
			secondWidth: second.offsetWidth
		};
		document.addEventListener('touchmove', onMouseMove);
		document.addEventListener('mousemove', onMouseMove);
		const remove = () => {
			document.removeEventListener('touchmove', onMouseMove);
			document.removeEventListener('mousemove', onMouseMove);
		};
		document.addEventListener('mouseup', remove);
		document.addEventListener('touchend', remove);
	}

	element.addEventListener('mousedown', onMouseDown);
	element.addEventListener('touchstart', onMouseDown);
}

export const SplitPane: FC<SplitPaneProps> = ({left, right, leftPercent}) => {
	return (
		<div style={{width: '100%', height: '100%', display: 'flex'}}>
			<div style={{overflow: 'auto', height: '100%', width: leftPercent+'%'}}>
				{left}
			</div>
			<Separator ref={(el) => splitter(el)}/>
			<div style={{position: 'relative', height: '100%', width: (100-leftPercent) + '%', overflow:'auto'}}>
				{right}
			</div>
		</div>
	)
}


