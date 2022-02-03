package repo_sync

import (
	"log"
)

// HandleRepository determines whether a directory with the repository name does exist.
// If it does, it checks out its main/master branch and updates it.
// Otherwise, it clones it.
func HandleRepository(organization, repository, output string) error {
	log.Println("handling repo:", repository)
	//TODO
	return nil
}
