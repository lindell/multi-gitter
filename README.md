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

## Install

### Manual binary install
Find the binary for your operating system from the [release page](https://github.com/lindell/multi-gitter/releases) and download it.

### Automatic binary install
To automatically install the latest version
```bash
curl -s https://raw.githubusercontent.com/lindell/multi-gitter/master/install.sh | sh
```

### From source
You can also install from source with `go get`, this is not recommended for most cases.
```bash
go get github.com/lindell/multi-gitter
```

## Usage

* [run](#-usage-of-run) Clones multiple repostories, run a script in that directory, and creates a PR with those changes.
* [merge](#-usage-of-merge) Merge pull requests.
* [status](#-usage-of-status) Get the status of pull requests.


### <img alt="run" src="docs/img/fa/rabbit-fast.svg" height="40" valign="middle" /> Usage of `run`

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY_NAME will be set to the name of the repository currently being executed by the script.

```
Usage:
  multi-gitter run [script path] [flags]

Flags:
      --author-email string     If set, this fields will be used as the email of the committer
      --author-name string      If set, this fields will be used as the name of the committer
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -d, --dry-run                 Run without pushing changes or creating pull requests
  -M, --max-reviewers int       If this value is set, reviewers will be randomized
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -o, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -P, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -p, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -u, --user strings         The name of a user. All repositories owned by that user will be used.
```


### <img alt="merge" src="docs/img/fa/code-merge.svg" height="40" valign="middle" /> Usage of `merge`
Merge pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter merge [flags]

Flags:
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -o, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -P, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -p, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -u, --user strings         The name of a user. All repositories owned by that user will be used.
```


### <img alt="status" src="docs/img/fa/tasks.svg" height="40" valign="middle" /> Usage of `status`
Get the status of all pull requests with a specified branch name in an organization.
```
Usage:
  multi-gitter status [flags]

Flags:
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")

Global Flags:
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -o, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -P, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -p, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -u, --user strings         The name of a user. All repositories owned by that user will be used.
```



## Example scripts

### general

<details>
  <summary>Replace text in all files</summary>

```bash
#!/bin/bash

# Assuming you are using gnu sed, if you are running this on a mac, please see https://stackoverflow.com/questions/4247068/sed-command-with-i-option-failing-on-mac-but-works-on-linux

find ./ -type f -exec sed -i -e 's/apple/orange/g' {} \;
```
</details>

### go

<details>
  <summary>Fix linting problems in all your go repositories</summary>

```bash
#!/bin/bash

golangci-lint run ./... --fix
```
</details>

<details>
  <summary>Updates a go module to a new (patch/minor) version</summary>

```bash
#!/bin/bash

### Change these values ###
MODULE=github.com/go-git/go-git/v5
VERSION=v5.1.0

# Check if the module already exist, abort if it does not
go list -m $MODULE &> /dev/null
status_code=$?
if [ $status_code -ne 0 ]; then
    echo "Module \"$MODULE\" does not exist"
    exit 1
fi

go get $MODULE@$VERSION
```
</details>

### node

<details>
  <summary>Updates a npm dependency if it does exist</summary>

```bash
#!/bin/bash

### Change these values ###
PACKAGE=webpack
VERSION=4.43.0

if [ ! -f "package.json" ]; then
    echo "package.json does not exist"
    exit 1
fi

# Check if the package already exist (without having to install all packages first), abort if it does not
current_version=`jq ".dependencies[\"$PACKAGE\"]" package.json`
if [ "$current_version" == "null" ];
then
    echo "Package \"$PACKAGE\" does not exist"
    exit 2
fi

npm install --save $PACKAGE@$VERSION
```
</details>

