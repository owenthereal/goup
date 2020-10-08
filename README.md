# goup

`goup` (pronounced Go Up) is a simple Go installer and version manager.

It is notorious that an operating system's package manager takes time to update Go to the latest version, particularly on Linux distros.
`goup` makes Go installation and version management easy, without tying yourself to the Go release cycle of your package manager.
You can download the latest version (or any version) of Go with one command.
Besides, `goup` does not inject junks into your shell, like other version managers: it only exports the Go bin directory to yoru PATH environment.

`goup` is written in Go and is heavily inspired by [golang/dl](https://github.com/golang/dl) & [getgo](https://github.com/golang/tools/tree/master/cmd/getgo).

## How it works

* `goup init` outputs a file (`$HOME/.go/env`) that exports goup's (`$HOME/.go/bin`) and Go's (`$HOME/.go/current/bin`) bin directory in your PATH environment variable.
* `goup update` downloads specified version of Go and symlinks downloaded version to `$HOME/.go/current`.
* `goup show` shows the current installed Go version.

## Installation

```
curl -sSf https://raw.githubusercontent.com/jingweno/goup/master/install.sh | sh

```

You need goup's bin directory (`$HOME/.go/bin`) and Go's bin directory (`$HOME/.go/current/bin`)
in your PATH environment variable. Add the following to your shell startup script:

```
echo 'source $HOME/.go/env' > ~/.bashrc # Equivalent of adding export PATH="$HOME/.go/bin":"$HOME/.go/current/bin:$PATH" to ~/.bashrc

```

## Quick Start

```
$ goup update
Downloaded   0.0% (    16384 / 121149509 bytes) ...
Downloaded   6.9% (  8404928 / 121149509 bytes) ...
Downloaded  17.3% ( 20987744 / 121149509 bytes) ...
Downloaded  33.5% ( 40533712 / 121149509 bytes) ...
Downloaded  48.5% ( 58736192 / 121149509 bytes) ...
Downloaded  66.0% ( 79920544 / 121149509 bytes) ...
Downloaded  84.0% (101743872 / 121149509 bytes) ...
Downloaded  94.3% (114244768 / 121149509 bytes) ...
Downloaded 100.0% (121149509 / 121149509 bytes)
INFO[0010] Unpacking /home/owen/.go/go1.15.2/go1.15.2.linux-amd64.tar.gz ...
INFO[0022] Success: go1.15.2 downloaded in /home/owen/.go/go1.15.2
INFO[0022] Activated go1.15.2
$ goup show
go1.15.2
$ go env GOROOT
/home/owen/.go/go1.15.2
$ go version
go version go1.15.2 linux/amd64

$ goup update tip
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
INFO[0297] Activated gotip
$ goup show
gotip
$ go env GOROOT
/home/owen/.go/gotip
$ go version
go version devel +5d13781 Thu Oct 8 00:28:09 2020 +0000 linux/amd64
```

## License

[Apache 2.0](https://github.com/jingweno/goup/blob/master/LICENSE)
