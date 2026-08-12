package formatter

import (
	"encoding/csv"
	"io"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/balajz/pgxcli/internal/perrors"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type TableFormatter struct {
	rows         int
	table        *tablewriter.Table
	streamWriter *csv.Writer

	tableConfig *config.TableConfig
}

func NewTableFormatter(w io.Writer, tableConfig *config.TableConfig) *TableFormatter {
	t := tablewriter.NewTable(w,
		tablewriter.WithRenderer(renderer.NewColorized(GetTableStyle(tableConfig))),
	)
	return &TableFormatter{
		table:       t,
		tableConfig: tableConfig,
	}
}

func (p *TableFormatter) Column(_ io.Writer, cols []string) error {
	p.table.Header(cols)
	return nil
}

func (p *TableFormatter) Iter(_, ew io.Writer, row []string) error {
	if p.streamWriter != nil {
		if err := p.streamWriter.Write(row); err != nil {
			return perrors.Wrap(err, perrors.WithMessage("failed to stream row"))
		}
		p.streamWriter.Flush()
		if err := p.streamWriter.Error(); err != nil {
			return perrors.Wrap(err, perrors.WithMessage("failed to flush streamed row"))
		}
		p.rows++
		return nil
	}

	if p.table == nil {
		return nil
	}

	if err := p.table.Append(row); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to append row to table"))
	}

	p.rows++
	return nil
}

// StartStreaming abandons the buffered table once it grows too large and
// emits the retained rows as tab-separated records. This keeps the query path
// bounded while preserving all values, including values containing tabs or
// newlines through csv quoting.
func (p *TableFormatter) StartStreaming(w io.Writer, cols []string, bufferedRows [][]string) error {
	p.table = nil
	p.streamWriter = csv.NewWriter(w)
	p.streamWriter.Comma = '\t'
	if err := p.streamWriter.Write(cols); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to stream header"))
	}
	for _, row := range bufferedRows {
		if err := p.streamWriter.Write(row); err != nil {
			return perrors.Wrap(err, perrors.WithMessage("failed to stream buffered row"))
		}
	}
	p.streamWriter.Flush()
	if err := p.streamWriter.Error(); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to flush streamed table"))
	}
	return nil
}

func (p *TableFormatter) Caption(w io.Writer, caption string) error {
	if p.table == nil {
		return nil
	}

	captionColor := getCaptionColor(p.tableConfig.Color.Caption)
	CC := tw.Caption{
		Text: color.New(captionColor).Sprint(caption),
		Spot: tw.SpotBottomLeft,
	}

	p.table.Caption(CC)
	return nil
}

func (p *TableFormatter) Render(_ io.Writer, _ int) error {
	if p.streamWriter != nil {
		p.streamWriter.Flush()
		if err := p.streamWriter.Error(); err != nil {
			return perrors.Wrap(err, perrors.WithMessage("failed to flush streamed table"))
		}
		return nil
	}
	if err := p.table.Render(); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to render table"))
	}
	return nil
}

func (p *TableFormatter) Done(_ io.Writer) error {
	p.table = nil
	p.streamWriter = nil
	return nil
}
