---
name: editing-model-diagrams
description: Creates, edits, reviews, and regenerates C4 architecture diagrams written with goa.design/model and the mdl CLI. Use when changing Model DSL, model.go or views.go files, system landscape, context, container, component, dynamic, or deployment views, element relationships, boundaries, layout, or generated SVG diagrams.
---

# Editing Model diagrams

Produce diagrams whose source model is architecturally true and whose rendered
views communicate that model without ambiguous ownership.

## Prioritize truth, readability, then polish

Use this order when design goals conflict:

1. Represent the real ownership, dependencies, directions, and runtime
   behavior.
2. Make that reality readable through the right abstraction level, focused
   views, deliberate layout, and clear labels.
3. Improve visual balance and polish without changing or hiding architectural
   meaning.

Never omit, reverse, reparent, or relabel architecture merely to make a diagram
look cleaner. If a truthful view is unreadable, reduce its question, remove
out-of-scope elements, or split it into additional truthful views while keeping
the main view representative of the whole system.

## Workflow

1. Read the model definitions, the affected view, and imported model packages.
2. State the one question the view answers and the expected elements,
   relationships, and boundaries.
3. Identify the owner of every element and the source-model relationship for
   every intended edge.
4. Choose the C4 view level that answers the question.
5. Preserve published view keys unless the requested change intentionally
   renames or removes an output.
6. Edit Model DSL source. Do not edit generated SVG or JSON output directly.
7. Regenerate every affected view.
8. Inspect the rendered diagram, not only the compiling DSL.

## Preserve ownership

- A `SoftwareSystem` owns the `Container` elements declared inside it.
- A `Container` owns the `Component` elements declared inside it.
- `Add`, `AddDefault`, `AddAll`, and imported packages change view membership;
  they do not change element ownership.
- A deployment node may contain infrastructure nodes, child deployment nodes,
  and container instances. A container instance represents deployment of its
  referenced container; it does not transfer software ownership.
- Do not redefine or reparent an element to make a layout easier. Correct the
  model first, then select and arrange the view.
- A view cannot create a relationship absent from the source model. Do not
  invent or reverse an edge to complete a desired narrative; report or correct
  the model contract when evidence supports that change.
- When a repository inventory reports a missing service or element, verify its
  ownership and typed callers before adding it. Model the real container and
  relationships; do not add a name-only placeholder merely to satisfy a check.

## Verify Goa service coverage

For repositories whose services are defined by a Goa system design, use the
Goa model validator as the service-inventory contract:

- Verify the system design calls
  `goa.design/plugins/v3/model/dsl.Model(<model package>, <system name>)`.
  A `goa.design/model` dependency or an MDL model package alone does not enable
  this validation.
- Map every Goa service to its owned model container. Use `ModelContainer` when
  a human-readable container name does not exactly match the plugin's naming
  format; do not rename model elements or rely on fuzzy string matching.
- Use `ModelNone` only when the Goa service is deliberately outside the
  architecture model's scope, and document why. Never use it merely to make
  generation pass.
- Use `ModelComplete` only when every in-scope model container must correspond
  to a Goa service. Omit it when the model intentionally includes workers,
  infrastructure, data stores, or other containers that are not Goa services.
- Run the repository's Goa generation or model-validation command after model
  changes. MDL rendering does not execute the Goa service-to-container check.
- Treat plugin failures as architecture drift. Verify the service's ownership
  and behavior, add or correct the real container, and then add the explicit
  service mapping.

Service coverage and view membership are separate. Every owned service must be
present in the source model, and the published view set should make each
architecturally relevant service visible in at least one purposeful view. Do
not force every service into every view or create an inventory-only diagram;
split the architecture into focused views and use documentation or a generated
catalog for exhaustive inventory.

## Enforce truthful boundaries

Treat a rendered boundary as an ownership statement.

- System boundaries must not overlap other system boundaries.
- A system boundary may contain only containers and descendants owned by that
  software system.
- Container boundaries must not overlap other container boundaries.
- A container boundary may contain only components owned by that container.
- A sibling container stays outside another container's boundary, even when
  both belong to the same software system.
- A person, external software system, external container, infrastructure node,
  or any other element not owned by a boundary stays outside that boundary.
- Nested boundaries must follow the model hierarchy: container inside its
  owning system and component inside its owning container.
- Relationships may cross boundaries; their endpoints may not be moved across
  boundaries to shorten lines.
- Boundary labels must remain visible and unambiguous.

