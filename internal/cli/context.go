// Package cli contains the command-line interface for pgxcli.
package cli

import (
	"github.com/balajz/pgxcli/internal/app"
	"github.com/balajz/pgxcli/internal/cliio"
	"github.com/balajz/pgxcli/internal/config"
	"github.com/balajz/pgxcli/internal/database"
	"github.com/balajz/pgxcli/internal/logger"
)

// CliContext holds the dependencies for cli.
type CliContext struct { //revive:disable suggested context name would be misunderstood to context.Context
	// config holds the global configuration for the cli
	config *config.Config

	// Logger is used for logging messages and errors
	Logger *logger.Logger

	// Printer is used for outputting messages to the user
	Printer cliio.Printer

	// Client is the database client used to interact with the Postgres database
	Client *database.Client

	// App is the application layer that contains the business logic of pgxcli
	// App orchestrates the execution of commands and interacts with the database client
	// printer to perform operations and display results.
	App app.Application
}
