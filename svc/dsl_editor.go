package svc

import (
	"context"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"goa.design/model/svc/editor"
	geneditor "goa.design/model/svc/gen/dsl_editor"
	gentypes "goa.design/model/svc/gen/types"
)

// Update the DSL for the given package, compile it and return the
// corresponding JSON if successful
func (svc *Service) UpdateDSL(ctx context.Context, p *gentypes.PackageFile) error {
	_, err := parser.ParseFile(token.NewFileSet(), fpath(p.Locator), p.Content, parser.ParseComments)
	if err != nil {
		return geneditor.MakeCompilationFailed(logAndReturn(ctx, err, "failed to parse DSL"))
	}
	if err := os.MkdirAll(filepath.Dir(fpath(p.Locator)), 0755); err != nil {
		return logAndReturn(ctx, err, "failed to create directory %s", p.Locator.Dir)
	}
	content, _ := format.Source([]byte(p.Content))
	if err := os.WriteFile(fpath(p.Locator), content, 0644); err != nil {
		return logAndReturn(ctx, err, "failed to write file %s", p.Locator.Filename)
	}
	return nil
}

// UpsertSystem updates the DSL for the given system or adds the DSL if it does not exist.
func (svc *Service) UpsertSystem(ctx context.Context, p *geneditor.System) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, systemDSL(p), editor.SoftwareSystemKind, p.Name)
}

// Create or update a person in the model
func (svc *Service) UpsertPerson(ctx context.Context, p *geneditor.Person) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, personDSL(p), editor.PersonKind, p.Name)
}

// Create or update a container in the model
func (svc *Service) UpsertContainer(ctx context.Context, p *geneditor.Container) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, containerDSL(p), editor.ContainerKind, p.SystemName+"/"+p.Name)
}

// Create or update a component in the model
func (svc *Service) UpsertComponent(ctx context.Context, p *geneditor.Component) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, componentDSL(p), editor.ComponentKind, p.SystemName+"/"+p.ContainerName+"/"+p.Name)
}

// Create or update a relationship in the model
func (svc *Service) UpsertRelationship(ctx context.Context, p *geneditor.Relationship) (*gentypes.PackageFile, error) {
	name := "Uses"
	if p.SourceKind == "Person" {
		if p.DestinationKind == "Person" {
			name = "InteractsWith"
		}
	} else if p.DestinationKind == "Person" {
		name = "Delivers"
	}
	data := &RelationshipData{Relationship: p, RelationName: name}
	ed := editor.NewEditor(p.Locator.Repository, p.Locator.Dir)
	f, err := ed.UpsertRelationship(editor.ElementKind(p.SourceKind), p.SourcePath, p.DestinationPath, relationshipDSL(data))
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return f, nil
}

// Create or update a landscape view in the model
func (svc *Service) UpsertLandscapeView(ctx context.Context, v *geneditor.LandscapeView) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, v.Locator.Repository, v.Locator.Dir, landscapeViewDSL(v), editor.LandscapeViewKind, "", v.Key)
}

// Create or update a system context view in the model
func (svc *Service) UpsertSystemContextView(ctx context.Context, v *geneditor.SystemContextView) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, v.Locator.Repository, v.Locator.Dir, systemContextViewDSL(v), editor.SystemContextViewKind, "", "", v.Key)
}

// Create or update a container view in the model
func (svc *Service) UpsertContainerView(ctx context.Context, v *geneditor.ContainerView) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, v.Locator.Repository, v.Locator.Dir, containerViewDSL(v), editor.ContainerViewKind, "", "", v.Key)
}

// Create or update a component view in the model
func (svc *Service) UpsertComponentView(ctx context.Context, v *geneditor.ComponentView) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, v.Locator.Repository, v.Locator.Dir, componentViewDSL(v), editor.ComponentViewKind, "", "", v.Key)
}

// Create or update an element style in the model
func (svc *Service) UpserElementStyle(ctx context.Context, e *geneditor.ElementStyle) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, e.Locator.Repository, e.Locator.Dir, elementStyleDSL(e), editor.ElementStyleKind, "", e.Tag)
}

// Create or update a relationship style in the model
func (svc *Service) UpsertRelationshipStyle(ctx context.Context, r *geneditor.RelationshipStyle) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, r.Locator.Repository, r.Locator.Dir, relationshipStyleDSL(r), editor.RelationshipStyleKind, "", r.Tag)
}

// DeleteSystem implements DeleteSystem.
func (svc *Service) DeleteSystem(ctx context.Context, p *geneditor.DeleteSystemPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.SoftwareSystemKind, p.SystemName)
}

// DeletePerson implements DeletePerson.
func (svc *Service) DeletePerson(ctx context.Context, p *geneditor.DeletePersonPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.PersonKind, p.PersonName)
}

// DeleteContainer implements DeleteContainer.
func (svc *Service) DeleteContainer(ctx context.Context, p *geneditor.DeleteContainerPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.ContainerKind, p.SystemName+"/"+p.ContainerName)
}

// DeleteComponent implements DeleteComponent.
func (svc *Service) DeleteComponent(ctx context.Context, p *geneditor.DeleteComponentPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.ComponentKind, p.SystemName+"/"+p.ContainerName+"/"+p.ComponentName)
}

// DeleteRelationship implements DeleteRelationship.
func (svc *Service) DeleteRelationship(ctx context.Context, p *geneditor.DeleteRelationshipPayload) (*gentypes.PackageFile, error) {
	edit := editor.NewEditor(p.Repository, p.Dir)
	f, err := edit.DeleteRelationship(editor.ElementKind(p.SourceKind), p.SourcePath, p.DestinationPath)
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return f, nil
}

// Delete an existing landscape view from the model
func (svc *Service) DeleteLandscapeView(ctx context.Context, p *geneditor.DeleteLandscapeViewPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.LandscapeViewKind, p.Key)
}

// Delete an existing system context view from the model
func (svc *Service) DeleteSystemContextView(ctx context.Context, p *geneditor.DeleteSystemContextViewPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.SystemContextViewKind, p.Key)
}

// Delete an existing container view from the model
func (svc *Service) DeleteContainerView(ctx context.Context, p *geneditor.DeleteContainerViewPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.ContainerViewKind, p.Key)
}

// Delete an existing component view from the model
func (svc *Service) DeleteComponentView(ctx context.Context, p *geneditor.DeleteComponentViewPayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.ComponentViewKind, p.Key)
}

// Delete an existing element style from the model
func (svc *Service) DeleteElementStyle(ctx context.Context, p *geneditor.DeleteElementStylePayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.ElementStyleKind, p.Tag)
}

// Delete an existing relationship style from the model
func (svc *Service) DeleteRelationshipStyle(ctx context.Context, p *geneditor.DeleteRelationshipStylePayload) (*gentypes.PackageFile, error) {
	return deleteElement(ctx, p.Repository, p.Dir, editor.RelationshipStyleKind, p.Tag)
}

func upsertElement(ctx context.Context, repo, dir, code string, kind editor.ElementKind, path string, args ...string) (*gentypes.PackageFile, error) {
	edit := editor.NewEditor(repo, dir)
	f, err := edit.UpsertElement(code, kind, path, args...)
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return f, nil
}

func deleteElement(ctx context.Context, repo, dir string, kind editor.ElementKind, key string) (*gentypes.PackageFile, error) {
	edit := editor.NewEditor(repo, dir)
	f, err := edit.DeleteElement(kind, key)
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return f, nil
}
