<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/img/logo-dark-mode.svg" />
    <img alt="Multi-gitter logo" src="docs/img/logo.svg" height="80" height="80" />
  </picture>
</h1>

<div align="center">
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ABuilding"><img alt="Go build status" src="https://github.com/lindell/multi-gitter/workflows/Building/badge.svg?branch=master" /></a>
  <a href="https://github.com/lindell/multi-gitter/actions?query=branch%3Amaster+workflow%3ATesting"><img alt="Go test status" src="https://github.com/lindell/multi-gitter/workflows/Testing/badge.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/lindell/multi-gitter"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/lindell/multi-gitter" /></a>
  <a href="https://securityscorecards.dev/viewer/?uri=github.com/lindell/multi-gitter"><img alt="OpenSSF Scorecard" src="https://api.securityscorecards.dev/projects/github.com/lindell/multi-gitter/badge" /></a>
</div>
<br>

*multi-gitter* allows you to make changes in multiple repositories simultaneously. This is achieved by running a script or program in the context of multiple repositories. If any changes are made, a pull request is created that can be merged manually by the set reviewers, or automatically by multi-gitter when CI pipelines have completed successfully.

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
You might want to test your changes before creating commits. The `--dry-run` flag provides an easy way to test without actually making any modifications. It works well when setting the log level to `debug`, with `--log-level=debug`, to also print the changes that would have been made.
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
go install github.com/lindell/multi-gitter@latest
```

## Token

To use multi-gitter, a token that is allowed to list repositories and create pull requests is needed. This token can either be set in the `GITHUB_TOKEN`, `GITLAB_TOKEN`, `GITEA_TOKEN` environment variable, or by using the `--token` flag.

### GitHub

[How to generate a GitHub personal access token (classic)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic). Make sure to give it `repo` permissions.

### GitLab

[How to generate a GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html). Make sure to give to it the `api` permission.

### Gitea

In Gitea, access tokens can be generated under Settings -> Applications -> Manage Access Tokens

## Config file

All configuration in multi-gitter can be done through command line flags, configuration files or a combination of both. If you want to use a configuration file, simply use the `--config=./path/to/config.yaml` option. Multi-gitter will also read from the file `~/.multi-gitter/config` and take and configuration from there. The priority of configs are first flags, then defined config file and lastly the static config file.



<details>
  <summary>All available run options</summary>

```yaml
# The username of the assignees to be added on the pull request.
assignees:
  - example

# Email of the committer. If not set, the global git config setting will be used.
author-email:

# Name of the committer. If not set, the global git config setting will be used.
author-name:

# The branch which the changes will be based on.
base-branch:

# Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# The temporary directory where the repositories will be cloned. If not set, the default os temporary directory will be used.
clone-dir:

# Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
code-search:

# The commit message. Will default to title + body if none is set.
commit-message:

# The maximum number of concurrent runs.
concurrent: 1

# What should happen if the branch already exist.
# Available values:
#   skip: Skip making any changes to the existing branch and do not create a new pull request.
#   replace: Replace the existing content of the branch by force pushing any new changes, then reuse any existing pull request, or create a new one if none exist.
conflict-strategy: skip

# Create pull request(s) as draft.
draft: false

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

# Labels to be added to any created pull request.
labels:
  - example

# The file where all logs should be printed to. "-" means stdout.
log-file: "-"

# The formatting of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# If this value is set, reviewers will be randomized.
max-reviewers: 0

# If this value is set, team reviewers will be randomized
max-team-reviewers: 0

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# Don't use any terminal formatting when printing the output.
plain-output: false

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta
platform: github

# The body of the commit message. Will default to everything but the first line of the commit message if none is set.
pr-body:

# The title of the PR. Will default to the first line of the commit message if none is set.
pr-title:

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# Skip pull request and only push the feature branch.
push-only: false

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# Exclude repositories that match with a given Regular Expression
repo-exclude:

# Include repositories that match with a given Regular Expression
repo-include:

# Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
repo-search:

# The username of the reviewers to be added on the pull request.
reviewers:
  - example

# Skip repositories which are forks.
skip-forks: false

# Skip pull request and directly push to the branch.
skip-pr: false

# Skip changes on specified repositories, the name is including the owner of repository in the format "ownerName/repoName".
skip-repo:
  - example

# Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
ssh-auth: false

# Github team names of the reviewers, in format: 'org/team'
team-reviewers:
  - example

# The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
token:

# The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
topic:
  - example

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
# Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
code-search:

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

# The formatting of the logs. Available values: text, json, json-pretty.
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

