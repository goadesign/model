package mdlsvc

import (
	"context"
	"io"

	genmodule "goa.design/model/mdlsvc/gen/module"
)

// List the model modules in the current Go workspace
func (svc *Service) ListModules(ctx context.Context) (res []*genmodule.Module, err error) {
	panic("not implemented")
}

// WebSocket endpoint for subscribing to updates to the model
func (svc *Service) Subscribe(ctx context.Context, mod *genmodule.Module, stream genmodule.SubscribeServerStream) (err error) {
	panic("not implemented")
}

// Stream the model JSON, see https://pkg.go.dev/goa.design/model/model#Model
func (svc *Service) GetModel(ctx context.Context, mod *genmodule.Module) (body io.ReadCloser, err error) {
	panic("not implemented")
}

// Stream the model DSL, save it, compile it and return the corresponding JSON
func (svc *Service) Compile(ctx context.Context, mod *genmodule.Module, rc io.ReadCloser) (res genmodule.Model, err error) {
	panic("not implemented")
}

// Stream the model layout JSON saved in the SVG
func (svc *Service) GetLayout(ctx context.Context, mod *genmodule.Module) (body io.ReadCloser, err error) {
	panic("not implemented")
}

// Save the SVG streamed in the request body
func (svc *Service) WriteDiagram(ctx context.Context, mod *genmodule.Module, rc io.ReadCloser) (err error) {
	panic("not implemented")
}
