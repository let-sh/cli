package dev

import (
	"github.com/let-sh/cli/handler/dev/tunnel"
	"github.com/let-sh/cli/info"
)

func StartClient(endpoint, local string) {
	client := tunnel.Client{
		Remote: endpoint,
		UpstreamMap: map[string]string{
			"": local,
		},
		Token:            info.Credentials.Token,
		StrictForwarding: false,
	}

	if err := client.Connect(); err != nil {
		return
	}
}
