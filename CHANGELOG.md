# Changelog

## [Unreleased]

## [0.2.1] - 2026-05-24

## Credits: CockroachDB for code references.

### Refactored
- **Removed custom results** Removed the custom results and switch to the database/Rows interface for scalability and maintainability, helps to supports future features like streaming results and exporting results.


## [0.2.0] - 2026-05-22

### Added
- **Modernized UI**: Completely redesigned REPL interface with styles and better visual hierarchy.
- **Orca Banner**: Colorful ASCII orca banner with gradient styling on startup.
- **Report Issue link**: Clickable link in the status bar to report issues on GitHub and banner.
- **Loading Spinner**: Visual spinner indicator that displays during query execution.

### Fixed
- **Version update**: Updated version string.
- **fix interactive page**: fix the interactive page to show only form when width is lesser.

### Refactored
- **UI Component Architecture**: Extracted child components (Input, Status, Spinner) into separate modules for better maintainability.
- **Style Management**: Separated style definitions into dedicated source files for cleaner organization.

## [0.1.2] - 2026-05-20

### Added
- **Query Cancellation**: Press `Ctrl+C` during a running query to cancel it immediately via PostgreSQL's out-of-band cancel signal.

### Fixed
- **hardcoded style**: Removed the hardcoded "monokai" style from the syntax highlighter.

## [0.1.1] - 2026-05-18

### Added
- **External Editor Support**: Launch your favorite editor with `Ctrl+E` to write complex queries.
- **Standard PostgreSQL Environment Variables**: Support for `PGHOST`, `PGPORT`, and `PGDATABASE` for easier connection management.
- **Modernized UI**: Refreshed REPL interface.
- **Dialer Timeout**: Added timeout support for database connections to improve responsiveness on network issues.
- **Error Colorization**: Improved error reporting with colorized output for better readability.

### Fixed
- **Secure Password Input**: Passwords are no longer echoed to the terminal during input.
- **Improved Password Retry Flow**: Added clear error messaging and retry prompts for incorrect passwords.

### Refactored
- **ParsePrompt Optimization**: Refactored prompt parsing to use a single-pass construction, significantly reducing memory allocations.
- **Keyword Management**: Migrated PostgreSQL keywords to the `completer` module to support future context-aware completion.

## [0.1.0] - 2026-04-30

### The First Release

The first release of pgxcli, a command-line interface for PostgreSQL inspired by pgcli, with a focus on performance and extensibility.

### Added
- Syntax-aware SQL splitting via jackc's sqlsplit
- JSONL-based history system
- configurable syntax highlighting
- Support for multiple table formats
- Linux packages: deb, rpm, apk, archlinux
- Windows MSI installer
