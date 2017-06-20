# Smartling

A CLI interface for the [Smartling Translation API](help.smartling.com/v1.0/reference) in Go.

# Example usages

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
Smartling. *Note single quotes, it's not shell expansion syntax.`

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

