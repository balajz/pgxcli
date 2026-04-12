package cli

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/balaji01-4d/pgxcli/internal/app"
	"github.com/balaji01-4d/pgxcli/internal/config"
	"github.com/balaji01-4d/pgxcli/internal/database"
	"github.com/balaji01-4d/pgxcli/internal/logger"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func NewRootCmd(ctx context.Context, cliCtx *CliContext) *cobra.Command {
	var debugFlag debugFlag
	var hostFlag hostFlag
	var portFlag portFlag
	var dbNameFlag dbNameFlag
	var usernameFlag usernameFlag
	var neverPromptFlag neverPromptFlag
	var forcePromptFlag forcePromptFlag

	rootCmd := &cobra.Command{
		Use:     "pgxcli [DBNAME] [USERNAME]",
		Short:   "Interactive PostgreSQL command-line client for querying and managing databases.",
		Version: version,
		Args:    cobra.MaximumNArgs(2), // Database name and username are optional example: pgxcli mydb myuser

		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cliCtx.config = cfg

			logger, err := logger.InitLogger(bool(debugFlag), cfg.Main.LogFile)
			if err != nil {
				return err
			}

			cliCtx.Logger = logger

			return nil
		},

		PreRunE: func(cmd *cobra.Command, args []string) error {
			var argDB string
			var argUser string
			if len(args) > 0 {
				argDB = args[0]
			}
			if len(args) > 1 {
				argUser = args[1]
			}
			dbName, userName := resolveDBAndUser(string(dbNameFlag), string(usernameFlag), argDB, argUser)
			if userName != "" {
				userName = getUserFromEnv()
				if userName == "" {
					currentUser, err := user.Current()
					if err != nil {
						return err
					}
					userName = currentUser.Username
					if strings.Contains(userName, "\\") {
						userName = userName[strings.LastIndex(userName, "\\")+1:]
					}
				}
			}
			if dbName == "" {
				dbName = userName
			}

			postgres := database.New(cliCtx.Logger.Logger)
			cliCtx.Client = postgres

			var connector database.Connector
			var connErr error
			if strings.Contains(dbName, "://") || strings.Contains(dbName, "=") {
				connector, connErr = database.NewPGConnectorFromConnString(dbName)
				if connErr != nil {
					cliCtx.Logger.Error("Invalid Connection string", "error", connErr)
					return connErr
				}
			} else {
				cliCtx.Logger.Debug("using field-based connection",
					"host", hostFlag,
					"port", portFlag,
					"database", dbName,
					"user", userName,
				)
				var password string
				if bool(neverPromptFlag) {
					password = getPasswordFromEnv()
				}

				if bool(forcePromptFlag) && password == "" {
					// Force prompt for password
					// TODO: Implement secure passowrd input
					var pwd string
					fmt.Print("Password: ")
					_, err := fmt.Scanln(&pwd)
					if err != nil {
						return err
					}
					password = pwd
				}

				connector, connErr = database.NewPGConnectorFromFields(
					string(hostFlag),
					dbName,
					userName,
					password,
					uint16(portFlag),
				)
				if connErr != nil {
					cliCtx.Logger.Error("Failed to create connector", "error", connErr)
					return connErr
				}

				cliCtx.Logger.Debug("Attempting database connection")
				connErr = cliCtx.Client.Connect(ctx, connector)
				if connErr != nil {
					if shouldAskForPassword(connErr, bool(neverPromptFlag)) {
						cliCtx.Logger.Debug("Connection failed, prompting for password")
						var pwd string
						fmt.Print("Password: ")
						_, err := fmt.Scanln(&pwd)
						if err != nil {
							return err
						}
						connector.UpdatePassword(pwd)
						connRetryErr := cliCtx.Client.Connect(ctx, connector)
						if connRetryErr != nil {
							cliCtx.Logger.Error("Connection retry failed", "error", connRetryErr)
							return connRetryErr
						}
					} else {
						cliCtx.Logger.Error("Failed to connect to database", "error", connErr)
						return connErr
					}
				}
			}
			if !cliCtx.Client.IsConnected() {
				err := fmt.Errorf("failed to connect to database")
				cliCtx.Logger.Error("Failed to connect to database", "error", err)
				return err
			}

			if err := cliCtx.Client.Ping(ctx); err != nil {
				cliCtx.Logger.Error("Failed to ping database", "error", err)
				return err
			}

			app, err := app.NewPgxCLI(cliCtx.config, cliCtx.Printer, cliCtx.Client.Logger)
			if err != nil {
				cliCtx.Logger.Error("Failed to initialize app", "error", err)
				return err
			}
			cliCtx.App = app
			return nil
		},

		RunE: func(_ *cobra.Command, _ []string) error {
			if cliCtx.App == nil {
				cliCtx.Logger.Error("Application context not initialized")
				return fmt.Errorf("application context not initialized")
			}
			cliCtx.App.Start(ctx, cliCtx.Client)
			return nil
		},

		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cliCtx.Logger != nil {
				if err := cliCtx.Logger.Close(); err != nil {
					return err
				}
			}
			if cliCtx.Client != nil {
				if err := cliCtx.Client.Close(ctx); err != nil {
					return err
				}
			}
			if cliCtx.App != nil {
				if err := cliCtx.App.Close(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	// deactivating of the -h shorthand flag, so that it can be used in the host flag
	rootCmd.PersistentFlags().BoolP("help", "", false, "Print usage")
	_ = rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "use --help")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true

	debugFlag.bind(rootCmd)
	hostFlag.bind(rootCmd)
	portFlag.bind(rootCmd)
	dbNameFlag.bind(rootCmd)
	usernameFlag.bind(rootCmd)
	neverPromptFlag.bind(rootCmd)
	forcePromptFlag.bind(rootCmd)

	return rootCmd
}

// getUserFromEnv gets username from environment variables
// support for pgcli specific environment variable
func getUserFromEnv() string {
	if userEnv := os.Getenv("PGXUSER"); userEnv != "" {
		return userEnv
	}
	if userEnv := os.Getenv("PGUSER"); userEnv != "" {
		return userEnv
	}
	return ""
}

func getPasswordFromEnv() string {
	if passEnv := os.Getenv("PGXPASSWORD"); passEnv != "" {
		return passEnv
	}
	if passEnv := os.Getenv("PGPASSWORD"); passEnv != "" {
		return passEnv
	}
	return ""
}

// when database is given as flag then the next argument as user
func resolveDBAndUser(dbnameOpt, userOpt, argDB, argUser string) (string, string) {
	// Case: cmd -d database user
	if dbnameOpt != "" && argDB != "" && argUser == "" {
		return dbnameOpt, argDB
	}

	database := firstNonEmpty(dbnameOpt, argDB)
	user := firstNonEmpty(userOpt, argUser)
	return database, user
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
