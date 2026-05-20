package renderer

import (
	"testing"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		in    config.TableColor
		field ColorField
		want  color.Attribute
	}{
		{
			name:  "default header uses cyan",
			in:    config.FgDefault,
			field: ColorHeader,
			want:  color.FgCyan,
		},
		{
			name:  "default column uses white",
			in:    config.FgDefault,
			field: ColorColumn,
			want:  color.FgWhite,
		},
		{
			name:  "known explicit color uses mapped value",
			in:    config.FgHiMagenta,
			field: ColorCaption,
			want:  color.FgHiMagenta,
		},
		{
			name:  "unknown color falls back to field default",
			in:    config.TableColor("unknown"),
			field: ColorHeader,
			want:  color.FgCyan,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := resolveColor(tc.in, tc.field)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDefaultColorForField(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		field ColorField
		want  color.Attribute
	}{
		{
			name:  "header",
			field: ColorHeader,
			want:  color.FgCyan,
		},
		{
			name:  "column",
			field: ColorColumn,
			want:  color.FgWhite,
		},
		{
			name:  "caption",
			field: ColorCaption,
			want:  color.FgWhite,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, defaultColorForField(tc.field))
		})
	}
}

func TestResolveStyle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   config.TableStyle
		want tw.BorderStyle
	}{
		{
			name: "known style",
			in:   config.StyleRounded,
			want: tw.StyleRounded,
		},
		{
			name: "unknown style falls back to default",
			in:   config.TableStyle("unknown"),
			want: tw.StyleDefault,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, resolveStyle(tc.in))
		})
	}
}

func TestGetTableStyle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		headerColor config.TableColor
		columnColor config.TableColor
		tableStyle  config.TableStyle
		wantHeader  color.Attribute
		wantColumn  color.Attribute
		wantSymbols tw.Symbols
	}{
		{
			name:        "uses explicit config values",
			headerColor: config.FgGreen,
			columnColor: config.FgYellow,
			tableStyle:  config.StyleDouble,
			wantHeader:  color.FgGreen,
			wantColumn:  color.FgYellow,
			wantSymbols: tw.NewSymbols(tw.StyleDouble),
		},
		{
			name:        "unknown values fall back to defaults",
			headerColor: config.TableColor("bad-color"),
			columnColor: config.TableColor("bad-color"),
			tableStyle:  config.TableStyle("bad-style"),
			wantHeader:  color.FgCyan,
			wantColumn:  color.FgWhite,
			wantSymbols: tw.NewSymbols(tw.StyleDefault),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{}
			cfg.Table.Color.Header = tc.headerColor
			cfg.Table.Color.Column = tc.columnColor
			cfg.Table.Style = tc.tableStyle

			got := GetTableStyle(cfg)
			require.Len(t, got.Header.FG, 1)
			require.Len(t, got.Column.FG, 1)
			assert.Equal(t, tc.wantHeader, got.Header.FG[0])
			assert.Equal(t, tc.wantColumn, got.Column.FG[0])
			assert.Equal(t, tc.wantSymbols, got.Symbols)
		})
	}
}
