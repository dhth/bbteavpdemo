// Copy of https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go
package main

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().

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
)

type model struct {
	content                    string
	ready                      bool
	viewport                   viewport.Model
	useHighPerformanceRenderer bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.HighPerformanceRendering = m.useHighPerformanceRenderer
			var contentHighlighted bytes.Buffer
			highlightErr := quick.Highlight(&contentHighlighted, m.content, "diff", "terminal16m", "monokai")
			if highlightErr != nil {
				m.viewport.SetContent(m.content)
			} else {
				prefix := fmt.Sprintf("This is bubbletea's VP. high performance rendering: %v\n\n", m.useHighPerformanceRenderer)
				m.viewport.SetContent(prefix + contentHighlighted.String())
			}
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

		if m.useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func bteaVPort(hp bool) {
	flag.Parse()
	// Load some text for our viewport
	content, err := os.ReadFile("diff.patch")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		model{
			content:                    string(content),
			useHighPerformanceRenderer: hp,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
