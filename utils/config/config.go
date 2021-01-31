package config

import (
	"encoding/json"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/types"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"time"
)

func init() {
	// check config dir exists
	home, _ := homedir.Dir()
	_, err := os.Stat(home + "/.let")
	if os.IsNotExist(err) {
		os.MkdirAll(home+"/.let", os.ModePerm)
	}

	_, err = os.Stat(home + "/.let/credentials.json")
	if os.IsNotExist(err) {
		f, _ := os.Create(home + "/.let/credentials.json")
		f.WriteString("{}")
	}

	_, err = os.Stat(home + "/.let/projects.json")
	if os.IsNotExist(err) {
		f, _ := os.Create(home + "/.let/projects.json")
		f.WriteString("{}")
	}

	_, err = os.Stat(home + "/.let/extra.json")
	if os.IsNotExist(err) {
		f, _ := os.Create(home + "/.let/extra.json")
		extra := types.Extra{NotifyUpgradeTime: time.Now()}
		f.WriteString(func() string {
			str, _ := json.Marshal(extra)
			return string(str)
		}())
	}

	// bootstrap configs
	credentialsFile, _ := ioutil.ReadFile(home + "/.let/credentials.json")
	err = json.Unmarshal(credentialsFile, &info.Credentials)
	if err != nil {
		log.Error(err)
	}
}

func Load() {}

func SetToken(token string) {
	info.Credentials.Token = token
	home, _ := homedir.Dir()

	file, _ := json.MarshalIndent(&info.Credentials, "", "  ")
	_ = ioutil.WriteFile(home+"/.let/credentials.json", file, 0644)
}

func GetLastUpdateNotifyTime() (latest time.Time) {
	home, _ := homedir.Dir()
	var extra types.Extra
	extrasFile, _ := ioutil.ReadFile(home + "/.let/extra.json")
	err := json.Unmarshal(extrasFile, &extra)
	if err != nil {
		os.Remove(home + "/.let/extra.json")
		f, _ := os.Create(home + "/.let/extra.json")
		extra := types.Extra{NotifyUpgradeTime: time.Now()}
		f.WriteString(func() string {
			str, _ := json.Marshal(extra)
			return string(str)
		}())
		return extra.NotifyUpgradeTime
	}

	latestUpdateNotifyTime := extra.NotifyUpgradeTime
	extra.NotifyUpgradeTime = time.Now()
	str, _ := json.Marshal(extra)
	err = ioutil.WriteFile(home+"/.let/extra.json", str, 0644)

	return latestUpdateNotifyTime
}
