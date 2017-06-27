# Example usages

**Table of content**

* [Introduction](#introduction)
* [Projects](#work-with-projects)
* Files
  * [Upload files](#upload-files)
  * [Download files](#download-files)
  * [List files](#list-files)
  * [Status](#status)
  * [Delete files](#delete-files)

## Introduction

All examples below assume that you configured credentials, project id and account id by `init` command. But don't forget about "global options" that allows you to override any configuration parameters. An example:
```
-c --config <file>      Config file in YAML format.
                         [default: smartling.yml]
-p --project <project>  Project ID to operate on.
                         This option overrides config value "project_id".
-a --account <account>  Account ID to operate on.
                         This option overrides config value "account_id".
--user <user>           User ID which will be used for authentication.
                         This option overrides config value "user_id".
--secret <secret>       Token Secret which will be used for authentication.
                         This option overrides config value "secret".
-s --short              Use short list output, usually outputs only first
                         column, e.g. file URI in case of files list.
-d --directory <dir>    Sets directory to operate on, usually, to store or to
                         read files.  Depends on command.  [default: .]
```

Also would recommend to review template of configuration file [smartling.yml.example](smartling.yml.example). It's very useful in case you have existing "localization project" locally and need to synchronize it with Smartling TMS.

## Work with projects

### List of all projects in your account

```
$ smartling-cli projects list
2f2xxxxx  Sitecore Connector           en-US
129xxxxx  Wordpress Connector          en
855xxxxx  Drupal Connector             en
7c7xxxxx  Word files                   en
```

### Information about project

```
$ smartling-cli projects info
ID       2f2xxxxx
ACCOUNT  a3exxxxx
NAME     Test site
LOCALE   en-US: English (United States)
STATUS   active
```

### Listing all target locales

Display all target project locales along with their description.
```
$ smartling-cli projects locales
zh-CN  Chinese (Simplified)     true
en-AU  English (Australia)      true
fr-FR  French (France)          true
de-DE  German (Germany)         true
it-IT  Italian (Italy)          true
ru-RU  Russian                  true
es     Spanish (International)  true

```

### Listing only locale IDs

Display short form of locales list. You can use `-s` or `--short` option

```
$ smartling-cli projects locales -s
zh-CN
nl-NL
de-DE
```

### Display only source locale

```
$ smartling-cli projects locales --source
en-US  English (United States)
```

Display short form of source locale.

```
$ smartling-cli projects locales --source -s
en-US
```

### Display only enabled target locales

Display only enable locales using extended output formatting.

```
$ smartling-cli projects locales --format='{{if .Enabled}}{{.LocaleID}}{{end}}\n'
zh-CN
nl-NL
de-DE
```

Read more about [Golang templates](https://golang.org/pkg/text/template).

## Upload files

### Simplest one-file upload

Upload file as is for translation with automatic file type detection.

```
$ smartling-cli files push my-file.txt
```

### One-file upload with URI

Pushes specified `my-file.txt` from the local directory into Smartling with
using remote URI `/my/super/file.txt`.

```
$ smartling-cli files push my-file.txt /my/super/file.txt
```

### Override file type

Pushes specified `README.md` from the local directory into Smartling as plaintext file.

```
$ smartling-cli files push README.md --type plaintext
README.md overwritten [3 strings 28 words]
```

### Upload files by mask

Uploads all `txt` files from local directory (and all subdirectories) into
Smartling. *Note* single quotes, it's not shell expansion syntax.

```
$ smartling-cli files push '**.txt'
```

### Branching \ versioning

Uploads all `txt` files from local directory (and all subdirectories) into
Smartling, prefixing each path with `testing` (use `-b` or `--branch` option)

```
$ smartling-cli files push '**.txt' -b 'testing'
testing/test.txt new [3 strings 28 words]
```

If your work directory is git workspace then you can configure CLI to detect the current branch name and use it as prefix.

```
$ smartling-cli files push '**.txt' --branch '@auto'
* 2017-06-27 16:20:54    [INFO] autodetected branch name: master
master/test.txt new [3 strings 28 words]
```

### Partial file type override via config file

Pushes `txt` and `md` files, but overrides `md` file type to `plaintext` using
configuration file (see below)

```
$ smartling-cli files push '**.{txt,md}'
```

#### smartling.yml

```yaml
# authentication parameters

files:
    "**.md":
        push:
            type: "plaintext"
```

CLI looks for configuration file in work directory. But you can always set explicit  path where to look for configuration file (`-c` or `--config` option)

```
$ smartling-cli -c smartling.yml files push '**.{txt,md}'
```


## Download files

### Download only source files

```
$ smartling-cli files pull **test.xml --source
downloaded files\test.xml 0%
downloaded test.xml 0%
downloaded master\test.xml 0%
```

### Download translated files files

Download translated files by mask and just for a couple locales

```
$ smartling-cli files pull '**test.xml' -l es -l fr-FR
downloaded test_es.xml 0%
downloaded files/test_es.xml 0%
downloaded master/test_es.xml 0%
downloaded files/test_fr-FR.xml 0%
downloaded test_fr-FR.xml 0%
downloaded master/test_fr-FR.xml 0%
```

### Use pipes

Get list of files and filter files only which have `_` or `-` in URI and then download original files

```
$ smartling-cli files list '**test.*' --short | grep '[_-]' | smartling-cli files pull --source
downloaded content/dam/geometrixx/portraits/yolanda_huggins.jpg-65359_de.xml 0%
downloaded content/geometrixx-media/en-14306_de.xml 0%
downloaded content/dam/geometrixx/portraits/scott_reynolds.jpg-65359_de.xml 0%
downloaded content/geometrixx-outdoors/en-96925_de.xml 0%
downloaded content/geometrixx/en-31001_de.xml 0%
```

### Copy files between projects

Download all `properties` files locally from the configured project

```
$ smartling-cli files pull **ep?.properties --source
downloaded files/ep1.properties 0%
downloaded files/ep5.properties 0%
downloaded files/ep2.properties 0%
```

Upload all local `properties` files into another project (`129xxxxx`)

```
$ smartling-cli files push **.properties -p 129xxxxx
files/ep1.properties new [1 strings 2 words]
files/ep2.properties new [1 strings 2 words]
files/ep5.properties new [1 strings 3 words]
```

Also you can specify another credentials for `push` command in case a target project sits in another account.

## List files

### List of all files in project

```
$ smartling-cli files list
/files/ep1.properties                 2016-08-04T08:18:14Z  javaProperties
/files/ep2.properties                 2016-08-04T08:30:44Z  javaProperties
/files/ep5.properties                 2016-08-10T14:04:45Z  javaProperties
/files/example.STRINGSDICT            2016-08-26T08:47:54Z  stringsdict
/files/Localizable.stringsdict        2016-09-15T13:13:37Z  stringsdict
/files/Localizable2.stringsdict       2016-09-23T16:22:49Z  stringsdict
....
```

### List files by mask

Get list of file URIs filtered by mask and return only file URIs
```
$ smartling-cli files list '**test.xml' --short
/files/placeholder_test.xml
/files/test.xml
master/test.xml
test.xml
```

### Custom output format

Use another format for rendering data in table and use more complex mask

```
$ smartling-cli files list '**{id,_}**.xml' --format='{{.FileType}}\t{{.FileURI}}\n'
xml  /files/placeholder_test.xml
xml  new_test_file.xml
xml  test-client-id.xml
xml  test_dropbox_variants_ns1.xml
xml  test_dropbox_variants_ns2.xml
xml  test_keys.xml
```

## Status

### Files status

Get full overview for your local and remote files

```
$ smartling-cli files status
config/locales/account.en.yml           en-US  missing  source  5     5
config/locales/account.en_be-BY.yml     be-BY  missing  0%      0     0
config/locales/account.en_de-DE.yml     de-DE  missing  100%    5     5
config/locales/account.en_en-AU.yml     en-AU  missing  0%      0     0
config/locales/account.en_es.yml        es     missing  0%      0     0
config/locales/account.en_es-ES.yml     es-ES  missing  100%    5     5
config/locales/account.en_fr-FR.yml     fr-FR  missing  100%    5     5
config/locales/account.en_it-IT.yml     it-IT  missing  0%      0     0
config/locales/account.en_nl-NL.yml     nl-NL  missing  100%    5     5
config/locales/account.en_ru-RU.yml     ru-RU  missing  0%      0     0
config/locales/account.en_zh-CN.yml     zh-CN  missing  0%      0     0
test2.json                              en-US  missing  source  1     4
test2_be-BY.json                        be-BY  missing  0%      0     0
test2_it-IT.json                        it-IT  missing  0%      0     0
test2_zh-CN.json                        zh-CN  missing  0%      0     0
test2_ru-RU.json                        ru-RU  missing  100%    1     4
test2_es.json                           es     missing  100%    1     4
test2_de-DE.json                        de-DE  missing  0%      0     0
test2_en-AU.json                        en-AU  missing  0%      0     0
test2_es-ES.json                        es-ES  missing  0%      0     0
test2_fr-FR.json                        fr-FR  missing  0%      0     0
....
```


## Delete files

### Delete files by mask

Delete all files from `master` branch

```
$ smartling-cli files delete 'master/**'
master/test.xml deleted
```

### Use pipes

Delete files listed in `files-list.txt` file

```
$ cat files-list.txt | smartling-cli files delete -
/files/placeholder_test.xml deleted
/files/test.xml deleted
test.xml deleted
```

### Delete all files in project

Let's say the default configured project is `129xxxxx`. But we need to delete **all** files in another project (`2a1xxxxxx`) but use configured credentials

```
$ smartling-cli files list -s -p 2a1xxxxxx | smartling-cli files delete - -p 2a1xxxxxx
27 Words To Learn Before You Visit Hawaii.docx deleted
....
```
