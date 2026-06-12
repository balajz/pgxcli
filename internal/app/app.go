// Package app is the application layer.
//
// It orchestrates the repl loop, command routing, result handling and history management.
// custom commands and edit command buffer are also handled here.
package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/balajz/pgxcli/internal/app/ui"
	"github.com/balajz/pgxcli/internal/cliio"
	"github.com/balajz/pgxcli/internal/config"
	"github.com/balajz/pgxcli/internal/database"
	"github.com/balajz/pgxcli/internal/perrors"
	compDB "github.com/balajz/pgxls/pkg/database"
)

// Application defines the interface for the main application logic.
type Application interface {
	// Start starts the main repl loop, reading input, executing commands and printing results until the user exits.
	Start(ctx context.Context) error

	// Close performs saving history before exiting.
	Close() error
}

var builtinsCommand = map[string]func() tea.Cmd{
	"\\clear": func() tea.Cmd { return tea.ClearScreen },
}

// pgxCLI is the main implementation of the Application interface.
type pgxCLI struct {
	model      *ui.Model
	program    *tea.Program
	Printer    cliio.Printer
	config     *config.Config
	logger     *slog.Logger
	client     *database.Client
	compWorker *compDB.Worker

	version string
}

func New(cfg *config.Config, printer cliio.Printer, logger *slog.Logger, client *database.Client, version string) (Application, error) {
	compWorker := compDB.NewWorker()
	compWorker.Start()

	return &pgxCLI{
		config:     cfg,
		logger:     logger,
		Printer:    printer,
		client:     client,
		version:    version,
		compWorker: compWorker,
	}, nil
}

func (p *pgxCLI) execute(ctx context.Context, query string) tea.Cmd {
	p.logger.Debug("received command", "command_length", len(query))

	if cmd, ok := builtinsCommand[query]; ok {
		p.logger.Debug("executing builtin command", "command", query)
		return p.withPrompt(cmd())
	}

	return func() tea.Msg {
		if msg, ok := p.runMeta(ctx, query); ok {
			return msg
		}

		return p.runQuery(ctx, query)
	}
}

func (p *pgxCLI) Start(ctx context.Context) error {
	p.printBanner(p.version)
	executeFunc := func(query string) tea.Cmd {
		return p.execute(ctx, query)
	}

	initialPrefix := p.client.ParsePrompt(p.config.Main.Prompt)
	m, err := ui.New(
		initialPrefix,
		p.config.Main.HistoryFile,
		string(p.config.Main.Style),
		p.version,
		executeFunc,
		p.Cancel,
		p.getCompletions(),
	)
	if err != nil {
		return err
	}

	p.model = m
	p.program = tea.NewProgram(p.model, tea.WithContext(ctx))

	if _, err := p.program.Run(); err != nil {
		return perrors.Wrap(
			err,
			perrors.WithMessage("failed to start UI"),
		)
	}

	return nil
}

func (p *pgxCLI) Cancel(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.client.Cancel(ctx)
}

func (p *pgxCLI) printBanner(version string) {
	lipgloss.Fprint(os.Stdout, ui.Banner(version)+"\n")
}

func (p *pgxCLI) nextPrompt() tea.Cmd {
	return func() tea.Msg {
		prefix := p.client.ParsePrompt(p.config.Main.Prompt)
		return ui.ReadyMsg{Prefix: prefix}
	}
}

func (p *pgxCLI) withPrompt(cmds ...tea.Cmd) tea.Cmd {
	all := make([]tea.Cmd, len(cmds)+1)
	copy(all, cmds)
	all[len(cmds)] = p.nextPrompt()
	return tea.Sequence(all...)
}

func (p *pgxCLI) printViaPager(str string) tea.Cmd {
	if p.Printer.ShouldUsePager(str) {
		cmd, ok := cliio.PagerCmd(str)
		if !ok {
			return ui.PrintCmd(str, ui.DefaultStyles().AppOutput)
		}
		return ui.ShowPagerCmd(cmd)
	}

	return ui.PrintCmd(str, ui.DefaultStyles().AppOutput)
}

func (p *pgxCLI) printError(err error) tea.Cmd {
	return ui.PrintErrCmd(err, ui.DefaultStyles().ErrorOutput)
}

func (p *pgxCLI) Close() error {
	p.logger.Info("closing application and saving history")
	if p.compWorker != nil {
		p.compWorker.Stop()
	}

	if p.model != nil {
		return p.model.Close()
	}
	return nil
}
