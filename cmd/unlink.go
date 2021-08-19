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
	"os"
	"strings"

	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils/cache"

	"github.com/spf13/cobra"
)

// unlinkCmd represents the unlink command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "UnLink domain from current project",
	Long: `UnLink domain from current project.
e.g.: lets unlink test.let.sh`,
	Run: func(cmd *cobra.Command, args []string) {
		//detectedType :=deploy.DetectProjectType()
		dir, _ := os.Getwd()
		p, err := cache.GetProjectInfo(dir)

		// if cache exists
		// todo: support query project
		if err != nil {
			log.Error(errors.New("please deploy first"))
			return
		}

		result, err := requests.Unlink(p.ID, strings.TrimSpace(args[0]))
		if err != nil {
			log.Error(err)
			return
		}

		if result == false {
			log.Error(errors.New("unlink failed"))
			return
		}
		log.Success("unlink success")
		return
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unlinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unlinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
