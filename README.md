# Stellaris Tool

A JSON converter for Stellaris save game files.

[![Go](https://github.com/ErikKalkoken/stellaris-tool/actions/workflows/go.yml/badge.svg)](https://github.com/ErikKalkoken/stellaris-tool/actions/workflows/go.yml)

## Description

This package contains the tool `sav2json` which can convert Stellaris save game files into JSON. We are releasing this tool for Linux, Windows and MAC (experimental).

## Installation

### Releases

You find the most current release of the tool for your platform on the [releases page](https://github.com/ErikKalkoken/stellaris-tool/releases).

Please download the tool for your platform directly from the release page and decompress it, e.g. with unzip.

The tool is a single executable.

To install it you need to copy the executable into a folder, which is already in your `PATH`, e.g. `~/.local/bin` on Linux.

### Build and install from repository

If you system has a go compiler you can install the tool directly with:

```sh
go install github.com/ErikKalkoken/stellaris-tool/cmd/sav2json@latest
```

## Support

Should you encounter a bug please feel free to open an issue with "Bug: ..." in the title. If the issue relates to a problem with a save file, please attach a copy of that save file to the issue, so we can debug it.

If you feel you are missing an important feature, please feel free to open an issue for it with "Feature request: ..." in the title.

## Contributions

We welcome any contributions to this open source project. Please open a PR with your suggested change.
