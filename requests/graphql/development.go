package graphql

import "context"

func StartDevelopment(projectID string) (m MutationStartDevelopment, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"projectID": projectID,
	})
	return m, err
}

func StopDevelopment(projectID string) (m MutationStopDevelopment, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"projectID": projectID,
	})
	return m, err
}
