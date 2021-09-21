package graphql

import (
	"context"
)

func SetPreference(name, value string) (m MutationSetPreference, err error) {
	err = NewClient().Mutate(context.Background(), &m, map[string]interface{}{
		"name":  name,
		"value": value,
	})
	return m, err
}

func GetAllPreference() (q QueryAllPreference, err error) {
	err = NewClient().Query(context.Background(), &q, nil)
	return q, err
}
