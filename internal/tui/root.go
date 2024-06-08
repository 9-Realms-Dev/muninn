package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	filepicker_view view = iota
	file_view
	response_view
)

const divisor = 4

type tuiState interface{}

type mainModel struct {
	help     help.Model
	loaded   bool
	focused  view
	states   []tuiState
	quitting bool
}

func (m *mainModel) Next() {
	if m.focused == response_view {
		m.focused = filepicker_view
	} else {
		m.focused++
	}
}

func (m *mainModel) Prev() {
	if m.focused == filepicker_view {
		m.focused = response_view
	} else {
		m.focused--
	}
}

func newModel(path string) mainModel {
	helpText := help.New()
	picker := initFilepicker(path)

	return mainModel{
		help: helpText,
		states: []tuiState{
			picker,
		},
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// TODO: Add http msg types later
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.states[m.focused], cmd = m.states[m.focused].(filepickerModel).Update(msg)
	return m, cmd
}

var passiveStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.HiddenBorder())

var activeStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

func (m mainModel) View() string {
	switch m.focused {
	case filepicker_view:
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			activeStyle.Render(m.states[0].(filepickerModel).View()),
			passiveStyle.Render(""),
			passiveStyle.Render(""),
		)
	default:
		return "something is wrong"
	}
}

func StartTui(path string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(newModel(path), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
