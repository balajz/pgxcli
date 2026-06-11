package config

import (
	"errors"
	"strings"
)

func Validate(cfg Config) error {
	var errs []error

	if cfg.Main.Prompt == "" {
		errs = append(errs, errors.New("prompt must not be empty"))
	}
	if cfg.Main.Style == "" {
		errs = append(errs, errors.New("style must not be empty"))
	} else if !cfg.Main.Style.isValid() {
		errs = append(errs, errors.New("style must be a valid syntax highlighting style"))
	}
	if cfg.Main.HistoryFile == "" {
		errs = append(errs, errors.New("history file path must not be empty"))
	}
	if cfg.Main.LogFile == "" {
		errs = append(errs, errors.New("log file path must not be empty"))
	}

	pagerMode := strings.ToLower(strings.TrimSpace(cfg.Main.Pager))
	if pagerMode == "" {
		errs = append(errs, errors.New("pager mode must not be empty"))
	} else {
		switch pagerMode {
		case "auto", "always", "never":
		default:
			errs = append(errs, errors.New("pager mode must be one of: auto, always, never"))
		}
	}
	onError := cfg.Main.OnError
	if onError == "" {
		errs = append(errs, errors.New("on_error action must not be empty"))
	} else if !onError.isValid() {
		errs = append(errs, errors.New("on_error action must be one of: STOP, RESUME"))
	}
	if !cfg.Table.Style.isValid() {
		errs = append(errs, errors.New("table style must be a valid style"))
	}

	if !cfg.Table.Color.Header.isValid() {
		errs = append(errs, errors.New("table color header must be a valid color"))
	}
	if !cfg.Table.Color.Column.isValid() {
		errs = append(errs, errors.New("table color column must be a valid color"))
	}
	if !cfg.Table.Color.Caption.isValid() {
		errs = append(errs, errors.New("table color caption must be a valid color"))
	}

	return errors.Join(errs...)
}
