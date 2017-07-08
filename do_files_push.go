package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesPush(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		project     = config.ProjectID
		file, _     = args["<file>"].(string)
		uri, useURI = args["<uri>"].(string)
		branch, _   = args["--branch"].(string)
		locales, _  = args["--locale"].([]string)
		authorize   = args["--authorize"].(bool)
		directory   = args["--directory"].(string)
		fileType, _ = args["--type"].(string)
	)

	if branch == "@auto" {
		var err error

		branch, err = getGitBranch()
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to autodetect branch name",
			)
		}

		logger.Infof("autodetected branch name: %s", branch)
	}

	if branch != "" {
		branch = strings.TrimSuffix(branch, "/") + "/"
	}

	patterns := []string{}

	if file != "" {
		patterns = append(patterns, file)
	} else {
		for pattern, section := range config.Files {
			if section.Push.Type != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	files := []string{}

	for _, pattern := range patterns {
		chunk, err := globFilesLocally(directory, pattern)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to find matching files to upload`,
				),

				`Check, that specified pattern is valid and refer to help for`+
					` more information about glob patterns.`,
			)
		}

		files = append(files, chunk...)
	}

	if len(files) == 0 {
		return NewError(
			fmt.Errorf(`no files found by specified patterns`),

			`Check command line pattern if any and configuration file for`+
				` more patterns to search for.`,
		)
	}

	if uri != "" && len(files) > 1 {
		return NewError(
			fmt.Errorf(
				`more than one file is matching speciifed pattern and <uri>`+
					` is specified too`,
			),

			`Either remove <uri> argument or make sure that only one file`+
				` is matching mask.`,
		)
	}

	for _, file := range files {
		if !useURI {
			uri = file
		}

		fileConfig, err := config.GetFileConfig(file)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to retrieve file specific configuration`,
				),

				``,
			)
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to read file contents "%s"`,
					file,
				),

				`Check that file exists and readable by current user.`,
			)
		}

		request := smartling.FileUploadRequest{
			File:               contents,
			Authorize:          authorize,
			LocalesToAuthorize: locales,
		}

		request.FileURI = branch + uri

		if fileConfig.Push.Type == "" {
			if fileType == "" {
				request.FileType = smartling.GetFileTypeByExtension(
					filepath.Ext(file),
				)

				if request.FileType == smartling.FileTypeUnknown {
					return NewError(
						fmt.Errorf(
							"unable to deduce file type from extension: %q",
							filepath.Ext(file),
						),

						`You need to specify file type via --type option.`,
					)
				}
			} else {
				request.FileType = smartling.FileType(fileType)
			}
		} else {
			request.FileType = smartling.FileType(fileConfig.Push.Type)
		}

		request.Smartling.Directives = fileConfig.Push.Directives

		response, err := client.UploadFile(project, request)

		if err != nil {
			return NewError(
				hierr.Errorf(
					err,
					`unable to upload file "%s"`,
					file,
				),

				`Check, that you have enough permissions to upload file to`+
					` the specified project`,
			)
		}

		status := "new"
		if response.Overwritten {
			status = "overwritten"
		}

		fmt.Printf(
			"%s (%s) %s [%d strings %d words]\n",
			branch+file,
			request.FileType,
			status,
			response.StringCount,
			response.WordCount,
		)
	}

	return nil
}
