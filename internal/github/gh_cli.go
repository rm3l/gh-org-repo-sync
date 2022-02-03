package github

import (
	"bytes"
	"fmt"
	"github.com/cli/safeexec"
	"os/exec"
)

func RunGhCliInDir(workingDir string, env []string, args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	ghPath, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("error while looking up the gh command: %w", err)
		return
	}
	cmd := exec.Command(ghPath, args...)
	cmd.Dir = workingDir //Not possible with the default gh.Exec function
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if env != nil {
		cmd.Env = env
	}
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh: %s. error: %w", stdErr.String(), err)
		return
	}
	return
}
