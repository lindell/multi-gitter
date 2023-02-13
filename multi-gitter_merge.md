## multi-gitter merge

Merge pull requests.

### Synopsis

Merge pull requests with a specified branch name in an organization and with specified conditions.

```
multi-gitter merge [flags]
```

### Options

```
  -g, --base-url string      Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -B, --branch string        The name of the branch where changes are committed. (default "multi-gitter-branch")
      --config string        Path of the config file.
      --fork                 Use pull requests made from forks instead of from the same repository.
      --fork-owner string    If set, use forks from the defined value instead of the logged in user.
  -G, --group strings        The name of a GitLab organization. All repositories in that group will be used.
  -h, --help                 help for merge
      --include-subgroups    Include GitLab subgroups when using the --group flag.
      --insecure             Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string      The file where all logs should be printed to. "-" means stdout. (default "-")
      --log-format string    The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string     The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
      --merge-type strings   The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed. (default [merge,squash,rebase])
  -O, --org strings          The name of a GitHub organization. All repositories in that organization will be used.
  -p, --platform string      The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings      The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings         The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --ssh-auth             Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string         The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
      --topic strings        The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings         The name of a user. All repositories owned by that user will be used.
  -u, --username string      The Bitbucket server username.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories.

