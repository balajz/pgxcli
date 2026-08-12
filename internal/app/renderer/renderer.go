package renderer

import (
	"io"

	"github.com/balajz/pgxcli/internal/app/renderer/formatter"
	"github.com/balajz/pgxcli/internal/config"
)

type Formatter interface {
	Column(w io.Writer, cols []string) error
	Iter(w, Ew io.Writer, row []string) error
	Render(w io.Writer, seenRows int) error
	Done(w io.Writer) error
}

// StreamingFormatter can switch from table buffering to bounded row output
// once a result exceeds the in-memory formatting limit.
type StreamingFormatter interface {
	Formatter
	StartStreaming(w io.Writer, cols []string, bufferedRows [][]string) error
}

const bufferedRowLimit = 500

func TableRender(cols []string, rowrowIter RowStrIter, caption string, w, Ew io.Writer, c *config.Config) error {
	tf := formatter.NewTableFormatter(w, &c.Table)
	return Render(w, Ew, tf, cols, rowrowIter)
}

func Render(w, Ew io.Writer, formatter Formatter, cols []string, row RowStrIter) error {
	if err := formatter.Column(w, cols); err != nil {
		return err
	}

	streamingFormatter, canStream := formatter.(StreamingFormatter)
	bufferedRows := make([][]string, 0, bufferedRowLimit)
	streaming := false
	nRows := 0
	for {
		r, err := row.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if canStream && !streaming && len(bufferedRows) >= bufferedRowLimit {
			if err := streamingFormatter.StartStreaming(w, cols, bufferedRows); err != nil {
				return err
			}
			bufferedRows = nil
			streaming = true
		}

		if err := formatter.Iter(w, Ew, r); err != nil {
			return err
		}
		if canStream && !streaming {
			bufferedRows = append(bufferedRows, append([]string(nil), r...))
		}
		nRows++
	}

	return formatter.Render(w, nRows)
}
