package main

import (
	"github.com/Smartling/api-sdk-go"
)

var (
	usePullFormat = func(config FileConfig) string {
		return config.Pull.Format
	}
)

var (
	usePushFormat = func(config FileConfig) string {
		return config.Push.Format
	}
)

func executeFileFormat(
	config Config,
	file smartling.File,
	fallback string,
	getter func(config FileConfig) string,
	data interface{},
) (string, error) {
	return executeFileURIFormat(config, file.FileURI, fallback, getter, data)
}

func executeFileURIFormat(
	config Config,
	fileURI string,
	fallback string,
	getter func(config FileConfig) string,
	data interface{},
) (string, error) {
	local, err := config.GetFileConfig(fileURI)
	if err != nil {
		return "", err
	}

	template := getter(local)

	if template == "" {
		template = fallback
	}

	format, err := compileFormat(template)
	if err != nil {
		return "", err
	}

	result, err := format.Execute(data)
	if err != nil {
		return "", err
	}

	return result, nil
}

