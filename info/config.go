package info

import (
	"encoding/json"
	"github.com/let-sh/cli/types"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Credentials credentials

type (
	credentials types.Credentials
)

func (c *credentials) LoadToken() string {
	if Credentials.Token == "" {
		home, _ := os.UserHomeDir()
		credentialsFile, _ := ioutil.ReadFile(home + "/.let/credentials.json")
		err := json.Unmarshal(credentialsFile, &Credentials)
		if err != nil {
			logrus.Debugln("load token error: %s", err.Error())
		}
	}
	return Credentials.Token
}

func (c *credentials) SetToken(token string) {
	c.Token = token
}
