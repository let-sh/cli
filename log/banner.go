package log

import (
	"github.com/theckman/yacspin"
	"time"
)

var S *yacspin.Spinner

func init() {
	cfg := yacspin.Config{
		Frequency: 50 * time.Millisecond,
		CharSet:   yacspin.CharSets[14],
		//Suffix:"steps",
		SuffixAutoColon: true,
		//Message:         "exporting data",
		StopCharacter: "✓",
		StopColors:    []string{"fgGreen"},
	}
	S, _ = yacspin.New(cfg)

	//s = spinner.New(spinner.CharSets[14], 50*time.Millisecond)
}

func BStart(suffix string) {
	S.Suffix(" " + suffix)
	//s.Suffix = " " + suffix
	S.Start()
}

func BUpdate(suffix string) {
	S.Suffix(suffix)
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
