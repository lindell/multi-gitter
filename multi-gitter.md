## multi-gitter

Multi gitter is a tool for making changes into multiple git repositories

### Synopsis

Multi gitter is a tool for making changes into multiple git repositories

### Options

```
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -h, --help                 help for multi-gitter
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -o, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -P, --platform string      The platform that is used. Available values: github, gitlab (default "github")
  -p, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -u, --user strings         The name of a user. All repositories owned by that user will be used.
```

### SEE ALSO

* [multi-gitter merge](multi-gitter_merge.md)	 - Merge pull requests.
* [multi-gitter run](multi-gitter_run.md)	 - Clones multiple repostories, run a script in that directory, and creates a PR with those changes.
* [multi-gitter status](multi-gitter_status.md)	 - Get the status of pull requests.
* [multi-gitter version](multi-gitter_version.md)	 - Get the version of multi-gitter.

