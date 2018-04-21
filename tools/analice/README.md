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
analice generate notice /path/to/github.com/user/project
```

### Generate [debian copyright file](https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/).
```sh
analice generate copyright /path/to/github.com/user/project
```

## Assumptions

* Root project uses `dep` as dependency management, `Gopkg.lock` need to be present in root directory.
* Every dependency is licensed under one license known to this tool.
* All dependencies need to be in the vendor folder including their `LICENSE` files.

## Known limitations

* Skips dependencies without license, either because they don't have one or it is absent in vendor folder.
* Skips dependencies which bury the license file in a non-standard location, e.g. by embedding it in `README`.
* Sometimes does not work for dependencies with heavy re-formatting of the `LICENSE` file, e.g. by using markdown, prefixing every line wth "> ", etc.
  Line breaks and spaces are OK though.
* A quick test showed this tool to work on ~50% of the projects in my `GOPATH`.