# Don't use any terminal formatting when printing the output.
plain-output: false

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
repo-search:

# Skip repositories which are forks.
skip-forks: false

# Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
ssh-auth: false

# The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
token:

# The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
topic:
  - example

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
# Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
code-search:

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

# The formatting of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# Don't use any terminal formatting when printing the output.
plain-output: false

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
repo-search:

# Skip repositories which are forks.
skip-forks: false

# Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
ssh-auth: false

# The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
token:

# The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
topic:
  - example

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
# Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
base-url:

# The name of the branch where changes are committed.
branch: multi-gitter-branch

# Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
code-search:

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

# The formatting of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# Don't use any terminal formatting when printing the output.
plain-output: false

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
repo-search:

# Skip repositories which are forks.
skip-forks: false

# Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
ssh-auth: false

# The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
token:

# The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
topic:
  - example

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
# Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
base-url:

# The temporary directory where the repositories will be cloned. If not set, the default os temporary directory will be used.
clone-dir:

# Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
code-search:

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

# The formatting of the logs. Available values: text, json, json-pretty.
log-format: text

# The level of logging that should be made. Available values: trace, debug, info, error.
log-level: info

# The name of a GitHub organization. All repositories in that organization will be used.
org:
  - example

# The file that the output of the script should be outputted to. "-" means stdout.
output: "-"

# Don't use any terminal formatting when printing the output.
plain-output: false

# The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta
platform: github

# The name, including owner of a GitLab project in the format "ownerName/repoName".
project:
  - group/project

# The name, including owner of a GitHub repository in the format "ownerName/repoName".
repo:
  - my-org/js-repo
  - other-org/python-repo

# Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
repo-search:

# Skip repositories which are forks.
skip-forks: false

# Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
ssh-auth: false

# The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
token:

# The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
topic:
  - example

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

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created.

When the script is invoked, these environment variables are set:
- REPOSITORY will be set to the name of the repository currently being executed
- DRY_RUN will be set =true, when running in with the --dry-run flag, otherwise it's absent

```
Usage:
  multi-gitter run [script path] [flags]

Flags:
  -a, --assignees strings          The username of the assignees to be added on the pull request.
      --author-email string        Email of the committer. If not set, the global git config setting will be used.
      --author-name string         Name of the committer. If not set, the global git config setting will be used.
      --base-branch string         The branch which the changes will be based on.
  -g, --base-url string            Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -B, --branch string              The name of the branch where changes are committed. (default "multi-gitter-branch")
      --clone-dir string           The temporary directory where the repositories will be cloned. If not set, the default os temporary directory will be used.
      --code-search fork:true      Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use fork:true to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
  -m, --commit-message string      The commit message. Will default to title + body if none is set.
  -C, --concurrent int             The maximum number of concurrent runs. (default 1)
      --config string              Path of the config file.
      --conflict-strategy string   What should happen if the branch already exist.
                                   Available values:
                                     skip: Skip making any changes to the existing branch and do not create a new pull request.
                                     replace: Replace the existing content of the branch by force pushing any new changes, then reuse any existing pull request, or create a new one if none exist.
                                    (default "skip")
      --draft                      Create pull request(s) as draft.
  -d, --dry-run                    Run without pushing changes or creating pull requests.
  -f, --fetch-depth int            Limit fetching to the specified number of commits. Set to 0 for no limit. (default 1)
      --fork                       Fork the repository instead of creating a new branch on the same owner.
      --fork-owner string          If set, make the fork to the defined value. Default behavior is for the fork to be on the logged in user.
      --git-type string            The type of git implementation to use.
                                   Available values:
                                     go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
                                     cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
                                    (default "go")
  -G, --group strings              The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups          Include GitLab subgroups when using the --group flag.
      --insecure                   Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
  -i, --interactive                Take manual decision before committing any change. Requires git to be installed.
      --labels strings             Labels to be added to any created pull request.
      --log-file string            The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string          The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string           The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -M, --max-reviewers int          If this value is set, reviewers will be randomized.
      --max-team-reviewers int     If this value is set, team reviewers will be randomized
  -O, --org strings                The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string              The file that the output of the script should be outputted to. "-" means stdout. (default "-")
      --plain-output               Don't use any terminal formatting when printing the output.
  -p, --platform string            The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta (default "github")
  -b, --pr-body string             The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string            The title of the PR. Will default to the first line of the commit message if none is set.
  -P, --project strings            The name, including owner of a GitLab project in the format "ownerName/repoName".
      --push-only                  Skip pull request and only push the feature branch.
  -R, --repo strings               The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --repo-exclude string        Exclude repositories that match with a given Regular Expression
      --repo-include string        Include repositories that match with a given Regular Expression
      --repo-search fork:true      Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use fork:true to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
  -r, --reviewers strings          The username of the reviewers to be added on the pull request.
      --skip-forks                 Skip repositories which are forks.
      --skip-pr                    Skip pull request and directly push to the branch.
  -s, --skip-repo strings          Skip changes on specified repositories, the name is including the owner of repository in the format "ownerName/repoName".
      --ssh-auth                   Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
      --team-reviewers strings     Github team names of the reviewers, in format: 'org/team'
  -T, --token string               The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
      --topic strings              The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings               The name of a user. All repositories owned by that user will be used.
  -u, --username string            The Bitbucket server username.
```


