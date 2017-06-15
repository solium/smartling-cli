package main

import (
	"github.com/Smartling/api-sdk-go"
)

func doFilesPull(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		directory = args["--directory"].(string)
		project   = args["--project"].(string)
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

	pool := NewThreadPool(config.Threads)

	for _, file := range files {
		// func closure required to pass different file objects to goroutines
		func(file smartling.File) {
			pool.Do(func() {
				downloadAllFileLocales(
					client,
					project,
					file,
					format,
					directory,
				)
			})
		}(file)
	}

	pool.Wait()

	return nil
}
