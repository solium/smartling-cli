package main

import (
	"strings"
	"text/template"

	"github.com/reconquest/hierr-go"
)

func CompileFormatOption(
	option string,
	value string,
) (*template.Template, error) {
	value = strings.NewReplacer(`\n`, "\n", `\t`, "\t").Replace(value)

	format, err := template.New(option).Parse(value)
	if err != nil {
		return nil, NewError(
			hierr.Errorf(
				err,
				"failed to compile template for %s option",
				option,
			),

			"Check template syntax accordingly to text/template "+
				"documentation:\n\thttps://golang.org/pkg/text/template/",
		)
	}

	return format, nil
}
