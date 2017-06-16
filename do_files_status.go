package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/Smartling/api-sdk-go"
)

func doFilesStatus(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		project   = args["--project"].(string)
		uri, _    = args["<uri>"].(string)
		directory = args["--directory"].(string)
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

	var (
		table = NewTableWriter(os.Stdout)
	)

	for _, file := range files {
		path, err := format.Execute(map[string]interface{}{
			"FileURI": file.FileURI,
		})
		if err != nil {
			return err
		}

		status, err := client.GetFileStatus(project, file.FileURI)
		if err != nil {
			return err
		}

		path = filepath.Join(directory, path)

		state := "source"
		if !isFileExists(path) {
			state = "missing"
		}

		writeFileStatus(table, map[string]string{
			"Path":     path,
			"State":    state,
			"Progress": "source",
			"Strings":  fmt.Sprint(status.TotalStringCount),
			"Words":    fmt.Sprint(status.TotalWordCount),
		})

		for _, locale := range status.Items {
			path, err := format.Execute(map[string]interface{}{
				"FileURI": file.FileURI,
				"Locale":  locale.LocaleID,
			})
			if err != nil {
				return err
			}

			path = filepath.Join(directory, path)

			state := "remote"
			if !isFileExists(path) {
				state = "missing"
			}

			writeFileStatus(table, map[string]string{
				"Path":   path,
				"Locale": locale.LocaleID,
				"State":  state,
				"Progress": fmt.Sprintf(
					"%d%%",
					int(
						100*
							float64(locale.CompletedStringCount)/
							float64(status.TotalStringCount),
					),
				),
				"Strings": fmt.Sprint(locale.CompletedStringCount),
				"Words":   fmt.Sprint(locale.CompletedWordCount),
			})
		}
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}

func writeFileStatus(table *tabwriter.Writer, row map[string]string) {
	fmt.Fprintf(
		table,
		"%s\t%s\t%s\t%s\t%s\t%s\n",
		row["Path"],
		row["Locale"],
		row["State"],
		row["Progress"],
		row["Strings"],
		row["Words"],
	)
}
