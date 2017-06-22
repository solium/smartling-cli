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
		project = config.ProjectID
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFilePullFormat
	}

	files, err := globFilesRemote(client, project, uri)
	if err != nil {
		return err
	}

	pool := NewThreadPool(config.Threads)

	for _, file := range files {
		// func closure required to pass different file objects to goroutines
		func(file smartling.File) {
			pool.Do(func() {
				err := downloadFileTranslations(client, config, args, file)

				if err != nil {
					logger.Error(err)
				}
			})
		}(file)
	}

	pool.Wait()

	return nil
}
