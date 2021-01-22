/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/let-sh/cli/handler/dev"
	"github.com/let-sh/cli/handler/dev/process"
	"github.com/let-sh/cli/log"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-ps"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development environment",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupCloseDevelopmentHandler()
		var endpoint string
		var ports []int
		if len(inputCommand) > 0 {
			// run server command
			cmdSlice := strings.Split(inputCommand, " ")
			currentCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
			go func() {
				l := logrus.New()
				customFormatter := new(logrus.TextFormatter)
				//customFormatter.
				customFormatter.DisableTimestamp = true
				customFormatter.EnvironmentOverrideColors = true
				customFormatter.ForceColors = false
				customFormatter.FullTimestamp = true
				customFormatter.DisableColors = false
				l.Formatter = customFormatter

				//stdout, err := currentCmd.StdoutPipe()
				//if err != nil {
				//	l.Error(err)
				//	return
				//}

				currentCmd.Stdin = os.Stdin
				currentCmd.Stdout = os.Stdout
				currentCmd.Stderr = os.Stderr
				// start the command after having set up the pipe
				if err := currentCmd.Start(); err != nil {
					l.Error(err)
					return
				}
				//
				//// read command's stdout line by line
				//in := bufio.NewScanner(stdout)
				//
				//for in.Scan() {
				//	l.Info(string(in.Bytes())) // write each line to your log, or anything you need
				//}
				//if err := in.Err(); err != nil {
				//	l.Errorf("error: %s", err)
				//}
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
				ports = process.GetPortByProcessID(currentCmd.Process.Pid)

				// get port by process child pid
				processes, err := ps.Processes()
				if err != nil {
					log.Error(err)
					return
				}
				for _, p := range processes {
					if p.PPid() == currentCmd.Process.Pid {
						ports = append(ports, process.GetPortByProcessID(p.Pid())...)
					}
				}

				time.Sleep(time.Second * 2)

				if len(ports) > 0 {
					break
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
					Label: "Please select a port to listen: ",
					Items: ports,
				}

				_, result, err := prompt.Run()
				if err != nil {
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

		log.Success("you can visit remotely at:\nhttp://" + inputRemoteEndpoint)
		dev.StartClient(inputRemoteEndpoint, endpoint)
	},
}

var inputRemoteEndpoint string
var inputLocalEndpoint string
var inputCommand string

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
		//if len(DeploymentID) > 0 {
		//	succeed, err := requests.CancelDeployment(DeploymentID)
		//	if err != nil {
		//		log.S.StopFail()
		//		log.Error(err)
		//		os.Exit(0)
		//	}
		//	if succeed {
		//		log.S.StopFail()
		//		log.Warning("Deployment canceled")
		//		os.Exit(0)
		//	} else {
		//		log.S.StopFail()
		//		log.Warning("Deployment cancellation failed")
		//		os.Exit(0)
		//	}
		//} else {
		//	log.S.StopFail()
		//	log.Warning("Deployment canceled")
		//	os.Exit(0)
		//}
		os.Exit(0)
	}()
}
