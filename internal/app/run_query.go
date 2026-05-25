package app

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/balajz/pgxcli/internal/app/renderer"
	"github.com/balajz/pgxcli/internal/database"
)

func (p *pgxCLI) handleQueryResult(r database.Rows, execDuration time.Duration) (cmd tea.Cmd, err error) {
	var s strings.Builder

	cols := renderer.GetColumnStrings(r, true)
	if len(cols) > 0 {
		rowIter := renderer.NewRowIter(r, true)
		if err := renderer.TableRender(cols, rowIter, "", &s, &s, p.config); err != nil {
			r.Close() // Ensure closed on error
			return nil, err
		}
	}

	// We must close the rows before reading the tag
	if closeErr := r.Close(); closeErr != nil {
		return nil, closeErr
	}

	tag, err := r.Tag()
	if err != nil {
		return nil, err
	}
	tagStr := tag.String()
	if tagStr == "" {
		tagStr = "OK"
	}

	output := s.String()
	if len(cols) == 0 {
		output = tagStr
	} else {
		output += tagStr
	}

	// Append timing info to the output
	timingInfo := fmt.Sprintf("\nTime %.3fs", execDuration.Seconds())
	output += timingInfo

	return p.printViaPager(output), nil
}
