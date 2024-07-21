package tui

import (
	"fmt"
	"github.com/9-Realms-Dev/muninn/internal/util"
	"os"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	filepicker_view view = iota
	file_view
	response_view
)

var (
	maxHeight = 100
	maxWidth  = 100
	divisor   = 4
)

var keymap = DefaultKeyMap()

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
	list := initList()

	return mainModel{
		help: helpText,
		states: []tuiState{
			picker,
			fileViewModel{list: list},
			responseViewModel{},
		},
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		util.Logger.Debug(msg)
		maxHeight = msg.Height
		maxWidth = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keymap.Next):
			m.Next()
		case key.Matches(msg, keymap.Prev):
			m.Prev()
		}
	case fileSelectedMsg:
		// Send update regardless if the file view is focused or not
		m.states[file_view], cmd = m.states[file_view].(fileViewModel).Update(msg)
		cmds = append(cmds, cmd)
	case httpRespMsg:
		m.states[response_view], cmd = m.states[response_view].(responseViewModel).Update(msg)
		cmds = append(cmds, cmd)
		m.focused = response_view
	}

	switch m.focused {
	case filepicker_view:
		m.states[m.focused], cmd = m.states[m.focused].(filepickerModel).Update(msg)
		cmds = append(cmds, cmd)
	case file_view:
		m.states[m.focused], cmd = m.states[m.focused].(fileViewModel).Update(msg)
		cmds = append(cmds, cmd)
	case response_view:
		m.states[m.focused], cmd = m.states[m.focused].(responseViewModel).Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

var rootStyle = lipgloss.NewStyle().
	Width(maxWidth+(maxWidth-40)).
	Height(util.SetDefaultHeight(maxHeight/3)).
	Padding(0, 2)

var passiveStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	Width(maxWidth / 2).
	Height(maxHeight / 3).
	MaxHeight(util.SetDefaultHeight(maxHeight)).
	BorderForeground(lipgloss.Color("#7aa2f7"))

var activeStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	Width(maxWidth / 2).
	Height(maxHeight / 3).
	MaxHeight(util.SetDefaultHeight(maxHeight)).
	BorderForeground(lipgloss.Color("#f69058"))

func (m mainModel) View() string {
	switch m.focused {
	case filepicker_view:
		return rootStyle.Render(lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeStyle.Render(m.states[filepicker_view].(filepickerModel).View()),
			passiveStyle.Render(m.states[file_view].(fileViewModel).View()),
			passiveStyle.Render(m.states[response_view].(responseViewModel).View()),
		))
	case file_view:
		return rootStyle.Render(lipgloss.JoinHorizontal(
			lipgloss.Top,
			passiveStyle.Render(m.states[filepicker_view].(filepickerModel).View()),
			activeStyle.Render(m.states[file_view].(fileViewModel).View()),
			passiveStyle.Render(m.states[response_view].(responseViewModel).View()),
		))
	case response_view:
		return rootStyle.Render(lipgloss.JoinHorizontal(
			lipgloss.Top,
			passiveStyle.Render(m.states[filepicker_view].(filepickerModel).View()),
			passiveStyle.Render(m.states[file_view].(fileViewModel).View()),
			activeStyle.Render(m.states[response_view].(responseViewModel).View()),
		))
	default:
		return "something is wrong"
	}
}

func StartTui(path string) {
	p := tea.NewProgram(newModel(path), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
