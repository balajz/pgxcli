# Copilot Instructions for pgxcli

## Build, Test, and Lint Commands

### Build
```bash
make build          # Build the application (output: bin/app)
make clean          # Remove compiled binaries
```

### Test
```bash
go test ./...                      # Run all tests
make test-integration              # Run integration tests (build tag: integration)
go test ./internal/config          # Test specific package
go test -v ./internal/database     # Verbose output
```

### Lint
```bash
make lint           # Run golangci-lint (requires golangci-lint binary)
make precommit      # Run lint AND test (equivalent to: lint test)
```

### Run
```bash
make build && ./bin/app mydb myuser          # Build and run
./bin/app --host localhost --user postgres  # With flags
```

## High-Level Architecture

pgxcli is an interactive PostgreSQL CLI written in Go. The application follows a layered architecture:

### Core Packages

- **`cmd/pgxcli/main.go`** – Entry point; initializes context, printer, and root CLI command.
- **`internal/cli`** – Cobra command definitions (flags parsing, connection setup, password prompting).
- **`internal/app`** – Application orchestration layer. `Application` manages the REPL loop, routes commands, renders results, and manages history.
- **`internal/app/ui`** – Bubble Tea REPL model, connection forms (huh), and UI components.
- **`internal/app/renderer`** – Result formatting and table rendering helpers.
- **`internal/database`** – PostgreSQL abstraction layer using `jackc/pgx`. Contains `Client` (connection mgmt) and special commands handler.
- **`internal/config`** – Configuration management (loads embedded defaults + user config from `~/.config/pgxcli/config.toml`).
- **`internal/logger`** – Structured logging via `log/slog`. Initializes file-based logger with debug flag support.
- **`internal/completer`** – SQL autocompletion engine; maintains database schema metadata.
- **`internal/parser`** – SQL statement splitter for multi-statement input.
- **`internal/cliio`** – Output printing abstraction (stdout/stderr wrapper).

### Data Flow

1. **cobra root** loads config/logger, initializes pager, and connects `database.Client`.
2. **app.Start()** builds the Bubble Tea model (`internal/app/ui`) with input, spinner, and status components.
3. User input is handled by Bubbline editline with history + syntax highlighting; multi-statement SQL is split by `parser.SplitSQLStatements`.
4. Builtins (e.g., `\clear`) are handled in app; pgSQL special commands (`\d`, `\q`, `\c`, `\conninfo`) go through `pgxspecial`.
5. SQL executes via `database.Client.ExecuteQuery`; `OnErrorAction` controls stop/resume for multi-statement input.
6. **renderer** formats results using table config; output flows through `cliio.Printer` with optional pager.
7. History is saved via editline on close (default `~/.pgxcli_history.jsonl` when `history_file = "default"`).

## Key Conventions

### Configuration & Initialization

- Configuration is validated on load (see `config.validate()`). Invalid values block startup.
- Paths support `"default"` as a sentinel value; `config.Load()` resolves it to OS-appropriate defaults (e.g., `~/.config/pgxcli/`).
- Logger accepts `debug bool` flag; enables debug-level slog output if true.

### Testing Patterns

- Tests use **stretchr/testify** for assertions and mocking.
- Database tests define mock doubles (`MockConn`, `MockRows`) in `mocks_test.go`; see `database/executor_test.go` for patterns.
- Config tests use `t.TempDir()` and env isolation (set `HOME`, `XDG_CONFIG_HOME`, etc. to temp dirs).
- Context passed as first argument to all db operations: `ctx` is `context.Background()` or passed from CLI.

### Error Handling & Logging

- All errors are wrapped with context using `fmt.Errorf(..., %w, err)`.
- Logs use structured `slog` with attribute key `"error"` (not `"err"`) for error values: `logger.Error("msg", "error", err)`.
- REPL loop catches errors mid-execution; logs and continues (doesn't exit).

### Special Commands & PostgreSQL Features

- Special pgSQL commands (e.g., `\d`, `\l`) are registered via **pgxspecial** package initialization (see `internal/database/special_commands.go`).
- The `OnErrorAction` config controls multi-statement behavior: `STOP` exits on error, `RESUME` continues.
- **Parser** splits multi-statement input; execution happens statement-by-statement.

### REPL Reader & History

- REPL is built on Bubble Tea v2 with Bubbline editline input (history, autocomplete, syntax highlighting).
- History is persisted by Bubbline to a JSONL file (default `~/.pgxcli_history.jsonl` when `history_file = "default"`).

### Rendering & Output

- Renderer lives in `internal/app/renderer`; table formatting uses `internal/app/renderer/formatter` with config-driven styles.
- Output goes through a `Printer` interface (wraps stdout/stderr) for testability.

### Go Version & Dependencies

- Go 1.26.4+ required.
- Key dependencies: pgx/v5 (database), cobra (CLI), viper (config), charm libraries (TUI), testify (testing).

### Naming Conventions

- Unexported (private) functions/vars use `camelCase` (e.g., `pgxCLI`, `newReader()`).
- Exported (public) types/funcs use `PascalCase` (e.g., `Application`, `Client`, `Logger`).
- Interfaces are typically named `Action` or end with `-er` (e.g., `Reader`, `Printer`, `Application`).
- Test files use `_test.go` suffix; test functions start with `Test` (standard Go convention).
