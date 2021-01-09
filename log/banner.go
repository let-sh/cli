package log

import (
	"github.com/briandowns/spinner"
	"time"
)

var s *spinner.Spinner

func init() {
	s = spinner.New(spinner.CharSets[14], 50*time.Millisecond)
}

func BStart(suffix string) {
	s.Suffix = " " + suffix
	s.Start()
}

func BUpdate(suffix string) {
	s.Suffix = " " + suffix
}

func BStop() {
	s.Stop()
}
