package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCommand(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	fmt.Println(filepath.Base(dir))
}
