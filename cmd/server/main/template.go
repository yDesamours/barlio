package main

import (
	"barlio/ui"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

const (
	base = "base"
)

type templateData map[string]interface{}

func newTemplateData() templateData {
	return templateData{
		"time": time.Now(),
	}
}

func (t templateData) Set(key string, value interface{}) {
	t[key] = value
}

type PageTemplate struct {
	Tmpl *template.Template
}

func appPage() (map[string]*PageTemplate, error) {
	templates := map[string]*PageTemplate{}
	baseTemplate := template.Must(template.ParseFS(ui.FILES, "html/base.tmpl.html"))
	baseTemplate = template.Must(baseTemplate.ParseFS(ui.FILES, "html/partials/*"))

	pages, _ := fs.Glob(ui.FILES, "html/pages/*.tmpl.html")

	for _, page := range pages {
		filename := filepath.Base(page)
		name := filename[:strings.Index(filename, ".")]

		ts, err := template.Must(baseTemplate.Clone()).ParseFS(ui.FILES, page)
		if err != nil {
			return nil, err
		}
		templates[name] = &PageTemplate{Tmpl: ts}
	}
	return templates, nil
}

func (tmpl *PageTemplate) Execute(out io.Writer, data templateData) error {
	return tmpl.Tmpl.ExecuteTemplate(out, base, data)
}
