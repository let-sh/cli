package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"log"
)

type SpinnerConfig struct {
	Prefix      string
	Suffix      string
	Placeholder string
	Default     bool
	SpinnerText string
}

type spinnerModel struct {
	spinner  spinner.Model
	config   SpinnerConfig
	quitting bool
	err      error
}

func Spinner(configs ...SpinnerConfig) bool {
	conf := SpinnerConfig{}
	if len(configs) > 0 {
		conf = configs[0]
	}

	ti := spinner.NewModel()

	m1 := &spinnerModel{
		spinner: ti,
		config:  conf,
		err:     nil,
	}
	p := tea.NewProgram(m1)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return false
}

func (m spinnerModel) initialModel() spinnerModel {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	return spinnerModel{spinner: s}
}

func (m spinnerModel) Init() tea.Cmd {
	return spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			fallthrough
		case "esc":
			fallthrough
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m spinnerModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	s := termenv.String(m.spinner.View()).Foreground(termenv.ANSIColor(int(MainColor))).String()
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", s)
	if m.quitting {
		return str + "\n"
	}
	return str
}
