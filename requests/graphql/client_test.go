package graphql

import (
	"context"
	"errors"
	"github.com/let-sh/cli/utils/config"
	"github.com/shurcooL/graphql"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	config.Load()
	var query struct {
		AllPreferences struct {
			Channel string
		}
	}
	err := Client.Query(context.Background(), &query, nil)

	var requestError *graphql.RequestError
	if errors.As(err, &requestError) {
		log.Println("request error: ", requestError.Error())
	}

	var graphqlError *graphql.GraphQLError
	if errors.As(err, &graphqlError) {
		log.Println("graphql error: ", graphqlError)
		return
	}

	if err != nil {
		t.Error(err)
		log.Println("error: ", err.Error())
	}
	t.Log(query)
}

func TestGetAllPreference(t *testing.T) {
	config.Load()
	pref, err := GetAllPreference()

	var requestError *graphql.RequestError
	if errors.As(err, &requestError) {
		log.Println("request error: ", requestError.Error())
	}

	var graphqlError *graphql.GraphQLError
	if errors.As(err, &graphqlError) {
		log.Println("graphql error: ", graphqlError)
		return
	}
	log.Println(pref)

}

func TestCombineQuery(t *testing.T) {
	config.Load()

	// pref, err := GetAllPreference()
	var query struct {
		QueryAllPreference
		QueryBuildTemplate
	}
	err := Client.Query(context.Background(), &query, map[string]interface{}{
		"type": graphql.String("gin"),
	})

	var requestError *graphql.RequestError
	if errors.As(err, &requestError) {
		log.Println("request error: ", requestError.Error())
	}

	var graphqlError *graphql.GraphQLError
	if errors.As(err, &graphqlError) {
		log.Println("graphql error: ", graphqlError)
		return
	}
	log.Println(query)
}
