package requests

import (
	"fmt"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	fmt.Println(GetLatestVersion("rc"))
}
