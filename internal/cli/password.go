package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/term"
)

// shouldAskForPassword decides whether to prompt for password based on error code and flags.
func shouldAskForPassword(err error, neverPrompt bool) bool {
	if neverPrompt {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "28P01" { // invalid_password
		return true
	}
	return false
}

var readPassword = term.ReadPassword

func promptPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := readPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}
