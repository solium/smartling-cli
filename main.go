package main

import (
	"fmt"
	"os"

	"github.com/Smartling/api-sdk-go"
	"github.com/docopt/docopt-go"
	"github.com/reconquest/hierr-go"
)

var version = "1.0"

var usage = `smartling - manage translation files using Smartling.

TBD.

Usage:
  smartling -h | --help
  smartling [options] projects list
  smartling [options] projects info <project>
  smartling [options] files list [-s] [--format=] <project>

Commands:
  projects               Used to access various project sub-commands.
    list                 Lists projects for current account.
    get                  Get project details about specific project.
                          Accepts project ID as <project> parameter.
  files                  Used to access various files sub-commands.
    list                 Lists files from specified project.
      -s --short         Output only file URI.
      --format <format>  Specifies format to use for file list output.
                          [default: $FILES_FORMAT_TEMPLATE]

Options:
  -h --help               Show this help.
  -c --config <file>      Config file in YAML format.
                           [default: $HOME/.config/smartling/config.yml]
  -p --project <project>  Project ID to operate on.
                           This option ovverides config value "project_id".
                           This option is not used by "projects" command.
  -a --account <account>  Account ID to operate on.
                           This option ovverides config value "account_id".
  --user <user>           User ID which will be used for authentication.
                           This option ovverides config value "user_id".
  --secret <secret>       Token Secret which will be used for authentication.
                           This option ovverides config value "secret".
  -s --short              Use short list output, usually outputs only first
                           column, e.g. file URI in case of files list.
`

func main() {
	const (
		defaultFilesListFormat = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	)

	usage := os.Expand(usage, func(key string) string {
		switch key {
		case "HOME":
			return os.Getenv(key)
		case "FILES_FORMAT_TEMPLATE":
			return defaultFilesListFormat
		}

		return key
	})

	args, err := docopt.Parse(usage, nil, true, "smartling "+version, false)
	if err != nil {
		panic(err)
	}

	if args["--format"] == nil {
		args["--format"] = defaultFilesListFormat
	}

	config, err := loadConfig(args)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	switch {
	case args["projects"].(bool):
		err = projects(config, args)

	case args["files"].(bool):
		err = files(config, args)

	}

	if err != nil {
		fmt.Println(err)
	}
}

func loadConfig(args map[string]interface{}) (Config, error) {
	path := args["--config"].(string)

	config, err := NewConfig(path)
	if err != nil {
		return config, NewError(
			hierr.Errorf(err, `failed to load configuration file "%s".`, path),
			`Check configuretion file contents according to documentation.`,
		)
	}

	if args["--user"] != nil {
		config.UserID = args["--user"].(string)
	}

	if args["--secret"] != nil {
		config.Secret = args["--secret"].(string)
	}

	if args["--account"] != nil {
		config.AccountID = args["--account"].(string)
	}

	if config.UserID == "" {
		return config, MissingConfigValueError{
			ConfigPath: config.path,
			ValueName:  "user ID",
			OptionName: "user",
			KeyName:    "user_id",
		}
	}

	if config.Secret == "" {
		return config, MissingConfigValueError{
			ConfigPath: config.path,
			ValueName:  "token secret",
			OptionName: "secret",
			KeyName:    "secret",
		}
	}

	return config, nil
}

func projects(config Config, args map[string]interface{}) error {
	client := smartling.NewClient(config.UserID, config.Secret)

	if config.AccountID == "" {
		return MissingConfigValueError{
			ConfigPath: config.path,
			ValueName:  "account ID",
			OptionName: "account",
			KeyName:    "account_id",
		}
	}

	switch {
	case args["list"].(bool):
		return doProjectsList(client, config, args)

	case args["info"].(bool):
		return doProjectsGet(client, config, args)
	}

	return nil
}

func files(config Config, args map[string]interface{}) error {
	client := smartling.NewClient(config.UserID, config.Secret)

	switch {
	case args["list"].(bool):
		return doFilesList(client, config, args)
	}

	return nil
}
