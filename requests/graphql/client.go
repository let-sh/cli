package graphql

import (
	"github.com/let-sh/cli/requests/http_client"
	"github.com/shurcooL/graphql"
	"net/http"
)

var Client = NewClient()

func NewClient() *graphql.Client {
	//retryClient := retryablehttp.NewClient()
	//retryClient.RetryMax = 3
	//httpClient := retryClient.StandardClient()
	//
	//rt := WithHeader(httpClient.Transport)
	//rt.Set("Authorization", "Bearer "+info.Credentials.LoadToken())
	//httpClient.Transport = rt
	//src := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: info.Credentials.LoadToken()},
	//)
	//httpClient := oauth2.NewClient(context.Background(), src)
	//
	//// set request timeout
	//httpClient.Timeout = 10 * time.Second

	return graphql.NewClient("https://api.let-sh.com/query", http_client.NewClient())
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
