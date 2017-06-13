package main

import (
	"io"
	"text/tabwriter"

	hierr "github.com/reconquest/hierr-go"
)

func NewTableWriter(target io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(target, 2, 4, 2, ' ', 0)
}

func RenderTable(writer *tabwriter.Writer) error {
	err := writer.Flush()
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to flush table to stdout",
		)
	}

	return nil
}
