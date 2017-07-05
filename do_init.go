package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/reconquest/hierr-go"
)

func doInit(config Config, args map[string]interface{}) error {
	fmt.Printf("Generating %s...\n", config.path)

	prompt := func(
		message string,
		value interface{},
		zero bool,
		input interface{},
	) {
		fmt.Print(message)

		if !zero {
			fmt.Printf(" [%s]: ", value)
		} else {
			fmt.Printf(": ")
		}

		fmt.Scanln(input)
	}

	var input Config

	prompt(
		"Enter User ID",
		config.UserID,
		config.UserID == "",
		&input.UserID,
	)

	if input.UserID != "" {
		config.UserID = input.UserID
	}

	prompt(
		"Enter Token Secret",
		config.Secret,
		config.Secret == "",
		&input.Secret,
	)

	if input.Secret != "" {
		config.Secret = input.Secret
	}

	prompt(
		"Enter Account ID (optional)",
		config.AccountID,
		config.AccountID == "",
		&input.AccountID,
	)

	if input.AccountID != "" {
		config.AccountID = input.AccountID
	}

	prompt(
		"Enter Project ID",
		config.ProjectID,
		config.ProjectID == "",
		&input.ProjectID,
	)

	if input.ProjectID != "" {
		config.ProjectID = input.ProjectID
	}

	var result bytes.Buffer
	err := configTemplate.Execute(&result, config)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to compile config template",
		)
	}

	if args["--dry-run"].(bool) {
		fmt.Println()
		fmt.Println("Dry run is specified, do not writing config.")
		fmt.Println("New config is dislpayed below.")
		fmt.Println()

		fmt.Println(result.String())
	} else {
		err = ioutil.WriteFile(config.path, result.Bytes(), 0644)
		if err != nil {
			return hierr.Errorf(
				err,
				"unable to write new config file",
			)
		}
	}

	return nil
}
