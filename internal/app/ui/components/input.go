package components

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/balajz/bubbline/computil"
	"github.com/balajz/bubbline/editline"
	"github.com/balajz/bubbline/history"
	"github.com/muesli/termenv"
)

var chromaFormatter = detectTerminalColorProfile()

// InputModel wraps the editline component.
type InputModel struct {
	Model       *editline.Model
	HistoryFile string
}

// NewInputModel creates and configures the input model.
func NewInputModel(prompt, historyFile string, style string, autoCompleter editline.AutoCompleteFn) (*InputModel, error) {
	el := editline.New(1, 1)
	el.Prompt = prompt

	if historyFile == "" || historyFile == "default" {
		historyFile = getHistoryFilePath()
	}

	if err := applyEditlineConfig(el, historyFile, style); err != nil {
		return nil, fmt.Errorf("applying input config: %w", err)
	}

	el.AutoComplete = autoCompleter

	return &InputModel{
		Model:       el,
		HistoryFile: historyFile,
	}, nil
}

func (m *InputModel) Init() tea.Cmd {
	return m.Model.Focus()
}

func (m *InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	newModel, cmd := m.Model.Update(msg)
	m.Model = newModel
	return *m, cmd
}

func (m *InputModel) View() string {
	return m.Model.View()
}

func (m *InputModel) SetSize(width, height int) {
	m.Model.SetSize(width, height)
}

func (m *InputModel) Value() string {
	return m.Model.Value()
}

func (m *InputModel) Reset() {
	m.Model.Reset()
}

func (m *InputModel) AddHistoryEntry(entry string) {
	m.Model.AddHistoryEntry(entry)
}

func (m *InputModel) SaveHistory() error {
	if m.HistoryFile == "" {
		return nil
	}
	return history.SaveHistory(m.Model.GetHistory(), m.HistoryFile)
}

func (m *InputModel) SetPrompt(prompt string) {
	m.Model.Prompt = prompt
}

func (m *InputModel) Prompt() string {
	return m.Model.Prompt
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

func postgresAutocomplete(pgKeywords []string) func(v [][]rune, line, col int) (string, editline.Completions) {
	return func(v [][]rune, line, col int) (string, editline.Completions) {
		word, wstart, wend := computil.FindWord(v, line, col)
		if word == "" {
			return "", nil
		}
		upperWord := strings.ToUpper(word)
		var matches []string
		for _, kw := range pgKeywords {
			if strings.HasPrefix(kw, upperWord) {
				matches = append(matches, kw)
			}
		}
		if len(matches) == 0 {
			return "", nil
		}
		return "", editline.SimpleWordsCompletion(matches, "Keywords", col, wstart, wend)
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

func applyEditlineConfig(el *editline.Model, historyFile string, style string) error {
	el.SetHighlighter(postgresHighlighter(style))
	el.SetExternalEditorEnabled(true, "sql")
	el.KeyMap.ExternalEdit = key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "edit query in external editor"),
	)

	el.CheckInputComplete = func(entireInput [][]rune, line, col int) bool {
		var sb strings.Builder
		for i, rline := range entireInput {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(string(rline))
		}
		input := strings.TrimSpace(sb.String())

		if input == "" {
			return true
		}

		if strings.HasPrefix(input, "\\") {
			return true
		}

		return strings.HasSuffix(input, ";")
	}

	entries, err := history.LoadHistory(historyFile)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	el.SetHistory(entries)
	return nil
}

func getHistoryFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".pgxcli_history.jsonl")
}
