package deploy

import (
	"testing"

	"github.com/let-sh/cli/utils/config"
)

func TestInitProject(t *testing.T) {
	config.Load()
	if err := InitProject("gin"); err != nil {
		t.Error(t)
	}
}
