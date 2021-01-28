package update

import (
	"errors"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils"
	"github.com/manifoldco/promptui"
	"strings"
)

func CheckUpdate() {
	if strings.Contains(info.Version, "beta") {
		if latest, err := requests.GetLatestVersion("beta"); info.Version != latest && err != nil {
			NotifyUpgrade("beta")
		}
		return
	}
	if strings.Contains(info.Version, "rc") {
		if latest, err := requests.GetLatestVersion("rc"); info.Version != latest && err != nil {
			NotifyUpgrade("rc")
		}
		return
	}

	if strings.Contains(info.Version, "development") {
		return
	}

	// else
	if latest, err := requests.GetLatestVersion("stable"); info.Version != latest && err != nil {
		NotifyUpgrade("stable")
	}
}

func NotifyUpgrade(channel string) {
	prompt := promptui.Prompt{
		Label:   "Detected new version of cli released, update now?[Y/n]",
		Default: "Y",
		Validate: func(input string) error {
			if utils.ItemExists([]string{"", "n", "N", "No", "Y", "y", "yes", "Yes"}, input) {
				return nil
			}
			return errors.New("no matching input")
		},
	}

	result, err := prompt.Run()
	if err != nil {
		log.Errorf("Prompt failed %v\n", err.Error())
		return
	}

	if utils.ItemExists([]string{"Y", "y", "yes", "Yes"}, result) {
		UpgradeCli(channel)
	}
}