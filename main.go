package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Smartling/api-sdk-go"
	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/hierr-go"
)

var version = "1.0"

var usage = `smartling - manage translation files using Smartling.

TBD.

Usage:
  smartling -h | --help
  smartling [options] [-v]... projects list
  smartling [options] [-v]... projects info -p=<project>
  smartling [options] [-v]... projects locales -p=<project> [--source] [--format=]
  smartling [options] [-v]... files list [--format=] [<uri>]
  smartling [options] [-v]... files pull [-l=]... [-d=] [--source] [<uri>]
  smartling [options] [-v]... files push <file> [<uri>] [(-z|-l=...)] [-b=] [-t=]
  smartling [options] [-v]... files status [-d=] [--format=] [<uri>]

All <uri> arguments support globbing with following patterns:
  > ** — matches any number of any chars;
  > *  — matches any number of chars except '/';
  > ?  — matches any single char except '/';
  > [xyz]   — matches 'x', 'y' or 'z' charachers;
  > [!xyz]  — matches not 'x', 'y' or 'z' charachers;
  > {a,b,c} — matches alternatives a, b or c;

Commands:
  projects                 Used to access various project sub-commands.
   list                    Lists projects for current account.
   info                    Get project details about specific project.
   locales                 Display list of target locales.
    -s --short             Display only target locale IDs.
    --format <format>      Use specified format for listing locales.
                            [format: $PROJECTS_LOCALES_FORMAT]
  files                    Used to access various files sub-commands.
   status <uri>            Shows file translation status.
    --format <format>      Specifies format to use for file status output.
                            [default: $FILE_STATUS_FORMAT]
   list <uri>              Lists files from specified project.
    -s --short             Output only file URI.
    --format <format>      Specifies format to use for file list output.
                            [default: $FILE_LIST_FORMAT]
   pull <uri>              Pulls specified files from server.
    --source               Pulls source file as well.
    -d --directory <dir>   Download all files to specified directory.
    --format <format>      Can be used to format path to downloaded files.
                            Note, that single file can be translated in
                            different locales, so format should include locale
                            to create several file paths.
                            [default: $FILE_PULL_FORMAT]
   push <file> <uri>       Uploads specified file into Smartling platform.
    -z --authorize         Automatically authorize all locales in specified
                            file. Incompatible with -l option.
    -l --locale <locale>   Authorize only specified locales.
    -b --branch <branch>   Prepend specified text to the file uri.
    -t --type <type>       Specifies file type which will be used instead of
                            automatically deduced from extension.


Options:
  -h --help               Show this help.
  -c --config <file>      Config file in YAML format.
                           [default: smartling.yml]
  -p --project <project>  Project ID to operate on.
                           This option ovverides config value "project_id".
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
  -z --authorize          Authorize all locales while pushing file.
                           Incompatible with -l option.
  -b --branch <branch>    Prepend specified value to the file URI.
  -t --type <type>        Specify file type. Depends on command.
  --threads <number>      If command can be executed concurrently, it will be
                           executed for at most <number> of threads.
                           [default: 4]
  -v --verbose            Sets verbosity level for logging messages. Specify
                           flag several time to increase verbosity. Useful
                           when debugging and investigating unexpected
                           behavior.
`

var (
	logger = lorg.NewLog()
)

const (
	defaultProjectsLocalesFormat = `{{.LocaleID}}\t{{.Description}}\t{{.Enabled}}\n`
	defaultFilesListFormat       = `{{.FileURI}}\t{{.LastUploaded}}\t{{.FileType}}\n`
	defaultFileStatusFormat      = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
	defaultFilePullFormat        = `{{name .FileURI}}{{with .Locale}}_{{.}}{{end}}{{ext .FileURI}}`
)

func main() {
	usage := os.Expand(usage, func(key string) string {
		switch key {
		case "FILE_LIST_FORMAT":
			return defaultFilesListFormat

		case "FILE_PULL_FORMAT":
			return defaultFilePullFormat

		case "FILE_STATUS_FORMAT":
			return defaultFileStatusFormat

		case "PROJECTS_LOCALES_FORMAT":
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

	switch args["--verbose"].(int) {
	case 1:
		logger.SetLevel(lorg.LevelInfo)

	case 2:
		logger.SetLevel(lorg.LevelDebug)
	}

	logger.SetFormat(lorg.NewFormat("* ${time} ${level:[%s]:right} %s"))
	logger.SetIndentLines(true)

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

	if args["--project"] != nil {
		config.ProjectID = args["--project"].(string)
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

	switch {
	case args["files"].(bool):
		fallthrough
	case args["projects"].(bool) && args["info"].(bool):
		if config.ProjectID == "" {
			return config, MissingConfigValueError{
				ConfigPath: config.path,
				ValueName:  "project ID",
				OptionName: "project",
				KeyName:    "project_id",
			}
		}
	}

	threads, err := strconv.ParseInt(args["--threads"].(string), 10, 0)
	if err != nil {
		return config, InvalidConfigValueError{
			ValueName:   "threads",
			Description: "should be valid integer number",
		}
	}

	if threads <= 0 {
		return config, InvalidConfigValueError{
			ValueName:   "threads",
			Description: "should be positive integer number",
		}
	}

	if config.Threads == 0 {
		config.Threads = int(threads)
	}

	return config, nil
}

func projects(config Config, args map[string]interface{}) error {
	client := smartling.NewClient(config.UserID, config.Secret)

	setLogger(client, logger, args["--verbose"].(int))

	switch {
	case args["list"].(bool):
		if config.AccountID == "" {
			return MissingConfigValueError{
				ConfigPath: config.path,
				ValueName:  "account ID",
				OptionName: "account",
				KeyName:    "account_id",
			}
		}
	}

	switch {
	case args["list"].(bool):
		return doProjectsList(client, config, args)

	case args["info"].(bool):
		return doProjectsInfo(client, config, args)

	case args["locales"].(bool):
		return doProjectsLocales(client, config, args)

	}

	return nil
}

func files(config Config, args map[string]interface{}) error {
	client := smartling.NewClient(config.UserID, config.Secret)

	setLogger(client, logger, args["--verbose"].(int))

	switch {
	case args["list"].(bool):
		return doFilesList(client, config, args)

	case args["pull"].(bool):
		return doFilesPull(client, config, args)

	case args["push"].(bool):
		return doFilesPush(client, config, args)

	case args["status"].(bool):
		return doFilesStatus(client, config, args)

	}

	return nil
}
