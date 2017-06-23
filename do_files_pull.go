package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
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

	var (
		err   error
		files []smartling.File
	)

	if uri == "-" {
		lines, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to read stdin",
			)
		}

		for _, line := range strings.Split(string(lines), "\n") {
			if line == "" {
				continue
			}

			files = append(files, smartling.File{
				FileURI: line,
			})
		}
	} else {
		files, err = globFilesRemote(client, project, uri)
		if err != nil {
			return err
		}
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
