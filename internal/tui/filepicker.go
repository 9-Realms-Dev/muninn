package tui

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type filepickerModel struct {
	picker   filepicker.Model
	selected string
	quitting bool
	err      error
	help     help.Model
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
		help:   help.New(),
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

	// Did the user select a file?
	if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selected = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.picker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selected = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

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
		s.WriteString("Selected file: \n" + m.picker.Styles.Selected.Render(m.selected))
	}
	s.WriteString("\n\n" + m.picker.View() + "\n")
	// TODO: Add a kep map for the filepicker settings
	return s.String()
}
