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
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mdp/qrterminal/v3"
	"github.com/muesli/termenv"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils/config"
	"github.com/matishsiao/goInfo"
	"github.com/spf13/cobra"
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
		// select login methods
		prompt := promptui.Select{
			Label: "Login Method",
			Items: []string{"GitHub", "WeChat"},
		}

		_, loginMethod, err := prompt.Run()

		if err != nil {
			if err == promptui.ErrInterrupt {
				os.Exit(-1)
				return
			}

			loginMethod = "GitHub"
		}

		log.BStart("loading...")

		// get ticket id
		tickeIDInterface, err := requests.GetJsonWithPath("https://api.let-sh.com/oauth/ticket_id", "data")
		if err != nil {
			log.Error(err)
			return
		}
		log.S.StopFail()

		switch loginMethod {
		case "GitHub":
			// open browser to login
			err = openBrowser("https://api.let-sh.com/oauth/login?method=github&client=cli&ticket_id=" + tickeIDInterface.
				String() + "&device=" + goInfo.GetInfo().OS + goInfo.GetInfo().Core)

			shortenedUrl, _ := requests.GenerateShortUrl("https://api.let-sh." +
				"com/oauth/login?method=github&client=cli&ticket_id=" + tickeIDInterface.
				String() + "&device=" + goInfo.GetInfo().OS + goInfo.GetInfo().Core)

			if shortenedUrl != "" {
				log.S.StopFail()
				fmt.Println(
					termenv.
						String("if your browser not opened automatically, please visit: ").
						Foreground(termenv.ColorProfile().Color("#808080")),

					termenv.
						String(shortenedUrl).
						Foreground(termenv.ColorProfile().Color("#808080")).Underline(),
				)
				log.BStart("redirecting to browser")
			}
		case "WeChat":

			shortenedUrl, _ := requests.GenerateShortUrl("https://api.let-sh.com" +
				"/oauth/login?method=wechat&client=cli&ticket_id=" + tickeIDInterface.
				String() + "&device=" + goInfo.GetInfo().OS + goInfo.GetInfo().Core)

			config := qrterminal.Config{
				Level:     qrterminal.L,
				Writer:    os.Stdout,
				BlackChar: qrterminal.BLACK,
				WhiteChar: qrterminal.WHITE,
				QuietZone: 0,
			}
			qrterminal.GenerateWithConfig(shortenedUrl, config)

			fmt.Println("\nplease use WeChat to scan the QR code above.\n")
		}

		// valid response
		start := time.Now()
		log.BUpdate("waiting for login result")
		//log.BUpdate("waiting for login result, you could also manually visit: " + "https://api.let-sh.com/oauth/login?method=github&client=cli&ticket_id=" + tickeIDInterface.String() + "&device=" + goInfo.GetInfo().OS + goInfo.GetInfo().Core)

		for {
			// Code to measure
			duration := time.Since(start)
			if duration >= time.Minute*1 {
				break
			}
			tokenInterface, err := requests.GetJsonWithPath("https://api.let-sh.com/oauth/ticket/"+tickeIDInterface.String(), "data.token")
			if err == nil && tokenInterface.String() != "" {
				// verify response
				config.SetToken(tokenInterface.String())
				log.S.StopMessage(" login succeed")
				log.BStop()
				return
			}
			time.Sleep(time.Second * 1)
		}
		log.S.StopFailMessage(" login timeout")
		log.S.StopFail()
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

func openBrowser(url string) (err error) {

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
		//log.Error(err)
		return err
	}
	return nil
}
