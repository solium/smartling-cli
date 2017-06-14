package main

import (
	"bytes"
	"path/filepath"

	"github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func doFilesPull(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		directory = args["--directory"].(string)
		project   = args["<project>"].(string)
		uri, _    = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFilePullFormat
	}

	format, err := CompileFormatOption(args)
	if err != nil {
		return err
	}

	files, err := globFiles(client, project, uri)
	if err != nil {
		return err
	}

	for _, file := range files {
		status, err := client.GetFileStatus(project, file.FileURI)
		if err != nil {
			return hierr.Errorf(
				err,
				`unable to retrieve file "%s" locales from project "%s"`,
				file.FileURI,
				project,
			)
		}

		for _, locale := range status.Items {
			buffer := &bytes.Buffer{}

			data := map[string]interface{}{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			}

			err = format.Execute(buffer, data)
			if err != nil {
				return FormatExecutionError{
					Cause:  err,
					Format: args["--format"].(string),
					Data:   file,
				}
			}

			path := filepath.Join(directory, buffer.String())

			err = downloadFile(client, project, file, locale.LocaleID, path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
