# Goup

`goup` (pronounced Go Up) is an elegant Go installer.

It is notorious that an operating system's package manager takes time to update Go to the latest version, particularly on Linux distros.
`goup` makes Go installation and version management easy.
You can download the latest version (or any version) of Go with one command, without tying yourself to the availability of Go in your package manager.
Besides, `goup` does not inject magic into your shell, like other version managers: it only exports Goup's bin directory and Go's bin directory to your PATH environment.

`goup` is written in Go. It is heavily inspired by [Rustup](https://rustup.rs/), [golang/dl](https://github.com/golang/dl) and [getgo](https://github.com/golang/tools/tree/master/cmd/getgo).

## Installation

```
curl -sSf https://raw.githubusercontent.com/owenthereal/goup/master/install.sh | sh

```

Install by skipping the confirmation prompt, e.g., for automation:

```
curl -sSf https://raw.githubusercontent.com/owenthereal/goup/master/install.sh | sh -s -- '--skip-prompt'

```

## Quick Start

```
$ goup install
Downloaded   0.0% (     3267 / 121149509 bytes) ...
Downloaded   4.5% (  5504416 / 121149509 bytes) ...
Downloaded  12.5% ( 15122320 / 121149509 bytes) ...
Downloaded  20.1% ( 24345872 / 121149509 bytes) ...
Downloaded  24.2% ( 29363162 / 121149509 bytes) ...
Downloaded  24.8% ( 30048032 / 121149509 bytes) ...
Downloaded  32.6% ( 39517376 / 121149509 bytes) ...
Downloaded  40.8% ( 49413840 / 121149509 bytes) ...
Downloaded  48.6% ( 58850880 / 121149509 bytes) ...
Downloaded  55.9% ( 67681792 / 121149509 bytes) ...
Downloaded  63.4% ( 76790656 / 121149509 bytes) ...
Downloaded  68.3% ( 82754960 / 121149509 bytes) ...
Downloaded  70.8% ( 85818784 / 121149509 bytes) ...
Downloaded  73.1% ( 88604000 / 121149509 bytes) ...
Downloaded  75.6% ( 91601680 / 121149509 bytes) ...
Downloaded  76.6% ( 92798304 / 121149509 bytes) ...
Downloaded  78.9% ( 95583536 / 121149509 bytes) ...
Downloaded  82.8% (100252368 / 121149509 bytes) ...
Downloaded  87.7% (106282208 / 121149509 bytes) ...
Downloaded  94.8% (114801248 / 121149509 bytes) ...
Downloaded 100.0% (121149509 / 121149509 bytes)
INFO[0030] Unpacking /home/owen/.go/go1.15.2/go1.15.2.linux-amd64.tar.gz ...
INFO[0043] Success: go1.15.2 downloaded in /home/owen/.go/go1.15.2
INFO[0043] Default Go is set to 'go1.15.2'
$ goup show
go1.15.2
$ go env GOROOT
/home/owen/.go/go1.15.2
$ go version
go version go1.15.2 linux/amd64

$ goup install tip
Cloning into '/home/owen/.go/gotip'...
remote: Counting objects: 10041, done
remote: Finding sources: 100% (10041/10041)
remote: Total 10041 (delta 1347), reused 6538 (delta 1347)
Receiving objects: 100% (10041/10041), 23.83 MiB | 3.16 MiB/s, done.
Resolving deltas: 100% (1347/1347), done.
Updating files: 100% (9212/9212), done.
INFO[0078] Updating the go development tree...
From https://go.googlesource.com/go
 * branch            master     -> FETCH_HEAD
HEAD is now at 5d13781 cmd/cgo: add more architectures to size maps
Building Go cmd/dist using /home/owen/.go/go1.15.2. (go1.15.2 linux/amd64)
Building Go toolchain1 using /home/owen/.go/go1.15.2.
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1.
Building Go toolchain2 using go_bootstrap and Go toolchain1.
Building Go toolchain3 using go_bootstrap and Go toolchain2.
Building packages and commands for linux/amd64.
---
Installed Go for linux/amd64 in /home/owen/.go/gotip
Installed commands in /home/owen/.go/gotip/bin
INFO[0297] Default Go is set to 'gotip'
$ goup show
gotip
$ go env GOROOT
/home/owen/.go/gotip
$ go version
go version devel +5d13781 Thu Oct 8 00:28:09 2020 +0000 linux/amd64
```

## How it works

* `goup install` downloads specified version of Go to`$HOME/.go/VERSION` and symlinks it to `$HOME/.go/current`.
* `goup show` shows the activated Go version located at `$HOME/.go/current`.
* `goup init` appends Goup's bin directory (`$HOME/.go/bin`) and Go's bin directory (`$HOME/.go/current/bin`) to your PATH environment variable.


## License

[Apache 2.0](https://github.com/owenthereal/goup/blob/master/LICENSE)
