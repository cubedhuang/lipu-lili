package app

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"

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
