package app

import (
	"context"
	"errors"

	"github.com/balaji01-4d/pgxcli/internal/app/commands"
	"github.com/balaji01-4d/pgxcli/internal/config"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

const (
	DefaultPrompt = `\u@\h:\d> `
	MaxLenPrompt  = 30
)

var builtinsCommand = map[string]func(){
	"clear": commands.ClearScreen,
}

var chromaFormatter = detectTerminalColorProfile()

type Reader interface {
	Read(prefix string, ctx context.Context) (string, error)
	History() []prompt.HistoryCommand
	SetAutocompleter(keywords []string)
}

type PgxReader struct {
	prompt prompt.Prompter
}

func NewPgxReader() (*PgxReader, error) {
	p, err := prompt.New()
	if err != nil {
		return nil, err
	}

	return &PgxReader{prompt: p}, nil
}

func (r *PgxReader) Read(prefix string, ctx context.Context) (string, error) {
	r.prompt.SetPrefix(prefix)
	text, err := r.prompt.Prompt(ctx)
	if err != nil && !errors.Is(err, prompt.ErrAborted) {
		return "", err
	}
	return text, nil
}

func (r *PgxReader) SetAutocompleter(keywords []string) {
	suggestions := make([]prompt.Suggestion, len(keywords))
	for i, kw := range keywords {
		suggestions[i] = prompt.Suggestion{Value: kw}
	}
	r.prompt.SetAutoCompleterContextual(prompt.AutoCompleteSimple(suggestions, true))
}

func (r *PgxReader) History() []prompt.HistoryCommand {
	return r.prompt.History()
}

func getSyntaxHighlighting(style string) (prompt.SyntaxHighlighter, error) {
	return prompt.SyntaxHighlighterChroma("PostgreSQL SQL dialect", chromaFormatter, style)
}

func applyReaderOptions(r *PgxReader, config *config.Config, histories []prompt.HistoryCommand) error {
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
