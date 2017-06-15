package main

import (
	"os"

	"github.com/Smartling/api-sdk-go"
)

func doFilesStatus(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		project = args["--project"].(string)
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFileStatusFormat
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
		//row := map[string]interface{}{
		//    "FileURI": file.FileURI,
		//    "Status": "untracked",
		//    "Locale":
		err := format.Execute(table, file)
		if err != nil {
			return FormatExecutionError{
				Cause:  err,
				Format: args["--format"].(string),
				Data:   file,
			}
		}
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}
