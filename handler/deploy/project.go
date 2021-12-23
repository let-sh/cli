package deploy

import (
	"fmt"
	"strings"

	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/let-sh/cli/ui"
	. "github.com/logrusorgru/aurora"
)

func (c *DeployContext) ConfirmProject() bool {
	log.S.StopFail()
	_, err := requests.QueryProject(c.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			// TODO: better error handler
			log.Error(err)
			return false
		}

		// let user check project info
		// pretty print current project info
		if ui.Radio(ui.RadioConfig{
			Prefix: fmt.Sprintf(
				"%s\nname: %s\ntype: %s\n%s",
				Index(51, "New project detected:"),
				c.Name,
				c.Type,
				Index(51, "\ncontinue to deploy?"),
			),
			RadioText: Index(51, "[Y/n]").String(),
		}) {
			// changing project config
			log.BStart("deploying")
			return true
		} else {
			return false
		}
	}
	log.BStart("deploying")
	return true
}
