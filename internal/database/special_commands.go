package database

import (
	"context"

	"github.com/balaji01-4d/pgxspecial"
	"github.com/balaji01-4d/pgxspecial/database"
	// Register built-in pgxspecial commands via package init side effects.
	_ "github.com/balaji01-4d/pgxspecial/dbcommands"
)

const (
	// Exit is the result kind for quit command actions.
	Exit pgxspecial.SpecialResultKind = 100 + iota
	// ChangeDB is the result kind for database switch command actions.
	ChangeDB
	// Conninfo is the result kind for connection info command actions.
	Conninfo
)

func init() {
	registerSpecialCommands()
}

func registerSpecialCommands() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:         "\\q",
		Syntax:      "\\q",
		Description: "Quit Pgxcli",
		Handler: func(_ context.Context, _ database.Queryer, _ string, _ bool) (pgxspecial.SpecialCommandResult, error) {
			return ExitAction{}, nil
		},
		CaseSensitive: true,
	})

	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:         "\\c",
		Syntax:      "\\c database_name",
		Description: "Change a new database",
		Handler: func(_ context.Context, _ database.Queryer, s string, _ bool) (pgxspecial.SpecialCommandResult, error) {
			return ChangeDbAction{Name: s}, nil
		},
		CaseSensitive: true,
		Alias:         []string{"\\connect"},
	})

	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:         "\\conninfo",
		Syntax:      "\\conninfo",
		Description: "Get connection details",
		Handler: func(_ context.Context, _ database.Queryer, _ string, _ bool) (pgxspecial.SpecialCommandResult, error) {
			return ConnInfoAction{}, nil
		},
		CaseSensitive: false,
	})
}

// ExitAction indicates that the REPL should terminate.
type ExitAction struct{}

// ResultKind returns the special result kind for ExitAction.
func (e ExitAction) ResultKind() pgxspecial.SpecialResultKind {
	return Exit
}

// ChangeDbAction carries target database name for \c / \connect commands.
type ChangeDbAction struct {
	Name string
}

// ResultKind returns the special result kind for ChangeDbAction.
func (c ChangeDbAction) ResultKind() pgxspecial.SpecialResultKind {
	return ChangeDB
}

// ConnInfoAction indicates that connection info should be displayed.
type ConnInfoAction struct{}

// ResultKind returns the special result kind for ConnInfoAction.
func (g ConnInfoAction) ResultKind() pgxspecial.SpecialResultKind {
	return Conninfo
}
