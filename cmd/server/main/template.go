package main

import (
	"barlio/cmd/server/model"
	"barlio/ui"
	"html/template"
	"io"
	"io/fs"
	"net/http"
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
	if strings.Compare("/", tmpl) == 0 {
		return t["home"]
	}
	tmpl = tmpl[strings.LastIndex(tmpl, "/")+1:]
	return t[tmpl]
}

func formatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

func (app *App) newTemplateData(user *model.User, r *http.Request) templateData {
	return templateData{
		"time":       time.Now(),
		"showHeader": true,
		"user":       user,
		"page":       r.URL.Path,
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
	funcMap := template.FuncMap{"formatdate": formatDate}

	baseTemplate := template.Must(template.ParseFS(ui.FILES, "web/html/base.tmpl.html"))
	baseTemplate = template.Must(baseTemplate.ParseFS(ui.FILES, "web/html/partials/*"))
	baseTemplate = baseTemplate.Funcs(funcMap)

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
