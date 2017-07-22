# Smartling

A CLI interface for the [Smartling Translation API](https://help.smartling.com/v1.0/reference) in Go.

# Building package

```
make <target>
```

Where target is:

* `deb` for building Debian packages:
   ```
   make deb
   ```

* `rpm` for building Fedora packages:
   ```
   make rpm
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

# Installation

### From sources
```
go get github.com/Smartling/smartling-cli
```

### Windows
```
curl --output smartling-cli.exe https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.windows.exe
```
### Mac OS
```
curl --output smartling-cli https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.darwin \
&& sudo chmod +x smartling-cli \
&& sudo mv smartling-cli /usr/local/bin/
```
### Linux
```
curl --output smartling-cli https://smartling-connectors-releases.s3.amazonaws.com/cli/smartling.linux \
&& sudo chmod +x smartling-cli \
&& sudo mv smartling-cli /usr/local/bin/
```

# Initial configuration

First of all, you need to create config file with authentication parameters
for your project. To ease process, there is `init` command:

```
smartling-cli init
```

# Example usages

Display all target project locales along with their description.
```
smartling-cli projects locales
```

Upload file as is for translation with automatic file type detection.
```
smartling-cli files push my-file.txt
```

Find more [example usages there](examples.md).

# Development

## Managing dependencies

Project uses [manul](https://github.com/kovetskiy/manul) vendoring tool,
which uses vendoring through git-submodules.

After adding any third-party libraries, you need to update vendoring:

```
manul -Ir
```
