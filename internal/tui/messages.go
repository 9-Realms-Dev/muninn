package tui

import (
	munnin "github.com/9-Realms-Dev/muninn-core"
	tea "github.com/charmbracelet/bubbletea"
)

type fileSelectedMsg struct {
	path string
}

func setFile(path string) tea.Cmd {
	return func() tea.Msg {
		return fileSelectedMsg{path: path}
	}
}

type httpErrMsg struct{}

type httpRespMsg struct {
	Response *munnin.HttpResponse
}

func sendRequest(pkg item) tea.Cmd {
	return func() tea.Msg {
		resp, err := munnin.SendHttpRequest(pkg.request)
		if err != nil {
			return httpErrMsg{}
		}

		return httpRespMsg{Response: resp}
	}
}
