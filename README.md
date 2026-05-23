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
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li><a href="#comparison-with-pgcli">Comparison with pgcli</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#development">Development</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

`pgxcli` is an interactive PostgreSQL command-line client built in Go, designed for a fast, and smooth REPL experience. It includes syntax highlighting, keyword autocompletion, command history, and support for PostgreSQL backslash commands.

Highlights:
* Interactive REPL with customizable prompt and style.
* SQL syntax highlighting.
* SQL keyword autocompletion.
* Persistent command history.
* PostgreSQL special backslash commands (`\d`, `\l`, `\dt`, `\q`).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PGXCLI VS PGCLI -->
## Comparison with pgcli

[Pgcli][pgcli-url] is a mature PostgreSQL CLI developed over many years, which has set the standard for interactive PostgreSQL clients.

**pgxcli** takes the simpler approach, focusing on speed and minmal setup. It is a singe Go binary with fast startup and TOML configuration. If you need a lightweight, It may be good fit. for a more feature-rich, established experience, pgcli remains the benchmark.

### Where pgxcli stands out:
#### Now 
* Single binary, no external runtime dependencies
* Fast startup and better performance

#### Planned
* Modern CLI Interface
* Streaming query results for large tables
* Browser based Table view via localhost
* Performance improvements for large tables
* Direct Table export to SQL INSERT statements, CSV, MD tables, Excel, and HTML.

<details>
  <summary><strong>Which one should I use?</strong></summary>

Right now, I would definitely choose pgcli. I think no explanation is needed.

That could change as pgxcli matures. I would really appreciate if you give pgxcli a try and share your feedback. If you want to contribute, that would be even better.

</details>

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

* `~/.config/pgxcli/config.toml` (or the OS-equivalent user config directory)

For configuration documentation, see the [docs][config-ref]

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

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [pgx][pgx-url]
* [Cobra][cobra-url]
* [Viper][viper-url]
* [go-pretty][go-pretty-url]
* [go-prompter][go-prompter-url]
* Inspired by [pgcli][pgcli-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[cli-shield]: https://img.shields.io/badge/CLI-%23000000?style=for-the-badge&logo=iterm2&logoColor=white
[pgxcli-shield]: https://img.shields.io/badge/Pgxcli-7B36ED?style=for-the-badge&logo=database&logoColor=white
[go-shield]: https://img.shields.io/badge/Go-%23000000?style=for-the-badge&logo=go&logoColor=white
[postgres-shield]: https://img.shields.io/badge/PostgreSQL-7B36ED.svg?style=for-the-badge&logoColor=white

[pgx-url]: https://github.com/jackc/pgx
[cobra-url]: https://github.com/spf13/cobra
[viper-url]: https://github.com/spf13/viper
[go-pretty-url]: https://github.com/jedib0t/go-pretty
[go-prompter-url]: https://github.com/jedib0t/go-prompter
[pgcli-url]: https://github.com/dbcli/pgcli
[cli-ref]: https://pgxcli.vercel.app/reference/cli-reference/
[config-ref]: https://pgxcli.vercel.app/guides/configuration/
[install-ref]: https://pgxcli.vercel.app/docs/guides/getting-started#installation
[releases-url]: https://github.com/Balaji01-4D/pgxcli/releases
