## gh-org-repo-sync

[![Build](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml)

[![CodeQL Analysis](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml)

> GitHub CLI extension to clone repositories in a GitHub Organization and update their default branch if already cloned

### Installation

This requires at least the 2.0.0+ version of the GitHub CLI to support extensions.

- Install the `gh` CLI. See [https://github.com/cli/cli#installation](https://github.com/cli/cli#installation) for further details.
- Install this extension:

```bash
gh extension install rm3l/gh-org-repo-sync
```

### Usage

```bash
gh org-repo-sync -h
Usage of ~/.local/share/gh/extensions/gh-org-repo-sync/gh-org-repo-sync:
  -batchSize int
        the number of elements to retrieve at once. Must not exceed 100 (default 50)
  -output string
        the output path (default ".")
  -protocol string
        the protocol to use for cloning. Possible values: system, ssh, https. (default "system")
```

### Working with the source code

Once you clone this repository, you can install a symbolic link to this in the `gh` CLI repositories cache.

```bash
cd gh-org-repo-sync

# Install the local version
gh extension install .

# At this point, you can start using it
gh org-repo-sync <my-organization>

# To see changes in the code as you develop, simply build and use the extension:
go build && gh org-repo-sync <my-organization>
```