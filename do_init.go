package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Smartling/api-sdk-go"
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
		"Enter Smartling API V2.0 User Identifier:",
		config.UserID,
		config.UserID == "",
		&input.UserID,
	)

	if input.UserID != "" {
		config.UserID = input.UserID
	}

	prompt(
		"Enter Smartling API V2.0 Token Secret:",
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

	fmt.Println("Testing connection to Smartling API...")

	client, err := createClient(config, args)
	if err != nil {
		return hierr.Errorf(
			err,
			"unable to create client for testing connection",
		)
	}

	err = client.Authenticate()
	if err != nil {
		if _, ok := err.(smartling.NotAuthorizedError); ok {
			return NewError(
				errors.New("Not authorized."),
				"Your credentials are invalid. Double check them and run "+
					"init command again.",
			)
		} else {
			return NewError(
				hierr.Errorf(err, "failure while testing connection"),
				"Contact developer for more info.",
			)
		}
	}

	fmt.Println("Connection is successfull.")

	if args["--dry-run"].(bool) {
		fmt.Println()
		fmt.Println("Dry run is specified, do not writing config.")
		fmt.Println("New config is displayed below.")
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
