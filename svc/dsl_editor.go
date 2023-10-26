package svc

import (
	"context"
	"os"
	"path/filepath"

	geneditor "goa.design/model/svc/gen/dsl_editor"
	gentypes "goa.design/model/svc/gen/types"
)

// Update the DSL for the given package, compile it and return the
// corresponding JSON if successful
func (svc *Service) UpdateDSL(ctx context.Context, p *gentypes.PackageFile) error {
	err := os.WriteFile(filepath.Join(p.Locator.Dir, p.Locator.Filename), []byte(p.Content), 0644)
	if err != nil {
		return logAndReturn(ctx, err, "failed to write file %s", p.Locator.Filename)
	}
	return nil
}

// UpsertSystem updates the DSL for the given system or adds the DSL if it does not exist.
func (svc *Service) UpsertSystem(ctx context.Context, p *geneditor.System) error {
	panic("not implemented")
}

// UpsertPerson implements UpsertPerson.
func (svc *Service) UpsertPerson(ctx context.Context, p *geneditor.Person) error {
	panic("not implemented")
}

// UpsertContainer implements UpsertContainer.
func (svc *Service) UpsertContainer(ctx context.Context, p *geneditor.Container) error {
	panic("not implemented")
}

// UpsertComponent implements UpsertComponent.
func (svc *Service) UpsertComponent(ctx context.Context, p *geneditor.Component) error {
	panic("not implemented")
}

// UpsertRelationship implements UpsertRelationship.
func (svc *Service) UpsertRelationship(ctx context.Context, p *geneditor.Relationship) error {
	panic("not implemented")
}

// DeleteSystem implements DeleteSystem.
func (svc *Service) DeleteSystem(ctx context.Context, p *geneditor.DeleteSystemPayload) error {
	panic("not implemented")
}

// DeletePerson implements DeletePerson.
func (svc *Service) DeletePerson(ctx context.Context, p *geneditor.DeletePersonPayload) error {
	panic("not implemented")
}

// DeleteContainer implements DeleteContainer.
func (svc *Service) DeleteContainer(ctx context.Context, p *geneditor.DeleteContainerPayload) error {
	panic("not implemented")
}

// DeleteComponent implements DeleteComponent.
func (svc *Service) DeleteComponent(ctx context.Context, p *geneditor.DeleteComponentPayload) error {
	panic("not implemented")
}

// DeleteRelationship implements DeleteRelationship.
func (svc *Service) DeleteRelationship(ctx context.Context, p *geneditor.DeleteRelationshipPayload) error {
	panic("not implemented")
}
