package ui

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"testing"
)

func TestCheckBox(t *testing.T) {
	for i := uint8(16); i <= 231; i++ {
		fmt.Println(i, aurora.Index(i, "pew-pew"), aurora.BgIndex(i, "pew-pew"))
	}
}
