package deploy

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/let-sh/cli/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (c *DeployContext) LoadLetJson() {
	_, err := os.Stat("let.json")
	if err == nil {
		jsonFile, err := os.Open("let.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			log.Error(err)
			return
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		configStr := string(byteValue)
		logrus.WithFields(logrus.Fields{"configFile": configStr}).Debugln("let.json")
		err = json.Unmarshal(byteValue, &c)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
}

func (c *DeployContext) LoadCliFlag(inputProjectName, inputProjectType string) {
	if inputProjectName != "" {
		c.Name = inputProjectName
	}

	if inputProjectType != "" {
		c.Type = inputProjectType
	}
}

func (c *DeployContext) LoadRegion(cmd *cobra.Command, inputCN bool) {
	var cn bool
	// user custom by json config
	if c.CN == nil {
		c.CN = &cn
	}

	if cmd.Flags().Changed("cn") {
		// user custom by cli flag
		c.CN = &inputCN
	}
}
