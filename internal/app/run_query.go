package app

import (
	"context"
	"fmt"
	"io"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/balajz/pgxcli/internal/app/renderer"
	"github.com/balajz/pgxcli/internal/app/ui"
	"github.com/balajz/pgxcli/internal/config"
	"github.com/balajz/pgxcli/internal/database"
	"github.com/balajz/pgxcli/internal/parser"
)

func (p *pgxCLI) runQuery(ctx context.Context, query string) tea.Msg {
	p.logger.Debug("executing query")
	stmts := parser.SplitSQLStatements(query)
	cmds := make([]tea.Cmd, 0, len(stmts))

StatementsLoop:
	for _, stmt := range stmts {
		p.logger.Debug("parsed statement", "statement", stmt)
		if stmt == "" || stmt == ";" {
			continue
		}

		start := time.Now()
		queryResult, _, err := p.client.ExecuteQuery(ctx, stmt, false)
		execDuration := time.Since(start)
		if err != nil {
			p.logger.Error("query execution failed", "error", err)
			cmds = append(cmds, p.printError(err))
			if p.config.Main.OnError == config.OnErrorStop {
				break StatementsLoop
			}
			continue
		}
		resultCmd, err := p.handleQueryResult(queryResult, execDuration)
		if err != nil {
			p.logger.Error("error handling query result", "error", err)
			cmds = append(cmds, p.printError(err))
			if p.config.Main.OnError == config.OnErrorStop {
				break StatementsLoop
			}
			continue
		}
		cmds = append(cmds, resultCmd)
	}

	return ui.ExecCmdMsg{Cmd: p.withPrompt(cmds...)}
}

func (p *pgxCLI) handleQueryResult(r database.Rows, execDuration time.Duration) (cmd tea.Cmd, err error) {
	cols := renderer.GetColumnStrings(r, true)
	return func() tea.Msg {
		streamErr := p.Printer.StreamViaPager(func(w io.Writer) error {
			if len(cols) > 0 {
				rowIter := renderer.NewRowIter(r, true)
				if err := renderer.TableRender(cols, rowIter, "", w, w, p.config); err != nil {
					_ = r.Close()
					return err
				}
			}

			// We must close the rows before reading the tag.
			if closeErr := r.Close(); closeErr != nil {
				return closeErr
			}

			tag, err := r.Tag()
			if err != nil {
				return err
			}
			tagStr := tag.String()
			if tagStr == "" {
				tagStr = "OK"
			}
			if _, err := fmt.Fprintln(w, tagStr); err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "Time %.3fs\n", execDuration.Seconds())
			return err
		})
		if streamErr != nil {
			p.logger.Error("error streaming query result", "error", streamErr)
			return ui.PrintErrCmd(streamErr, ui.DefaultStyles().ErrorOutput)
		}
		return nil
	}, nil
}
