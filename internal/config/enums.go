package config

// OnErrorAction controls multi-statement behavior after a statement fails.
type OnErrorAction string

const (
	// OnErrorResume continues executing remaining statements after an error.
	OnErrorResume OnErrorAction = "RESUME"
	// OnErrorStop stops execution immediately when a statement fails.
	OnErrorStop OnErrorAction = "STOP"
)

func (a OnErrorAction) isValid() bool {
	switch a {
	case OnErrorResume, OnErrorStop:
		return true
	default:
		return false
	}
}
