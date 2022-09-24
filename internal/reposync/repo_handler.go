package reposync

import (
	"context"
	"errors"
	"fmt"
	"github.com/cli/go-gh"
	"github.com/rm3l/gh-org-repo-sync/internal/cli"
	"github.com/rm3l/gh-org-repo-sync/internal/github"
	"log"
	"os"
	"path/filepath"
)

// CloneProtocol indicates the Git protocol to use for cloning.
// See the constants exported in this package for further details.
type CloneProtocol string

const (
	// SystemProtocol indicates whether to use the Git protocol configured in the GitHub CLI,
	// e.g., via the 'gh config set git_protocol' configuration command
	SystemProtocol CloneProtocol = "system"

	// SSHProtocol forces this extension to clone repositories via SSH.
	// As such, the Git remote will look like: git@github.com:org/repo.git
	SSHProtocol CloneProtocol = "ssh"

	// HTTPSProtocol  forces this extension to clone repositories via HTTPS.
	// As such, the Git remote will look like: https://github.com/org/repo.git
	HTTPSProtocol CloneProtocol = "https"
)

// HandleRepository determines whether a directory with the repository name does exist.
// If it does, it checks out its default branch and updates it locally.
// Otherwise, it clones it.
func HandleRepository(
	_ context.Context,
	dryRun bool,
	output,
	organization string,
	repositoryInfo github.RepositoryInfo,
	protocol CloneProtocol,
	force bool,
) error {
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
	return updateLocalClone(repoPath, organization, repositoryInfo, force)
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

func updateLocalClone(outputPath, organization string, repositoryInfo github.RepositoryInfo, force bool) error {
	repository := repositoryInfo.Name
	repoPath, err := safeAbsPath(outputPath)
	if err != nil {
		return err
	}
	err = fetchAllRemotes(repoPath, force)
	if err != nil {
		return err
	}
	if repositoryInfo.IsEmpty {
		log.Printf("[warn] skipped syncing empty repo: %s. Only remotes have been fetched\n", repoPath)
		return nil
	}
	args := []string{"repo", "sync", "--source", fmt.Sprintf("%s/%s", organization, repository)}
	if force {
		args = append(args, "--force")
	}
	_, _, err = github.RunGhCliInDir(repoPath, nil, args...)
	return err
}

func fetchAllRemotes(outputPath string, force bool) error {
	repoPath, err := safeAbsPath(outputPath)
	if err != nil {
		return err
	}
	args := []string{"fetch", "--all", "--prune", "--tags", "--recurse-submodules"}
	if force {
		args = append(args, "--force")
	}
	_, _, err = cli.RunCommandInDir("git", repoPath, nil, args...)
	return err
}

func safeAbsPath(p string) (string, error) {
	return filepath.Abs(filepath.FromSlash(p))
}
