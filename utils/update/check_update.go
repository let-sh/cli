package update

import (
	"errors"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/utils"
	"github.com/let-sh/cli/utils/config"
	"github.com/manifoldco/promptui"
	"strings"
	"time"
)

func CheckUpdate() {
	if time.Since(config.GetLastUpdateNotifyTime()) < time.Hour*24 {
		return
	}

	switch GetCurrentReleaseChannel() {
	case "beta":
		if latest, err := requests.GetLatestVersion("beta"); info.Version != latest && err == nil {
			NotifyUpgrade("beta")
		}
	case "rc":
		if latest, err := requests.GetLatestVersion("beta"); info.Version != latest && err == nil {
			NotifyUpgrade("beta")
		}
	case "dev":
		return
	default:
		latest, err := requests.GetLatestVersion("stable")
		if info.Version != latest && err == nil {
			NotifyUpgrade("stable")
		}
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
		UpgradeCli(false, "")
	}
}

func GetCurrentReleaseChannel() (channel string) {
	if strings.Contains(info.Version, "beta") {
		return "beta"
	}
	if strings.Contains(info.Version, "rc") {
		return "rc"
	}

	if strings.Contains(info.Version, "dev") {
		return "dev"
	}

	return "stable"
}
