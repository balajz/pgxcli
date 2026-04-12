package main

import (
	"context"
	"os"

	"github.com/balaji01-4d/pgxcli/internal/cli"
	"github.com/balaji01-4d/pgxcli/internal/cliio"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	printer := cliio.NewPgxPrinter(os.Stdout, os.Stderr)
	cliCtx := &cli.CliContext{Printer: printer}


	rootCmd := cli.NewRootCmd(
		ctx,
		cliCtx,
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

}
