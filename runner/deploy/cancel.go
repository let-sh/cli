package deploy

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		//DeleteFiles()
		os.Exit(0)
	}()
}
