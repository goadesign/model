---
name: editing-model-diagrams
description: Creates, edits, reviews, and regenerates C4 architecture diagrams written with goa.design/model and the mdl CLI. Use when changing Model DSL, model.go or views.go files, system landscape, context, container, component, dynamic, or deployment views, element relationships, boundaries, layout, or generated SVG diagrams.
---

# Editing Model diagrams

Produce diagrams whose source model is architecturally true and whose rendered
views communicate that model without ambiguous ownership.

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
- Start with `AutoLayout` in the direction of the primary flow. Use saved
  coordinates only for intentional refinements.

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

After rendering, verify:

- Every expected element appears once.
- Rendered node, edge, and boundary counts match the view's stated scope.
- The C4 abstraction level is consistent.
- All boundaries satisfy the ownership and non-overlap rules.
- External elements are outside internal boundaries.
- Relationship direction and labels are readable.
- Nodes, labels, arrows, and boundary titles do not overlap.
- Text stays inside its node or boundary.
- The complete diagram is visible at common viewport sizes.
- Generated files match the DSL and are included when the repository publishes
  rendered artifacts.

Run the repository's architecture-drift checks, tests, and formatting commands
after changing Go DSL.
