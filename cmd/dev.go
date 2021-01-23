/*
Copyright © 2021 Fred Liang <fred@oasis.ac>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"github.com/fatih/color"
	"github.com/let-sh/cli/handler/dev"
	"github.com/let-sh/cli/handler/dev/process"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var inputRemoteEndpoint string
var inputLocalEndpoint string
var inputCommand string
var processPids []int

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development environment",
	Long:  `Start development environment, let.sh cli will automatically export your service with development endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupCloseDevelopmentHandler()
		var command string
		var endpoint string
		var ports []int

		//detectedType :=deploy.DetectProjectType()

		if len(inputCommand) == 0 {
			// if current dir is not previous dir
			prompt := promptui.Prompt{
				Label: "Please enter your command to start service",
			}
			result, err := prompt.Run()
			if err != nil {
				log.Error(err)
				return
			}
			command = result
		} else {
			command = inputCommand
		}

		{
			// run server command
			cmdSlice := strings.Split(command, " ")
			currentCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
			go func() {
				// start the command after having set up the pipe
				currentCmd.Stdin = os.Stdin
				currentCmd.Stdout = os.Stdout
				currentCmd.Stderr = os.Stderr
				if err := currentCmd.Start(); err != nil {
					KillServiceProcess()
					log.Error(err)
					return
				}
			}()

			for {
				// wait for process to start
				if currentCmd.Process != nil {
					break
				}
			}

			log.BStart("let.sh: awaiting service local port binding")
			// awaiting port binding
			for i := 0; i < 10; i++ {
				// get local port by pid
				processPids = append(processPids, currentCmd.Process.Pid)

				FindAllChildrenProcess(currentCmd.Process.Pid)
				for _, p := range processPids {
					ports = utils.RemoveDuplicates(append(ports, process.GetPortByProcessID(p)...))
				}

				time.Sleep(time.Second * 2)

				if len(ports) > 0 {
					break
				}

				if _, err := ps.FindProcess(currentCmd.Process.Pid); err != nil {
					log.Error(errors.New("service process existed, please check logs above"))
					KillServiceProcess()
					return
				}

				if i == 9 {
					log.Error(errors.New("timeout waiting for service port, please check your service status"))
					KillServiceProcess()
					return
				}
			}
		}
		log.S.StopFail()

		if len(inputLocalEndpoint) == 0 {
			if len(ports) == 0 {
				log.Warning("please specify a local endpoint")
			}

			if len(ports) == 1 {
				endpoint = "localhost:" + strconv.Itoa(ports[0])
			}

			if len(ports) > 1 {
				// if current dir is not previous dir
				prompt := promptui.Select{
					Label: "Please select a port to listen",
					Items: ports,
				}

				_, result, err := prompt.Run()
				if err != nil {
					KillServiceProcess()
					log.Error(err)
					return
				}
				endpoint = "localhost:" + result
			}
		} else {
			endpoint = inputLocalEndpoint
		}

		if inputRemoteEndpoint == "" || endpoint == "" {
			log.Error(errors.New("currently under development"))
		}

		log.Success("you can visit remotely at: " + color.New(color.Bold).Sprint("http://"+inputRemoteEndpoint))
		dev.StartClient(inputRemoteEndpoint, endpoint)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	devCmd.Flags().StringVarP(&inputCommand, "command", "c", "", "command to serve service, e.g. `yarn start`, `go run main.go`")
	devCmd.Flags().StringVarP(&inputRemoteEndpoint, "remote", "r", "", "custom remote endpoint, e.g. 127.0.0.1")
	devCmd.Flags().StringVarP(&inputLocalEndpoint, "local", "l", "", "custom local upstream endpoint, e.g. 127.0.0.1")
}

func SetupCloseDevelopmentHandler() {
	// TODO: trigger stop tunnel
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		KillServiceProcess()
		log.Warning("exited development")
		os.Exit(0)
	}()
}

func KillServiceProcess() {
	for _, pid := range processPids {
		process.Kill(pid)
	}
}

func FindAllChildrenProcess(pid int) (exists bool, childrenPid int) {
	processes, err := ps.Processes()
	if err != nil {
		log.Error(err)
		KillServiceProcess()
		return
	}
	for _, p := range processes {
		if p.PPid() == pid {
			processPids = append(processPids, p.Pid())
			return FindAllChildrenProcess(p.Pid())
		}
	}
	return false, 0
}
