package main

import (
	"bytes"
	"fmt"
	"path/filepath"

	smartling "github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func downloadFileTranslations(
	client *smartling.Client,
	project string,
	file smartling.File,
	format *Format,
	directory string,
	original bool,
) error {
	status, err := client.GetFileStatus(project, file.FileURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to retrieve file "%s" locales from project "%s"`,
			file.FileURI,
			project,
		)
	}

	translations := status.Items

	if original {
		translations = append(translations, smartling.FileStatusTranslation{
			LocaleID: "",
		})
	}

	for _, locale := range translations {
		buffer := &bytes.Buffer{}

		data := map[string]interface{}{
			"FileURI": file.FileURI,
			"Locale":  locale.LocaleID,
		}

		err = format.Execute(buffer, data)
		if err != nil {
			return FormatExecutionError{
				Cause: err,
				Data:  file,
			}
		}

		path := filepath.Join(directory, buffer.String())

		err = downloadFile(client, project, file, locale.LocaleID, path)
		if err != nil {
			return err
		}

		fmt.Printf("downloaded %s\n", path)
	}

	return err
}
