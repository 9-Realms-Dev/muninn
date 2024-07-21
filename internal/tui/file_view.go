package tui

import (
	"fmt"

	munnin "github.com/9-Realms-Dev/muninn-core"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var fileViewKeyMap = DefaultFileViewKeyMap()

var (
	appStyle       = lipgloss.NewStyle().Padding(1, 2)
	listEmptyTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5dbffc")).
			Padding(0, 1)
	listTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1a1b26")).
			Background(lipgloss.Color("#c0caf5")).
			Padding(0, 1)
	listItemTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#c5c8c6", Dark: "#c5c8c6"})
	listActiveItemTitleStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#54ced6", Dark: "#54ced6"})
	filterPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#e0af68", Dark: "#e0af68"})
)

// Items
type item struct {
	request munnin.HttpRequest
}

func (i item) Title() string { return i.request.Title }
func (i item) Description() string {
	return fmt.Sprintf("%s: %s", i.request.Method, i.request.URL)
}
func (i item) FilterValue() string { return i.request.Title }

// File View
type fileViewModel struct {
	path     string
	requests []munnin.HttpRequest
	list     list.Model
}

func initList() list.Model {
	delegate := list.NewDefaultDelegate()

	delegate.Styles.NormalTitle = listItemTitleStyle
	delegate.Styles.NormalDesc = listItemTitleStyle

	delegate.Styles.SelectedTitle = listActiveItemTitleStyle
	delegate.Styles.SelectedDesc = listActiveItemTitleStyle

	// TODO: Setup reset list for window size changes
	initList := list.New([]list.Item{}, delegate, 50, 20)
	initList.Styles.Title = listEmptyTitle
	initList.Title = "Waiting for file..."

	return initList
}

func (m fileViewModel) SelectRequest() tea.Cmd {
	request, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}
	return sendRequest(request)
}

func (m fileViewModel) SelectFile(path string) (fileViewModel, tea.Cmd) {
	var cmd tea.Cmd
	if m.path != path {
		requests, err := munnin.ReadHttpFile(path)
		if err != nil {
			return m, tea.Quit
		}

		items := []list.Item{}
		for _, req := range requests {
			items = append(items, item{
				request: req,
			})
		}

		m.requests = requests
		m.list.Title = "Files requests"
		m.list.SetItems(items)

		// Style the list
		m.list.Styles.Title = listTitleStyle

		m.path = path
	}
	return m, cmd
}

func (m fileViewModel) Init() tea.Cmd {
	return nil
}

func (m fileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height - (divisor * 2))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fileViewKeyMap.Select):
			if m.list.FilterState() != list.Filtering {
				cmd = m.SelectRequest()
				return m, cmd
			}
		}
	case fileSelectedMsg:
		m, cmd = m.SelectFile(msg.path)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m fileViewModel) View() string {
	return m.list.View()
}
