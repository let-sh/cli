package deploy

import (
	"github.com/let-sh/cli/requests"
	"github.com/sirupsen/logrus"
)

func InitProject(projectName string) error {
	projectInfo, err := requests.QueryProject(projectName)
	// if project exists return
	// todo: catch other errors
	//if err == nil {
	//	return nil
	//}
	if err != nil {
		//pp.Println(err)
		return err
	}

	// project
	logrus.Debug("project info:", projectInfo)

	return nil
}
