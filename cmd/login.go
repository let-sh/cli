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
	"fmt"
	"github.com/let-sh/let.cli/log"
	"github.com/let-sh/let.cli/requests"
	"github.com/let-sh/let.cli/utils"
	"github.com/matishsiao/goInfo"
	"github.com/spf13/cobra"
	"os/exec"
	"runtime"
	"time"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to let.sh",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// loading
		log.BStart("redirecting to browser")

		// get ticket id
		tickeIDInterface, err := requests.GetJsonWithPath("https://api.let.sh/oauth/ticket_id", "data")
		if err != nil {
			log.Error(err)
			return
		}

		// open browser to login
		openBrowser("https://api.let.sh/oauth/login?method=github&client=cli&ticket_id=" + tickeIDInterface.String() + "&device=" + goInfo.GetInfo().OS + goInfo.GetInfo().Core)

		// valid response
		start := time.Now()
		log.BUpdate("waiting for login result")
		for {
			// Code to measure
			duration := time.Since(start)
			if duration >= time.Minute*1 {
				break
			}
			tokenInterface, err := requests.GetJsonWithPath("https://api.let.sh/oauth/ticket/"+tickeIDInterface.String(), "data.token")
			if err == nil && tokenInterface.String() != "" {
				// verify response
				utils.SetToken(tokenInterface.String())
				log.BStop()
				log.Success("login to let.sh succeed")
				return
			}

			time.Sleep(time.Second * 1)
		}

		log.BStop()
		log.Error(errors.New("login failed"))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}
