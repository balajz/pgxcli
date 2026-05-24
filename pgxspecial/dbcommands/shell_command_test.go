package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
)

func TestShellCommand_Success(t *testing.T) {
	ctx := context.Background()

	_, err := dbcommands.ShellCommand(ctx, nil, "echo hello", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestShellCommand_Failure(t *testing.T) {
	ctx := context.Background()

	// invalid command should error
	_, err := dbcommands.ShellCommand(ctx, nil, "commandthatdoesnotexist", false)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
