package types

type LetConfig struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
	//Build struct {
	//	Env struct {
	//		K string `json:"K,omitempty"`
	//	} `json:"env,omitempty"`
	//} `json:"build,omitempty"`

	// static dir
	Static   string `json:"static,omitempty"`
	Redirect []struct {
		Source      string `json:"source,omitempty"`
		Destination string `json:"destination,omitempty"`
		Type        int    `json:"type,omitempty"`
	} `json:"redirect,omitempty"`
	Rewrite  []struct {
		Source      string `json:"source,omitempty"`
		Destination string `json:"destination,omitempty"`
		Type        int    `json:"type,omitempty"`
	} `json:"rewrite,omitempty"`
	Link    []string      `json:"link,omitempty"`
}