### <img alt="merge" src="docs/img/fa/code-merge.svg" height="40" valign="middle" /> Usage of `merge`
Merge pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter merge [flags]

Flags:
  -g, --base-url string         Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
      --code-search fork:true   Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use fork:true to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
      --config string           Path of the config file.
      --fork                    Use pull requests made from forks instead of from the same repository.
      --fork-owner string       If set, use forks from the defined value instead of the logged in user.
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups       Include GitLab subgroups when using the --group flag.
      --insecure                Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string         The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string       The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
      --merge-type strings      The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed. (default [merge,squash,rebase])
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
      --plain-output            Don't use any terminal formatting when printing the output.
  -p, --platform string         The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta (default "github")
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --repo-search fork:true   Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use fork:true to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
      --skip-forks              Skip repositories which are forks.
      --ssh-auth                Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string            The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
      --topic strings           The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
  -u, --username string         The Bitbucket server username.
```


### <img alt="status" src="docs/img/fa/tasks.svg" height="40" valign="middle" /> Usage of `status`
Get the status of all pull requests with a specified branch name in an organization.
```
Usage:
  multi-gitter status [flags]

Flags:
  -g, --base-url string         Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
      --code-search fork:true   Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use fork:true to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
      --config string           Path of the config file.
      --fork                    Use pull requests made from forks instead of from the same repository.
      --fork-owner string       If set, use forks from the defined value instead of the logged in user.
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups       Include GitLab subgroups when using the --group flag.
      --insecure                Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string         The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string       The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string           The file that the output of the script should be outputted to. "-" means stdout. (default "-")
      --plain-output            Don't use any terminal formatting when printing the output.
  -p, --platform string         The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta (default "github")
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --repo-search fork:true   Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use fork:true to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
      --skip-forks              Skip repositories which are forks.
      --ssh-auth                Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string            The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
      --topic strings           The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
  -u, --username string         The Bitbucket server username.
```


### <img alt="close" src="docs/img/fa/times-hexagon.svg" height="40" valign="middle" /> Usage of `close`
Close pull requests with a specified branch name in an organization and with specified conditions.
```
Usage:
  multi-gitter close [flags]

Flags:
  -g, --base-url string         Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
      --code-search fork:true   Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use fork:true to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
      --config string           Path of the config file.
      --fork                    Use pull requests made from forks instead of from the same repository.
      --fork-owner string       If set, use forks from the defined value instead of the logged in user.
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups       Include GitLab subgroups when using the --group flag.
      --insecure                Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string         The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string       The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
      --plain-output            Don't use any terminal formatting when printing the output.
  -p, --platform string         The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta (default "github")
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --repo-search fork:true   Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use fork:true to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
      --skip-forks              Skip repositories which are forks.
      --ssh-auth                Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string            The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
      --topic strings           The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
  -u, --username string         The Bitbucket server username.
```


### <img alt="print" src="docs/img/fa/print.svg" height="40" valign="middle" /> Usage of `print`

This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. The output of each script run in each repo will be printed, by default to stdout and stderr, but it can be configured to write to files as well.

When the script is invoked, these environment variables are set:
- REPOSITORY will be set to the name of the repository currently being executed

```
Usage:
  multi-gitter print [script path] [flags]

