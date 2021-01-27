// +build !windows

package command

import (
	"github.com/creack/pty"
	"github.com/let-sh/cli/log"
	"github.com/logrusorgru/aurora"
	"github.com/segmentio/textio"
	"io"
	"os"
	"os/exec"
)

func RunCmd(cmd *exec.Cmd) {

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Error(err)
		return
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGWINCH)
	//go func() {
	//	for range ch {
	//		size, _ := pty.GetsizeFull(ptmx)
	//		pty.Setsize(ptmx, size)
	//
	//		//if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
	//		//	log.Errorf("error resizing pty: %s", err.Error())
	//		//}
	//	}
	//}()
	//ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	//oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()
	copyIndent(os.Stdout, ptmx)
}

func copyIndent(w io.Writer, r io.Reader) error {
	p := textio.NewPrefixWriter(w, aurora.Gray(5, "[log] ").String())

	// Copy data from an input stream into the PrefixWriter, all lines will
	// be prefixed with a '\t' character.
	if _, err := io.Copy(p, r); err != nil {
		return err
	}

	// Flushes any data buffered in the PrefixWriter, this is important in
	// case the last line was not terminated by a '\n' character.
	return p.Flush()
}
