# Contributing to pgxcli

First off — thank you for being here. Whether you're fixing a typo, reporting a bug, or building a feature, every contribution matters. pgxcli is a solo project and contributions from the community genuinely move the needle.

## Table of Contents

- [I Just Have a Quick Question](#i-just-have-a-quick-question)
- [What Should I Know Before Getting Started?](#what-should-i-know-before-getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Development Setup](#development-setup)
- [Style Guide](#style-guide)
- [AI Usage](#ai-usage)
- [Recognition](#recognition)

## I Just Have a Quick Question

GitHub Issues are for bugs and feature requests only — not questions. If you're unsure whether something is a bug, start a Discussion first.

- Open a [GitHub Discussion](https://github.com/balajz/pgxcli/discussions) for Q&A, ideas, and general conversation.
- Check the [docs](https://pgxcli.vercel.app/) — configuration, CLI reference, and guides live there.

## What Should I Know Before Getting Started?

pgxcli is a single Go binary — a PostgreSQL REPL built for speed and a smooth user experience. Understanding the high-level architecture will help you find your way around quickly.

### Package Structure

| Layer | Package | Purpose |
|---|---|---|
| Entry | `cmd/pgxcli/main.go` | Bootstrap context & printer |
| CLI | `internal/cli/` | Cobra commands, flags, connection params |
| Config | `internal/config/` | TOML loading, validation, paths |
| Logger | `internal/logger/` | `slog`-based file logging |
| App / REPL | `internal/app/` | Main loop, command routing, history |
| Database | `internal/database/` | pgx connection and query execution |
| Parser | `internal/parser/` | SQL splitting and query classification |
| UI | `internal/ui/` | Interactive connection forms (Charm TUI) |
| Completer | `internal/completer/` | SQL autocomplete suggestions |
| Output | `internal/cliio/` | `stdout`/`stderr` abstraction |

### Key Dependencies

| Dependency | Purpose | Role in pgxcli |
|---|---|---|
| [`pgx`](https://github.com/jackc/pgx) | PostgreSQL driver | Core — all DB connections and queries |
| [`cobra`](https://github.com/spf13/cobra) | CLI framework | Core — every command is built on this |
| [`viper`](https://github.com/spf13/viper) | Config management | Core — handles all config file parsing |
| [`tablewriter`](https://github.com/olekukonko/tablewriter) | Table rendering | UI — formats query results in the terminal |
| [`go-prompter`](https://github.com/jedib0t/go-prompter) | Interactive prompts | UI — handles REPL input and history |

## How Can I Contribute?

### Reporting Bugs

Before opening a bug report:

- Check [existing issues](https://github.com/balajz/pgxcli/issues) — if the bug is already reported and still open, leave a comment there instead of opening a new one.
- Try to reproduce the issue on the latest release.

When filing a bug, please include:

- **A clear, specific title.** "pgxcli crashes" is not helpful. "pgxcli panics when running `\d` on a table with no columns" is.
- **Steps to reproduce** — as precise as possible.
- **Expected vs. actual behavior.**
- **Your environment:** OS, architecture, pgxcli version (`pgxcli --version`), PostgreSQL version.
- **Relevant config** from `~/.config/pgxcli/config.toml`, if applicable.
- **A stack trace or terminal output**, pasted in a fenced code block.

Please use the [bug report template](https://github.com/balajz/pgxcli/issues/new?template=bug_report.md) when available.

### Suggesting Enhancements

Feature ideas are welcome. Before submitting one:

- Search [existing issues and discussions](https://github.com/balajz/pgxcli/issues) to avoid duplicates.
- Check the [roadmap](#roadmap--what-to-work-on) — it might already be planned.

A good enhancement request includes:

- **What problem it solves** — not just what it does.
- **How you'd expect it to work** from the user's perspective.
- **Alternatives you've considered**, and why this approach is better.
- **Examples from other tools**, if relevant (psql, pgcli, usql, etc.).

Please use the [feature request template](https://github.com/balajz/pgxcli/issues/new?template=feature_request.md) when available.

### Your First Code Contribution

Not sure where to start? Look for issues tagged:

- [`good first issue`](https://github.com/balajz/pgxcli/labels/good%20first%20issue) — small, well-scoped, good for getting familiar with the codebase.
- [`help wanted`](https://github.com/balajz/pgxcli/labels/help%20wanted) — meaningful contributions that don't require deep context.

Issues are sorted by comment count as a rough proxy for impact.

### Pull Requests

1. **Fork** the repo and create your branch from `main`.
2. **Set up your dev environment** (see [Development Setup](#development-setup) below).
3. **Make your changes.** Keep scope focused — one fix or feature per PR.
4. **Add or update tests** for whatever you changed.
5. **Run `golangci-lint run`** and resolve any new warnings before pushing.
6. **Open a PR** against `main` with a clear description of what changed and why. Link any related issues.
7. **Be responsive** — if a reviewer leaves feedback, engage with it promptly.

For roadmap-sized features (streaming results, browser table view, export formats), open a Discussion first to align on approach before writing significant code.

> **Note:** All status checks must pass before a PR is reviewed. If a check fails for an unrelated reason, leave a comment explaining why.

## Development Setup

**Prerequisites:** Go 1.26.4+, a running PostgreSQL instance (local or via Docker).

```bash
# Clone your fork
git clone https://github.com/<your-username>/pgxcli.git
cd pgxcli

# Install dependencies
go mod download

# Run from source
go run ./cmd/pgxcli --help

# Build the binary
go build -o pgxcli ./cmd/pgxcli

# Run tests
go test ./...

# Lint — install golangci-lint first: https://golangci-lint.run/usage/install/
golangci-lint run
```

For a quick local PostgreSQL instance via Docker:

```bash
docker run --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16
```

Then connect:

```bash
go run ./cmd/pgxcli --host localhost --port 5432 --user postgres --dbname postgres
```

## Style Guide

- Follow standard Go conventions — `gofmt` is non-negotiable.
- Run `golangci-lint run` before committing. Fix all new warnings.
- Prefer table-driven tests using `t.Run(...)`.
- All exported functions and types must have doc comments.
- Avoid unnecessary abstractions — pgxcli values clarity over cleverness.

## AI Usage

AI tools are welcome in your workflow — with some boundaries.

**We discourage** submitting AI-generated code directly into core Go source files (packages like `internal/database/`, `internal/app/`, `internal/parser/`, etc.). These are load-bearing parts of the codebase where correctness, clarity, and intentionality matter. AI-generated logic here tends to introduce subtle bugs or patterns that are hard to review and maintain.

**We actively encourage** using AI for:

- **Test cases** — generating table-driven tests, edge case coverage, and test scaffolding.
- **Docs & comments** — writing or improving godoc comments, README sections, and guides.
- **Release notes** — summarizing changelogs and formatting release content.
- **Understanding the codebase** — using AI to explore, explain, or map out unfamiliar parts before diving in.
- **Vulnerability hunting** — prompting AI to review code for security issues, edge cases, or misuse of APIs.
- **Non-critical tooling** — scripts, CI config, Dockerfiles, and other supporting files.

The rule of thumb: if it's going into a `.go` file that runs in production, write it yourself and understand every line. If it's helping you *around* that work — use whatever tools help you do it better.

## Recognition

Every merged contribution gets credited in the release notes. Significant contributors will be added to the `ACKNOWLEDGMENTS` section of the README.

Thank you for taking the time to improve pgxcli. 🙌
