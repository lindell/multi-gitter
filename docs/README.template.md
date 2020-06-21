<h1 align="center">
  🛠 multi-gitter
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3AGo"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Go/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
</div>
<br>

Multi-gitter is a tool that allows you to run a script or program for every repository in a GitHub organisation that will then be committed and a PR will be created.

The script can be both a shell script or a binary. If the script returns with a 0 exit code and has made changes to the directory, a PR will be created.

## Usage
```
{{.MainUsage}}
```

{{range .Commands}}
### {{.TitleExtra}} Usage of `{{.Name}}`
{{.Description}}
```
{{.Usage}}
```

{{end}}

## Example script

```bash
#!/bin/bash

FILE="README.md"
if [ ! -f "$FILE" ]; then
    exit 1
fi

echo "Some extra text" >> README.md
```
