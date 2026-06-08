package ui

import "charm.land/bubbles/v2/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	Quit      key.Binding
	Interrupt key.Binding
	Clear     key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "quit"),
		),
		Interrupt: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "interrupt"),
		),
		Clear: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc esc", "clear"),
		),
	}
}
