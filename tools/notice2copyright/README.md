# notice2copyright

## About

Parse a `NOTICE` file to generate the [debian copyright file](https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/).

## Build

```sh
go build
```

## Usage

```sh
notice2copyright path/to/github.com/user/project "MIT License (Expat)"
```

## Assumptions

* Every dependency is licensed under one license known to this tool.
* The project has a `NOTICE` file in the following format:
  ```text
  ...
  
  ========================================================================
  Apache License 2.0 (Apache-2.0)
  ========================================================================
  
  The following software have components provided under the terms of this license:
  
  - cobra (from https://github.com/spf13/cobra)

  ========================================================================
  MIT License (Expat)
  ========================================================================
  
  The following software have components provided under the terms of this license:
  
  - color (from https://github.com/fatih/color)

  ...
  ```
  The file needs to reside at the top-level directory of the project.
* Vendoring. All dependencies need to be in the vendor folder including their `LICENSE` files.
  A dependency referenced via `- cobra (from https://github.com/spf13/cobra)` in the `NOTICE`
  file has to be present in the vendor folder as `vendor/github.com/spf13/cobra)`.
