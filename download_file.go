package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	smartling "github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func downloadFile(
	client *smartling.Client,
	project string,
	uri string,
	locale string,
	path string,
) error {
	request := smartling.FileDownloadRequest{}
	request.FileURI = uri

	reader, err := client.DownloadFile(project, locale, request)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to download file "%s" from project "%s" (locale "%s")`,
			uri,
			project,
			locale,
		)
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create dirs hierarchy "%s" for downloaded file`,
			path,
		)
	}

	writer, err := os.Create(path)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to create output file "%s"`,
			path,
		)
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to write file contents into "%s"`,
			path,
		)
	}

	return nil
}

func downloadAllFileLocales(
	client *smartling.Client,
	project string,
	file smartling.File,
	format *Format,
	directory string,
) error {
	status, err := client.GetFileStatus(project, file.FileURI)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to retrieve file "%s" locales from project "%s"`,
			file.FileURI,
			project,
		)
	}

	for _, locale := range status.Items {
		buffer := &bytes.Buffer{}

		data := map[string]interface{}{
			"FileURI": file.FileURI,
			"Locale":  locale.LocaleID,
		}

		err = format.Execute(buffer, data)
		if err != nil {
			return FormatExecutionError{
				Cause: err,
				Data:  file,
			}
		}

		path := filepath.Join(directory, buffer.String())

		err = downloadFile(
			client,
			project,
			file.FileURI,
			locale.LocaleID,
			path,
		)
		if err != nil {
			return err
		}

		fmt.Printf("downloaded %s\n", path)
	}

	return err
}
