## multi-gitter merge

Merge pull requests.

### Synopsis

Merge pull requests with a specified branch name in an organization and with specified conditions.

```
multi-gitter merge [flags]
```

### Options

```
  -B, --branch string   The name of the branch where changes are committed. (default "multi-gitter-branch")
  -h, --help            help for merge
  -o, --org string      The name of the GitHub organization.
```

### Options inherited from parent commands

```
  -g, --gh-base-url string   Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.
  -T, --token string         The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories

