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
var docStyle = lipgloss.NewStyle().Margin(1, 2)

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
	// TODO: Setup reset list for window size changes
	initList := list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 20)
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
		m.path = path
	}
	return m, cmd
}

func (m fileViewModel) Init() tea.Cmd {
	return nil
}

func (m fileViewModel) Update(msg tea.Msg) (fileViewModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fileViewKeyMap.Select):
			cmd = m.SelectRequest()
			return m, cmd
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
