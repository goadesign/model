package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"goa.design/model/mdl"
)

// DefaultTemplate is the template used to render and serve diagrams by
// default.
const DefaultTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>{{ .Title }}</title>
	<style>
		{{ .CSS }}
	</style>
</head>
<body>
	<div id="diagram"></div>
	<div class="footer">
		<div class="title">
			{{ .Title }}
			<div class="description">
				{{ .Description }}
			</div>
			<div class="version">
				{{ .Version }}
			</div>
		</div>
		<div class="legend">
			<div class="legend-title">
				Legend <span id="toggle">≚</span>
			</div>
			<div id="legend-diagram" style="display:none"></div>
		</div>
	</div>
	<script src="http://localhost:35729/livereload.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
	<script>
		function renderSvg(src, id) {
			var diagram = document.getElementById(id);
			mermaidAPI.render("mermaid"+id, src, function(svgCode) { diagram.innerHTML = svgCode; });
		} 
		var mermaidAPI = mermaid.mermaidAPI;
		mermaidAPI.initialize({
			securityLevel: 'loose',
			theme: 'neutral',
			startOnLoad:false{{ if .MermaidConfig }},
			...{{ .MermaidConfig }}{{ end }}
		});
		var diagramSrc = ` + "`{{ .MermaidSource }}`;" + `
		renderSvg(diagramSrc, "diagram")
		var legendSrc = ` + "`{{ .MermaidLegendSource }}`;" + `
		renderSvg(legendSrc, "legend-diagram")
	</script>
	<script>
		var toggle = document.getElementById("toggle");
		var legend = document.getElementById("legend-diagram");
		toggle.addEventListener('click', function (event) {
			if (legend.style.display == "") {
				legend.style.display = "none";
				toggle.innerHTML = "≚";
			} else {
				legend.style.display = "";
				toggle.innerHTML = "≙";
			}
		});
	</script>
</body>
</html>
`

// DefaultCSS is the CSS used to render and serve diagrams by default.
const DefaultCSS = `
body {
	padding: 10px;
	font-family: Arial;
}

//-----------------
// Diagram elements
//-----------------

.element {
	font-family: Arial;
}

.element-title {
	font-weight: bold;
}

.element-technology {
	font-size: 70%;
	padding-bottom: 0.8em;
}

.element-description {
	font-size: 80%;
}

.relationship {
	font-family: Arial;
	background-color: white;
}

.relationship-label {
	font-size: 80%;
	font-weight: bold;
	color: #909090;
}

.relationship-technology {
	font-size: 70%;
	font-weight: bold;
	color: #909090;
}

//-----------------
// Footer
//-----------------

.footer {
	-ms-box-orient: horizontal;
	display: -webkit-box;
	display: -moz-box;
	display: -ms-flexbox;
	display: -moz-flex;
	display: -webkit-flex;
	display: flex;
	-webkit-align-items: flex-start;
}

.title {
	margin-top: 2em;
	font-size: 110%;
	font-weight: bold;
	vertical-align: top;
	display: inline-block;
}

.description {
	color: #A0A0A0;
	font-size: 90%;
}

.version {
	color: #A0A0A0;
	font-size: 80%;
}

.legend {
	margin-top: 2em;
	margin-left: 2em;
	padding-top: 1em;
	padding-bottom: 1em;
	border: 2px solid #A0A0A0;
	border-radius: 5px;
	vertical-align: top;
	display: inline-block;
}

.legend-title {
	border-bottom: 1px solid #A0A0A0;
	margin-bottom: 0.5em;
	padding-left: 0.5em;
	padding-right: 0.5em;
	padding-bottom: 0.5em;
}

#toggle {
	margin-left: 0.2em;
	padding: 0.1em;
	border: 1px solid #A0A0A0;
	border-radius: 5px;
	background-color: #606060;
	color: white;
	cursor: pointer;
}
`

// ViewData is the data structure used to render the HTML template for a
// given view.
type ViewData struct {
	// Key of view
	Key string
	// Title of view
	Title string
	// Description of view
	Description string
	// Version of design
	Version string
	// MermaidSource is the Mermaid diagram source code.
	MermaidSource template.JS
	// MermaidLegendSource is the Mermaid legend source code.
	MermaidLegendSource template.JS
	// MermaidConfig is the Mermaid config JSON.
	MermaidConfig template.JS
	// CSS rendered inline
	CSS template.CSS
}

// indexTmpl is the default Go template used to render views.
var indexTmpl = template.Must(template.New("view").Parse(DefaultTemplate))

// loadViews generates the views for the given Go package, loads and returns the
// results indexed by view keys.
func loadViews(pkg, out string, debug bool) (map[string]*mdl.RenderedView, error) {
	if err := gen(pkg, out, debug); err != nil {
		return nil, err
	}
	views := make(map[string]*mdl.RenderedView)
	err := filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var view mdl.RenderedView
		if err := json.Unmarshal(b, &view); err != nil {
			return err
		}
		views[view.Key] = &view
		return nil
	})
	if err != nil {
		return nil, err
	}
	return views, nil
}

// render generates the views and renders static pages from the results.
func render(pkg, config, out string, debug bool) error {
	views, err := loadViews(pkg, out, debug)
	if err != nil {
		return err
	}
	for _, view := range views {
		f, err := os.Create(filepath.Join(out, view.Key+".html"))
		if err != nil {
			return err
		}
		data := &ViewData{
			Title:               view.Title,
			Description:         view.Description,
			Version:             view.Version,
			MermaidSource:       template.JS(view.Mermaid),
			MermaidLegendSource: template.JS(view.Legend),
			MermaidConfig:       template.JS(config),
			CSS:                 template.CSS(DefaultCSS),
		}
		if err := indexTmpl.Execute(f, data); err != nil {
			return err
		}
	}
	return nil
}
