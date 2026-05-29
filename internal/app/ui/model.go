// Package ui contains the BubbleTea model for pgxcli's interactive prompt.
package ui

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/balajz/bubbline/editline"
	"github.com/balajz/pgxcli/internal/app/ui/components"
	"github.com/davecgh/go-spew/spew"
	"github.com/muesli/termenv"
)

var chromaFormatter = detectTerminalColorProfile()

type State int

const (
	StateInput State = iota
	StateExecuting
)

// ReadyMsg signals the ui that execution is done and it should prompt.
type ReadyMsg struct{ Prefix string }

// ExecCmdMsg is used to dispatch a batch/sequence of commands.
type ExecCmdMsg struct{ Cmd tea.Cmd }

// QuitRequestMsg signals that the app wants to quit.
type QuitRequestMsg struct{}

// ConfirmQuitMsg is used internally to finalize quitting.
type ConfirmQuitMsg struct{}

type cancel func(ctx context.Context) error

type execute func(query string) tea.Cmd

type Model struct {
	input         components.InputModel
	statusModel   components.StatusModel
	spinner       components.SpinnerModel
	isSpinning    bool
	width, height int
	state         State
	quitting      bool
	prevUserInput string
	version       string
	highlighter   func(string) string

	keys   KeyMap
	styles Styles

	// execute executes a query passed and return as ExecCmdMsg + ReadyMsg.
	execute execute
	cancel  cancel

	dump *os.File
}

func New(initialPrefix string, historyFile string, style string, version string, executeFunc execute, cancelFunc cancel, autocompleter editline.AutoCompleteFn) (*Model, error) {
	inputModel, err := components.NewInputModel(initialPrefix, historyFile, style, autocompleter)
	if err != nil {
		return nil, fmt.Errorf("creating input model: %w", err)
	}

	styles := DefaultStyles()

	statusModel := components.NewStatusModel(version, issueLink)
	statusModel.SeparatorStyle = styles.InputSeparator
	statusModel.StatusBarStyle = styles.StatusBar

	spinnerModel := components.NewSpinnerModel(styles.Spinner, styles.SpinnerCaption)

	var dump *os.File
	if _, ok := os.LookupEnv("PGXCLI_DEBUG"); ok {
		dump, err = os.OpenFile("pgxcli_messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("opening debug log: %w", err)
		}
	}

	return &Model{
		input:       *inputModel,
		statusModel: statusModel,
		spinner:     spinnerModel,
		version:     version,
		state:       StateInput,
		keys:        DefaultKeyMap(),
		styles:      styles,
		execute:     executeFunc,
		cancel:      cancelFunc,
		dump:        dump,
		highlighter: postgresHighlighter(style),
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.input.Init(),
		m.statusModel.Init(),
		m.spinner.Init(),
	)
}

