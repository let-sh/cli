package config

import (
	"encoding/json"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
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
