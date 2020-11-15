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
$ multi-gitter run ./my-script.sh -o my-org -m "Commit message" -B branch-name
```

### Run code through interpreter
If you are running an interpreted language or similar, it's important to specify the path as an absolute value (since the script will be run in the context of each repository). Using the `$PWD` variable helps with this.
```bash
$ multi-gitter run "python $PWD/run.py" -o my-org -m "Commit message" -B branch-name
$ multi-gitter run "node $PWD/script.js" -R repo1 -r repo2 -m "Commit message" -B branch-name
$ multi-gitter run "go run $PWD/main.go" -u my-user -m "Commit message" -B branch-name
```

### Test before live run
You might want to test your changes before creating commits. The `--dry-run` provides an easy way to test without actually making any modifications. It works well with setting the log level to `debug` with `--log-level=debug` to also print the changes that would have been made.
```
$ multi-gitter run ./script.sh --dry-run --log-level=debug -o my-org -m "Commit message" -B branch-name
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

```bash
{{.Body}}
```
</details>
{{end}}{{end}}
