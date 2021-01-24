package update

import (
	"fmt"
	"os"
	"os/exec"
)

func UpgradeCli(channel string) {
	c := exec.Command("sh", "-c", "curl install.let.sh.cn | bash")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Printf("commad run failed with %s\n", err)
	}
}
