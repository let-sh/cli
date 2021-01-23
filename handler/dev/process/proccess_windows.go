// +build windows

package process

import (
	"fmt"
	"github.com/let-sh/cli/log"
	"os/exec"
	"strconv"
	"strings"
)

func GetPortByProcessID(pid int) []int {
	cmd := exec.Command("sh", "-c", "lsof -nP -iTCP -sTCP:LISTEN | grep "+strconv.Itoa(pid))
	out, err := cmd.Output()
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("%s", out)

	var ports []int
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			break
		}
		spaces := strings.Fields(line)

		splited := strings.Split(spaces[8], ":")
		port, err := strconv.Atoi(splited[1])
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
