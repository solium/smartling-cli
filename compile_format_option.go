package main

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/reconquest/hierr-go"
)

func CompileFormatOption(
	args map[string]interface{},
) (*template.Template, error) {
	value := args["--format"].(string)

	value = strings.NewReplacer(`\n`, "\n", `\t`, "\t").Replace(value)

	funcs := template.FuncMap{
		"name": func(path string) string {
			return strings.TrimSuffix(path, filepath.Ext(path))
		},

		"ext": func(path string) string {
			return filepath.Ext(path)
		},
	}

	format, err := template.New("format").Funcs(funcs).Parse(value)
	if err != nil {
		return nil, NewError(
			hierr.Errorf(
				err,
				"failed to compile template for --format option",
			),

			"Check template syntax accordingly to text/template "+
				"documentation:\n\thttps://golang.org/pkg/text/template/",
		)
	}

	return format, nil
}
