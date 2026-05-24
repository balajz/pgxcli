# pgxspecial

`pgxspecial` is a Go library that provides an API to execute PostgreSQL meta-commands (a.k.a. “special” or “backslash” commands), modeled on the behavior of tools like `psql` and inspired by the Python library [pgspecial](https://github.com/dbcli/pgspecial).

## Features

- Execute `psql`-style backslash commands directly from Go code  
- Get structured metadata about databases: tables, types, functions, schemas, roles — not just raw SQL results  
- Works with `pgx/v5` and `pgxpool` (or any adapter implementing the included DB interface)  

## Installation

```bash
go get github.com/balajz/pgxspecial
```

## Basic Usage (Go API)

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/balajz/pgxspecial"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    ctx := context.Background()
    pool, err := pgxpool.New(ctx, "postgres://user:password@localhost:5432/database?sslmode=disable")
    if err != nil {
        log.Fatalf("Unable to connect: %v\n", err)
    }
    defer pool.Close()

    // Execute a special command
    // Returns a SpecialCommandResult interface
    res, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, pool, "\l")
    if err != nil {
        log.Fatalf("Special command error: %v\n", err)
    }

    if isSpecial {
        // Handle the result based on its type
        switch r := res.(type) {
        case pgxspecial.RowResult:
            // Standard result (like \l, \dt) - wraps pgx.Rows
            fmt.Println("Rows returned:")
            // Iterate r.Rows ...
            r.Rows.Close()

        case pgxspecial.DescribeTableListResult:
            // Complex result (like \d my_table)
            for _, table := range r.Results {
                fmt.Printf("Table: %s\n", table.TableMetaData.TypedTableOf) // Example
                // Access Columns, Data, and TableMetaData
            }
        }
    }
}
```

## Supported Commands

| Cmd            | Syntax               | Description                                  |
| -------------- | -------------------- | -------------------------------------------- |
| `\l` (`\list`) | `\l[+] [pattern]`    | List databases                               |
| `\d`           | `\d[+] [pattern]`    | List or describe tables, views and sequences |
| `DESCRIBE`     | `DESCRIBE [pattern]` | List or describe tables, views and sequences |
| `\dT`          | `\dT[+] [pattern]`   | List data types                              |
| `\ddp`         | `\ddp [pattern]`     | List default access privilege settings       |
| `\dD`          | `\dD[+] [pattern]`   | List or describe domains                     |
| `\dx`          | `\dx[+] [pattern]`   | List extensions                              |
| `\dE`          | `\dE[+] [pattern]`   | List foreign tables                          |
| `\df`          | `\df[+] [pattern]`   | List functions                               |
| `\dt`          | `\dt[+] [pattern]`   | List tables                                  |
| `\dv`          | `\dv[+] [pattern]`   | List views                                   |
| `\dm`          | `\dm[+] [pattern]`   | List materialized views                      |
| `\ds`          | `\ds[+] [pattern]`   | List sequences                               |
| `\di`          | `\di[+] [pattern]`   | List indexes                                 |
| `\dp` (`\z`)   | `\dp [pattern]`      | List privileges                              |
| `\du`          | `\du[+] [pattern]`   | List roles                                   |
| `\dn`          | `\dn[+] [pattern]`   | List schemas                                 |
| `\db`          | `\db[+] [pattern]`   | List tablespaces                             |
| `\!`           | `\! command`         | Execute a shell command                      |
| `\sf`          | `\sf[+] FUNCNAME`    | Show a function's definition                 |


## Result Types

The library now uses a polymorphic result type `SpecialCommandResult` to handle different kinds of output:

1.  **`RowResult`**: Wraps standard `pgx.Rows`. Used by list commands like `\l`, `\dt`, `\du`.
2.  **`DescribeTableListResult`**: Returned by `\d [pattern]`. Contains a list of `DescribeTableResult` structs, each with:
    *   `Columns`: Header names
    *   `Data`: Grid data (rows)
    *   `TableMetaData`: Footer info (Indexes, Constraints, Triggers, etc.)
3. **`ExtensionVerboseListResult`**: Returned by `\dx+ [pattern`, Contains a list of `ExtensionVerboseResult` structs, each with:
   *    `Name`: Extension name
   *    `Description`: Extension's description

## Contributing

Contributions are welcome!
Feel free to open issues or submit pull requests for bug fixes, new commands, improved tests, or documentation enhancements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
