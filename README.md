# 🔍Livegrep-CLI🔍

Livegrep-CLI is a command line interface for the
[Livegrep](https://github.com/livegrep/livegrep) tool. ⚡

## Why?
Livegrep is a popular code-search tool used by many organizations to search
across their code. A public instance (that searches within the Linux kernel
source code) can be found at [https://livegrep.com/](https://livegrep.com). I
couldn't find a CLI tool that provided the set of configurable options I wanted
so I decided to write one myself. :)

This is my first Golang project of any substance, so please
excuse any mistakes or non-idiomatic code. Please open an issue
or file a PR if you spot any bugs or would like to suggest any
improvements.

## Use

By default `livegrep-cli` will run against the public Livegrep instance at
[https://livegrep.com/](https://livegrep.com). This can be altered by setting
the `LIVEGREP_URL` environment variable.

```
$ LIVEGREP_URL=livegrep.com ./livegrep-cli $query
```

Command line flags are intended to be at least mostly
compatible with familiar flags from `grep` and `ag` [The Silver
Searcher](https://github.com/ggreer/the_silver_searcher).

Supported environment variables:
- `LIVEGREP_HOST` sets the url that Livegrep should use (by default
[livegrep.com](https://livegrep.com).
- `LIVEGREP_USE_HTTPS` determines whether Livegrep should use https
(By default HTTPS is enabled).
- `LIVEGREP_UNIX_SOCKET` allows you to proxy traffic through a
local Unix socket.


Copyright 2018 Isaac Diamond. Released under the MIT license.
