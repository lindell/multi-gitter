<h1 align="center">
  <img alt="" src="docs/img/logo.svg" height="80" />
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ABuilding"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Building/badge.svg?branch=master" /></a>
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ATesting"><img alt="Go test status" src="https://github.com/lindell/multi-gitter/workflows/Testing/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
</div>
<br>

*multi-gitter* allows you to make changes in multiple repositories simultaneously. This is achieved by running a script or program in the context of multiple repositories. If any changes are made, a pull request is created that can be merged manually by the set reviewers, or automatically by multi-gitter when CI pipelines has completed successfully.

Are you a bash-guru or simply prefer your scripting in Node.js? It doesn't matter, since multi-gitter support any type of script or program. **If you can script it to run in one place, you can run it in all your repositories with one command!**

### Some examples:
* Syncing a file (like a PR-template)
* Programmatic refactoring
* Updating a dependency
* Automatically fixing linting issues
* Search and replace
* Anything else you are able to script!

## Demo

![Gif](docs/img/demo.gif)

## Example

### Run with file
```bash
$ multi-gitter run ./my-script.sh -O my-org -m "Commit message" -B branch-name
```

Make sure the script has execution permissions before running it (`chmod +x ./my-script.sh`)

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

### Homebrew
If you are using Mac or Linux, [Homebrew](https://brew.sh/) is an easy way of installing multi-gitter.
```bash
brew install lindell/multi-gitter/multi-gitter
```

### Manual binary install
Find the binary for your operating system from the [release page](https://github.com/lindell/multi-gitter/releases) and download it.

### Automatic binary install
To automatically install the latest version
```bash
curl -s https://raw.githubusercontent.com/lindell/multi-gitter/master/install.sh | sh
```

### From source
You can also install from source with `go install`, this is not recommended for most cases.
```bash
go install github.com/lindell/multi-gitter
```

## Token

To use multi-gitter, a token that is allowed to list repositories and create pull requests is needed. This token can either be set in the `GITHUB_TOKEN`, `GITLAB_TOKEN`, `GITEA_TOKEN` environment variable, or by using the `--token` flag.

### GitHub
[How to generate a GitHub personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token). Make sure to give to `repo` permissions.

### GitLab

[How to generate a GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html). Make sure to give to it the `api` permission.

### Gitea

In Gitea, access tokens can be generated under Settings -> Applications -> Manage Access Tokens 

## Config file

All configuration in multi-gitter can be done through command line flags, configuration files or a mix of both. If you want to use a configuration file, simply use the `--config=./path/to/config.yaml`. Multi-gitter will also read from the file `~/.multi-gitter/config` and take and configuration from there. The priority of configs are first flags, then defined config file and lastly the static config file.



<details>
  <summary>All available run options</summary>

```yaml
# Email of the committer. If not set, the global git config setting will be used.
author-email:

# Name of the committer. If not set, the global git config setting will be used.
author-name:

# The branch which the changes will be based on.
base-branch:

# Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# The commit message. Will default to title + body if none is set.
commit-message:

# The maximum number of concurrent runs.
concurrent: 1

# Run without pushing changes or creating pull requests.
dry-run: false

# Limit fetching to the specified number of commits. Set to 0 for no limit.
fetch-depth: 1

# Fork the repository instead of creating a new branch on the same owner.
fork: false

# If set, make the fork to the defined value. Default behavior is for the fork to be on the logged in user.
fork-owner:

# The type of git implementation to use.
# Available values:
#   go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
#   cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
git-type: go

# The name of a GitLab organization. All repositories in that group will be used.
group:
  - example

# Include GitLab subgroups when using the --group flag.
include-subgroups: false

# Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
insecure: false

# Take manual decision before committing any change. Requires git to be installed.
interactive: false

# The file where all logs should be printed to. "-" means stdout.
log-file: "-"

# The formating of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# If this value is set, reviewers will be randomized.
max-reviewers: 0

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.
platform: github

# The body of the commit message. Will default to everything but the first line of the commit message if none is set.
pr-body:

# The title of the PR. Will default to the first line of the commit message if none is set.
pr-title:

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# The username of the reviewers to be added on the pull request.
reviewers:
  - example

# Skip pull request and directly push to the branch.
skip-pr: false

# The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
token:

# The name of a user. All repositories owned by that user will be used.
user:
  - example

# The Bitbucket server username.
username:
```
</details>


<details>
  <summary>All available merge options</summary>

```yaml
# Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use pull requests made from forks instead of from the same repository.
fork: false

# If set, use forks from the defined value instead of the logged in user.
fork-owner:

# The name of a GitLab organization. All repositories in that group will be used.
group:
  - example

# Include GitLab subgroups when using the --group flag.
include-subgroups: false

# Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
insecure: false

# The file where all logs should be printed to. "-" means stdout.
log-file: "-"

# The formating of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed.
merge-type:
  - merge
  - squash
  - rebase

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
token:

# The name of a user. All repositories owned by that user will be used.
user:
  - example

# The Bitbucket server username.
username:
```
</details>


<details>
  <summary>All available status options</summary>

```yaml
# Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use pull requests made from forks instead of from the same repository.
fork: false

# If set, use forks from the defined value instead of the logged in user.
fork-owner:

# The name of a GitLab organization. All repositories in that group will be used.
group:
  - example

# Include GitLab subgroups when using the --group flag.
include-subgroups: false

# Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
insecure: false

# The file where all logs should be printed to. "-" means stdout.
log-file: "-"

# The formating of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
token:

# The name of a user. All repositories owned by that user will be used.
user:
  - example

# The Bitbucket server username.
username:
```
</details>


<details>
  <summary>All available close options</summary>

```yaml
# Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use pull requests made from forks instead of from the same repository.
fork: false

# If set, use forks from the defined value instead of the logged in user.
fork-owner:

# The name of a GitLab organization. All repositories in that group will be used.
group:
  - example

# Include GitLab subgroups when using the --group flag.
include-subgroups: false

# Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
insecure: false

# The file where all logs should be printed to. "-" means stdout.
log-file: "-"

# The formating of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
token:

# The name of a user. All repositories owned by that user will be used.
user:
  - example

# The Bitbucket server username.
username:
```
</details>


<details>
  <summary>All available print options</summary>

```yaml
# Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
base-url:

# The maximum number of concurrent runs.
concurrent: 1

# The file that the output of the script should be outputted to. "-" means stderr.
error-output: "-"

# Limit fetching to the specified number of commits. Set to 0 for no limit.
fetch-depth: 1

# The type of git implementation to use.
# Available values:
#   go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
#   cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
git-type: go

# The name of a GitLab organization. All repositories in that group will be used.
group:
  - example

# Include GitLab subgroups when using the --group flag.
include-subgroups: false

# Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
insecure: false

# The file where all logs should be printed to. "-" means stdout.
log-file:

# The formating of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
token:

# The name of a user. All repositories owned by that user will be used.
user:
  - example

# The Bitbucket server username.
username:
```
</details>


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
  -g, --base-url string         Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -C, --concurrent int          The maximum number of concurrent runs. (default 1)
      --config string           Path of the config file.
  -d, --dry-run                 Run without pushing changes or creating pull requests.
  -f, --fetch-depth int         Limit fetching to the specified number of commits. Set to 0 for no limit. (default 1)
      --fork                    Fork the repository instead of creating a new branch on the same owner.
      --fork-owner string       If set, make the fork to the defined value. Default behavior is for the fork to be on the logged in user.
      --git-type string         The type of git implementation to use.
                                Available values:
                                  go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
                                  cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
                                 (default "go")
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups       Include GitLab subgroups when using the --group flag.
      --insecure                Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
  -i, --interactive             Take manual decision before committing any change. Requires git to be installed.
      --log-file string         The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string       The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -M, --max-reviewers int       If this value is set, reviewers will be randomized.
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string           The file that the output of the script should be outputted to. "-" means stdout. (default "-")
  -p, --platform string         The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.
      --skip-pr                 Skip pull request and directly push to the branch.
  -T, --token string            The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
  -u, --username string         The Bitbucket server username.
```


### <img alt="merge" src="docs/img/fa/code-merge.svg" height="40" valign="middle" /> Usage of `merge`
Merge pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter merge [flags]

Flags:
  -g, --base-url string      Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
      --config string        Path of the config file.
      --fork                 Use pull requests made from forks instead of from the same repository.
      --fork-owner string    If set, use forks from the defined value instead of the logged in user.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups    Include GitLab subgroups when using the --group flag.
      --insecure             Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string      The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string    The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
      --merge-type strings   The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed. (default [merge,squash,rebase])
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string      The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
  -u, --username string      The Bitbucket server username.
```


### <img alt="status" src="docs/img/fa/tasks.svg" height="40" valign="middle" /> Usage of `status`
Get the status of all pull requests with a specified branch name in an organization.
```
Usage:
  multi-gitter status [flags]

Flags:
  -g, --base-url string     Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string       The name of the branch where changes are committed. (default "multi-gitter-branch")
      --config string       Path of the config file.
      --fork                Use pull requests made from forks instead of from the same repository.
      --fork-owner string   If set, use forks from the defined value instead of the logged in user.
  -G, --group strings       The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups   Include GitLab subgroups when using the --group flag.
      --insecure            Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string     The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string   The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string    The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings         The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string       The file that the output of the script should be outputted to. "-" means stdout. (default "-")
  -p, --platform string     The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings     The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings        The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -T, --token string        The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings        The name of a user. All repositories owned by that user will be used.
  -u, --username string     The Bitbucket server username.
```


### <img alt="close" src="docs/img/fa/times-hexagon.svg" height="40" valign="middle" /> Usage of `close`
Close pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter close [flags]

Flags:
  -g, --base-url string     Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string       The name of the branch where changes are committed. (default "multi-gitter-branch")
      --config string       Path of the config file.
      --fork                Use pull requests made from forks instead of from the same repository.
      --fork-owner string   If set, use forks from the defined value instead of the logged in user.
  -G, --group strings       The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups   Include GitLab subgroups when using the --group flag.
      --insecure            Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string     The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string   The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string    The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings         The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string     The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings     The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings        The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -T, --token string        The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings        The name of a user. All repositories owned by that user will be used.
  -u, --username string     The Bitbucket server username.
```


### <img alt="print" src="docs/img/fa/print.svg" height="40" valign="middle" /> Usage of `print`

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. The output of each script run in each repo will be printed, by default to stdout and stderr, but it can be configured to files as well.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.

```
Usage:
  multi-gitter print [script path] [flags]

Flags:
  -g, --base-url string       Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -C, --concurrent int        The maximum number of concurrent runs. (default 1)
      --config string         Path of the config file.
  -E, --error-output string   The file that the output of the script should be outputted to. "-" means stderr. (default "-")
  -f, --fetch-depth int       Limit fetching to the specified number of commits. Set to 0 for no limit. (default 1)
      --git-type string       The type of git implementation to use.
                              Available values:
                                go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
                                cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
                               (default "go")
  -G, --group strings         The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups     Include GitLab subgroups when using the --group flag.
      --insecure              Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string       The file where all logs should be printed to. "-" means stdout.
      --log-format string     The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string      The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings           The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string         The file that the output of the script should be outputted to. "-" means stdout. (default "-")
  -p, --platform string       The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings       The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings          The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -T, --token string          The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings          The name of a user. All repositories owned by that user will be used.
  -u, --username string       The Bitbucket server username.
```



## Example scripts

### general

<details>
  <summary>Replace a file if it exist</summary>

```sh
#!/bin/bash

REPLACE_FILE=~/test/pull_request_template.md # The file that should replace the file in the repo, must be an absolute path
FILE=.github/pull_request_template.md # Relative from any repos root

# Don't replace this file if it does not already exist in the repo
if [ ! -f "$FILE" ]; then
    exit 1
fi

cp $REPLACE_FILE $FILE
```
</details>

<details>
  <summary>Replace text in all files</summary>

```sh
#!/bin/bash

# Assuming you are using gnu sed, if you are running this on a mac, please see https://stackoverflow.com/questions/4247068/sed-command-with-i-option-failing-on-mac-but-works-on-linux

find ./ -type f -exec sed -i -e 's/apple/orange/g' {} \;
```
</details>

### go

<details>
  <summary>Fix linting problems in all your go repositories</summary>

```sh
#!/bin/bash

golangci-lint run ./... --fix
```
</details>

<details>
  <summary>Updates a go module to a new (patch/minor) version</summary>

```sh
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

```sh
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

<details>
  <summary>Simple replace using node</summary>

```js
const { readFile, writeFile } = require("fs").promises;

async function replace() {
  let data = await readFile("./README.md", "utf8");
  data = data.replace("apple", "orange");
  await writeFile("./README.md", data, "utf8");
}

replace();
```
</details>


Do you have a nice script that might be useful to others? Please create a PR that adds it to the [examples folder](/examples).