Flags:
  -g, --base-url string         Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
      --clone-dir string        The temporary directory where the repositories will be cloned. If not set, the default os temporary directory will be used.
      --code-search fork:true   Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use fork:true to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.
  -C, --concurrent int          The maximum number of concurrent runs. (default 1)
      --config string           Path of the config file.
  -E, --error-output string     The file that the output of the script should be outputted to. "-" means stderr. (default "-")
  -f, --fetch-depth int         Limit fetching to the specified number of commits. Set to 0 for no limit. (default 1)
      --git-type string         The type of git implementation to use.
                                Available values:
                                  go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
                                  cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
                                 (default "go")
  -G, --group strings           The name of a GitLab organization. All repositories in that group will be used.
      --include-subgroups       Include GitLab subgroups when using the --group flag.
      --insecure                Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string         The file where all logs should be printed to. "-" means stdout.
      --log-format string       The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string        The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings             The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string           The file that the output of the script should be outputted to. "-" means stdout. (default "-")
      --plain-output            Don't use any terminal formatting when printing the output.
  -p, --platform string         The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud. Note: bitbucket_cloud is in Beta (default "github")
  -P, --project strings         The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings            The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --repo-search fork:true   Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use fork:true to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.
      --skip-forks              Skip repositories which are forks.
      --ssh-auth                Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string            The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD environment variable.
      --topic strings           The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings            The name of a user. All repositories owned by that user will be used.
  -u, --username string         The Bitbucket server username.
```



## Example scripts

### general

<details>
  <summary>Clone all repositories locally while maintaining their group folder structure</summary>

```sh
#!/bin/bash

# This script should be used with the print command.
mkdir -p ~/multi-gitter/$REPOSITORY
cp -r . ~/multi-gitter/$REPOSITORY
```
</details>

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
  <summary>Replace all instances of empty interface with any</summary>

```sh
#!/bin/bash

gofmt -r 'interface{} -> any' -w **/*.go
```
</details>

<details>
  <summary>Fix the ioutil deprecation</summary>

```sh
#!/bin/bash

gofmt -w -r 'ioutil.Discard -> io.Discard' .
gofmt -w -r 'ioutil.NopCloser -> io.NopCloser' .
gofmt -w -r 'ioutil.ReadAll -> io.ReadAll' .
gofmt -w -r 'ioutil.ReadFile -> os.ReadFile' .
gofmt -w -r 'ioutil.TempDir -> os.MkdirTemp' .
gofmt -w -r 'ioutil.TempFile -> os.CreateTemp' .
gofmt -w -r 'ioutil.WriteFile -> os.WriteFile' .
gofmt -w -r 'ioutil.ReadDir -> os.ReadDir ' . # (note: returns a slice of os.DirEntry rather than a slice of fs.FileInfo)

goimports -w .
```
</details>

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

<details>
  <summary>Upgrade Go version in go modules</summary>

```sh
#!/bin/bash

go mod edit -go 1.18
go mod tidy
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


<details>

<summary> Bitbucket Cloud </summary>

_note: bitbucket cloud support is currently in Beta_

In order to use bitbucket cloud you will need to create and use an [App Password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/). The app password you create needs sufficient permissions so ensure you grant it Read and Write access to projects, repositories and pull requests and at least Read access to your account and workspace membership.

You will need to configure the bitbucket workspace using the `org` option for multi-gitter for the repositories you want to make changes to e.g. `multi-gitter run examples/go/upgrade-go-version.sh -u your_username --org "your_workspace"`

### Example
Here is an example of using the command line options to run a script from the `examples/` directory and make pull-requests for a few repositories in a specified workspace.
```shell
export BITBUCKET_CLOUD_APP_PASSWORD="your_app_password"
multi-gitter run examples/go/upgrade-go-version.sh -u your_username --org "your_workspace" --repo "your_first_repository,your_second_repository" --platform bitbucket_cloud -m "your_commit_message" -B your_branch_name
```

### Bitbucket Cloud Limitations
Currently, we add the repositories default reviewers as a reviewer for any pull-request you create. If you want to specify specific reviewers, you will need to add them using their `UUID` instead of their username since bitbucket does not allow us to look up a `UUID` using their username. [This article has more information about where you can get a users UUID.](https://community.atlassian.com/t5/Bitbucket-articles/Retrieve-the-Atlassian-Account-ID-AAID-in-bitbucket-org/ba-p/2471787)

We don't support specifying specific projects for bitbucket cloud yet, you should still be able to make changes to the repositories you want but certain functionality, like forking, does not work as well until we implement that feature.
Using `fork: true` is currently experimental within multi-gitter for Bitbucket Cloud, and will be addressed in future updates.
Here are the known limitations:
- The forked repository will appear in the user-provided workspace/orgName, with the given repo name, but will appear in a random project within that workspace.
- Using `git-type: cmd` is required for Bitbucket Cloud forking for now until better support is added, as `git-type: go` causes inconsistent behavior(intermittent unauthorized errors).

We also only support modifying a single workspace, any additional workspaces passed into the multi-gitter `org` option will be ignored after the first value.

We also have noticed the performance is slower with larger workspaces and we expect to resolve this when we add support for projects to make filtering repositories by project faster.

</details>