## multi-gitter print

Clones multiple repositories, run a script in that directory, and prints the output of each run.

### Synopsis


This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. The output of each script run in each repo will be printed, by default to stdout and stderr, but it can be configured to files as well.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.


```
multi-gitter print [script path] [flags]
```

### Options

```
  -g, --base-url string       Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -C, --concurrent int        The maximum number of concurrent runs (default 1)
  -E, --error-output string   The file that the output of the script should be outputted to. "-" means stderr (default "-")
  -f, --fetch-depth int       Limit fetching to the specified number of commits. Set to 0 for no limit (default 1)
  -G, --group strings         The name of a GitLab organization. All repositories in that group will be used.
  -h, --help                  help for print
      --log-file string       The file where all logs should be printed to. "-" means stdout
      --log-format string     The formating of the logs. Available values: text, json, json-pretty (default "text")
  -L, --log-level string      The level of logging that should be made. Available values: trace, debug, info, error (default "info")
  -O, --org strings           The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string         The file that the output of the script should be outputted to. "-" means stdout (default "-")
  -p, --platform string       The platform that is used. Available values: github, gitlab, gitea (default "github")
  -P, --project strings       The name, including owner of a GitLab project in the format "ownerName/repoName"
  -R, --repo strings          The name, including owner of a GitHub repository in the format "ownerName/repoName"
  -T, --token string          The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.
  -U, --user strings          The name of a user. All repositories owned by that user will be used.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories

