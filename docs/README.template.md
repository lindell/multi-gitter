multi-gitter
----

Multi-gitter is a tool that allows you to run a script or program for every repository in a GitHub organisation that will then be committed and a PR will be created.

The script can be both a shell script or a binary. If the script returns with a 0 exit code and has made changes to the directory, a PR will be created.

### Usage
```
{{.Usage}}
```

## Example script

```bash
#!/bin/bash

FILE="README.md"
if [ ! -f "$FILE" ]; then
    exit 1
fi

echo "Some extra text" >> README.md
```
