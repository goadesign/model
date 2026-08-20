# Diagram geometry

MDL uses one path from model data to SVG:

1. Draw each node once so the browser can measure its real text and shape.
2. Build a measured diagram with exact node, boundary-title, and edge-label
   sizes.
3. Ask ELK to place the whole diagram, including nested boundaries and
   relationship labels, in one run.
4. Check the complete result.
5. Draw the checked positions without moving them.

`geometry.ts` defines the data passed between these steps. Drawing code accepts
only `ValidatedLayout`, which can be created only by the geometry checks.

## Automatic and manual layout

`automatic-layout.ts` owns automatic positions. It passes measured sizes and
declared ownership to ELK and copies ELK's result without repairs.

`manual-layout.ts` owns saved positions. It checks one complete saved layout
against the current measured diagram. It never fills missing manual values with
automatic ones.

The two paths meet only after validation.

## What validation checks

Validation rejects:

- missing, old, or duplicate element IDs;
- invalid numbers or changed measured sizes;
- node, boundary, title, label, and route collisions;
- false boundary ownership or partial boundary overlap;
- incomplete routes or endpoints that miss their nodes;
- routes through unrelated nodes or boundaries; and
- any geometry outside the view.

All problems from one pass are reported together. There is no retry, grid
fallback, label search, or post-layout movement.

## Headless SVG generation

`mdl svg` loads `headless.html`, not the React editor. The page renders one
requested view and posts a typed result containing its view ID and model
digest. Go accepts only the matching result and replaces the SVG atomically.
One browser process and one tab are reused across all requested views.
