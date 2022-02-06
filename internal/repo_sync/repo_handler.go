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
func HandleRepository(dryRun bool, output, organization string, repositoryInfo github.RepositoryInfo, protocol CloneProtocol) error {
	repository := repositoryInfo.Name
	repoPath, err := safeAbsPath(fmt.Sprintf("%s/%s", output, repository))
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
			if err := clone(repoPath, organization, repository, protocol); err != nil {
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
	return updateLocalClone(repoPath, organization, repositoryInfo)
}

func clone(output, organization string, repository string, protocol CloneProtocol) error {
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
	repoPath, err := safeAbsPath(output)
	if err != nil {
		return err
	}
	args := []string{"repo", "clone", repoUrl, repoPath}
	_, stdErr, err := gh.Exec(args...)
	if stdErrString := stdErr.String(); stdErrString != "" {
		fmt.Println(stdErrString)
	}
	return err
}

func updateLocalClone(outputPath, organization string, repositoryInfo github.RepositoryInfo) error {
	repository := repositoryInfo.Name
	repoPath, err := safeAbsPath(outputPath)
	if err != nil {
		return err
	}
	err = fetchAllRemotes(repoPath)
	if err != nil {
		log.Println("[warn]", err)
	}
	if repositoryInfo.IsEmpty {
		log.Printf("[warn] skipped syncing empty repo: %s. Only remotes have been fetched\n", repoPath)
		return nil
	}
	args := []string{"repo", "sync", "--source", fmt.Sprintf("%s/%s", organization, repository)}
	_, _, err = github.RunGhCliInDir(repoPath, nil, args...)
	return err
}

func fetchAllRemotes(outputPath string) error {
	repoPath, err := safeAbsPath(outputPath)
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	remotes, err := r.Remotes()
	if err != nil {
		return err
	}
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
	return nil
}

func safeAbsPath(p string) (string, error) {
	return filepath.Abs(filepath.FromSlash(p))
}
