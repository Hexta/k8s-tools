package duckdb

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Hexta/k8s-tools/internal/format"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	ctx         context.Context
	dataDir     string
}

const (
	TaPlaceholder = "SQL "
)

func RunTUI(ctx context.Context, dataDir string) {
	p := tea.NewProgram(initialModel(ctx, dataDir))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel(ctx context.Context, dataDir string) model {
	ta := textarea.New()
	ta.Placeholder = TaPlaceholder
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 0
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(10, 5)
	vp.Style.AlignVertical(lipgloss.Bottom)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		ctx:         ctx,
		dataDir:     dataDir,
	}
}
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			result := m.handleKeyEnter()

			query := formatSQL(m.textarea.Value())

			m.messages = append(m.messages, fmt.Sprintf("%s\n%s", query, result))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 1
		m.viewport.Height = msg.Height - 5
		m.textarea.SetWidth(msg.Width - 5)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) handleKeyEnter() string {
	result, err := Query(m.ctx, m.dataDir, m.textarea.Value())
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	return format.Table(result, format.Options{})
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func formatSQL(sql string) string {
	var b bytes.Buffer
	err := quick.Highlight(&b, sql, "sql", "terminal16m", "monokai")
	if err != nil {
		log.Errorf("Failed to highlight SQL: %v", err)
		return ""
	}

	return b.String()
}
