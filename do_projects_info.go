package main

import (
	"fmt"
	"os"

	"github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func doProjectsInfo(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	project := args["--project"].(string)

	details, err := client.GetProjectDetails(project)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to get project "%s" details`,
			project,
		)
	}

	table := NewTableWriter(os.Stdout)

	status := "active"

	if details.Archived {
		status = "archived"
	}

	info := [][]interface{}{
		{"ID", details.ProjectID},
		{"ACCOUNT", details.AccountUID},
		{"NAME", details.ProjectName},
		{
			"LOCALE",
			details.SourceLocaleID + ": " + details.SourceLocaleDescription,
		},
		{"STATUS", status},
	}

	for _, row := range info {
		fmt.Fprintf(
			table,
			"%s\t%s\n",
			row...,
		)
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}
