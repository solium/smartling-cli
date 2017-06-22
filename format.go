package main

import (
	"bytes"
	"text/template"
)

type Format struct {
	*template.Template

	Source string
}

func (format *Format) Execute(data interface{}) (string, error) {
	buffer := &bytes.Buffer{}

	err := format.Template.Execute(buffer, data)
	if err != nil {
		return "", FormatExecutionError{
			Cause:  err,
			Format: format.Source,
			Data:   data,
		}
	}

	return buffer.String(), nil
}
