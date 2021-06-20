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

{{range .Commands}}
{{if .YAMLExample}}
<details>
  <summary>All available {{.Name}} options</summary>

```yaml
{{ .YAMLExample }}
```
</details>
{{end}}{{end}}

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

## Example scripts
{{range .ExampleCategories}}
### {{.Name}}
{{range .Examples}}
<details>
  <summary>{{.Title}}</summary>

```{{.Type}}
{{.Body}}
```
</details>
{{end}}{{end}}

Do you have a nice script that might be useful to others? Please create a PR that adds it to the [examples folder](/examples).
