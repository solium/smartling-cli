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
  smartling [options] files list <project> [-s] [--format=] 
  smartling [options] files pull <project> [<uri>] [-l=] [-d=]
  smartling [options] files status <project> [<uri>] [--format=] 

Commands:
  projects                Used to access various project sub-commands.
   list                   Lists projects for current account.
   get <project>          Get project details about specific project.
                           Accepts project ID as <project> parameter.
  files                    Used to access various files sub-commands.
   status <project>       Shows file translation status.
    --format <format>     Specifies format to use for file status output.
                           [default: $FILE_STATUS_FORMAT]
   list <project>         Lists files from specified project.
    -s --short            Output only file URI.
    --format <format>     Specifies format to use for file list output.
                           [default: $FILE_LIST_FORMAT]
   pull <project> <uri>   Pulls specified files from server. URI supports
                           following globbing patterns:
                            > ** — matches any number of any chars;
                            > *  — matches any number of chars except '/';
                            > ?  — matches any single char except '/';
                            > [xyz]   — matches 'x', 'y' or 'z' charachers;
                            > [!xyz]  — matches not 'x', 'y' or 'z' charachers;
                            > {a,b,c} — matches alternatives a, b or c;
    -d --directory <dir>  Download all files to specified directory.
	--format <format>     Can be used to format path to downloaded files. Note,
	                       that single file can be translated in different
						   locales, so format should include locale to create
						   several file paths.
						   [default: $FILE_PULL_FORMAT]

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
  -l --locale <locale>    Sets locale to filter by or operate upon. Depends on
                           command.
  -d --directory <dir>    Sets directory to operate on, usually, to store or to
                           read files.  Depends on command.  [default: .]
`

const (
	defaultFilesListFormat  = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	defaultFileStatusFormat = `{{.FileURI}}\t{{.Locale}}\t{{.Status}}\t{{.Progress}}\n`
	defaultFilePullFormat   = `{{name .FileURI}}@{{.Locale}}{{ext .FileURI}}`
)

func main() {

	usage := os.Expand(usage, func(key string) string {
		switch key {
		case "HOME":
			return os.Getenv(key)

		case "FILE_LIST_FORMAT":
			return defaultFilesListFormat

		case "FILE_PULL_FORMAT":
			return defaultFilePullFormat

		case "FILE_STATUS_FORMAT":
			return defaultFileStatusFormat
		}

		return key
	})

	args, err := docopt.Parse(usage, nil, true, "smartling "+version, false)
	if err != nil {
		panic(err)
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

	case args["pull"].(bool):
		return doFilesPull(client, config, args)

	case args["status"].(bool):
		return doFilesStatus(client, config, args)

	}

	return nil
}
