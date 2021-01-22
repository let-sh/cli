package dev

import (
	"fmt"
	"github.com/weaveworks/procspy"
	"testing"
)

func TestStartClient(t *testing.T) {
	StartClient("", "")
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
