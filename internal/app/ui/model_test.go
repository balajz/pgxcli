package ui

import (
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/require"
)

func TestModel_StaleSequenceTimeoutDoesNotResetExecutingState(t *testing.T) {
	m := &Model{
		state:      StateExecuting,
		pendingSeq: 7,
	}

	m.Update(seqTimeoutMsg{seq: 7})

	require.Equal(t, StateExecuting, m.state)
}

func TestModel_SequenceTimeoutResetsPendingState(t *testing.T) {
	m := &Model{
		state:      StatePendingClear,
		pendingSeq: 7,
	}

	m.Update(seqTimeoutMsg{seq: 7})

	require.Equal(t, StateInput, m.state)
}

func TestModel_SecondClearInvalidatesSequenceTimeout(t *testing.T) {
	m, err := New("", filepath.Join(t.TempDir(), "history"), "", "test", nil, nil, nil)
	require.NoError(t, err)
	m.state = StatePendingClear
	m.pendingSeq = 7
	m.statusModel.Message = "Press ESC again to clear"

	_, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})

	require.Equal(t, StateInput, m.state)
	require.NotEqual(t, 7, m.pendingSeq)
	require.Empty(t, m.statusModel.Message)
	m.Update(seqTimeoutMsg{seq: 7})
	require.Equal(t, StateInput, m.state)
}

func TestModel_AbandoningPendingSequenceResetsStatus(t *testing.T) {
	m, err := New("", filepath.Join(t.TempDir(), "history"), "", "test", nil, nil, nil)
	require.NoError(t, err)
	m.state = StatePendingQuit
	m.pendingSeq = 7
	m.statusModel.Message = "Press Ctrl+C again to quit"

	_, _ = m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})

	require.Equal(t, StateInput, m.state)
	require.NotEqual(t, 7, m.pendingSeq)
	require.Empty(t, m.statusModel.Message)
	m.Update(seqTimeoutMsg{seq: 7})
	require.Equal(t, StateInput, m.state)
}