If automatic or saved layout violates these rules, first verify ownership and
view membership. Then reduce or split the view before using intentional manual
positions. Never accept false containment as a visual compromise.

## Choose the right view

- Use a system landscape view for people and software systems across the
  enterprise or domain.
- Use a system context view for one software system, its users, and external
  systems it directly interacts with.
- Use a container view for the containers owned by one software system plus
  directly related people and external systems or containers.
- Use a component view for the components owned by one container plus directly
  related external elements.
- Use a dynamic view for an ordered runtime interaction, not static ownership.
- Use a deployment view for runtime placement in environments and nodes.

Prefer a small view with one clear question over one diagram that exposes every
known element and relationship.

Every view must communicate architecture or behavior through meaningful
relationships, boundaries, dependencies, lifecycle, or runtime flow. Do not
create a view whose sole purpose is listing elements. When readers need an
element inventory, use documentation or a generated catalog; keep diagrams
focused on how the elements work together.

Use separate views when readers need different flows or levels of detail. Each
view must still answer its own architectural question.

Split a view when it combines independent ownership or runtime questions and
its canonical relationship labels cannot be routed without collisions. Move
each complete question into a purposefully titled view; do not shorten,
disconnect, or hide the relationships merely to retain one output.

### Make the main view a system summary

When a published diagram set has multiple views, designate one stable main or
overview view:

- The main view must summarize the entire system: its entry points, major
  capabilities, owning services, shared runtime or data services, and key
  external dependencies or execution paths.
- Prefer showing every owned service when their relationships remain readable
  and meaningful.
- If all services make the overview unreadable, show representative owners
  from every major subsystem and the relationships that connect those
  subsystems. Put omitted service-level detail in focused secondary views.
- Do not let the main view describe only one feature, runtime path, subsystem,
  or user journey. A reader who sees only the main view should still understand
  the system's complete architectural shape and how its major parts work
  together.
- Preserve the main view's published key and filename. Refine its scope rather
  than replacing it with a narrowly focused view.

## Author the DSL

- Give every view a concrete purpose in its description.
- Use variables or stable element paths for references; do not select elements
  by incidental rendered text.
- Preserve stable view keys and output filenames used by documentation or
  publishing. When intentionally removing a key, update references and remove
  its stale generated output because regeneration does not prove stale files
  disappeared.
- Use `AddDefault`, `Add`, and `Remove` deliberately. Avoid `AddAll` when it
  obscures the view's question.
- `SelectRelationships` is available only in a `SystemLandscapeView`. In other
  view types, curate membership and use `Unlink` for relationships that do not
  answer the view's question.
- `Unlink` hides a real source-model relationship from one view; it does not
  mean the relationship is absent. Never unlink solely to improve layout.
- Before every `Unlink`, ask whether a reader seeing both endpoints without the
  edge could reasonably infer that no relationship exists. If so, keep and
  arrange the edge, remove an out-of-scope endpoint, or split the view.
- A narrowly titled dynamic flow may omit relationships that are outside that
  exact runtime interaction. Make the scope explicit in the title and
  description, and ensure another purposeful view or authoritative
  documentation communicates any omitted relationship that matters to system
  understanding.
- Main and overview views must retain the key architecturally significant
  relationships among their visible elements. Do not make an overview look
  simpler by disconnecting services that materially depend on one another.
- `NoRelationship` removes every relationship to and from that element after
  view finalization, including explicitly linked relationships. Use it only
  when one element is intentionally isolated within an otherwise meaningful
  view, not to turn the whole view into a listing or as a general edge filter.
- In a `DynamicView`, `Link(source, destination, description)` selects an
  existing source-model relationship. The description must exactly match the
  canonical relationship description; it is not a display-label override.
- Linked dynamic-view elements may also render other model relationships among
  those elements. Compare the rendered edge count with the intended links and
  use `Unlink` for every incidental relationship.
- Supply the exact canonical description to `Unlink`, even when only one
  relationship exists between the source and destination.
- Describe relationships with domain actions such as "Publishes alarm state"
  or "Retrieves schedules", not vague labels such as "Uses".
- Keep relationship direction consistent with the runtime call, event, or data
  flow.
- Coalesce duplicate relationships only when one label truthfully represents
  them.
- Do not weaken or shorten canonical element metadata merely to make a crowded
  diagram fit. Prefer a smaller view or a layout/rendering correction.
