Contributing Guide
----

We welcome dontributions in multiple forms. This file describes some of the ways how you can help out in various ways.

## Bug fixes or new feature

You would like to help out with code? Great! Before you get started, please read the following guidelines:

### General guidelines for Pull Requests

💬 Create a ticket discussing the change and get a confirmation that this change will be merged if built before starting your pull request. This is true for most pull requests, except for small changes such as typos.

🚧 Clearly state the state of you pull requests. If the PR is not finished but you want feedback. Mark it as a draft PR.

1️⃣ Limit the pull requests to solve one specific task. 

⬇️ Create a new commit for any changes after the pull request was created. Changes will be squashed once merged. But to make the review process easier, never rewrite the history.

### Get started coding

The only dependency needed for multi-gitter development is `go`. By forking and cloning the repo you are ready to go.

### Test your code

All tests can be run with `go test ./...`. These tests will also run on multiple platforms once you push it.

### Linting

This project is linted with golangci-lint, and the [configuration file exist in the root of this project](/.golangci.yml). Linting is never perfect, and if something is not linting correctly, please discuss it in the PR.

Some rules might seem superfluous, like comments on everything that is exported, but to encourage the addition of as many useful comments as possible, this repository imposes the rule that everything that is exported should have comments.

### Docs

*Don't change the README.md file*. If you want to make changes to the README.md file. Please take a look in `./docs/README.template.md` since the README.md file is generated based on it. If you want to change something like the description of a flag, that change can be made directly in the code and docs will be updated automatically.

### Structure

* **cmd**: This is where the CLI is created. Everything that has to do with parsing flags or similair is done here, but no actual execution logic.
* **docs**: Documentation.
* **examples**: Example scripts that can be used together with multi-gitter.
* **internal**
  * **domain**: This folder contains basic structures that can be used by multiple packages.
  * **git**: All implementations of git.
  * **multi-gitter**: The main logic of multi-gitter. This is the code that glues everything together.
  * **scm**: Source control system implementations such as GitHub/GitLab/etc.
* **test**: Integration tests.
* **tools**: Tools for CI/CD or for development.
