package ui

import (
	"errors"
	"github.com/let-sh/cli/log"
	"github.com/theckman/yacspin"
	"time"
)

type spinner struct {
	Spinner *yacspin.Spinner
}

var Spinner = spinner{}

func (s *spinner) Start(message string) {
	if s.Spinner == nil {
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
		s.Spinner, _ = yacspin.New(cfg)
	}

	s.Spinner.Message(" " + message)
	s.Spinner.Start()
}

func (s *spinner) Update(message string) {
	s.Spinner.Message(" " + message)
}

func (s *spinner) Stop() {
	s.Spinner.StopFail()
}

func (s *spinner) Success(message string) {
	s.Spinner.StopFail()
	log.Success(message)
}

func (s *spinner) Failed(message string) {
	s.Spinner.StopFail()
	log.Error(errors.New(message))
}
