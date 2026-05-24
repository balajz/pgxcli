package dbcommands

import (
	"context"
	"os"
	"os/exec"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
	"github.com/google/shlex"
)

func init() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "\\!",
		Description:   "Execute a shell command.",
		Syntax:        "\\! command",
		Handler:       ShellCommand,
		CaseSensitive: true,
	})
}

func ShellCommand(ctx context.Context, db database.Queryer, args string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	parts, err := shlex.Split(args)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	return nil, err
}
