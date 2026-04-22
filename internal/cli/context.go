// Package cli contains the command-line interface for pgxcli.
package cli

import (
	"github.com/balaji01-4d/pgxcli/internal/app"
	"github.com/balaji01-4d/pgxcli/internal/cliio"
	"github.com/balaji01-4d/pgxcli/internal/config"
	"github.com/balaji01-4d/pgxcli/internal/database"
	"github.com/balaji01-4d/pgxcli/internal/logger"
)

// CliContext holds the dependencies for cli.
type CliContext struct { //revive:disable suggested context name would be misunderstood to context.Context
	config  *config.Config
	Logger  *logger.Logger
	Printer cliio.Printer
	Client  *database.Client

	App app.Application
}
