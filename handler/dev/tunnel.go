package dev

import (
	"github.com/let-sh/cli/handler/dev/tunnel"
	"github.com/let-sh/cli/info"
)

func StartClient(endpoint, local, exposeFqdn string) {
	client := tunnel.Client{
		Remote: endpoint,
		UpstreamMap: map[string]string{
			exposeFqdn: local,
		},
		Token:            info.Credentials.LoadToken(),
		StrictForwarding: false,
	}

	if err := client.Connect(); err != nil {
		return
	}
}
