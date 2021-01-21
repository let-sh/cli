package dev

import (
	"github.com/let-sh/cli/handler/dev/tunnel"
)

func StartClient(endpoint, local string) {
	client := tunnel.Client{
		Remote: endpoint, // server address, i.e. 127.0.0.1:8000
		UpstreamMap: map[string]string{
			"": local,
		}, // upstream server(local service), http://127.0.0.1:3000
		//Token:            token, // authentication token
		//StrictForwarding: strictForwarding, // forward only to the upstream URLs specified
	}

	if err := client.Connect(); err != nil {
		return
	}
}