//nolint:gocyclo
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	var cmds []tea.Cmd

	// Send WindowSize to children
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		var smCmd tea.Cmd
		m.statusModel, smCmd = m.statusModel.Update(msg)
		cmds = append(cmds, smCmd)
		m.updateInputSize()
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {

	case QuitRequestMsg:
		m.quitting = true
		return m, func() tea.Msg {
			return ConfirmQuitMsg{}
		}

	case ConfirmQuitMsg:
		return m, tea.Quit

	case ReadyMsg:
		m.state = StateInput
		m.isSpinning = false
		if msg.Prefix != "" {
			m.input.SetPrompt(msg.Prefix)
		}
		return m, nil

	case ExecCmdMsg:
		return m, msg.Cmd

	case spinner.TickMsg:
		if m.isSpinning {
			var smCmd tea.Cmd
			m.spinner, smCmd = m.spinner.Update(msg)
			return m, smCmd
		}
		return m, nil

	case editline.InputCompleteMsg:
		if m.state == StateExecuting {
			return m, nil
		}
		return m.handleInput()

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, func() tea.Msg {
				return QuitRequestMsg{}
			}
		}
		if key.Matches(msg, m.keys.Interrupt) {
			if m.state == StateExecuting {
				m.state = StateInput
				cancelFn := m.cancel
				return m, func() tea.Msg {
					if err := cancelFn(context.Background()); err != nil {
						return ExecCmdMsg{Cmd: PrintErrCmd(err, m.styles.ErrorOutput)}
					}
					return nil
				}
			}

			m.input.Reset()
			return m, nil
		}
	}

	// Route to input only if in input state, avoiding capturing keystrokes while executing.
	if m.state == StateInput {
		var nextCmd tea.Cmd
		m.input, nextCmd = m.input.Update(msg)
		cmds = append(cmds, nextCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleInput() (tea.Model, tea.Cmd) {
	input := m.input.Value()
	trimmed := strings.TrimSpace(input)

	if trimmed == "" {
		m.input.Reset()
		return m, tea.Sequence(
			m.printUserInput(m.styles.UserInput.Render(m.input.Prompt()), ""),
			func() tea.Msg {
				return ReadyMsg{Prefix: m.input.Prompt()}
			},
		)
	}

	m.prevUserInput = input
	m.state = StateExecuting
	m.isSpinning = true
	m.input.AddHistoryEntry(input)
	m.input.Reset()

	return m, tea.Sequence(
		m.printUserInput(m.styles.UserInput.Render(m.input.Prompt()), input),
		m.spinner.Tick(),
		m.execute(trimmed),
	)
}

func (m *Model) printUserInput(prefix, input string) tea.Cmd {
	return func() tea.Msg {
		var highlightedInput string
		if input != "" {
			highlightedInput = m.highlighter(input)
		}

		// used to separate previous user input from the current one with half straight line.
		line := strings.Repeat("─", m.width/2)

		userContent := lipgloss.JoinHorizontal(lipgloss.Left, prefix, highlightedInput)
		userContent = lipgloss.JoinVertical(lipgloss.Top, m.styles.UserInputSepartor.Render(line), userContent)

		return ExecCmdMsg{Cmd: tea.Printf("%s", userContent)}
	}
}

func (m *Model) updateInputSize() {
	if m.width == 0 || m.height == 0 {
		return
	}
	// Calculate available height for input
	h := m.height - m.statusModel.StaticHeight()
	if h < 1 {
		h = 1
	}
	m.input.SetSize(m.width, h)
}

func (m *Model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	// Don't render until we know the terminal size.
	if m.width == 0 {
		return tea.NewView("")
	}

	var baseView string

	if m.state == StateExecuting {
		baseView = m.spinner.View()
	} else {
		separator := m.statusModel.SeparatorStyle.Render(strings.Repeat("─", m.width))

		baseView = lipgloss.JoinVertical(
			lipgloss.Top,
			separator,
			m.input.View(),
			m.statusModel.View(),
		)
	}

	return tea.NewView(baseView)
}

func (m *Model) saveHistory() error {
	return m.input.SaveHistory()
}

func (m *Model) Close() error {
	if m.dump != nil {
		_ = m.dump.Close()
	}
	if err := m.saveHistory(); err != nil {
		return fmt.Errorf("saving history: %w", err)
	}
	return nil
}

// PrintCmd returns a command that prints formatted text.
func PrintCmd(text string, style lipgloss.Style) tea.Cmd {
	formattedInteraction := lipgloss.Sprintf(
		"%s\n",
		style.Render(text),
	)
	return tea.Printf("%s", formattedInteraction)
}

// PrintErrCmd returns a command that prints a formatted error.
func PrintErrCmd(err error, style lipgloss.Style) tea.Cmd {
	formattedError := lipgloss.Sprintf(
		"%s\n",
		style.Render("✗ "+err.Error()),
	)
	return tea.Printf("%s", formattedError)
}

// ShowPagerCmd returns a command to execute a pager process.
func ShowPagerCmd(cmd *exec.Cmd) tea.Cmd {
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return nil
		}
		return nil
	})
}

func postgresHighlighter(style string) func(string) string {
	return func(s string) string {
		var buf bytes.Buffer
		if err := quick.Highlight(&buf, s, "postgresql", chromaFormatter, style); err != nil {
			return s
		}
		return buf.String()
	}
}

func detectTerminalColorProfile() string {
	switch termenv.ColorProfile() {
	case termenv.TrueColor:
		return "terminal16m"
	case termenv.ANSI256:
		return "terminal256"
	case termenv.ANSI:
		return "terminal16"
	default:
		return "noop"
	}
}
