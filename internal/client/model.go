package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/marczahn/person/internal/i18n"
	"github.com/marczahn/person/internal/server"
)

// inputMode determines how the client input is interpreted by the server.
type inputMode int

const (
	modeSpeech      inputMode = iota // plain text
	modeAction                       // *action*
	modeEnvironment                  // ~environment
)

var modeLabels = [...]string{"speech", "action", "environment"}

func (m inputMode) String() string { return modeLabels[m] }

func (m inputMode) next() inputMode {
	return (m + 1) % inputMode(len(modeLabels))
}

// serverMsg wraps a ServerMessage for the Bubbletea message loop.
type serverMsg server.ServerMessage

// errMsg signals a connection error.
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

var (
	thoughtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))  // green
	triggerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))   // gray
	modeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))  // blue
	dividerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// Model is the Bubbletea model for the client TUI.
type Model struct {
	conn     *Connection
	ctx      context.Context
	mode     inputMode
	thoughts []string // rendered thought lines
	input    textinput.Model
	viewport viewport.Model
	width    int
	height   int
	ready    bool
	err      error
}

// NewModel creates a new TUI model connected to the given server connection.
func NewModel(ctx context.Context, conn *Connection) Model {
	ti := textinput.New()
	ti.Placeholder = i18n.T().Client.PlaceholderSpeech
	ti.Focus()

	return Model{
		conn:  conn,
		ctx:   ctx,
		input: ti,
	}
}

// Init starts listening for server messages.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.waitForMessage(),
	)
}

// Update handles Bubbletea messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			m.mode = m.mode.next()
			m.updatePlaceholder()
			return m, nil
		case tea.KeyEnter:
			return m.handleSubmit()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputHeight := 3 // input line + divider + mode indicator
		vpHeight := m.height - inputHeight
		if vpHeight < 1 {
			vpHeight = 1
		}
		if !m.ready {
			m.viewport = viewport.New(m.width, vpHeight)
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = vpHeight
		}
		m.refreshViewport()
		return m, nil

	case serverMsg:
		m.addThought(server.ServerMessage(msg))
		m.refreshViewport()
		return m, m.waitForMessage()

	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	// Forward to text input.
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	return m, tea.Batch(cmds...)
}

// View renders the TUI.
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf(i18n.T().Client.ConnectionError+"\n", m.err)
	}
	if !m.ready {
		return i18n.T().Client.Connecting + "\n"
	}

	divider := dividerStyle.Render(strings.Repeat("─", m.width))
	modeIndicator := modeStyle.Render(fmt.Sprintf("[%s]", m.mode))

	return fmt.Sprintf("%s\n%s\n%s %s",
		m.viewport.View(),
		divider,
		modeIndicator,
		m.input.View(),
	)
}

func (m *Model) handleSubmit() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return m, nil
	}

	msg := server.ClientMessage{
		Type:    m.mode.String(),
		Content: text,
	}

	m.input.SetValue("")

	return m, func() tea.Msg {
		if err := m.conn.Send(m.ctx, msg); err != nil {
			return errMsg{err: err}
		}
		return nil
	}
}

func (m *Model) addThought(msg server.ServerMessage) {
	var line string
	ts := msg.Timestamp.Format("15:04:05")
	if msg.Trigger != "" {
		line = fmt.Sprintf("%s %s (trigger: %s)",
			triggerStyle.Render(ts),
			thoughtStyle.Render(msg.Content),
			triggerStyle.Render(msg.Trigger),
		)
	} else {
		line = fmt.Sprintf("%s %s",
			triggerStyle.Render(ts),
			thoughtStyle.Render(msg.Content),
		)
	}
	m.thoughts = append(m.thoughts, line)
}

func (m *Model) refreshViewport() {
	if !m.ready {
		return
	}
	content := strings.Join(m.thoughts, "\n")
	if m.width > 0 {
		content = ansi.Wordwrap(content, m.width, "")
	}
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

func (m *Model) updatePlaceholder() {
	tr := i18n.T()
	switch m.mode {
	case modeSpeech:
		m.input.Placeholder = tr.Client.PlaceholderSpeech
	case modeAction:
		m.input.Placeholder = tr.Client.PlaceholderAction
	case modeEnvironment:
		m.input.Placeholder = tr.Client.PlaceholderEnvironment
	}
}

func (m Model) waitForMessage() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.conn.Messages()
		if !ok {
			return errMsg{err: fmt.Errorf("%s", i18n.T().Client.ConnectionClosed)}
		}
		return serverMsg(msg)
	}
}
