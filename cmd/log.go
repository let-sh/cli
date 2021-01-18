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
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show latest logs",
	Long: `Show latest logs under current project.

e.g. 
"lets logs --tail 10"         print latest 10 line logs
"lets logs -p hello-world"    print latest logs under project hello-world
`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := requests.QueryDeployments("gin", 1)
		if err != nil {
			log.Error(err)
		}
		fmt.Println(data)
		fmt.Println("log called")
	},
}

var inputLines int
var logInputProjectName string

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dir, _ := os.Getwd()

	logCmd.Flags().IntVarP(&inputLines, "lines", "l", 10, "latest lines of logs")
	logCmd.Flags().StringVarP(&logInputProjectName, "project", "p", filepath.Base(dir), "project name, e.g. react")
}
