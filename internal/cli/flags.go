package cli

import "github.com/spf13/cobra"

// hostFlag refers to -h / --host for database host.
type hostFlag string

func (f *hostFlag) bind(cmd *cobra.Command) {
	cmd.Flags().StringVarP((*string)(f), "host", "h", "", "host address of the postgres database")
}

// portFlag refers to -p / --port for database port, default is 5432
type portFlag uint16

func (f *portFlag) bind(cmd *cobra.Command) {
	cmd.Flags().Uint16VarP((*uint16)(f), "port", "p", 5432, "port number at which the postgres server is listening")
}

// usernameFlag refers to -u / -U / --user / --username for database user
// It binds two flags.
type usernameFlag string

func (f *usernameFlag) bind(cmd *cobra.Command) {
	cmd.Flags().StringVarP((*string)(f), "username", "u", "", "Username to connect to the postgres database.")
	cmd.Flags().StringVarP((*string)(f), "user", "U", "", "Username to connect to the postgres database.")
}

// dbNameFlage refers to -d / --dbname for database name.
type dbNameFlag string

func (f *dbNameFlag) bind(cmd *cobra.Command) {
	cmd.Flags().StringVarP((*string)(f), "dbname", "d", "", "database name to connect to.")
}

// forcePromptFlage refers to -W / --password for prompting for password.
type forcePromptFlag bool

func (f *forcePromptFlag) bind(cmd *cobra.Command) {
	cmd.Flags().BoolVarP((*bool)(f), "password", "W", false, "Force password prompt")
}


// neverPromptFlag refers to -w / --no-password for enforcing no prompt for passowrd.
type neverPromptFlag bool

func (f *neverPromptFlag) bind(cmd *cobra.Command) {
	cmd.Flags().BoolVarP((*bool)(f), "no-password", "w", false, "never prompt for the password")
}


// debugFlag use to debug application.
type debugFlag bool

func (f *debugFlag) bind(cmd *cobra.Command) {
	cmd.Flags().BoolVar((*bool)(f), "debug", false, "Enable debug mode for verbose logging.")
}


// interactiveConnFlag launches a form for filling the database connection parameters.
type interactiveConnFlag bool

func (f *interactiveConnFlag) bind(cmd *cobra.Command) {
	cmd.Flags().BoolVarP((*bool)(f), "interactive", "i", false, "Interactive connection mode")
}
