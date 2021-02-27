package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

type RadioConfig struct {
	Prefix      string
	Suffix      string
	Placeholder string
	Default     bool
	RadioText   string
}

type radioModel struct {
	textInput textinput.Model
	config    RadioConfig
	err       error
}

func Radio(configs ...RadioConfig) bool {
	conf := RadioConfig{}
	if len(configs) > 0 {
		conf = configs[0]
	}

	ti := textinput.NewModel()
	ti.Placeholder = conf.Placeholder
	ti.Focus()

	ti.Prompt = ""
	m1 := &radioModel{
		textInput: ti,
		config:    conf,
		err:       nil,
	}
	p := tea.NewProgram(m1)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

	// y
	if strings.Contains(strings.ToLower(m1.Value()), "y") {
		return true
	}

	// n
	if strings.Contains(strings.ToLower(m1.Value()), "n") {
		return false
	}

	// default
	if conf.Default {
		return true
	} else {
		return false
	}
}

func (m *radioModel) Init() tea.Cmd {
	if len(m.config.RadioText) == 0 {
		if m.config.Default {
			m.config.RadioText = "[Y/n]"
		} else {
			m.config.RadioText = "[y/N]"
		}
	}
	return textinput.Blink
}

// Value return the input string
func (m *radioModel) Value() string {
	return m.textInput.Value()
}

type radioErrMsg error

func (m *radioModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyEsc:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case radioErrMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	val := m.textInput.Value()
	if len(val) > 1 {
		logrus.Debugf(val)
	}

	return m, cmd
}

func (m radioModel) View() string {
	var s string
	s += m.config.Prefix

	if m.config.Default {
		s += " " + m.config.RadioText + ": "
	} else {
		s += " " + m.config.RadioText + ": "
	}

	s += m.textInput.View()
	s += m.config.Suffix
	return s + "\n"
}
