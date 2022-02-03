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
	"fmt"
	"github.com/let-sh/cli/handler/deploy"
	"github.com/let-sh/cli/ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/let-sh/cli/handler/dev"
	c "github.com/let-sh/cli/handler/dev/command"
	"github.com/let-sh/cli/handler/dev/process"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils"
	"github.com/let-sh/cli/utils/cache"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
)

var inputRemoteEndpoint string
var inputLocalEndpoint string
var inputCommand string
var processPids []int
var forceLocal bool

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development environment",
	Long:  `Start development environment, let.sh cli will automatically export your service with development endpoint`,
	Run: func(cmd *cobra.Command, args []string) {

		var command string
		var localEndpoint string
		var ports []int

		var deploymentCtx deploy.DeployContext
		detectedType := deploymentCtx.DetectProjectType()
		logrus.Debug("detected project type: ", detectedType)

		deploymentCtx.LoadLetJson()

		dir, _ := os.Getwd()
		p, err := cache.GetProjectInfo(dir)

		// if cache exists
		if err == nil {
			command = p.ServeCommand
		} else {
			p.Name = filepath.Base(dir)
			p.Dir = dir
		}

		// if user not specified command
		// TODO: add default command placeholders for different project types
		if detectedType == "static" {
			// TODO: trigger to reverse port
			freePort, err := GetFreePort()
			if err != nil {
				log.Error(err)
				return
			}
			inputLocalEndpoint = "localhost:" + cast.ToString(freePort)
			go ServeStaticFiles(deploymentCtx.Static, freePort)
		} else {
			if len(inputCommand) == 0 {
				if len(command) == 0 {
					//detectedType

					defaultCommand := func() string {
						switch detectedType {
						case "gin":
							return "go run main.go"
						case "react":
							return "yarn dev"
						case "vue":
							return "yarn dev"
						default:
							return ""
						}
					}()
					result, err := ui.InputArea(ui.InputAreaConfig{
						Layout:             "Please enter your command to start service: ",
						DefaultPlaceholder: defaultCommand,
						PlaceHolders:       []string{defaultCommand},
					})
					if err != nil {
						log.Error(err)
						return
					}
					command = result
				}
			} else {
				command = inputCommand
			}
		}

		p.ServeCommand = command
		cache.SaveProjectInfo(p)
		SetupCloseDevelopmentHandler(p.ID)

		defer KillServiceProcess(p.ID)

		// request to start tunnel
		remoteEndpoint := inputRemoteEndpoint
		var result struct {
			RemotePort    int    `json:"remotePort,omitempty"`
			RemoteAddress string `json:"remoteAddress,omitempty"`
			Fqdn          string `json:"fqdn,omitempty"`
		}
		if !forceLocal {
			result, err = requests.StartDevelopment(p.ID)
			if err != nil {
				log.Error(err)
				return
			}

			// using wss://
			if result.RemotePort == 443 {
				remoteEndpoint = "wss://" + result.RemoteAddress + ":" + strconv.Itoa(result.RemotePort)
			} else {
				remoteEndpoint = "ws://" + result.RemoteAddress + ":" + strconv.Itoa(result.RemotePort)
			}

		}

		if detectedType != "static" {
			ui.Spinner.Start("starting development service")

			// run server command
			cmdSlice := strings.Split(command, " ")
			currentCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

			go c.RunCmd(currentCmd)

			for {
				// wait for process to start
				if currentCmd.Process != nil {
					break
				}
			}
			ui.Spinner.Stop()

			// awaiting port binding
			for i := 0; i < 15; i++ {
				// get local port by pid
				processPids = append(processPids, currentCmd.Process.Pid)

				FindAllChildrenProcess(currentCmd.Process.Pid)
				for _, p := range processPids {
					ports = utils.RemoveDuplicates(append(ports, process.GetPortByProcessID(p)...))
				}

				time.Sleep(time.Second)

				if len(ports) > 0 {
					break
				}

				if _, err := ps.FindProcess(currentCmd.Process.Pid); err != nil {
					log.Error(errors.New("service process existed, please check logs above"))
					KillServiceProcess(p.ID)
					return
				}

				if i == 9 {
					log.Error(errors.New("timeout waiting for service port, please check your service status"))
					KillServiceProcess(p.ID)
					return
				}
			}
		}

		// if no input local port
		// detect service port
		if len(inputLocalEndpoint) == 0 {
			if len(ports) == 0 {
				log.Warning("please specify a local endpoint")
			}

			if len(ports) == 1 {
				localEndpoint = "localhost:" + strconv.Itoa(ports[0])
			}

			if len(ports) > 1 {
				// if current dir is not previous dir
				prompt := promptui.Select{
					Label: "Please select a port to listen",
					Items: ports,
				}

				_, result, err := prompt.Run()
				if err != nil {
					KillServiceProcess(p.ID)
					log.Error(err)
					return
				}
				localEndpoint = "localhost:" + result
			}
		} else {
			localEndpoint = inputLocalEndpoint
		}

		// if remote or local endpoint not exists
		if remoteEndpoint == "" || localEndpoint == "" {
			log.Error(errors.New("currently under development"))
		}

		if len(result.Fqdn) == 0 {
			log.Error(errors.New("missing public visit fqdn"))
		}

		fmt.Println("\n"+aurora.BrightCyan("[msg]").Bold().String(), "you can visit remotely at: "+aurora.Bold("https://"+result.Fqdn).String()+"\n\r")
		openBrowser("https://" + result.Fqdn)
		dev.StartClient(remoteEndpoint, localEndpoint, result.Fqdn)
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
	devCmd.Flags().StringVarP(&inputCommand, "command", "c", "", "command to serve service, e.g. 'yarn start', 'go run main.go'")

	devCmd.Flags().StringVarP(&inputRemoteEndpoint, "remote", "r", "", "custom remote endpoint, e.g. remote.example.com:3000")
	deployCmd.Flags().MarkHidden("remote")

	devCmd.Flags().StringVarP(&inputLocalEndpoint, "local", "l", "", "custom local upstream endpoint, e.g. 127.0.0.1:3000")

	devCmd.Flags().BoolVarP(&forceLocal, "force", "f", false, "force local test development")
	deployCmd.Flags().MarkHidden("force")
}

func SetupCloseDevelopmentHandler(projectID string) {
	// TODO: trigger stop tunnel
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ui.Spinner.Start("exiting development mode")
		KillServiceProcess(projectID)

		ui.Spinner.Stop()
		log.Warning("exited development")
		os.Exit(0)
	}()
}

func KillServiceProcess(projectID string) {
	// kill process

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		requests.StopDevelopment(projectID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, pid := range processPids {
			process.Kill(pid)
		}
	}()
	wg.Wait()
}

func FindAllChildrenProcess(pid int) (exists bool, childrenPid int) {
	processes, err := ps.Processes()
	if err != nil {
		log.Error(err)
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

func ServeStaticFiles(dir string, port int) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	err := http.ListenAndServe(":"+cast.ToString(port), nil)
	if err != nil {
		log.Error(err)
		return
	}
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
