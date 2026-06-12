<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
![CLI][cli-shield]
![pgxcli][pgxcli-shield]

![Go][go-shield]
![PostgreSQL][postgres-shield]
<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://pgxcli.vercel.app/">
    <img src="https://res.cloudinary.com/dsdupsv2g/image/upload/q_auto/f_auto/v1779504361/home_lrj6gh.png" width="100%"/>
  </a>
  <h1 align="center">pgxcli</h1>
  <p align="center">
    Interactive PostgreSQL command-line client written in Go.
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#comparison-with-pgcli">Comparison with pgcli</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#thanks">Thanks</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

`pgxcli` is an interactive PostgreSQL command-line client built in Go, designed for a fast, and smooth REPL experience. It includes syntax highlighting, keyword autocompletion, command history, and support for PostgreSQL backslash commands.

Highlights:
* Interactive REPL with customizable prompt and style.
* Rich syntax highlighting.
* Context aware autocompletion.
* Persistent command history.
* PostgreSQL special backslash commands (`\d`, `\l`, `\dt`, `\q`).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PGXCLI VS PGCLI -->
## Comparison with pgcli

[Pgcli][pgcli-url] is a mature PostgreSQL CLI developed over many years, which has set the standard for interactive PostgreSQL clients.

### Where pgxcli stands out:
#### Now 
* Single binary, no external runtime dependencies
* Fast startup and better performance
* Modern CLI Interface
* Performance improvements for large tables

#### Planned
* Streaming query results for large tables
* Browser based Table view via localhost
* Direct Table export to SQL INSERT statements, CSV, MD tables, Excel, and HTML.

For more details, see the [Comparison with pgcli][comparison].

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Installation

`pgxcli` is distributed as a single binary for Linux, macOS, and Windows. You can download the latest version from the [releases page][releases-url].

#### Linux
Download the `.deb`, `.rpm`, `.apk`, or `.pkg.tar.zst` for your architecture and install it using your package manager.
```bash
# Example for Debian / Ubuntu
sudo dpkg -i pgxcli_*_linux_amd64.deb
```

#### macOS
Download the `.tar.gz` for your architecture (Intel or Apple Silicon), extract it, and move the binary to your `PATH`.
```bash
tar -xzf pgxcli_*_darwin_arm64.tar.gz
sudo mv pgxcli /usr/local/bin/
```

#### Windows
Download the `.msi` or `.zip` from the [releases page][releases-url]. The installer will automatically add `pgxcli` to your `PATH`.

#### Docker
You can run `pgxcli` directly via Docker using the pre-built images:
```bash
docker run -it --rm ghcr.io/balajz/pgxcli:v0.3.0 postgres://user:pass@host:5432/db
```

#### Go Install
If you have Go installed, you can install directly:
```bash
go install github.com/balajz/pgxcli/cmd/pgxcli@latest
```

For more details, see the [Installation Guide][install-ref].

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

```sh
# positional arguments
pgxcli mydb myuser

# flags
pgxcli --host localhost --port 5432 --user postgres --dbname postgres

# connection URI
pgxcli postgres://user:password@localhost:5432/dbname

# interactive connection form
pgxcli -i
```
<img src="https://res.cloudinary.com/dsdupsv2g/image/upload/q_auto/f_auto/v1777298704/5_h2fxui.png" alt="pgxcli flags screenshot" width="100%"/>

For full flag documentation, see the [docs][cli-ref].

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONFIGURATION -->
## Configuration

On first run, a config file is created at:

* `$XDG_CONFIG_HOME/pgxcli/config.toml` (or the OS-equivalent user config directory)

For configuration documentation, see the [configuration docs][config-ref].
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

We welcome contributions of all kinds — bug reports, feature requests, documentation improvements, and code contributions. Whether you're fixing a typo, reporting an issue, or building a feature, every contribution helps make pgxcli better.

To get started, please read our [Contributing Guide](CONTRIBUTING.md) for:
- How to report bugs and suggest features
- Development setup and local testing
- Our code style and testing conventions
- Guidelines for pull requests

Thank you for being part of the community! 🙌

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- THANKS -->
## Thanks

* [pgx][pgx-url]
* [bubbletea][bubbletea-url]
* [Cobra][cobra-url]
* [Viper][viper-url]
* Inspired by [pgcli][pgcli-url]
* [sqls][sqls-url]
* [CockroachDB][cockroachdb-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[cli-shield]: https://img.shields.io/badge/CLI-%23000000?style=for-the-badge&logo=iterm2&logoColor=white
[pgxcli-shield]: https://img.shields.io/badge/Pgxcli-7B36ED?style=for-the-badge&logo=database&logoColor=white
[go-shield]: https://img.shields.io/badge/Go-%23000000?style=for-the-badge&logo=go&logoColor=white
[postgres-shield]: https://img.shields.io/badge/PostgreSQL-7B36ED.svg?style=for-the-badge&logoColor=white
[comparison]: https://pgxcli.vercel.app/docs/guides/comparison-with-pgcli
[pgx-url]: https://github.com/jackc/pgx
[bubbletea-url]: https://github.com/charmbracelet/bubbletea
[cobra-url]: https://github.com/spf13/cobra
[viper-url]: https://github.com/spf13/viper
[pgcli-url]: https://github.com/dbcli/pgcli
[sqls-url]: https://github.com/sqls-server/sqls
[cockroachdb-url]: https://github.com/cockroachdb/cockroach
[cli-ref]: https://pgxcli.vercel.app/docs/reference/cli-reference/
[config-ref]: https://pgxcli.vercel.app/docs/guides/configuration
[install-ref]: https://pgxcli.vercel.app/docs/guides/getting-started#installation
[releases-url]: https://github.com/balajz/pgxcli/releases
