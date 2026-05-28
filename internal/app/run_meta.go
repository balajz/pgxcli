package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/balajz/pgxcli/internal/app/renderer"
	"github.com/balajz/pgxcli/internal/app/ui"
	"github.com/balajz/pgxcli/internal/database"
	"github.com/balajz/pgxcli/internal/database/result"
	"github.com/balajz/pgxcli/pgxspecial"
)

func (p *pgxCLI) handleSpecialCommand(ctx context.Context, res pgxspecial.SpecialCommandResult, client *database.Client, execTime time.Duration) tea.Cmd {
	timingInfo := fmt.Sprintf("\nTime %.3fs", execTime.Seconds())

	switch msg := res.(type) {

	case database.ExitAction:
		return func() tea.Msg {
			return ui.QuitRequestMsg{}
		}

	case database.ChangeDbAction:
		if msg.Name != "" {
			if err := client.ChangeDatabase(ctx, msg.Name); err != nil {
				return p.withPrompt(p.printError(err))
			}
		}
		out := fmt.Sprintf(
			"You are now connected to database %q as user %q",
			client.GetDatabase(),
			client.GetUser(),
		)

		go func() {
			if err := p.client.Cache(p.compWorker); err != nil {
				p.logger.Error("failed to refresh database cache for completion", "error", err)
			}
		}()

		return p.withPrompt(p.printViaPager(out + timingInfo))

	case database.ConnInfoAction:
		var host string
		if strings.HasPrefix(client.GetHost(), "/") {
			host = fmt.Sprintf("Socket %q", client.GetHost())
		} else {
			host = fmt.Sprintf("Host %q", client.GetHost())
		}

		var port string
		if client.GetPort() == 0 {
			port = "None"
		} else {
			port = strconv.Itoa(int(client.GetPort()))
		}

		out := fmt.Sprintf(
			"You are connected to database %q as user %q on %s at port %s",
			client.GetDatabase(), client.GetUser(), host, port,
		)
		return p.withPrompt(p.printViaPager(out + timingInfo))

	case result.SpecialRow:
		table, err := renderer.RowsResult(res, p.config)
		if err != nil {
			return p.withPrompt(p.printError(err))
		}
		return p.withPrompt(p.printViaPager(table + timingInfo))

	case pgxspecial.DescribeTableListResult:
		tables, err := renderer.DescribeTableResult(res, p.config)
		if err != nil {
			return p.withPrompt(p.printError(err))
		}
		return p.withPrompt(p.printViaPager(tables + timingInfo))

	case pgxspecial.ExtensionVerboseListResult:
		tables, err := renderer.ExtensionVerboseResult(res, p.config)
		if err != nil {
			return p.withPrompt(p.printError(err))
		}
		return p.withPrompt(p.printViaPager(tables + timingInfo))

	default:
		return p.nextPrompt()
	}
}
