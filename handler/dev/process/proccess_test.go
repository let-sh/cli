package process

import (
	"fmt"
	"testing"
)

func TestGetPortByProcessID(t *testing.T) {
	fmt.Println(GetPortByProcessID(82763))
}
