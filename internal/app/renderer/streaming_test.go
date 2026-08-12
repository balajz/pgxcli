package renderer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/balajz/pgxcli/internal/app/renderer/formatter"
	"github.com/balajz/pgxcli/internal/config"
)

func TestTableRenderSwitchesToBoundedStreamingOutput(t *testing.T) {
	rows := make([][]string, bufferedRowLimit+1)
	for i := range rows {
		rows[i] = []string{fmt.Sprintf("row-%d", i), "value"}
	}

	var out strings.Builder
	err := Render(&out, &out, formatter.NewTableFormatter(&out, &config.TableConfig{}), []string{"name", "value"}, NewRowSliceIter(rows))
	if err != nil {
		t.Fatal(err)
	}

	result := out.String()
	if !strings.Contains(result, "row-0\tvalue\n") {
		t.Fatalf("expected streaming TSV output, got %q", result[len(result)-min(len(result), 200):])
	}
	if !strings.Contains(result, fmt.Sprintf("row-%d\tvalue\n", bufferedRowLimit)) {
		t.Fatalf("expected final streamed row record in output, got %q", result[len(result)-min(len(result), 200):])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
