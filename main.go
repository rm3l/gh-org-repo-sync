package main

import (
	"flag"
	"fmt"
	"github.com/rm3l/gh-org-repo-sync/internal/github"
	"github.com/rm3l/gh-org-repo-sync/internal/repo_sync"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultBatchSize = 50

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <organization> [-batchSize int] [protocol ssh|https|system] [-output /path]", os.Args[0])
		os.Exit(1)
	}

	var query string
	var batchSize int
	var output string
	var protocol string
	flag.StringVar(&query, "query", "",
		`GitHub search query, to filter the Organization repositories. Example: "language:Go stars:>10 pushed:>2010-11-12"
See https://bit.ly/3HurHe3 for more details on the search syntax
`)
	flag.IntVar(&batchSize, "batchSize", defaultBatchSize,
		"the number of elements to retrieve at once. Must not exceed 100")
	flag.StringVar(&protocol, "protocol", string(repo_sync.SystemProtocol),
		fmt.Sprintf("the protocol to use for cloning. Possible values: %s, %s, %s.", repo_sync.SystemProtocol,
			repo_sync.SSHProtocol, repo_sync.HTTPSProtocol))
	flag.StringVar(&output, "output", ".", "the output path")
	if os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help" {
		flag.Parse()
	} else {
		// Ignore errors since flag.CommandLine is set for ExitOnError.
		_ = flag.CommandLine.Parse(os.Args[2:])
	}

	organization := os.Args[1]

	if batchSize <= 0 || batchSize > 100 {
		fmt.Printf("invalid batch size (%d). Must be strictly higher than 0 and less than 100",
			batchSize)
		os.Exit(1)
	}
	cloneProtocol := repo_sync.CloneProtocol(strings.ToLower(protocol))

	log.Println("trying to list repos in the following organization:", organization)
	repositories, err := github.GetOrganizationRepos(organization, query, batchSize)
	if err != nil {
		log.Fatal(err)
	}
	nbRepos := len(repositories)
	log.Println("found", nbRepos, "repositories")
	if nbRepos == 0 {
		os.Exit(0)
	}

	var wg sync.WaitGroup
	wg.Add(nbRepos)
	for _, repository := range repositories {
		go func(repo string) {
			defer wg.Done()
			err := repo_sync.HandleRepository(output, organization, repo, cloneProtocol)
			if err != nil {
				log.Println("an error occurred while handling repo", repo, err)
			}
		}(repository)
	}
	wg.Wait()
	log.Println("done handling", nbRepos, "repositories in", time.Now().Sub(start))
}
