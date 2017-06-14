package main

import (
	"github.com/Smartling/api-sdk-go"
	"github.com/gobwas/glob"
	hierr "github.com/reconquest/hierr-go"
)

func globFiles(
	client *smartling.Client,
	project string,
	uri string,
) ([]smartling.File, error) {
	if uri == "" {
		uri = "**"
	}

	pattern, err := glob.Compile(uri, '/')
	if err != nil {
		return nil, NewError(
			err,
			"Search file URI is malformed. Check out help for more "+
				"information about search patterns.",
		)
	}

	request := smartling.FilesListRequest{}

	files, err := client.ListAllFiles(project, request)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`unable to list files in project "%s"`,
			project,
		)
	}

	result := []smartling.File{}

	for _, file := range files {
		if pattern.Match(file.FileURI) {
			result = append(result, file)
		}
	}

	return result, nil
}
