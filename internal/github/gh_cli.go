package github

import (
	"bytes"
	"github.com/rm3l/gh-org-repo-sync/internal/cli"
)

// RunGhCliInDir runs any gh command in the specified working directory,
// because this is not possible to do with the default gh.Exec function
func RunGhCliInDir(workingDir string, env []string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	return cli.RunCommandInDir("gh", workingDir, env, args...)
}
