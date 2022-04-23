package graphql

import "context"

func Link(projectID string, hostname string) (m MutationLink, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"projectID": projectID,
		"hostname":  hostname,
	})
	return m, err
}

func Unlink(projectID string, hostname string) (m MutationUnlink, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"projectID": projectID,
		"hostname":  hostname,
	})
	return m, err
}
