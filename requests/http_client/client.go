package http_client

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/let-sh/cli/info"
	"io/ioutil"
	"log"
	"net/http"
)

func NewClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	retryClient.RetryMax = 3
	httpClient := retryClient.StandardClient()

	rt := WithHeader(httpClient.Transport)
	rt.Set("Authorization", "Bearer "+info.Credentials.Token)
	httpClient.Transport = rt
	//src := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: info.Credentials.LoadToken()},
	//)
	//httpClient := oauth2.NewClient(context.Background(), src)
	//
	//// set request timeout
	//httpClient.Timeout = 10 * time.Second

	return httpClient
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
