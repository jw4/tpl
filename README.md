# tpl

A simple CLI tool that applies JSON data to Go templates.

## Installation

```sh
go install github.com/jw4/tpl@latest
```

## Usage

```sh
tpl -d data.json template.tpl > output.txt
```

- `-d` specifies the JSON data file
- Template file is provided as a positional argument
- Output is written to stdout

The JSON data is parsed into a `map[string]any` and made available to the Go `text/template` for rendering.
