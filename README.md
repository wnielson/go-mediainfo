```
███╗   ███╗███████╗██████╗ ██╗ █████╗ ██╗███╗   ██╗███████╗ ██████╗
████╗ ████║██╔════╝██╔══██╗██║██╔══██╗██║████╗  ██║██╔════╝██╔═══██╗
██╔████╔██║█████╗  ██║  ██║██║███████║██║██╔██╗ ██║█████╗  ██║   ██║
██║╚██╔╝██║██╔══╝  ██║  ██║██║██╔══██║██║██║╚██╗██║██╔══╝  ██║   ██║
██║ ╚═╝ ██║███████╗██████╔╝██║██║  ██║██║██║ ╚████║██║     ╚██████╔╝
╚═╝     ╚═╝╚══════╝╚═════╝ ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═╝      ╚═════╝
```

Go rewrite of MediaInfo CLI.

## Installation

- Homebrew (macOS):

```sh
brew tap autobrr/go-mediainfo https://github.com/autobrr/go-mediainfo
brew install --cask autobrr/go-mediainfo/mediainfo
```

Note: Homebrew’s official MediaInfo install can conflict (`media-info` / `mediainfo`). This cask also links `go-mediainfo` so you can run both:

```sh
go-mediainfo version
```

- Go install (requires Go toolchain):

```sh
go install github.com/autobrr/go-mediainfo/cmd/mediainfo@latest
```

- Latest release (one-liner, Linux x86_64):
  - Replace `linux_amd64` with `linux_arm64`, `darwin_amd64`, or `darwin_arm64` as needed.

```sh
curl -sL "$(curl -s https://api.github.com/repos/autobrr/go-mediainfo/releases/latest | grep browser_download_url | grep linux_amd64 | cut -d\" -f4)" | tar -xz -C /usr/local/bin
```

## Usage

Recommended (likely what you want):

```sh
mediainfo /path/to/file
```

```sh
mediainfo /path/to/file --output=JSON
mediainfo /path/to/file --full
mediainfo /path/to/dir
mediainfo --info-parameters
mediainfo update
mediainfo version
```

Path is required (file or directory).

Default output: text.
Footer includes `ReportBy : go-mediainfo - vX.Y.Z`.

## Options

- `-f, --full` (complete report)
- `--output=...` (TEXT/JSON/XML/OLDXML/HTML/CSV/EBUCore/EBUCore_JSON/PBCore/PBCore2/Graph_Svg/Graph_Dot)
- `--language=...` (output language)
- `--logfile=...` (write output to a file)
- `--bom` (write UTF-8 BOM on Windows)
- `--help`, `--help-output`
- `--info-parameters`

## Commands

- `update` (self-update this binary; release builds only)
- `version` (print go-mediainfo version)
