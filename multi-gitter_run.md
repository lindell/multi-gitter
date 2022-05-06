## multi-gitter run

Clones multiple repositories, run a script in that directory, and creates a PR with those changes.

### Synopsis


This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.


```
multi-gitter run [script path] [flags]
```

### Options

```
  -a, --assignees strings          The username of the assignees to be added on the pull request.
      --author-email string        Email of the committer. If not set, the global git config setting will be used.
      --author-name string         Name of the committer. If not set, the global git config setting will be used.
      --base-branch string         The branch which the changes will be based on.
  -g, --base-url string            Base URL of the GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.
  -B, --branch string              The name of the branch where changes are committed. (default "multi-gitter-branch")
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
  -h, --help                       help for run
      --include-subgroups          Include GitLab subgroups when using the --group flag.
      --insecure                   Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
  -i, --interactive                Take manual decision before committing any change. Requires git to be installed.
      --log-file string            The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string          The formating of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string           The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -M, --max-reviewers int          If this value is set, reviewers will be randomized.
  -O, --org strings                The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string              The file that the output of the script should be outputted to. "-" means stdout. (default "-")
  -p, --platform string            The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -b, --pr-body string             The body of the commit message. Will default to everything but the first line of the commit message if none is set.
  -t, --pr-title string            The title of the PR. Will default to the first line of the commit message if none is set.
  -P, --project strings            The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings               The name, including owner of a GitHub repository in the format "ownerName/repoName".
  -r, --reviewers strings          The username of the reviewers to be added on the pull request.
      --skip-pr                    Skip pull request and directly push to the branch.
  -s, --skip-repo strings          Skip changes on specified repositories, the name is including the owner of repository in the format "ownerName/repoName".
      --ssh-auth                   Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string               The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
  -U, --user strings               The name of a user. All repositories owned by that user will be used.
  -u, --username string            The Bitbucket server username.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories.

