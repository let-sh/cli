package graphql

import "context"

func Deploy(projectType, projectName, config, channel string, cn bool) (m MutationDeploy, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"type":    projectType,
		"name":    projectName,
		"config":  config,
		"channel": channel,
		"cn":      cn,
	})
	return m, err
}

func DeployWithCheckRunID(projectType, projectName, config, channel string, cn bool,
	checkRunID int64) (m MutationDeployWithCheckRunID, err error) {
	NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"type":       projectType,
		"name":       projectName,
		"config":     config,
		"channel":    channel,
		"cn":         cn,
		"checkRunID": checkRunID,
	})
	return m, err
}
