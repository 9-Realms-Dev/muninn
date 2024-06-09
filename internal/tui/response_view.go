package tui

import (
	"fmt"
	"net/http"
	"strings"

	munnin "github.com/9-Realms-Dev/muninn-core"
	munninFormat "github.com/9-Realms-Dev/muninn-core/formats"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
	defaultWidth  = 30
	defaultHeight = 30
)

type responseViewModel struct {
	response *http.Response
	viewport viewport.Model
	ready    bool
	content  string
	err      error
}

func (m responseViewModel) RenderResponse(pkg munnin.HttpResponse) (responseViewModel, tea.Cmd) {
	// TODO: Proper error handling
	if pkg.Error != nil {
		return m, nil
	}

	resp := pkg.Response
	json, err := munninFormat.FormatJSONResponse(resp)
	if err != nil {
		return m, nil
	}

	m.response = resp
	m.content = json.CliRender()

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	m.viewport = viewport.New(defaultWidth*2, defaultHeight-verticalMarginHeight)
	m.viewport.YPosition = headerHeight
	m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
	m.viewport.SetContent(m.content)
	m.ready = true

	return m, nil
}

func (m responseViewModel) Init() tea.Cmd {
	return nil
}

func (m responseViewModel) Update(msg tea.Msg) (responseViewModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case httpRespMsg:
		m, cmd = m.RenderResponse(*msg.Response)
		cmds = append(cmds, cmd)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m responseViewModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m responseViewModel) headerView() string {
	title := titleStyle.Render("Json Response")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m responseViewModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
