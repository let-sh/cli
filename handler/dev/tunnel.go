package dev

import (
	"github.com/let-sh/cli/handler/dev/tunnel"
	"github.com/let-sh/cli/utils"
)

func StartClient(endpoint, local string) {
	client := tunnel.Client{
		Remote: endpoint,
		UpstreamMap: map[string]string{
			"": local,
		},
		Token:            utils.Credentials.Token,
		StrictForwarding: false,
	}

	if err := client.Connect(); err != nil {
		return
	}
}
