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
	"github.com/let-sh/cli/log"
	"github.com/spf13/cobra"
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
		if inputRemoteEndpoint == "" || inputLocalEndpoint == "" {
			log.Error(errors.New("currently under development"))
		}
		log.Success("you can visit http://" + inputRemoteEndpoint)
		dev.StartClient(inputRemoteEndpoint, inputLocalEndpoint)
	},
}

var inputRemoteEndpoint string
var inputLocalEndpoint string

func init() {
	rootCmd.AddCommand(devCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	devCmd.Flags().StringVarP(&inputRemoteEndpoint, "remote", "r", "", "custom remote endpoint, e.g. 127.0.0.1")
	devCmd.Flags().StringVarP(&inputLocalEndpoint, "local", "l", "", "custom local upstream endpoint, e.g. 127.0.0.1")
}
