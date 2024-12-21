package app

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	texttemplate "text/template"

	"github.com/cubedhuang/lipu-lili/internal/models"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

func compileTemplates(path string, funcMap template.FuncMap) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	tmpl := template.New("").Funcs(funcMap)

	files, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("failed to find templates: %w", err)
	}

	for _, filename := range files {
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, fmt.Errorf("failed to minify file %s: %w", filename, err)
		}

		tmpl, err = tmpl.Parse(string(mb))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", filename, err)
		}
	}

	return tmpl, nil
}

const sitemapTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>http://lili.nimi.li/</loc>
		<changefreq>monthly</changefreq>
		<priority>1.0</priority>
	</url>
{{- range . -}}
	<url>
		<loc>http://lili.nimi.li/{{ .Id }}</loc>
		<changefreq>monthly</changefreq>
		<priority>{{ pagePriority . }}</priority>
	</url>
{{- end -}}
</urlset>`

func compileSiteMapTemplate() (*texttemplate.Template, error) {
	tmpl, err := texttemplate.New("sitemap").Funcs(template.FuncMap{
		"pagePriority": func(word models.WordData) string {
			usage := word.GetUsage()
			if usage == -1 {
				return "0.0"
			}
			priority := float64(usage) / 100
			return fmt.Sprintf("%.2f", priority)
		},
	}).Parse(sitemapTmpl)

	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap template: %w", err)
	}

	return tmpl, nil
}
