package main

import (
	"fmt"
)

const formatOptionHelp = `
That command supports advanced formatting via --format flag with full
support of Golang templates (https://golang.org/pkg/text/template).
`

const authenticationOptionsHelp = `
  --user <user>
    Specify user ID for authentication.

  --secret <secret>
    Specify secret token for authentication.

  -a --account <account>
    Specify account ID.
`

const initHelp = `smartling init — create config file interactively.

Walk down common config file parameters and fill them through dialog.

Init process will inspect if config file already exists and if it is, it will
be loaded as default values, so init can be used sequentially without config
is lost.

Options like --user, --secret, --account and --project can be used to specify
config values prior dialog:

  smartling init --user=your_user_id

Also, --dry-run option can be used to just look at resulting config without
overwritting anything:

  smartling init --dry-run

By default, smartling.yml file in the local directory will be used as target
config file, but it can be overriden by using --config option:

  smartling init --config=/path/to/project/smartling.yml


Available options:
  -c --config <file>
    Specify config file to operate on. Default: smartling.yml

  --dry-run
    Do not overwrite config file, only output to stdout.

Default config values can be passed via following options:` +
	authenticationOptionsHelp + `
  -p --project <project>
    Specify default project.
`

const projectsListHelp = `smartling projects list — list projects from account.

Command will list projects from specified account in tabular format with
following information:

  > Project ID
  > Project Description
  > Project Source Locale ID

Only project IDs will be listed if --short option is specified.

Note, that you should specify account ID either in config file or via --account
option to be able to see projects list.


Available options:
  -s --short
    List only project IDs.
` + authenticationOptionsHelp

const projectsInfoHelp = `smartling projects info — show detailed project info.

Displays detailed information for specific project.

Project should be specified either in config or via --project option.

Available options:` + authenticationOptionsHelp

const projectsLocalesHelp = `smartling projects locales — list target locales.

Lists target locales from specified project.

To list only locale IDs --short option can be used.
` + formatOptionHelp + `
Following variables are available:

  > .LocaleID — target locale ID to translate into;
  > .Description — human-readable locale description;
  > .Enabled — true/false specifying is locale active or not;

Available options:
  -s --short
    List only locale IDs.

  --format
    Use specific output format instead of default.
` + authenticationOptionsHelp

func showHelp(args map[string]interface{}) {
	switch {
	case args["init"].(bool):
		fmt.Print(initHelp)

	case args["projects"].(bool):
		switch {
		case args["list"].(bool):
			fmt.Print(projectsListHelp)

		case args["info"].(bool):
			fmt.Print(projectsInfoHelp)

		case args["locales"].(bool):
			fmt.Print(projectsLocalesHelp)
		}

	default:
		fmt.Print(usage)
	}
}
