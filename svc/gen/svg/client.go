// Code generated by goa v3.13.2, DO NOT EDIT.
//
// SVG client
//
// Command:
// $ goa gen goa.design/model/svc/design -o svc/

package svg

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	types "goa.design/model/svc/gen/types"
)

// Client is the "SVG" service client.
type Client struct {
	LoadEndpoint goa.Endpoint
	SaveEndpoint goa.Endpoint
}

// NewClient initializes a "SVG" service client given the endpoints.
func NewClient(load, save goa.Endpoint) *Client {
	return &Client{
		LoadEndpoint: load,
		SaveEndpoint: save,
	}
}

// Load calls the "Load" endpoint of the "SVG" service.
// Load may return the following errors:
//   - "NotFound" (type *goa.ServiceError): File not found
//   - error: internal error
func (c *Client) Load(ctx context.Context, p *types.FileLocator) (res SVG, err error) {
	var ires any
	ires, err = c.LoadEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(SVG), nil
}

// Save calls the "Save" endpoint of the "SVG" service.
func (c *Client) Save(ctx context.Context, p *SavePayload) (err error) {
	_, err = c.SaveEndpoint(ctx, p)
	return
}