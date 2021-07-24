package graphql

import (
	"context"
	"github.com/let-sh/cli/info"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	"time"
)

var Client = NewClient()

func NewClient() *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: info.Credentials.LoadToken()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	// set request timeout
	httpClient.Timeout = 10 * time.Second

	return graphql.NewClient("https://api.let-sh.com/query", httpClient)
}
