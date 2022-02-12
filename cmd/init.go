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
	"io/ioutil"
	"os"
	"strings"

	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils/download"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init your project via let.sh",
	Long: `Init your project via let.sh
e.g.: 
    lets init react
    lets init react new-react-project
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectType := strings.TrimSpace(args[0])
		if projectType == "" {
			log.Warning("Please specify the project type")
			return
		}

		currentDir, _ := os.Getwd()

		// check whether user custom folderName
		var folderName = projectType
		if len(args) > 1 {
			folderName = strings.TrimSpace(args[1])
		}

		// check template exists
		if _, err := requests.GetTemplate(projectType); err != nil {
			log.Error(err)
			return
		}

		log.BStart(" checking latest type") // validate project type
		tempDir, _ := ioutil.TempDir("", "init")
		defer os.RemoveAll(tempDir)
		if err := download.DownloadFile(
			fmt.Sprintf("%s/%s.zip", tempDir, projectType),
			fmt.Sprintf("http://github.com/let-sh/example/releases/latest/download/%s.zip", projectType),
		); err != nil {
			log.Error(err)
			return
		}

		log.BUpdate("downloading project template")
		logrus.Debug("download: ", fmt.Sprintf("%s/%s.zip", tempDir, projectType))
		if err := archiver.Unarchive(fmt.Sprintf("%s/%s.zip", tempDir, projectType), tempDir); err != nil {
			log.Error(err)
			return
		}
		// mv to current folder
		err := os.Rename(fmt.Sprintf("%s/%s", tempDir, projectType), fmt.Sprintf("%s/%s", currentDir, folderName))
		if err != nil {
			log.Error(errors.New("cannot init project to current folder: " + err.Error()))
			//logrus.Debug("current project dir: ", pwd)
			return
		}

		log.S.StopMessage(
			" Init succeeded\n\n" +
				"You could directly visit " + folderName + " folder by \n" +
				"> " + log.CyanUnderline("cd "+folderName),
		)
		log.BStop()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
