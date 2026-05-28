package pgxspecial

// Export returns a list of all registered commands.
func Export() []CommandExport {
	var cmds []CommandExport
	for key, cmd := range commandRegistry {
		// key contains the command name with the alias
		// \quit, '\q, \exit has the same command, only differs by key/alias
		cmd := New(key, cmd.Syntax, cmd.Description) 
		cmds = append(cmds, cmd)
	}
	return cmds
}
