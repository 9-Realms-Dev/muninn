package tui

import "github.com/charmbracelet/bubbles/key"

type MainKeyMap struct {
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k MainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k MainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev}, // first column
		{k.Quit},         // second column
	}
}

func DefaultKeyMap() MainKeyMap {
	return MainKeyMap{
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next window"),
		),
		Prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous window"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "previous window"),
		),
	}
}

type FileViewKeyMap struct {
	Select key.Binding
}

func (k FileViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select}
}

func (k FileViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select}, // first column
		{},         // second column
	}
}

func DefaultFileViewKeyMap() FileViewKeyMap {
	return FileViewKeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select file"),
		),
	}
}

type ResponseViewKeyMap struct {
	CopyBody key.Binding
}

func (k ResponseViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CopyBody}
}

func (k ResponseViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CopyBody}, // first column
		{},           // second column
	}
}

func DefaultResponseViewKeyMap() ResponseViewKeyMap {
	return ResponseViewKeyMap{
		CopyBody: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy body"),
		),
	}
}
