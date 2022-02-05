# gh-org-repo-sync

[![Build](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml)
[![CodeQL Analysis](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml)

> GitHub CLI extension to clone repositories in a GitHub Organization and update their default branch if already cloned

## Installation

- Install the `gh` CLI. See [https://github.com/cli/cli#installation](https://github.com/cli/cli#installation) for further details.
- Run `gh auth login` to authenticate with your GitHub account. Alternatively, the CLI will respect the `GITHUB_TOKEN` [environment variable](https://cli.github.com/manual/gh_help_environment).
- Install this extension:

```bash
gh extension install rm3l/gh-org-repo-sync
```

## Usage

```bash
❯ gh org-repo-sync -h

Usage: gh org-repo-sync <organization> [options]
Options: 
  -batchSize int
        the number of elements to retrieve at once. Must not exceed 100 (default 50)
  -output string
        the output path (default ".")
  -protocol string
        the protocol to use for cloning. Possible values: system, ssh, https. (default "system")
  -query string
        GitHub search query, to filter the Organization repositories. Example: "language:Go stars:>10 pushed:>2010-11-12"
        See https://bit.ly/3HurHe3 for more details on the search syntax
```

### Working with the source code

Clone the repository and install the local version.

```bash
cd gh-org-repo-sync

# Install the local version
gh extension install .

# At this point, you can start using it
gh org-repo-sync <my-organization>

# To see changes in the code as you develop, simply build and use the extension:
go build && gh org-repo-sync <my-organization>
```
