## multi-gitter run

Clones multiple repostories, run a script in that directory, and creates a PR with those changes.

### Synopsis


This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY_NAME will be set to the name of the repository currently being executed by the script.


```
multi-gitter run [script path] [flags]
```

### Options

```
      --author-email string     If set, this fields will be used as the email of the committer
      --author-name string      If set, this fields will be used as the name of the committer
  -B, --branch string           The name of the branch where changes are committed. (default "multi-gitter-branch")
  -m, --commit-message string   The commit message. Will default to title + body if none is set.
  -d, --dry-run                 Run without pushing changes or creating pull requests
  -h, --help                    help for run
  -M, --max-reviewers int       If this value is set, reviewers will be randomized
  -b, --pr-body string          The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string         The title of the PR. Will default to the first line of the commit message if none is set.
  -r, --reviewers strings       The username of the reviewers to be added on the pull request.
```

### Options inherited from parent commands

```
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

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories

