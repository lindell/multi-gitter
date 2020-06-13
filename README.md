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

### Usage
```
Usage:
  multi-gitter run [script path] [flags]

Flags:
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -R, --max-reviewers int       If this value is set, reviewers will be randomized
  -o, --org string              The name of the GitHub organization.
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. (default "https://api.github.com/")
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
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
