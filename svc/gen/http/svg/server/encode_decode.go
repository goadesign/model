// Code generated by goa v3.13.2, DO NOT EDIT.
//
// SVG HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/model/svc/design -o svc/

package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	svg "goa.design/model/svc/gen/svg"
)

// EncodeLoadResponse returns an encoder for responses returned by the SVG Load
// endpoint.
func EncodeLoadResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(svg.SVG)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeLoadRequest returns a decoder for requests sent to the SVG Load
// endpoint.
func DecodeLoadRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			filename   string
			repository string
			dir        string
			err        error
		)
		filename = r.URL.Query().Get("filename")
		if filename == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Filename", "query string"))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("Filename", filename, "\\.go$"))
		repository = r.URL.Query().Get("repo")
		if repository == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Repository", "query string"))
		}
		if utf8.RuneCountInString(repository) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("Repository", repository, utf8.RuneCountInString(repository), 1, true))
		}
		dir = r.URL.Query().Get("dir")
		if dir == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Dir", "query string"))
		}
		if utf8.RuneCountInString(dir) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("Dir", dir, utf8.RuneCountInString(dir), 1, true))
		}
		if err != nil {
			return nil, err
		}
		payload := NewLoadFileLocator(filename, repository, dir)

		return payload, nil
	}
}

// EncodeLoadError returns an encoder for errors returned by the Load SVG
// endpoint.
func EncodeLoadError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewLoadNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeSaveResponse returns an encoder for responses returned by the SVG Save
// endpoint.
func EncodeSaveResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// DecodeSaveRequest returns a decoder for requests sent to the SVG Save
// endpoint.
func DecodeSaveRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body SaveRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateSaveRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewSavePayload(&body)

		return payload, nil
	}
}