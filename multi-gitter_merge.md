## multi-gitter merge

Merge pull requests.

### Synopsis

Merge pull requests with a specified branch name in an organization and with specified conditions.

```
multi-gitter merge [flags]
```

### Options

```
  -g, --base-url string      Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
      --config string        Path of the config file.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -h, --help                 help for merge
      --include-subgroups    Include GitLab subgroups when using the --group flag.
      --log-file string      The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string    The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
      --merge-type strings   The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed. (default [merge,squash,rebase])
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string      The platform that is used. Available values: github, gitlab, gitea. (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -T, --token string         The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories.

