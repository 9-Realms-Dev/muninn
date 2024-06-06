package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	filepicker view = iota
	file_view
	response_view
)

type tuiState interface{}

type mainModel struct {
	help     help.Model
	loaded   bool
	focused  view
	states   []tuiState
	quitting bool
}

func newModel() mainModel {
	helpText := help.New()
	return mainModel{
		help: helpText,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m mainModel) View() string {
	return "Hello, world!"
}

func StartTui() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
