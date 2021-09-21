package cache

import (
	"encoding/json"
	"fmt"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/types"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
)

var ProjectsInfo = make(types.ProjectsInfo)

func init() {
	home, _ := homedir.Dir()

	if _, err := os.Stat(home + "/.let/"); err != nil {
		os.MkdirAll(home+"/.let/", os.ModePerm)
	}

	_, err := os.Stat(home + "/.let/projects.json")
	if err != nil {
		if _, err = os.Stat(home + "/.let/projects.json"); err != nil {
			f, _ := os.Create(home + "/.let/projects.json")
			f.WriteString("{}")
		}
	} else {
		projectsInfoFile, _ := ioutil.ReadFile(home + "/.let/projects.json")
		err := json.Unmarshal(projectsInfoFile, &ProjectsInfo)
		if err != nil {
			ioutil.WriteFile(home+"/.let/projects.json", []byte("{}"), 0644)
			log.Errorf("load project error: %s", err.Error())
			return
		}
	}
}

func SaveProjectInfo(projectInfo types.Project) error {
	ProjectsInfo[projectInfo.Name] = projectInfo

	// Convert golang object back to byte
	byteValue, err := json.Marshal(ProjectsInfo)
	if err != nil {
		log.Error(err)
		return err
	}

	// Write back to file
	home, _ := homedir.Dir()
	err = ioutil.WriteFile(home+"/.let/projects.json", byteValue, 0644)
	return nil
}

func GetProjectInfo(dir string) (project types.Project, err error) {
	for _, v := range ProjectsInfo {
		if v.Dir == dir {
			return v, nil
		}
	}
	return project, fmt.Errorf("project not found: %v", dir)
}
