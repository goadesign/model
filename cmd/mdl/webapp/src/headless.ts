// This page renders one requested view for mdl svg. It does not start React,
// routing, editor shortcuts, saved-layout loading, or LiveReload.

import "@fortawesome/fontawesome-free/css/all.css";
import "./fonts.css";
import "./style.css";

import {buildGraphView, LayoutDirection} from "./graph-view/graph";
import {LayoutOptions} from "./graph-view/layout";
import {parseView} from "./parseModel";

interface RenderRequest {
	viewId: string;
	modelDigest: string;
	options: LayoutOptions;
}

type RenderResult =
	| {
		status: "complete";
		viewId: string;
		modelDigest: string;
		svg: string;
	}
	| {
		status: "error";
		viewId: string;
		modelDigest: string;
		error: string;
	};

/** run renders one view and reports one result to the local MDL process. */
async function run(): Promise<void> {
	let request: RenderRequest;
	try {
		request = readRequest();
	} catch (error) {
		throw new Error(`invalid headless request: ${errorMessage(error)}`);
	}

	try {
		const response = await fetch("/data/model.json");
		if (!response.ok) {
			throw new Error(`load model: HTTP ${response.status}`);
		}
		const model = await response.json();
		const graph = parseView(model, {}, request.viewId);
		if (!graph) {
			throw new Error(`model does not contain view ${request.viewId}`);
		}
		const svg = buildGraphView(graph);
		document.body.append(svg);
		await graph.autoLayout(request.options);
		await sendResult({
			status: "complete",
			viewId: request.viewId,
			modelDigest: request.modelDigest,
			svg: graph.exportSVG(),
		});
	} catch (error) {
		await sendResult({
			status: "error",
			viewId: request.viewId,
			modelDigest: request.modelDigest,
			error: errorMessage(error),
		});
	}
}

/** readRequest reads the view, model digest, and explicit layout choices. */
function readRequest(): RenderRequest {
	const params = new URLSearchParams(document.location.search);
	const viewId = params.get("view") ?? "";
	const modelDigest = params.get("digest") ?? "";
	if (!viewId || !modelDigest) {
		throw new Error("view and digest are required");
	}
	const direction = readDirection(params.get("direction"));
	return {
		viewId,
		modelDigest,
		options: {
			direction,
			compactLayout: params.get("compact") === "true",
		},
	};
}

/** readDirection accepts only the four directions supported by the Model DSL. */
function readDirection(value: string | null): LayoutDirection | undefined {
	switch (value) {
		case null:
		case "":
			return undefined;
		case "UP":
		case "DOWN":
		case "LEFT":
		case "RIGHT":
			return value;
		default:
			throw new Error(`invalid direction ${value}`);
	}
}

/** sendResult sends the complete typed result to the waiting Go request. */
async function sendResult(result: RenderResult): Promise<void> {
	const response = await fetch("/headless/result", {
		method: "POST",
		headers: {"Content-Type": "application/json"},
		body: JSON.stringify(result),
	});
	if (!response.ok) {
		throw new Error(`report result: HTTP ${response.status} ${await response.text()}`);
	}
}

/** errorMessage includes the stack so CLI failures name the browser code. */
function errorMessage(error: unknown): string {
	return error instanceof Error ? error.stack ?? error.message : String(error);
}

void run();
