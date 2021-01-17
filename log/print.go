package log

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Error(err error) {
	S.StopFail()

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s.\n", red("[error]"), err.Error())
	os.Exit(-1)
}

func Success(msg string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s.\n", green("[success]"), msg)
}

func Warning(msg string) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	fmt.Printf("%s %s.\n", yellow("[warn]"), msg)
}
