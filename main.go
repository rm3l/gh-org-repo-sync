package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rm3l/gh-org-repo-sync/internal/github"
	"github.com/rm3l/gh-org-repo-sync/internal/reposync"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"strings"
	"time"
)

const defaultBatchSize = 50

var nbRepos int

func main() {
	start := time.Now()
	defer func() {
		log.Println("[info] done handling", nbRepos, "repositories in", time.Since(start))
	}()

	var dryRun bool
	var query string
	var batchSize int
	var output string
	var protocol string
	var force bool
	flag.BoolVar(&dryRun, "dry-run", false,
		`dry run mode. to display the repos that will get cloned or updated, 
without actually performing those actions`)
	flag.StringVar(&query, "query", "",
		`GitHub search query, to filter the Organization repositories.
Example: "language:Go stars:>10 pushed:>2010-11-12"
See https://bit.ly/3HurHe3 for more details on the search syntax`)
	flag.IntVar(&batchSize, "batch-size", defaultBatchSize,
		"the number of elements to retrieve at once. Must not exceed 100")
	flag.StringVar(&protocol, "protocol", string(reposync.SystemProtocol),
		fmt.Sprintf("the protocol to use for cloning. Possible values: %s, %s, %s.", reposync.SystemProtocol,
			reposync.SSHProtocol, reposync.HTTPSProtocol))
	flag.StringVar(&output, "output", ".", "the output path")
	flag.BoolVar(&force, "force", false,
		`whether to force sync repositories.
Caution: this will hard-reset the branch of the destination repository to match the source repository.`)

	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "Usage: gh org-repo-sync <organization> [options]")
		fmt.Println("Options: ")
		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "missing organization")
		flag.Usage()
		os.Exit(1)
	}

	organization := os.Args[1]

	if organization == "-h" || organization == "-help" || organization == "--help" {
		flag.Usage()
		os.Exit(1)
	} else {
		// Ignore errors since flag.CommandLine is set for ExitOnError.
		_ = flag.CommandLine.Parse(os.Args[2:])
	}

	if batchSize <= 0 || batchSize > 100 {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, "invalid batch size (%d). Must be strictly higher than 0 and less than 100",
			batchSize)
		os.Exit(1)
	}
	cloneProtocol := reposync.CloneProtocol(strings.ToLower(protocol))

	repositories, err := github.GetOrganizationRepos(organization, query, batchSize)
	if err != nil {
		log.Fatal(err)
	}
	nbRepos = len(repositories)
	log.Println("[debug] found", nbRepos, "repositories")
	if nbRepos == 0 {
		return
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, repository := range repositories {
		g.Go(func(repo github.RepositoryInfo) func() error {
			return func() error {
				err := reposync.HandleRepository(ctx, dryRun, output, organization, repo, cloneProtocol, force)
				if err != nil {
					return fmt.Errorf("error while handling repo %q: %w", repo.Name, err)
				}
				return nil
			}
		}(repository))
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
