package tui

import (
	"errors"
	"strings"
	"time"

	"github.com/9-Realms-Dev/muninn/internal/util"
	"github.com/charmbracelet/bubbles/filepicker"
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

func (m filepickerModel) Init(path string) tea.Cmd {
	util.Logger.Debug("Running filepickerModel init")
	return m.picker.Init()
}

func (m filepickerModel) Update(msg tea.Msg) (filepickerModel, tea.Cmd) {
	var cmd tea.Cmd
	var fileSelectCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.picker.Height = msg.Height / 4
	}
	m.picker, cmd = m.picker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selected = path

		// Send the cmd for updating the file view
		fileSelectCmd = setFile(m.selected)
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.picker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selected = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(cmd, fileSelectCmd)
}

func (m filepickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("Select the .http file you which to work on\n")
	if m.err != nil {
		s.WriteString(m.picker.Styles.DisabledFile.Render(m.err.Error()))
	}
	s.WriteString("\n\n" + m.picker.View() + "\n")
	// TODO: Add a kep map for the filepicker settings
	return s.String()
}
