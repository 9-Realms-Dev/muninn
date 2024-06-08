package tui

import (
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type filepickerModel struct {
	picker   filepicker.Model
	selected string
	quitting bool
	err      error
}

func initFilepicker(path string) filepickerModel {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".http"}
	fp.CurrentDirectory = path
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.AutoHeight = false

	return filepickerModel{
		picker: fp,
	}
}

func (m filepickerModel) Init() tea.Cmd {
	return nil
}

func (m filepickerModel) Update(msg tea.Msg) (filepickerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.picker.Height = msg.Height - (divisor * 2)
	}
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m filepickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.picker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selected == "" {
		s.WriteString("Pick a file: " + m.picker.CurrentDirectory)
	} else {
		s.WriteString("Selected file: " + m.picker.Styles.Selected.Render(m.selected))
	}
	s.WriteString("\n\n" + m.picker.View() + "\n")
	return s.String()
}
