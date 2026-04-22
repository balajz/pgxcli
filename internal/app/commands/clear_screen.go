// Package commands provides implementations for various commands
// that can be executed in the pgxCLI application.
package commands

import "fmt"

const escCode = "\033[2J\033[H\033[3J"

// ClearScreen clears the terminal screen
// by printing the appropriate escape code.
func ClearScreen() {
	fmt.Print(escCode)
}
