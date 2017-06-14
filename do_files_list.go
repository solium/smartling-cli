package main

import (
	"fmt"
	"os"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesList(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		short   = args["--short"].(bool)
		project = args["<project>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFilesListFormat
	}

	format, err := CompileFormatOption(args)
	if err != nil {
		return err
	}

	var (
		request = smartling.FilesListRequest{}
	)

	files, err := client.ListAllFiles(project, request)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to list files from project "%s"`,
			project,
		)
	}

	table := NewTableWriter(os.Stdout)

	for _, file := range files {
		if short {
			fmt.Fprintf(table, "%s\n", file.FileURI)
		} else {
			err := format.Execute(table, file)
			if err != nil {
				return FormatExecutionError{
					Cause:  err,
					Format: args["--format"].(string),
					Data:   file,
				}
			}
		}
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}
