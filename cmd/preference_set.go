package cmd

import (
	"errors"
	"strings"

	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/spf13/cobra"
)

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

// preferenceSetCmd represents the preference command
var preferenceSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set you personal preferences",
	Long: `Set your personal preferences

e.g. lets pref set default_channel dev
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Warning(`error input, you could try as below:

e.g. lets pref set channel dev`)
		}
		value, err := requests.SetPreference(strings.TrimSpace(args[0]), strings.TrimSpace(args[1]))
		if err != nil || !value {
			log.Error(errors.New("cannot set preference: " + err.Error()))
			return
		}

		log.Success("set preference: " + args[0] + "=" + args[1])
	},
}

func init() {
	preferenceCmd.AddCommand(preferenceSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// preferenceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// preferenceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
