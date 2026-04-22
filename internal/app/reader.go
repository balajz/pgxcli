package app

import (
	"context"
	"errors"

	"github.com/balaji01-4d/pgxcli/internal/app/commands"
	"github.com/balaji01-4d/pgxcli/internal/config"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

// Prompt rendering defaults for the interactive reader.
const (
	DefaultPrompt = `\u@\h:\d> `
	MaxLenPrompt  = 30
)

var builtinsCommand = map[string]func(){
	"clear": commands.ClearScreen,
}

var chromaFormatter = detectTerminalColorProfile()

// Reader abstracts user input, history, and completion behavior for the REPL.
type Reader interface {
	Read(ctx context.Context, prefix string) (string, error)
	History() []prompt.HistoryCommand
	SetAutocompleter(keywords []string)
}

type pgxReader struct {
	prompt prompt.Prompter
}

func newReader() (*pgxReader, error) {
	p, err := prompt.New()
	if err != nil {
		return nil, err
	}

	return &pgxReader{prompt: p}, nil
}

func (r *pgxReader) Read(ctx context.Context, prefix string) (string, error) {
	r.prompt.SetPrefix(prefix)
	text, err := r.prompt.Prompt(ctx)
	if err != nil && !errors.Is(err, prompt.ErrAborted) {
		return "", err
	}
	return text, nil
}

// SetAutocompleter configures SQL keyword suggestions for interactive input.
func (r *pgxReader) SetAutocompleter(keywords []string) {
	suggestions := make([]prompt.Suggestion, len(keywords))
	for i, kw := range keywords {
		suggestions[i] = prompt.Suggestion{Value: kw}
	}
	r.prompt.SetAutoCompleterContextual(prompt.AutoCompleteSimple(suggestions, true))
}

// History returns the current prompt history entries.
func (r *pgxReader) History() []prompt.HistoryCommand {
	return r.prompt.History()
}

func getSyntaxHighlighting(style string) (prompt.SyntaxHighlighter, error) {
	return prompt.SyntaxHighlighterChroma("PostgreSQL SQL dialect", chromaFormatter, style)
}

func applyReaderOptions(r *pgxReader, config *config.Config, histories []prompt.HistoryCommand) error {
	highlighter, err := getSyntaxHighlighting(config.Main.Style)
	if err != nil {
		return err
	}
	r.prompt.SetSyntaxHighlighter(highlighter)
	r.prompt.SetHistory(histories)
	return nil
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
		return "noop" // Chroma's no-op formatter
	}
}
