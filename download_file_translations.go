package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func downloadFileTranslations(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
	file smartling.File,
) error {
	var (
		project   = config.ProjectID
		directory = args["--directory"].(string)
		source    = args["--source"].(bool)
		locales   = args["--locale"].([]string)

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

	var translations []smartling.FileStatusTranslation

	if source {
		translations = []smartling.FileStatusTranslation{
			{LocaleID: ""},
		}
	} else {
		translations = status.Items
	}

	for _, locale := range translations {
		if len(locales) > 0 {
			if !hasLocaleInList(locale.LocaleID, locales) {
				continue
			}
		}

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

func hasLocaleInList(locale string, locales []string) bool {
	for _, filter := range locales {
		if strings.ToLower(filter) == strings.ToLower(locale) {
			return true
		}
	}

	return false
}
