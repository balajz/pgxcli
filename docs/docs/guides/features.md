---
title: Features
description: An overview of what pgxcli brings to your PostgreSQL workflow.
keywords: [postgresql syntax highlighting, sql autocompletion, psql features, postgres auto complete, chroma]
sidebar_position: 6
---

pgxcli is an interactive PostgreSQL REPL built in Go. This page highlights the core features that make it fast, readable, and comfortable to use.

---

## Syntax Highlighting

pgxcli applies real-time syntax highlighting to your SQL as you type, using [Chroma](https://github.com/alecthomas/chroma) — the same engine that powers many static site generators and code tools.

![pgxcli syntax highlighting — monokai and onedark themes](/img/syntax.png)

Colors are applied to keywords, identifiers, string literals, operators, and comments, so long queries remain easy to read at a glance.

pgxcli automatically detects your terminal's color depth — TrueColor, 256-color, or 16-color — and picks the right formatter for your environment.

---

## Autocompletion

pgxcli offers context-aware SQL completion. Press `Tab` to cycle through suggestions for keywords, tables, schemas, and columns based on cursor position (for example after `FROM`, `JOIN`, or `alias.`).

![pgxcli autocompletion in action](/img/completion.gif)

---

## External Editor Support

For complex queries, you can launch your favorite text editor directly from the REPL.

Press `Ctrl+E` to open the current query in the editor defined by your `$EDITOR` environment variable (e.g., `vim`, `nano`, `code`). When you save and exit the editor, the query is automatically brought back to the prompt, ready for execution.

---

## Interactive UI Elements

pgxcli provides a modern and responsive user interface, including:
- **Loading Spinner:** A visual spinner indicator that displays during query execution, letting you know your query is running.
- **Orca Banner:** A colorful ASCII orca banner with gradient styling on startup.
- **Issue Reporting:** A clickable link in the status bar to easily report issues on GitHub.

---