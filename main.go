package main

import (
	"flag"
	"fmt"
	"github.com/rm3l/gh-org-clone/internal/github"
	"log"
	"os"
)

const defaultBatchSize = 50

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <organization> [-batchSize 10] [-output /path]", os.Args[0])
		os.Exit(1)
	}
	organization := os.Args[1]

	var batchSize int
	var output string
	flag.IntVar(&batchSize, "batchSize", defaultBatchSize,
		"the number of elements to retrieve at once. Must not exceed 100")
	flag.StringVar(&output, "output", ".", "the output path")
	// Ignore errors; CommandLine is set for ExitOnError.
	_ = flag.CommandLine.Parse(os.Args[2:])

	if batchSize <= 0 || batchSize > 100 {
		fmt.Printf("invalid batch size (%d). Must be strictly higher than 0 and less than 100",
			batchSize)
		os.Exit(1)
	}

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
	log.Println("found", len(repositories), "repositories")
	if len(repositories) == 0 {
		os.Exit(0)
	}

	//TODO Now for each determine whether the directory in the output path exists.
	// if it does, checkout its main/master branch and update it.
	// otherwise, clone it
}
