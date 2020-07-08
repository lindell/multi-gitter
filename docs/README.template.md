<h1 align="center">
  <img alt="" src="docs/img/logo.svg" height="80" />
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3AGo"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Go/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
</div>
<br>

multi-gitter is a tool that allows you to run a script or program for multiple repositories. Changes done by the script will be made as a pull request or similar and can then be merged manually by the set reviewers, or automatically by multi-gitter when CI pipelines has completed successfully.

multi-gitter currently supports GitHub and GitLab where you can run the script on all repositories in an organization, group, user or specify individual repositories. For each repository, the script will run in the context of the root folder, and if any changes is done to the filesystem together with an exit code of 0, the changes will be committed and pushed as a pull/merge request.

## Usage
{{range .Commands}}
* [{{ .Name }}](#-usage-of-{{ .Name }}) {{ .Short }}{{end}}

{{range .Commands}}
### <img alt="{{.Name}}" src="{{.ImageIcon}}" height="40" valign="middle" /> Usage of `{{.Name}}`
{{.Long}}
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
