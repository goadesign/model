package mdlsvc

import (
	"context"
	"io"

	gensvg "goa.design/model/mdlsvc/gen/svg"
)

// Stream the model layout JSON saved in the SVG
func (svc *Service) Load(ctx context.Context, file *gensvg.Filename) (io.ReadCloser, error) {
	panic("not implemented")
}

// Save the SVG streamed in the request body
func (svc *Service) Save(ctx context.Context, file *gensvg.Filename, rc io.ReadCloser) error {
	panic("not implemented")
}
