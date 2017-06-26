# Smartling

A CLI interface for the [Smartling Translation API](https://help.smartling.com/v1.0/reference) in Go.

# Installation

### From sources
```
go get github.com/Smartling/smartling-cli
```

# Building package

```
make <target>
```

Where target is:

* `deb` for building deb-packages:
   ```
   make deb
   ```

Specific settings can be set in built-time:

*VERSION*:

```
make VERSION=<version-string> <target>
```

*MAINTAINER*:

```
make MAINTAINER=<maintainer> <target>
```

An executable named `smartling-cli` should become available in your
`$GOPATH/bin`.

### Windows
```
curl --output smartling-cli.exe https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.windows.exe
```
### Mac OS
```
curl --output smartling-cli https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.darwin
sudo chmod +x smartling-cli
sudo mv smartling-cli /usr/local/bin/
```
### Linux
```
curl --output smartling-cli https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.linux
sudo chmod +x smartling-cli
sudo mv smartling-cli /usr/local/bin/
```

# Initial configuration

First of all, you need to create config file with authentication parameters
for your project. To ease process, there is `init` command:

```
smartling init
```

# Example usages

## Listing project locales

### Listing all target locales

Display all target project locales along with their description.

```
smartling projects locales -p <project-id>
```

### Listing only locale IDs

Display short form of locales list.

```
smartling projects locales -p <project-id> -s
```

### Display only source locale

```
smartling projects locales -p <project-id> --source
```

### Display only source locale ID

Display short form of source locale.

```
smartling projects locales -p <project-id> --source -s
```

### Display only enabled target locales

Dislay only enable locales using extended output formatting.

```
smartling projects locales -p <project-id> --format='{{if .Enabled}}{{.LocaleID}}{{end}}\n'
```

## Uploading files

### Simplest one-file upload

Upload file as is for translation with automatic file type detection.

```
smartling files push -p <project-id> my-file.txt
```

### One-file upload with URI

Pushes specified `my-file.txt` from the local directory into Smartling with
using remote URI `/my/super/file.txt`.

```
smartling files push -p <project-id> my-file.txt /my/super/file.txt
```

### Override file type

Pushes specified `README.md` from the local directory into Smartling with
as it is plaintext file.

```
smartling files push -p <project-id> README.md --type plaintext
```

### Upload files by mask

Uploads all `txt` files from local directory (and all subdirectories) into
Smartling. *Note* single quotes, it's not shell expansion syntax.

```
smartling files push -p <project-id> '**.txt'
```

### Upload files into specific branch

Uploads all `txt` files from local directory (and all subdirectories) into
Smartling, prefixing each path with `testing/`

```
smartling files push -p <project-id> '**.txt' -b 'testing/'
```

### Partial file type override via config file

Pushes `txt` and `md` files, but overrides `md` file type to `plaintext` using
configuration file (see below).

```
smartling -c smartling.yml files push -p <project-id> '**.{txt,md}'
```

#### smartling.yml

```yaml
# authentication parameters

files:
    "**.md":
        push:
            type: "plaintext"
```
