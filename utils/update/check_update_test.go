package update

import (
	"github.com/let-sh/cli/info"
	"testing"
)

func TestCheckUpdate(t *testing.T) {
	info.Version = "0.0.67"
	CheckUpdate()
}
