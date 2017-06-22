package main

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/reconquest/hierr-go"
)

var (
	// because we can have specific formats for different file types defined
	// in config file, we need to cache templates to prevent overhead in
	// runtime
	compiledFormatsCache = map[string]*Format{}
)

func compileFormat(definition string) (*Format, error) {
	if format, ok := compiledFormatsCache[definition]; ok {
		return format, nil
	}

	definition = strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
	).Replace(definition)

	funcs := template.FuncMap{
		"name": func(path string) string {
			return strings.TrimSuffix(path, filepath.Ext(path))
		},

		"ext": func(path string) string {
			return filepath.Ext(path)
		},
	}

	var (
		format Format
		err    error
	)

	format.Source = definition
	format.Template, err = template.New("format").Funcs(funcs).Option(
		"missingkey=error",
	).Parse(
		definition,
	)
	if err != nil {
		return nil, NewError(
			hierr.Errorf(
				err,
				"failed to compile format template",
			),

			"Check template syntax accordingly to text/template "+
				"documentation:\n\thttps://golang.org/pkg/text/template/",
		)
	}

	compiledFormatsCache[definition] = &format

	return &format, nil
}
