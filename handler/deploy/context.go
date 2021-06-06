package deploy

import "github.com/let-sh/cli/types"

type DeployContext struct {
	types.LetConfig
	Channel string `json:"channel,omitempty"`
}
