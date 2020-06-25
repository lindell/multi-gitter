<h1 align="center">
  <img alt="" src="docs/img/logo.svg" height="80" />
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3AGo"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Go/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
</div>
<br>

Multi-gitter is a tool that allows you to run a script or program for every repository in a GitHub organisation that will then be committed and a PR will be created.

The script can be both a shell script or a binary. If the script returns with a 0 exit code and has made changes to the directory, a PR will be created.

## Usage

* [run](#-usage-of-run) Clones multiple repostories, run a script in that directory, and creates a PR with those changes.
* [merge](#-usage-of-merge) Merge pull requests.
* [status](#-usage-of-status) Get the status of pull requests.


### <img alt="run" src="docs/img/fa/rabbit-fast.svg" height="40" valign="middle" /> Usage of `run`
Run will clone down multiple repositories. For each of those repositories, the script will be run. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.
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
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -L, --log-level string     The level of logging that should be made. Available values: debug, info, error (default "info")
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
```


### <img alt="merge" src="docs/img/fa/code-merge.svg" height="40" valign="middle" /> Usage of `merge`
Merge pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter merge [flags]

Flags:
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")
  -o, --org string      The name of the GitHub organization.

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -L, --log-level string     The level of logging that should be made. Available values: debug, info, error (default "info")
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
```


### <img alt="status" src="docs/img/fa/tasks.svg" height="40" valign="middle" /> Usage of `status`
Get the status of all pull requests with a specified branch name in an organization.
```
Usage:
  multi-gitter status [flags]

Flags:
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")
  -o, --org string      The name of the GitHub organization.

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -L, --log-level string     The level of logging that should be made. Available values: debug, info, error (default "info")
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
