# analice

## About

Get an overview of OSS licenses used in your `go` project. 

## Build

```sh
go build
```

## Usage

### Generate `NOTICE` file
```sh
analice generate notice github.com/user/project github.com/user/project/pkg --tests --gooses=linux,windows
```

### Generate [debian copyright file](https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/).
```sh
analice generate notice github.com/user/project --tests --gooses=linux
```

## How it works



## Assumptions

* Every dependency is licensed under one license known to this tool.

## Known limitations

* Skips dependencies without license, either because they don't have one or it is absent in vendor folder.
* Skips dependencies which bury the license file in a non-standard location, e.g. by embedding it in `README`.
* Sometimes does not work for dependencies with heavy re-formatting of the `LICENSE` file, e.g. by using markdown, prefixing every line wth "> ", etc.
  Line breaks and spaces are OK though.
