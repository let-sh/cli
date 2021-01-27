// +build windows

package command

import (
	"github.com/let-sh/cli/log"
	"os"
	"os/exec"
)

func RunCmd(cmd *exec.Cmd) {
	// start the command after having set up the pipe
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return
	}
}
