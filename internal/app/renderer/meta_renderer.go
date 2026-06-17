package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/balajz/pgxcli/internal/perrors"
	"github.com/balajz/pgxcli/pgxspecial"
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

// Table renders the provided data as a table.
// Currently been used by the special results / meta commands only.
func Table(data Data, w io.Writer, c *config.Config) error {
	t := tablewriter.NewTable(w, tablewriter.WithRenderer(renderer.NewColorized(GetTableStyle(c))))
	rows, err := data.Rows()
	if err != nil {
		return err
	}

	t.Header(data.Columns())
	if err := t.Bulk(rows); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to bulk append rows to table"))
	}

	if captionText := data.Caption(); captionText != "" {
		captionColor := getCaptionColor(c.Table.Color.Caption)
		caption := tw.Caption{
			Text: color.New(captionColor).Sprint(captionText),
			Spot: tw.SpotBottomLeft,
		}
		t.Caption(caption)
	}
	if err := t.Render(); err != nil {
		return perrors.Wrap(err, perrors.WithMessage("failed to render table"))
	}
	return nil
}

type rowsTableResult interface {
	pgxspecial.SpecialCommandResult
	Columns() []string
	Data() [][]any
}

type staticData struct {
	columns []string
	rows    [][]any
	caption string
}

func (d staticData) Columns() []string {
	return d.columns
}

func (d staticData) Rows() ([][]any, error) {
	return d.rows, nil
}

func (d staticData) Caption() string {
	return d.caption
}

// RowsResult renders row-based special command output.
func RowsResult(result pgxspecial.SpecialCommandResult, c *config.Config) (string, error) {
	resultRows, ok := result.(rowsTableResult)
	if !ok {
		return "", fmt.Errorf("invalid row result type")
	}

	return renderData(staticData{columns: resultRows.Columns(), rows: resultRows.Data()}, c)
}

// DescribeTableResult renders each describe-table section.
func DescribeTableResult(result pgxspecial.SpecialCommandResult, c *config.Config) (string, error) {
	describeTableResult, ok := result.(pgxspecial.DescribeTableListResult)
	if !ok {
		return "", fmt.Errorf("invalid describe table result type")
	}

	out := make([]string, 0, len(describeTableResult.Results))

	for _, tableDesc := range describeTableResult.Results {
		rendered, err := renderTableDescription(tableDesc, c)
		if err != nil {
			return "", err
		}
		out = append(out, rendered)
	}
	return strings.Join(out, "\n"), nil
}

// ExtensionVerboseResult renders each verbose extension result.
func ExtensionVerboseResult(result pgxspecial.SpecialCommandResult, c *config.Config) (string, error) {
	extResult, ok := result.(pgxspecial.ExtensionVerboseListResult)
	if !ok {
		return "", fmt.Errorf("invalid extension verbose result type")
	}
	out := make([]string, 0, len(extResult.Results))

	for _, ext := range extResult.Results {
		rendered, err := renderExtensionVerbose(ext, c)
		if err != nil {
			return "", err
		}
		out = append(out, rendered)
	}
	return strings.Join(out, "\n"), nil
}

func renderExtensionVerbose(ext pgxspecial.ExtensionVerboseResult, c *config.Config) (string, error) {
	rows := make([][]any, 0, len(ext.Description))
	for _, objDesc := range ext.Description {
		rows = append(rows, []any{objDesc})
	}

	return renderData(staticData{
		columns: []string{"Object Description"},
		rows:    rows,
		caption: ext.Name,
	}, c)
}

func renderTableDescription(result pgxspecial.DescribeTableResult, c *config.Config) (string, error) {
	rows := make([][]any, 0, len(result.Data))
	for _, values := range result.Data {
		row := make([]any, len(values))
		for i, v := range values {
			row[i] = v
		}
		rows = append(rows, row)
	}

	return renderData(staticData{
		columns: result.Columns,
		rows:    rows,
		caption: renderTableFooter(result.TableMetaData),
	}, c)
}

func renderTableFooter(meta pgxspecial.TableFooterMeta) string {
	var sb strings.Builder

	writeList := func(title string, v []string) {
		if len(v) == 0 {
			return
		}
		sb.WriteString(title)
		sb.WriteByte('\n')
		for _, s := range v {
			sb.WriteString("    ")
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
	}

	writeValue := func(title string, v *string) {
		if v == nil || *v == "" {
			return
		}
		sb.WriteString(title)
		sb.WriteString(*v)
		sb.WriteByte('\n')
	}

	writeBool := func(title string, v *bool) {
		if v == nil {
			return
		}
		sb.WriteString(title)
		if *v {
			sb.WriteString("yes\n")
		} else {
			sb.WriteString("no\n")
		}
	}

	writeList("Indexes:", meta.Indexes)
	writeList("Check constraints:", meta.CheckConstraints)
	writeList("Foreign-key constraints:", meta.ForeignKeys)
	writeList("Referenced by:", meta.ReferencedBy)
	writeValue("View definition:\n", meta.ViewDefinition)

	writeList("Rules:", meta.RulesEnabled)
	writeList("Disabled rules:", meta.RulesDisabled)
	writeList("Rules firing always:", meta.RulesAlways)
	writeList("Rules firing on replica only:", meta.RulesReplica)

	writeList("Triggers:", meta.TriggersEnabled)
	writeList("Disabled triggers:", meta.TriggersDisabled)
	writeList("Triggers firing always:", meta.TriggersAlways)
	writeList("Triggers firing on replica only:", meta.TriggersReplica)

	writeList("Partition of:", meta.PartitionOf)
	writeList("Partition constraint:", meta.PartitionConstraints)
	writeValue("Partition key: ", meta.PartitionKey)
	writeList("Partitions:", meta.Partitions)
	writeValue("", meta.PartitionsSummary)

	writeList("Inherits:", meta.Inherits)
	writeList("Child tables:", meta.ChildTables)
	writeValue("", meta.ChildTablesSummary)
	writeValue("Typed table of type: ", meta.TypedTableOf)
	writeBool("Has OIDs: ", meta.HasOIDs)
	writeValue("Options: ", meta.Options)
	writeValue("Server: ", meta.Server)
	writeValue("FDW Options: ", meta.FDWOptions)
	writeValue("Owned by: ", meta.OwnedBy)

	return sb.String()
}

func renderData(data Data, c *config.Config) (string, error) {
	var sb strings.Builder
	if err := Table(data, &sb, c); err != nil {
		return "", err
	}
	return sb.String(), nil
}
