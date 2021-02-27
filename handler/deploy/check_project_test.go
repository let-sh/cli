package deploy

import (
	"github.com/let-sh/cli/utils/config"
	"testing"
)

func TestInitProject(t *testing.T) {
	config.Load()
	if err := InitProject("gin"); err != nil {
		t.Error(t)
	}
}
