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
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, editor.SoftwareSystemKind, p.Name, systemDSL(p))
}

// UpsertPerson implements UpsertPerson.
func (svc *Service) UpsertPerson(ctx context.Context, p *geneditor.Person) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, editor.PersonKind, p.Name, personDSL(p))
}

// UpsertContainer implements UpsertContainer.
func (svc *Service) UpsertContainer(ctx context.Context, p *geneditor.Container) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, editor.ContainerKind, p.Name, containerDSL(p))
}

// UpsertComponent implements UpsertComponent.
func (svc *Service) UpsertComponent(ctx context.Context, p *geneditor.Component) (*gentypes.PackageFile, error) {
	return upsertElement(ctx, p.Locator.Repository, p.Locator.Dir, editor.ComponentKind, p.Name, componentDSL(p))
}

// UpsertRelationship implements UpsertRelationship.
func (svc *Service) UpsertRelationship(ctx context.Context, p *geneditor.Relationship) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update a landscape view in the model
func (svc *Service) UpsertLandscapeView(ctx context.Context, v *geneditor.LandscapeView) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update a system context view in the model
func (svc *Service) UpsertSystemContextView(ctx context.Context, v *geneditor.SystemContextView) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update a container view in the model
func (svc *Service) UpsertContainerView(ctx context.Context, v *geneditor.ContainerView) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update a component view in the model
func (svc *Service) UpsertComponentView(ctx context.Context, v *geneditor.ComponentView) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update an element style in the model
func (svc *Service) UpserElementStyle(ctx context.Context, v *geneditor.ElementStyle) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Create or update a relationship style in the model
func (svc *Service) UpsertRelationshipStyle(ctx context.Context, v *geneditor.RelationshipStyle) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// DeleteSystem implements DeleteSystem.
func (svc *Service) DeleteSystem(ctx context.Context, p *geneditor.DeleteSystemPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// DeletePerson implements DeletePerson.
func (svc *Service) DeletePerson(ctx context.Context, p *geneditor.DeletePersonPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// DeleteContainer implements DeleteContainer.
func (svc *Service) DeleteContainer(ctx context.Context, p *geneditor.DeleteContainerPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// DeleteComponent implements DeleteComponent.
func (svc *Service) DeleteComponent(ctx context.Context, p *geneditor.DeleteComponentPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// DeleteRelationship implements DeleteRelationship.
func (svc *Service) DeleteRelationship(ctx context.Context, p *geneditor.DeleteRelationshipPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing landscape view from the model
func (svc *Service) DeleteLandscapeView(ctx context.Context, p *geneditor.DeleteLandscapeViewPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing system context view from the model
func (svc *Service) DeleteSystemContextView(ctx context.Context, p *geneditor.DeleteSystemContextViewPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing container view from the model
func (svc *Service) DeleteContainerView(ctx context.Context, p *geneditor.DeleteContainerViewPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing component view from the model
func (svc *Service) DeleteComponentView(ctx context.Context, p *geneditor.DeleteComponentViewPayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing element style from the model
func (svc *Service) DeleteElementStyle(ctx context.Context, p *geneditor.DeleteElementStylePayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

// Delete an existing relationship style from the model
func (svc *Service) DeleteRelationshipStyle(ctx context.Context, p *geneditor.DeleteRelationshipStylePayload) (*gentypes.PackageFile, error) {
	panic("not implemented")
}

func upsertElement(ctx context.Context, repo, dir string, kind editor.ElementKind, elementPath, code string) (*gentypes.PackageFile, error) {
	edit := editor.NewEditor(repo, dir)
	f, err := edit.UpsertElement(kind, elementPath, code)
	if err != nil {
		return nil, logAndReturn(ctx, err)
	}
	return f, nil
}
