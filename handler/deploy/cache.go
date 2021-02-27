package deploy

import "github.com/let-sh/cli/utils/cache"

func (c *DeployContext) LoadProjectInfoCache() {
	if i, err := cache.GetProjectInfo(c.Name); err == nil {
		c.LetConfig.Name = i.Name
		c.LetConfig.Type = i.Type
	}
}
