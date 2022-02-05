package repo_sync

import (
	"errors"
	"fmt"
	"github.com/cli/go-gh"
	"github.com/go-git/go-git/v5"
	"github.com/rm3l/gh-org-repo-sync/internal/github"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type CloneProtocol string

const (
	SystemProtocol CloneProtocol = "system"
	SSHProtocol    CloneProtocol = "ssh"
	HTTPSProtocol  CloneProtocol = "https"
)

// HandleRepository determines whether a directory with the repository name does exist.
// If it does, it checks out its default branch and updates it locally.
// Otherwise, it clones it.
func HandleRepository(dryRun bool, output, organization, repository string, protocol CloneProtocol) error {
	repoPath, err := filepath.Abs(filepath.FromSlash(fmt.Sprintf("%s/%s", output, repository)))
	if err != nil {
		return err
	}
	info, err := os.Stat(repoPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if dryRun {
				fmt.Printf("=> %s/%s: new clone in '%s'\n", organization, repository, repoPath)
				return nil
			}
			log.Println("[debug] cloning repo because local folder not found:", repoPath)
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
	if dryRun {
		fmt.Printf("=> %s/%s: update in '%s'\n", organization, repository, repoPath)
		return nil
	}
	log.Println("[debug] updating local clone for repo:", repoPath)
	return updateLocalClone(output, organization, repository)
}

func clone(output, organization, repository string, protocol CloneProtocol) error {
	repoPath := fmt.Sprintf("%s/%s", output, repository)
	var repoUrl string
	if protocol == SystemProtocol {
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
	err := fetchAllRemotes(repoPath)
	if err != nil {
		log.Println("[warn]", err)
	}
	args := []string{"repo", "sync", "--source", fmt.Sprintf("%s/%s", organization, repository)}
	_, _, err = github.RunGhCliInDir(repoPath, nil, args...)
	return err
}

func fetchAllRemotes(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	var errorReturned error
	if err != nil {
		errorReturned = err
	} else {
		remotes, err := r.Remotes()
		if err != nil {
			errorReturned = err
		} else {
			var wg sync.WaitGroup
			wg.Add(len(remotes))
			for _, remote := range remotes {
				go func(rem *git.Remote) {
					defer wg.Done()
					log.Printf("[debug] fetching remote '%s' in %s", rem.Config().Name, repoPath)
					_ = rem.Fetch(&git.FetchOptions{Depth: 0, Tags: git.AllTags})
				}(remote)
			}
			wg.Wait()
		}
	}
	return errorReturned
}
