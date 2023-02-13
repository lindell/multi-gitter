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
  -g, --base-url string       Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket.
  -C, --concurrent int        The maximum number of concurrent runs. (default 1)
      --config string         Path of the config file.
  -E, --error-output string   The file that the output of the script should be outputted to. "-" means stderr. (default "-")
  -f, --fetch-depth int       Limit fetching to the specified number of commits. Set to 0 for no limit. (default 1)
      --git-type string       The type of git implementation to use.
                              Available values:
                                go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
                                cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
                               (default "go")
  -G, --group strings         The name of a GitLab organization. All repositories in that group will be used.
  -h, --help                  help for print
      --include-subgroups     Include GitLab subgroups when using the --group flag.
      --insecure              Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.
      --log-file string       The file where all logs should be printed to. "-" means stdout.
      --log-format string     The formatting of the logs. Available values: text, json, json-pretty. (default "text")
  -L, --log-level string      The level of logging that should be made. Available values: trace, debug, info, error. (default "info")
  -O, --org strings           The name of a GitHub organization. All repositories in that organization will be used.
  -o, --output string         The file that the output of the script should be outputted to. "-" means stdout. (default "-")
  -p, --platform string       The platform that is used. Available values: github, gitlab, gitea, bitbucket_server. (default "github")
  -P, --project strings       The name, including owner of a GitLab project in the format "ownerName/repoName".
  -R, --repo strings          The name, including owner of a GitHub repository in the format "ownerName/repoName".
      --ssh-auth              Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.
  -T, --token string          The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.
      --topic strings         The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.
  -U, --user strings          The name of a user. All repositories owned by that user will be used.
  -u, --username string       The Bitbucket server username.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories.

