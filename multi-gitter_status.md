## multi-gitter status

Get the status of pull requests.

### Synopsis

Get the status of all pull requests with a specified branch name in an organization.

```
multi-gitter status [flags]
```

### Options

```
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")
  -h, --help            help for status
```

### Options inherited from parent commands

```
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -o, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
  -u, --user strings         The name of a GitHub user. All repositories owned by that user will be used.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories

