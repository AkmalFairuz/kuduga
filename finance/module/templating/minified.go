package templating

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"html/template"
)

func LoadMinified(name, content string) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	tmpl := template.New(name)
	minified, err := m.String("text/html", content)
	if err != nil {
		return nil, err
	}
	if _, err := tmpl.Parse(minified); err != nil {
		return nil, err
	}
	return tmpl, nil
}
