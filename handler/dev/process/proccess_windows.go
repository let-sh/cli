// +build windows

package process

import (
	"bytes"
	"fmt"
	"github.com/let-sh/cli/log"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func GetPortByProcessID(pid int) []int {
	first := exec.Command("netstat", "-a", "-n", "-o")
	second := exec.Command("find", "\""+strconv.Itoa(pid)+"\"")

	// http://golang.org/pkg/io/#Pipe

	reader, writer := io.Pipe()

	// push first command output to writer
	first.Stdout = writer

	// read from first command output
	second.Stdin = reader

	// prepare a buffer to capture the output
	// after second command finished executing
	var buff bytes.Buffer
	second.Stdout = &buff

	first.Start()
	second.Start()
	first.Wait()
	writer.Close()
	second.Wait()

	out := buff.String()
	fmt.Printf("%s", out)

	var ports []int
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			break
		}
		spaces := strings.Fields(line)

		split := strings.Split(spaces[2], ":")
		port, err := strconv.Atoi(split[1])
		if err != nil {
			log.Error(err)
			return ports
		}
		ports = append(ports, port)
	}
	return ports
}

func Kill(pid int) {
	p := process.Process{Pid: int32(pid)}
	p.Kill()
}
