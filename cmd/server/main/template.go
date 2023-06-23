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

type templateMap map[string]*PageTemplate

func (t templateMap) Get(tmpl string) *PageTemplate {
	if strings.EqualFold("/", tmpl) {
		return t["home"]
	}
	tmpl = tmpl[1:]
	return t[tmpl]
}

func (app *App) newTemplateData() templateData {
	return templateData{
		"time":       time.Now(),
		"showHeader": true,
	}
}

func (t templateData) Set(key string, value interface{}) {
	t[key] = value
}

type PageTemplate struct {
	Tmpl *template.Template
}

func appPage() (templateMap, error) {
	templates := templateMap{}
	baseTemplate := template.Must(template.ParseFS(ui.FILES, "web/html/base.tmpl.html"))
	baseTemplate = template.Must(baseTemplate.ParseFS(ui.FILES, "web/html/partials/*"))

	pages, err := fs.Glob(ui.FILES, "web/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

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

func (tmpl *PageTemplate) ExecuteTemplate(out io.Writer, name string, data templateData) error {
	return tmpl.Tmpl.ExecuteTemplate(out, name, data)
}

func mailTemplates() (templateMap, error) {
	templates := map[string]*PageTemplate{}
	pages, err := fs.Glob(ui.FILES, "mail/*")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		base := filepath.Base(page)
		name := base[:strings.Index(base, ".")]

		tmpl := template.Must(template.ParseFS(ui.FILES, page))
		templates[name] = &PageTemplate{Tmpl: tmpl}
	}
	return templates, nil
}
