package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/let-sh/cli/utils"
	"github.com/sirupsen/logrus"
	"os"
)

type InputAreaConfig struct {
	Layout             string   `json:"layout"`
	DefaultValue       string   `json:"default_value"`
	DefaultPlaceholder string   `json:"default_placeholder,omitempty"`
	PlaceHolders       []string `json:"placeholders,omitempty"`
}

func InputArea(config ...InputAreaConfig) (string, error) {
	conf := InputAreaConfig{}
	if len(config) > 0 {
		conf = config[0]
	}

	ti := textinput.NewModel()
	ti.Placeholder = conf.DefaultPlaceholder
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 36

	m := &inputAreaModal{
		textInput:    ti,
		layoutConfig: conf,
		err:          nil,
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		logrus.Fatal(err)
		return "", err
	}

	return m.Value(), nil
}

type inputAreaModal struct {
	textInput    textinput.Model
	layoutConfig InputAreaConfig
	err          error
}

type errMsg error

func (m inputAreaModal) Init() tea.Cmd {
	return textinput.Blink
}

// Value return the input string
func (m *inputAreaModal) Value() string {
	return m.textInput.Value()
}

func (m *inputAreaModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// prompt
	if matched := utils.MatchPrefix(m.layoutConfig.PlaceHolders, m.textInput.Value()); matched != "" {
		m.textInput.Placeholder = utils.MatchPrefix(m.layoutConfig.PlaceHolders, m.textInput.Value())
	}

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
		case tea.KeyRight:
			prompt := utils.MatchPrefix(m.layoutConfig.PlaceHolders, m.textInput.Value())
			m.textInput.SetValue(prompt)
			m.textInput.SetCursor(len(prompt))
		case tea.KeyTab:
			prompt := utils.MatchPrefix(m.layoutConfig.PlaceHolders, m.textInput.Value())
			m.textInput.SetValue(prompt)
			m.textInput.SetCursor(len(prompt))
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputAreaModal) View() string {
	return fmt.Sprintf(
		m.layoutConfig.Layout+"%s\n",
		m.textInput.View(),
	)
}
