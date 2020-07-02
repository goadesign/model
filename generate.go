package docs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen/openapi"
)

// init registers the plugin generator function.
func init() {
	codegen.RegisterPlugin("structurizr", "gen", nil, Generate)
}

// Generate produces the documentation JSON file.
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, docsFile(r))
		}
	}
	return files, nil
}

func docsFile(r *expr.RootExpr) *codegen.File {
	docs := &data{
		API:         apiDocs(r.API),
		Services:    servicesDocs(r),
		Definitions: openapi.Definitions,
	}
	jsonPath := filepath.Join(codegen.Gendir, "docs.json")
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		// goa does not delete files in the top-level gen folder.
		// https://github.com/goadesign/goa/pull/2194
		// The plugin must delete docs.json so that the generator does not append
		// to the existing docs.json.
		if err := os.Remove(jsonPath); err != nil {
			panic(err)
		}
	}
	jsonSection := &codegen.SectionTemplate{
		Name:    "docs",
		FuncMap: template.FuncMap{"toJSON": toJSON},
		Source:  "{{ toJSON .}}",
		Data:    docs,
	}
	return &codegen.File{
		Path:             jsonPath,
		SectionTemplates: []*codegen.SectionTemplate{jsonSection},
	}
}

func apiDocs(api *expr.APIExpr) *apiData {
	data := &apiData{
		Name:        api.Name,
		Title:       api.Title,
		Description: api.Description,
		Version:     api.Version,
		Terms:       api.TermsOfService,
	}
	if len(api.Servers) > 0 {
		data.Servers = make(map[string]*serverData, len(api.Servers))
		for _, s := range api.Servers {
			data.Servers[s.Name] = generateServer(s)
		}
	}
	if c := api.Contact; c != nil {
		data.Contact = &contactData{c.Name, c.Email, c.URL}
	}
	if l := api.License; l != nil {
		data.License = &licenseData{l.Name, l.URL}
	}
	if d := api.Docs; d != nil {
		data.Docs = &docsData{d.Description, d.URL}
	}
	data.Requirements = make([]*requirementData, len(api.Requirements))
	for i, req := range api.Requirements {
		data.Requirements[i] = generateRequirement(req)
	}

	return data
}

func servicesDocs(r *expr.RootExpr) map[string]*serviceData {
	svcs := make(map[string]*serviceData, len(r.Services))
	for _, svc := range r.Services {
		n := svc.Name
		svcs[n] = &serviceData{
			Name:        n,
			Description: svc.Description,
		}

		svcs[n].Methods = make(map[string]*methodData, len(svc.Methods))
		for _, meth := range svc.Methods {
			svcs[n].Methods[meth.Name] = generateMethod(r.API, meth)
		}

		svcs[n].Requirements = make([]*requirementData, len(svc.Requirements))
		for i, req := range svc.Requirements {
			svcs[n].Requirements[i] = generateRequirement(req)
		}
	}
	return svcs
}

func generateServer(s *expr.ServerExpr) *serverData {
	data := &serverData{
		Name:        s.Name,
		Description: s.Description,
		Services:    s.Services,
	}
	if len(s.Hosts) > 0 {
		data.Hosts = make(map[string]*hostData)
		for _, h := range s.Hosts {
			data.Hosts[h.Name] = &hostData{
				Name:        h.Name,
				ServerName:  h.ServerName,
				Description: h.Description,
			}
			if len(h.URIs) > 0 {
				data.Hosts[h.Name].URIs = make([]string, len(h.URIs))
				for i, u := range h.URIs {
					data.Hosts[h.Name].URIs[i] = string(u)
				}
			}
			if o := expr.AsObject(h.Variables.Type); o != nil {
				data.Hosts[h.Name].Variables = make([]*variableData, len(*o))
				for i, na := range *o {
					var def string
					if na.Attribute.DefaultValue != nil {
						def = fmt.Sprintf("%v", na.Attribute.DefaultValue)
					}
					var e []string
					if na.Attribute.Validation != nil && len(na.Attribute.Validation.Values) > 0 {
						e = make([]string, len(na.Attribute.Validation.Values))
						for j, v := range na.Attribute.Validation.Values {
							e[j] = fmt.Sprintf("%v", v)
						}
					}
					data.Hosts[h.Name].Variables[i] = &variableData{na.Name, def, e}
				}
			}
		}
	}
	return data
}

func generateRequirement(req *expr.SecurityExpr) *requirementData {
	r := &requirementData{Scopes: req.Scopes}
	if len(req.Schemes) > 0 {
		r.Schemes = make([]*schemeData, len(req.Schemes))
		for i, sch := range req.Schemes {
			r.Schemes[i] = &schemeData{
				Type:        sch.Type(),
				Description: sch.Description,
				Name:        sch.Name,
				In:          sch.In,
				Scheme:      sch.SchemeName,
			}
			if len(sch.Flows) > 0 {
				r.Schemes[i].Flows = make([]*flowData, len(sch.Flows))
				for j, f := range sch.Flows {
					r.Schemes[i].Flows[j] = &flowData{f.Type(), f.AuthorizationURL, f.TokenURL, f.RefreshURL}
				}
			}
		}
	}
	return r
}

func generateMethod(api *expr.APIExpr, meth *expr.MethodExpr) *methodData {
	m := &methodData{
		Name:        meth.Name,
		Description: meth.Description,
		Payload:     generatePayload(api, meth.Payload, meth.IsPayloadStreaming()),
		Result:      generatePayload(api, meth.Result, meth.Stream == expr.BidirectionalStreamKind || meth.Stream == expr.ServerStreamKind),
	}
	m.Errors = make(map[string]*errorData, len(meth.Errors))
	for _, er := range meth.Errors {
		m.Errors[er.Name] = generateError(api, er)
	}
	m.Requirements = make([]*requirementData, len(meth.Requirements))
	for i, req := range meth.Requirements {
		m.Requirements[i] = generateRequirement(req)
	}
	return m
}

func generatePayload(api *expr.APIExpr, att *expr.AttributeExpr, streaming bool) *payloadData {
	schema := openapi.AttributeTypeSchema(api, att)
	return &payloadData{
		Type:      schema,
		Example:   att.Example(api.Random()),
		Streaming: streaming,
	}
}

func generateError(api *expr.APIExpr, er *expr.ErrorExpr) *errorData {
	_, temporary := er.AttributeExpr.Meta["goa:error:temporary"]
	_, timeout := er.AttributeExpr.Meta["goa:error:timeout"]
	_, fault := er.AttributeExpr.Meta["goa:error:fault"]
	return &errorData{
		Name:        er.Name,
		Description: er.Description,
		Type:        openapi.AttributeTypeSchema(api, er.AttributeExpr),
		Temporary:   temporary,
		Timeout:     timeout,
		Fault:       fault,
	}
}

func generateScheme(sch *service.SchemeData) *schemeData {
	s := &schemeData{
		Type:   sch.Type,
		Name:   sch.Name,
		In:     sch.In,
		Scheme: sch.SchemeName,
	}
	if len(sch.Flows) > 0 {
		s.Flows = make([]*flowData, len(sch.Flows))
		for i, f := range sch.Flows {
			s.Flows[i] = &flowData{f.Type(), f.AuthorizationURL, f.TokenURL, f.RefreshURL}
		}
	}
	return s
}

func toJSON(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		panic("openapi: " + err.Error()) // bug
	}
	return string(b)
}
