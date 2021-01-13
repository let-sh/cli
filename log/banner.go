package log

import (
	"github.com/theckman/yacspin"
	"time"
)

var S *yacspin.Spinner

func init() {

	//s = spinner.New(spinner.CharSets[14], 50*time.Millisecond)
}

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

func BUpdate(message string) {
	S.Message(" " + message)
}

func BStop() {
	S.Stop()
}

func BPause() {
	S.Pause()
}

func BUnpause() {
	S.Unpause()
}
