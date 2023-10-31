package svc

import (
	"context"
	"os"
	"path/filepath"

	gensvg "goa.design/model/svc/gen/svg"
	gentypes "goa.design/model/svc/gen/types"
)

// Returns the SVG for the given file.
func (svc *Service) Load(ctx context.Context, file *gentypes.FileLocator) (gensvg.SVG, error) {
	data, err := os.ReadFile(fpath(file))
	if err != nil {
		if os.IsNotExist(err) {
			return "", gensvg.MakeNotFound(err)
		}
		return "", err
	}
	return gensvg.SVG(data), nil
}

// Save the SVG to the given file.
func (svc *Service) Save(ctx context.Context, payload *gensvg.SavePayload) error {
	fpath := filepath.Join(payload.Repository, payload.Dir, payload.Filename)
	return os.WriteFile(fpath, []byte(payload.SVG), 0644)
}
