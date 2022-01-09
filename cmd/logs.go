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
	"github.com/let-sh/cli/ui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// logsCmd represents the log command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show latest logs",
	Long: `Show latest logs under current project.

e.g. 
"lets logs --tail 10"         print latest 10 line logs
"lets logs -p hello-world"    print latest logs under project 'hello-world'
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented yet")
		//data, err := requests.QueryDeployments("gin", 1)
		//if err != nil {
		//	log.Error(err)
		//}
		//fmt.Println(data)
		//deploy.InitProject("gin")
		//ui.Radio(ui.RadioConfig{
		//	Prefix: "type: gin",
		//})
		//ui.InputArea(ui.InputAreaConfig{
		//	Layout:             "key",
		//	DefaultValue:       "v",
		//	DefaultPlaceholder: "value",
		//	PlaceHolders:       []string{"value", "test"},
		//})
		//
		//fmt.Println("log called")
		detectedType := "gin"
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
		ui.InputArea(ui.InputAreaConfig{
			Layout:             "Please enter your command to start service:",
			DefaultPlaceholder: defaultCommand,
			PlaceHolders:       []string{defaultCommand},
		})
	},
}

var inputLines int
var logInputProjectName string

func init() {
	rootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dir, _ := os.Getwd()

	logsCmd.Flags().IntVarP(&inputLines, "lines", "l", 10, "latest lines of logs")
	logsCmd.Flags().StringVarP(&logInputProjectName, "project", "p", filepath.Base(dir), "project name, e.g. react")
}