- Use `AutoLayout` as the default complete placement. It measures rendered
  content, places nodes, boundaries, routes, and labels together, and rejects
  invalid geometry instead of saving a partial result.
- Treat an `mdl svg` geometry failure as a real model, view-scope, saved-layout,
  or MDL defect. Do not work around it by unlinking relationships, shortening
  truthful text, retrying with guessed spacing, or accepting a partly rendered
  file.
- Always inspect every rendered view visually. Compilation, successful
  rendering, and non-overlapping geometry are not evidence that the view
  communicates well.
- Use saved coordinates from the MDL visual editor only for intentional
  refinement. A manual layout must contain every current element, complete
  relationship routes, and placed labels; stale or partial saved geometry is
  invalid.

## Regenerate and inspect

Render all affected views from the repository root, using the repository's
pinned `go tool mdl` invocation when available:

```bash
mdl svg <model-package> -all -dir <output-directory>
# or: go tool mdl svg <model-package> -all -dir <output-directory>
```

For interactive layout refinement:

```bash
mdl serve <model-package> -dir <output-directory>
```

### Arrange with the MDL visual editor

Start with `AutoLayout`, then visually review every view in the generated set.
Keep the automatic result when its hierarchy, spacing, labels, and edge routing
communicate the view's question clearly. Use the editor only when deliberate
placement would improve that communication. When a rendered view has excessive
whitespace, weak visual hierarchy, or avoidable edge crossings:

1. Confirm the view contains only relationships that answer its architectural
   question. Do not unlink a real relationship merely because its label or
   route is difficult to place. If too many in-scope relationships remain,
   split the question before positioning. Omit an incidental relationship only
   when the title and description make that scope clear and another purposeful
   view or authoritative documentation preserves the relevant fact.
2. Run `mdl serve` for the model package and output directory. If DSL changes
   while the editor is running, verify the displayed node and edge counts
   changed; restart `mdl serve` when it still shows the previously compiled
   model.
3. Select the affected view in the editor.
4. Arrange nodes, relationship labels, and boundaries so the primary
   architectural flow is apparent before reading every label.
5. Keep every boundary truthful while moving elements: only owned descendants
   may sit inside it, and sibling or external elements must remain outside.
6. Reset or reposition stale edge bend points and labels after moving nodes or
   changing membership. Saved routes from an earlier layout must not leave
   lines outside boundaries, unnecessary detours, or detached labels. If MDL
   rejects a stale or incomplete saved layout, regenerate that whole affected
   view or deliberately migrate every element and route together. Never mix
   old manual positions with newly guessed automatic values.
7. Save through the MDL editor so it records supported layout coordinates. Do
   not hand-edit generated SVG or JSON layout data.
8. Confirm the expected SVG's timestamp or content changed, wait for the write
   to finish, then reload the view from disk. Verify node coordinates and edge
   vertices survived before accepting the layout.
9. Reopen the persisted SVG at fitted viewport scale and verify the saved
   nodes, labels, arrows, and boundary titles.

For MDL renderer or layout changes, render the full repository view set at
least three times and compare the SVG files byte for byte. Also run independent
model packages concurrently. Any changed bytes between identical runs, port
collision, timeout, partial file, or cross-view result is a tool defect.

Always review the main view in the editor and arrange it deliberately whenever
that improves the whole-system summary. Review every secondary view at fitted
viewport scale and arrange it as needed. Do not use manual positioning to
compensate for excessive scope or an incorrect model; split or correct the view
first.

After rendering, verify:

- Every expected element appears once.
- Across the published view set, every architecturally relevant owned service
  appears in at least one purposeful view or has an explicit documented reason
  to remain model-only.
- The main view represents every major subsystem and capability, even when
  detailed services are delegated to secondary views.
- Rendered node, edge, and boundary counts match the view's stated scope.
- Every `Unlink` has been reviewed against the source relationship, the view's
  stated scope, and the inference a reader may draw from its omission.
- The C4 abstraction level is consistent.
- All boundaries satisfy the ownership and non-overlap rules.
- External elements are outside internal boundaries.
- Relationship direction and labels are readable.
- Nodes, labels, arrows, and boundary titles do not overlap.
- Text stays inside its node or boundary.
- The complete structure is visible at common viewport sizes. Labels in an
  honestly dense overview may require zoom, but focused secondary views must
  make its major flows readable without hiding real relationships.
- Generated files match the DSL and are included when the repository publishes
  rendered artifacts.

Run the repository's architecture-drift checks, tests, and formatting commands
after changing Go DSL.
