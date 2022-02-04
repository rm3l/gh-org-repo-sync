package repo_sync

import (
	"errors"
	"fmt"
	"github.com/cli/go-gh"
	"github.com/rm3l/gh-org-repo-sync/internal/github"
	"log"
	"os"
)

type CloneProtocol string

const (
	DefaultProtocol CloneProtocol = "system"
	SSHProtocol     CloneProtocol = "ssh"
	HTTPSProtocol   CloneProtocol = "https"
)

// HandleRepository determines whether a directory with the repository name does exist.
// If it does, it checks out its default branch and updates it locally.
// Otherwise, it clones it.
func HandleRepository(output, organization, repository string, protocol CloneProtocol) error {
	repoPath := fmt.Sprintf("%s/%s", output, repository)
	info, err := os.Stat(repoPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("cloning repo because local folder not found:", repoPath)
			if err := clone(output, organization, repository, protocol); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("expected folder for repository '%s'", repoPath)
	}
	log.Println("updating local clone for repo:", repoPath)
	return updateLocalClone(output, organization, repository)
}

func clone(output, organization, repository string, protocol CloneProtocol) error {
	repoPath := fmt.Sprintf("%s/%s", output, repository)
	var repoUrl string
	if protocol == DefaultProtocol {
		repoUrl = fmt.Sprintf("%s/%s", organization, repository)
	} else if protocol == SSHProtocol {
		repoUrl = fmt.Sprintf("git@github.com:%s/%s.git", organization, repository)
	} else if protocol == HTTPSProtocol {
		repoUrl = fmt.Sprintf("https://github.com/%s/%s.git", organization, repository)
	} else {
		return fmt.Errorf("unknown protocol for cloning: %s", protocol)
	}
	args := []string{"repo", "clone", repoUrl, repoPath}
	_, stdErr, err := gh.Exec(args...)
	if stdErrString := stdErr.String(); stdErrString != "" {
		fmt.Println(stdErrString)
	}
	return err
}

func updateLocalClone(output, organization, repository string) error {
	repoPath := fmt.Sprintf("%s/%s", output, repository)
	args := []string{"repo", "sync", "--source", fmt.Sprintf("%s/%s", organization, repository)}
	_, _, err := github.RunGhCliInDir(repoPath, nil, args...)
	return err
}
