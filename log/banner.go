package log

import (
	"github.com/let-sh/cli/log/sentry"
	"github.com/theckman/yacspin"
	"time"
)

var S *yacspin.Spinner

func init() {
	if S == nil {
		cfg := yacspin.Config{
			Frequency: 50 * time.Millisecond,
			CharSet:   yacspin.CharSets[14],
			Colors:    []string{"fgCyan"},
			//Suffix:"steps",
			SuffixAutoColon: true,
			//Message:         "exporting data",
			StopCharacter: "✓",
			StopColors:    []string{"fgGreen"},
		}
		S, _ = yacspin.New(cfg)
	}

	sentry.Init()
}

// deprecated
func BStart(message string) {
	cfg := yacspin.Config{
		Frequency: 50 * time.Millisecond,
		CharSet:   yacspin.CharSets[14],
		Colors:    []string{"fgCyan"},
		//Suffix:"steps",
		SuffixAutoColon: true,
		//Message:         "exporting data",
		StopCharacter: "✓",
		StopColors:    []string{"fgGreen"},
	}
	S, _ = yacspin.New(cfg)
	S.Message(" " + message)
	S.Start()
}

// deprecated
func BUpdate(message string) {
	S.Message(" " + message)
}

// deprecated
func BStop() {
	S.Stop()
}

// deprecated
func BPause() {
	S.Pause()
}

// deprecated
func BUnpause() {
	S.Unpause()
}
