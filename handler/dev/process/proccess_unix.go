// +build !windows

package process

import (
	"github.com/let-sh/cli/log"
	"github.com/shirou/gopsutil/v3/process"
	"os/exec"
	"strconv"
	"strings"
)

func GetPortByProcessID(pid int) []int {
	cmd := exec.Command("sh", "-c", "lsof -nP -iTCP -sTCP:LISTEN | grep "+strconv.Itoa(pid))
	out, _ := cmd.Output()

	var ports []int
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			break
		}
		spaces := strings.Fields(line)

		split := strings.Split(spaces[8], ":")
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
