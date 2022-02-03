package main

import (
	"flag"
	"fmt"
	"github.com/rm3l/gh-org-clone/internal/github"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <organization> [-output /path]", os.Args[0])
		os.Exit(1)
	}
	organization := os.Args[1]

	var output string
	flag.StringVar(&output, "output", ".", "the output path. Defaults to the current dir")
	flag.Parse()

	currentUser, err := github.GetUser()
	if err != nil {
		log.Println("could not determine current user:", err)
	} else {
		log.Println("running extension as", currentUser)
	}

	log.Println("trying to list repos in the following organization:", organization)
	repositories, err := github.GetOrganizationRepos(organization)
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
