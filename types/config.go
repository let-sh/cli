package types

type Credentials struct {
	Token string `json:"token"`
}

// project-name: info
type ProjectsInfo map[string]Project

type Project struct {
	Name string `json:"name"`

	// local dir
	Dir string `json:"dir"`

	Type         string `json:"type"`
	ServeCommand string `json:"serve_command"`
}
