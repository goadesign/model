// Manual layout has its own entry point. It never fills missing positions from
// automatic layout and never changes the positions a user authored.

import {
	MeasuredDiagram,
	ResolvedLayout,
	ValidatedLayout,
	validateResolvedLayout,
} from "./geometry";

/**
 * resolveManual checks one complete saved layout against the current measured
 * diagram. Missing or old IDs, incomplete routes, and collisions are errors.
 */
export function resolveManual(
	measured: MeasuredDiagram,
	completeSnapshot: ResolvedLayout,
): ValidatedLayout {
	return validateResolvedLayout(measured, completeSnapshot);
}
