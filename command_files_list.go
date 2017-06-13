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

	format, err := CompileFormatOption("--format", args["--format"].(string))
	if err != nil {
		return err
	}

	var (
		request = smartling.FilesListRequest{}
		result  = []smartling.FileStatus{}
	)

	for {
		files, err := client.ListFiles(project, request)
		if err != nil {
			return hierr.Errorf(
				err,
				`unable to list files from project "%s"`,
				project,
			)
		}

		result = append(result, files.Items...)

		if request.Cursor.Offset+len(files.Items) < files.TotalCount {
			request.Cursor.Offset = len(files.Items)
		} else {
			break
		}
	}

	table := NewTableWriter(os.Stdout)

	for _, file := range result {
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
