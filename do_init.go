package main

import (
	"fmt"
	"io/ioutil"

	hierr "github.com/reconquest/hierr-go"

	yaml "gopkg.in/yaml.v2"
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

	result, err := yaml.Marshal(config)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to encode new config representation",
		)
	}

	if args["--dry-run"].(bool) {
		fmt.Println()
		fmt.Println("Dry run is specified, do not writing config.")
		fmt.Println("New config is dislpayed below.")
		fmt.Println()

		fmt.Println(string(result))
	}

	err = ioutil.WriteFile(config.path, result, 0644)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to write new config file",
		)
	}

	return nil
}
