package deploy

import (
	"github.com/let-sh/cli/requests/graphql"
	"github.com/let-sh/cli/types"
)

type DeployContext struct {
	types.LetConfig
	Channel          string `json:"channel,omitempty"`
	PreDeployRequest struct {
		graphql.QueryCheckDeployCapability
		graphql.QueryBuildTemplate
		graphql.QueryStsToken
		graphql.QueryPreference
	} `json:"-"`
}
