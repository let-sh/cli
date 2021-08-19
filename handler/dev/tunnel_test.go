package dev

import (
	"fmt"
	"testing"

	"github.com/weaveworks/procspy"
)

func TestStartClient(t *testing.T) {
	// StartClient("", "")
}

func TestNet(t *testing.T) {
	lookupProcesses := true
	cs, err := procspy.Connections(lookupProcesses)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TCP Connections:\n")
	for c := cs.Next(); c != nil; c = cs.Next() {
		fmt.Printf(" - %v\n", c)
	}
}
