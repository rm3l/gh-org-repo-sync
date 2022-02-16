package cli

import (
	"bytes"
	"fmt"
	"github.com/cli/safeexec"
	"os"
	"os/exec"
)

// RunCommandInDir runs any command in the specified working directory
func RunCommandInDir(executable string, workingDir string, env []string, args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	executablePath, err := safeexec.LookPath(executable)
	if err != nil {
		err = fmt.Errorf("error while looking up the command specified: %s: %w", executablePath, err)
		return
	}
	cmd := exec.Command(executablePath, args...)
	cmd.Dir = workingDir
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if env != nil {
		cmd.Env = env
	}
	err = cmd.Run()
	fmt.Print(stdOut.String())
	_, _ = fmt.Fprint(os.Stderr, stdErr.String())
	if err != nil {
		err = fmt.Errorf("failed to run command: %s. error: %w", stdErr.String(), err)
		return
	}
	return
}
