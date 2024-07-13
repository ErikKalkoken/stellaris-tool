# Stellaris Tool

A tool for converting Stellaris save games into JSON.

[![Go](https://github.com/ErikKalkoken/stellaris-tool/actions/workflows/go.yml/badge.svg)](https://github.com/ErikKalkoken/stellaris-tool/actions/workflows/go.yml)

## Description

This package contains the tool `sav2json` which converts the contents of Stellaris save games into JSON. The tool can be downloaded directly for Windows, Linux and macOS or build from source for many other platforms. The tool is written in Go and has no build dependencies.

## Installation

### Latest release

You find the latest release for your platform on the [releases page](https://github.com/ErikKalkoken/stellaris-tool/releases).

Please download the package for your respective platform directly from the releases page and decompress it, e.g. with unzip or tar. The tool is a single executable and can be run directly.

To install it you need to copy the executable into a folder, which is already in your `PATH`, e.g. `~/.local/bin` on Linux.

### Build and install from repository

If you system has a go compiler you can install the tool directly with:

```sh
go install github.com/ErikKalkoken/stellaris-tool/cmd/sav2json@latest
```

## Usage

The usage is as follows:

```plain
Usage: sav2json [options] <inputfile>:

sav2json converts a Stellaris save game into JSON.

Options:
  -d string
        destination directory for output files (default ".")
  -k    keep original data files
  -s    create output files in same directory as source files
  -v    show the current version
```

You can always print the current usage of the tool with: `sav2json -h`.

> [!TIP]
> The location of the Stellaris save game files various by platform and installation method. Please see the official [Stellaris Wiki](https://stellaris.paradoxwikis.com/Save-game_editing) on how to find them.

## Support

Should you encounter a bug please feel free to open an issue with "Bug: ..." in the title. If the issue relates to a problem with a save file, please attach a copy of that save file to the issue, so we can debug it.

If you feel you are missing an important feature, please feel free to open an issue for it with "Feature request: ..." in the title.

## Contributions

We welcome any contributions to this open source project.
