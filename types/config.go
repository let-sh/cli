package types

import "time"

type Credentials struct {
	Token string `json:"token"`
}

// project-name: info
type ProjectsInfo map[string]Project

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// local dir
	Dir string `json:"dir"`

	Type         string `json:"type"`
	ServeCommand string `json:"serve_command"`
}

type Extra struct {
	NotifyUpgradeTime time.Time `json:"notify"`
}
