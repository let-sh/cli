package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Error(err error) {
	red := color.New(color.FgRed).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s.\n", red("[error]"), err.Error())
}

func Success(msg string) {
	green := color.New(color.FgGreen).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s.\n", green("[success]"), msg)
}

func Warning(msg string) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s %s.\n", yellow("[warn]"), msg)
}
