package main

import (
	"fmt"
	"os"

	"github.com/Smartling/api-sdk-go"
)

func doFilesList(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		short   = args["--short"].(bool)
		project = args["<project>"].(string)
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFilesListFormat
	}

	format, err := CompileFormatOption(args)
	if err != nil {
		return err
	}

	files, err := globFiles(client, project, uri)
	if err != nil {
		return err
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
