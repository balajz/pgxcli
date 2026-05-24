package formatter

import (
	"io"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type TableFormatter struct {
	rows  int
	table *tablewriter.Table

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
	if p.table == nil {
		return nil
	}

	if err := p.table.Append(row); err != nil {
		return err
	}

	p.rows++
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
	return p.table.Render()
}

func (p *TableFormatter) Done(_ io.Writer) error {
	p.table = nil
	return nil
}
