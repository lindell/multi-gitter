<h1 align="center">
  <img alt="" src="docs/img/logo.svg" height="80" />
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ABuilding"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Building/badge.svg?branch=master" /></a>
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ATesting"><img alt="Go test status" src="https://github.com/lindell/multi-gitter/workflows/Testing/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
</div>
<br>

*multi-gitter* allows you to make changes in multiple repositories simultaneously. This is archived by running a script or program in the context of all repositories and if any changes are made, a pull request is created that can be merged manually by the set reviewers, or automatically by multi-gitter when CI pipelines has completed successfully.

It currently supports GitHub and GitLab where you can run it on all repositories in an organization, group, user or specify individual repositories. For each repository, the script will run in the context of the root folder, and if any changes is done to the filesystem together with an exit code of 0, the changes will be committed and pushed as a pull/merge request.

## Demo

![Gif](docs/img/demo.gif)

## Example

### Run with file
```bash
$ multi-gitter run ./my-script.sh -O my-org -m "Commit message" -B branch-name
```

### Run code through interpreter
If you are running an interpreted language or similar, it's important to specify the path as an absolute value (since the script will be run in the context of each repository). Using the `$PWD` variable helps with this.
```bash
$ multi-gitter run "python $PWD/run.py" -O my-org -m "Commit message" -B branch-name
$ multi-gitter run "node $PWD/script.js" -R repo1 -R repo2 -m "Commit message" -B branch-name
$ multi-gitter run "go run $PWD/main.go" -U my-user -m "Commit message" -B branch-name
```

### Test before live run
You might want to test your changes before creating commits. The `--dry-run` provides an easy way to test without actually making any modifications. It works well with setting the log level to `debug` with `--log-level=debug` to also print the changes that would have been made.
```
$ multi-gitter run ./script.sh --dry-run --log-level=debug -O my-org -m "Commit message" -B branch-name
```

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

* [run](#-usage-of-run) Clones multiple repositories, run a script in that directory, and creates a PR with those changes.
* [merge](#-usage-of-merge) Merge pull requests.
* [status](#-usage-of-status) Get the status of pull requests.
* [close](#-usage-of-close) Close pull requests.
* [print](#-usage-of-print) Clones multiple repositories, run a script in that directory, and prints the output of each run.


### <img alt="run" src="docs/img/fa/rabbit-fast.svg" height="40" valign="middle" /> Usage of `run`

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.

```
Usage:
  multi-gitter run [script path] [flags]

Flags:
      --author-email string     Email of the committer. If not set, the global git config setting will be used.
      --author-name string      Name of the committer. If not set, the global git config setting will be used.
      --base-branch string      The branch which the changes will be based on.
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -C, --concurrent int          The maximum number of concurrent runs (default 1)
  -d, --dry-run                 Run without pushing changes or creating pull requests
  -g, --gh-base-url string      Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --log-file string         The file where all logs should be printed to. "-" means stdout (default "-")
      --log-format string       The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -M, --max-reviewers int       If this value is set, reviewers will be randomized
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string           The file that the output of the script should be outputted to. "-" means stdout (default "-")
  -p, --platform string         The platform that is used. Available values: github, gitlab (default "github")
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.
  -T, --token string            The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
```


### <img alt="merge" src="docs/img/fa/code-merge.svg" height="40" valign="middle" /> Usage of `merge`
Merge pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter merge [flags]

Flags:
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
      --log-file string      The file where all logs should be printed to. "-" means stdout (default "-")
      --log-format string    The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
```


### <img alt="status" src="docs/img/fa/tasks.svg" height="40" valign="middle" /> Usage of `status`
Get the status of all pull requests with a specified branch name in an organization.
```
Usage:
  multi-gitter status [flags]

Flags:
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
      --log-file string      The file where all logs should be printed to. "-" means stdout (default "-")
      --log-format string    The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string        The file that the output of the script should be outputted to. "-" means stdout (default "-")
  -p, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
```


### <img alt="close" src="docs/img/fa/times-hexagon.svg" height="40" valign="middle" /> Usage of `close`
Close pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter close [flags]

Flags:
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
      --log-file string      The file where all logs should be printed to. "-" means stdout (default "-")
      --log-format string    The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
```


### <img alt="print" src="docs/img/fa/print.svg" height="40" valign="middle" /> Usage of `print`

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. The output of each script run in each repo will be printed, by default to stdout and stderr, but it can be configured to files as well.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.

```
Usage:
  multi-gitter print [script path] [flags]

Flags:
  -C, --concurrent int        The maximum number of concurrent runs (default 1)
  -E, --error-output string   The file that the output of the script should be outputted to. "-" means stderr (default "-")
  -g, --gh-base-url string    Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings         The name of a GitLab organization. All repositories in that group will be used.
      --log-file string       The file where all logs should be printed to. "-" means stdout
      --log-format string     The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string      The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -O, --org strings           The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string         The file that the output of the script should be outputted to. "-" means stdout (default "-")
  -p, --platform string       The platform that is used. Available values: github, gitlab (default "github")
  -P, --project strings       The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings          The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string          The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings          The name of a user. All repositories owned by that user will be used.
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

