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
	"os"

	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/ui"
	"github.com/let-sh/cli/utils/config"
	"github.com/let-sh/cli/utils/update"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var Debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lets",
	Short: aurora.Index(ui.MainColor, "λ").String() + "️ Launch your app with just one command",
	Long:  aurora.Index(ui.MainColor, "λ").String() + " let.sh helps you test, preview and launch your app",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("exit with error: ", err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, update.CheckUpdate)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.SetVersionTemplate(info.Version)
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "debugging cli command")
	rootCmd.PersistentFlags().MarkHidden("debug")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	//rootCmd.PersistentFlags().String("token", "", "let.sh access token")
	rootCmd.PersistentFlags().StringVarP(&info.Credentials.Token, "token", "", "", "specify the let.sh access token, ")
	//
	//if token, err := rootCmd.PersistentFlags().GetString("token"); err != nil {
	//	if len(token) > 0 {
	//		info.Credentials.SetToken(token)
	//	}
	//}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.Load()

	if Debug || info.Version == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}
}
