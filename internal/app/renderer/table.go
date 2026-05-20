package renderer

import (
	"io"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type Data interface {
	Columns() []string
	Rows() ([][]any, error)
	Caption() string
}

func Table(data Data, w io.Writer, c *config.Config) error {
	t := tablewriter.NewTable(w, tablewriter.WithRenderer(renderer.NewColorized(GetTableStyle(c))))
	rows, err := data.Rows()
	if err != nil {
		return err
	}

	t.Header(data.Columns())
	if err := t.Bulk(rows); err != nil {
		return err
	}

	if captionText := data.Caption(); captionText != "" {
		captionColor := getCaptionColor(c.Table.Color.Caption)
		caption := tw.Caption{
			Text: color.New(captionColor).Sprint(captionText),
			Spot: tw.SpotBottomLeft,
		}
		t.Caption(caption)
	}
	return t.Render()
}
