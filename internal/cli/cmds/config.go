package cmds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/balajz/pgxcli/internal/app"
	"github.com/balajz/pgxcli/internal/config"
	"github.com/spf13/cobra"
)

func NewConfigCmd(ctx *app.AppContext) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "manage configuration",
	}

	configCmd.AddCommand(pathCmd(ctx))
	configCmd.AddCommand(checkCmd(ctx))
	return configCmd
}

func pathCmd(ctx *app.AppContext) *cobra.Command {
	pathCmd := &cobra.Command{
		Use:   "path",
		Short: "return the config path",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), ctx.Config.Path)
			return nil
		},
	}

	return pathCmd
}

func checkCmd(ctx *app.AppContext) *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "validate the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Validate(*ctx.Config)
			if err == nil {
				return nil
			}

			var sb strings.Builder

			if uw, ok := err.(interface { Unwrap() []error }); ok {
				fmt.Fprintf(&sb, "%s\n", "Configuration check failed")

				for i, e := range uw.Unwrap() {
					fmt.Fprintf(&sb, "%d. %s\n", i+1, e)
				}

				return errors.New(sb.String())
			}
			return err
		},
	}

	return checkCmd
}
