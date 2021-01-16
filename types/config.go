package types

type Credentials struct {
	Token string `json:"token"`
}

// project-name: info
type ProjectsInfo map[string]Project

type Project struct {
	Name string `json:"name"`
	Dir  string `json:"dir"`
	Type string `json:"type"`
}
