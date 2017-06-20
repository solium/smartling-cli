package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesPush(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		project     = args["--project"].(string)
		file        = args["<file>"].(string)
		uri, useURI = args["<uri>"].(string)
		branch, _   = args["--branch"].(string)
		locales, _  = args["--locale"].([]string)
		authorize   = args["--authorize"].(bool)
		directory   = args["--directory"].(string)
		fileType, _ = args["--type"].(string)
	)

	files, err := globFilesLocally(directory, file)
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
			FileURI:            branch + uri,
			Authorize:          authorize,
			LocalesToAuthorize: locales,
		}

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
			"%s %s [strings %d words %d]\n",
			branch+file,
			status,
			response.StringCount,
			response.WordCount,
		)
	}

	return nil
}
