# gh-org-repo-sync

[![Build](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/build.yml)
[![Lint](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/lint.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/lint.yml)
[![CodeQL Analysis](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/rm3l/gh-org-repo-sync/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rm3l/gh-org-repo-sync)](https://goreportcard.com/report/github.com/rm3l/gh-org-repo-sync)

> GitHub CLI extension to clone all repositories in an Organization, with the ability to filter via search queries.
> If a local clone already exists, it fetches all remotes and pulls changes from the default branch.

## Installation

- Install the `gh` CLI. See [https://github.com/cli/cli#installation](https://github.com/cli/cli#installation) for further details.
- If not done already, also install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git).
- To be able to clone private repos you have access to inside the Organization, authenticate with your GitHub account by running `gh auth login`. Alternatively, the CLI will respect the `GITHUB_TOKEN` [environment variable](https://cli.github.com/manual/gh_help_environment).
- Install the extension:

```bash
gh extension install rm3l/gh-org-repo-sync
```

## Usage

```bash
❯ gh org-repo-sync -h

Usage: gh org-repo-sync <organization> [options]
Options: 
  -batch-size int
        the number of elements to retrieve at once. Must not exceed 100 (default 50)
  -dry-run
        dry run mode. to display the repositories that will get cloned or updated, 
        without actually performing those actions
  -force
        whether to force sync repositories.
        Caution: this will hard-reset the branch of the destination repository to match the source repository.
  -output string
        the output path (default ".")
  -protocol string
        the protocol to use for cloning. Possible values: system, ssh, https. (default "system")
  -query string
        GitHub search query, to filter the Organization repositories.
        Example: "language:Java stars:>10 pushed:>2010-11-12"
        See https://bit.ly/3HurHe3 for more details on the search syntax
```

## Working with the source code

1. Clone the repository:

```
git clone https://github.com/rm3l/gh-org-repo-sync
cd gh-org-repo-sync
```

2. Download and install [Go](https://go.dev/doc/install) to build the project.
   Or if you are already using the [asdf](https://asdf-vm.com/) version manager, you can just run `asdf install` to install all the necessary tools (declared in the [.tool-versions](.tool-versions) file).

3. Install the local version of this extension; `gh` symlinks to your local source code directory.

```bash
# Install the local version
gh extension install .
```

4. At this point, you can start using it:

```bash
gh org-repo-sync <my-organization>
```

5. To see changes in the code as you develop, simply build and use the extension

```bash
go build && gh org-repo-sync <my-organization>
```

## Contribution Guidelines

Contributions and issue reporting are more than welcome. So to help out, do feel free to fork this repo and open up a pull request.
I'll review and merge your changes as quickly as possible.

You can use [GitHub issues](https://github.com/rm3l/gh-org-repo-sync/issues) to report bugs.
However, please make sure your description is clear enough and has sufficient instructions to be able to reproduce the issue.

## Developed by

* Armel Soro
    * [keybase.io/rm3l](https://keybase.io/rm3l)
    * [rm3l.org](https://rm3l.org) - &lt;armel+gh-org-repo-sync@rm3l.org&gt; - [@rm3l](https://twitter.com/rm3l)
    * [paypal.me/rm3l](https://paypal.me/rm3l)
    * [coinbase.com/rm3l](https://www.coinbase.com/rm3l)

## License

    The MIT License (MIT)

    Copyright (c) 2022 Armel Soro

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
