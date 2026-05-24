package pgxspecial

import (
	"context"
	"fmt"
	"strings"

	"github.com/balajz/pgxcli/pgxspecial/database"
)

// SpecialHandler defines the signature for a special command handler.
//
// A SpecialHandler is invoked with the parsed command arguments and an optional
// verbose flag. The provided db is used to execute queries as needed.
//
// The returned pgx.Rows, if non-nil, is passed back to the caller for consumption.
// Any error returned indicates command execution failure.
type SpecialHandler func(ctx context.Context, db database.Queryer, args string, verbose bool) (SpecialCommandResult, error)

// commandRegistry stores all registered special commands, indexed by command
// name and aliases. Command lookup is performed against this registry during
// special command execution.
var commandRegistry = map[string]SpecialCommand{}

// RegisterCommand registers a special command and its aliases in the command registry.
//
// The command name and aliases are normalized based on the CaseSensitive flag.
// If CaseSensitive is false, command keys are stored in lowercase, making lookup
// case-insensitive.
//
// If multiple aliases are provided, each alias is registered to reference the
// same command definition.
//
// RegisterCommand does not perform validation for duplicate command names or
// aliases; later registrations will overwrite existing entries with the same key.
func RegisterCommand(cmdRegistry SpecialCommandRegistry) {

	normalize := func(s string) string {
		if cmdRegistry.CaseSensitive {
			return s
		}
		return strings.ToLower(s)
	}

	cmd := SpecialCommand{
		Cmd:           cmdRegistry.Cmd,
		Description:   cmdRegistry.Description,
		Syntax:        cmdRegistry.Syntax,
		CaseSensitive: cmdRegistry.CaseSensitive,
		Handler:       cmdRegistry.Handler,
	}

	commandRegistry[normalize(cmdRegistry.Cmd)] = cmd

	for _, alias := range cmdRegistry.Alias {
		commandRegistry[normalize(alias)] = cmd
	}
}

// ExecuteSpecialCommand parses and executes a special command using the command registry.
// Syntax: \command[+] [args]
//
// \command - actual special command
//
// \command[+] - actual special command with verbose mode
//
// A special command is identified by a leading backslash (`\`). If the input does not
// start with a backslash, ExecuteSpecialCommand returns (nil, false, nil) to indicate
// that the input should be treated as a normal query.
//
// The first whitespace-delimited token is treated as the command name. If the command
// name ends with a plus sign (`+`), verbose mode is enabled and the suffix is removed
// before command lookup. The remaining input is passed to the command handler as
// arguments.
//
// The provided Queryer is used by the command handler to execute any required queries.
// Return values:
//   - SpecialCommandResult: the result returned by the command handler, if any
//   - bool: true if the input was recognized as a special command, even if execution failed
//   - error: non-nil if the command is unknown or execution fails
//
// An error is returned if the command is not found in the registry or if the command
// handler returns an error.
func ExecuteSpecialCommand(ctx context.Context, queryer database.Queryer, specialCommand string) (SpecialCommandResult, bool, error) {
	if !strings.HasPrefix(specialCommand, "\\") {
		return nil, false, nil
	}

	checkVerbose := func(cmd string) (string, bool) {
		suff := "+"
		return strings.TrimSuffix(cmd, suff), strings.HasSuffix(cmd, suff)
	}

	fields := strings.Fields(specialCommand)
	cmd := fields[0]
	args := strings.TrimSpace(strings.TrimPrefix(specialCommand, cmd))

	cmd, verbose := checkVerbose(cmd)

	command, ok := commandRegistry[cmd]
	if !ok {
		return nil, true, fmt.Errorf("unknown command: %s", cmd)
	}
	res, err := command.Handler(ctx, queryer, args, verbose)
	if err != nil {
		return nil, true, err
	}
	return res, true, nil
}
