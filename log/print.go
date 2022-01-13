package log

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/getsentry/sentry-go"
	"os"
)

func Errorf(template string, formatString ...string) {

	err := errors.New(fmt.Sprintf(template, formatString))

	S.StopFail()
	sentry.CaptureException(err)
	red := color.New(color.BgRed, color.FgBlack).SprintFunc()
	fmt.Printf("%s %s.\n", red(" error "), err.Error())
	os.Exit(-1)
}

func Error(err error) {
	S.StopFail()
	sentry.CaptureException(err)
	red := color.New(color.BgRed, color.FgBlack).SprintFunc()
	fmt.Printf("%s %s.\n", red(" error "), err.Error())
	os.Exit(-1)
}

func Success(msg string) {
	green := color.New(color.BgGreen, color.FgBlack).SprintFunc()
	fmt.Printf("%s %s\n", green(" success "), msg)
}

func Warning(msg string) {
	yellow := color.New(color.BgHiYellow, color.FgBlack).SprintFunc()
	fmt.Printf("%s %s.\n", yellow(" warn "), msg)
}
