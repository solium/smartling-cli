package main

import (
	"fmt"
	"path/filepath"

	smartling "github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func downloadFileTranslations(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
	file smartling.File,
) error {
	var (
		project   = args["--project"].(string)
		directory = args["--directory"].(string)
		source    = args["--source"].(bool)

		defaultFormat, _ = args["--format"].(string)
	)

	if defaultFormat == "" {
		defaultFormat = defaultFileStatusFormat
	}

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

	if source {
		translations = append(translations, smartling.FileStatusTranslation{
			LocaleID: "",
		})
	}

	for _, locale := range translations {
		path, err := executeFileFormat(
			config,
			file,
			defaultFormat,
			usePullFormat,
			map[string]interface{}{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			},
		)
		if err != nil {
			return err
		}

		path = filepath.Join(directory, path)

		err = downloadFile(client, project, file, locale.LocaleID, path)
		if err != nil {
			return err
		}

		fmt.Printf("downloaded %s\n", path)
	}

	return err
}
