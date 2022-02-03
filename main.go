package main

import (
	"flag"
	"fmt"
	"github.com/rm3l/gh-org-clone/internal/github"
	"github.com/rm3l/gh-org-clone/internal/repo_sync"
	"log"
	"os"
	"strings"
	"sync"
)

const defaultBatchSize = 50

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <organization> [-batchSize 10] [protocol ssh|https] [-output /path]", os.Args[0])
		os.Exit(1)
	}
	organization := os.Args[1]

	var batchSize int
	var output string
	var protocol string
	flag.IntVar(&batchSize, "batchSize", defaultBatchSize,
		"the number of elements to retrieve at once. Must not exceed 100")
	flag.StringVar(&protocol, "protocol", string(repo_sync.DefaultProtocol), "the protocol to use for cloning")
	flag.StringVar(&output, "output", ".", "the output path")
	// Ignore errors; CommandLine is set for ExitOnError.
	_ = flag.CommandLine.Parse(os.Args[2:])

	if batchSize <= 0 || batchSize > 100 {
		fmt.Printf("invalid batch size (%d). Must be strictly higher than 0 and less than 100",
			batchSize)
		os.Exit(1)
	}
	cloneProtocol := repo_sync.CloneProtocol(strings.ToLower(protocol))

	currentUser, err := github.GetUser()
	if err != nil {
		log.Println("could not determine current user:", err)
	} else {
		log.Println("running extension as", currentUser)
	}

	log.Println("trying to list repos in the following organization:", organization)
	repositories, err := github.GetOrganizationRepos(organization, batchSize)
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
	log.Println("done handling", nbRepos, "repositories!")
}
