# goup

`goup` is a simple Go installer. 
It is notorious that an operating system's package manager takes a long time to update Go to the latest version, particularly on Linux distros.
`goup` makes Go installation and version management easier.
Besides, `goup` does not inject junks into your shell, like other version managers.

`goup` is heavily inspired by [golang/dl](https://github.com/golang/dl) and [getgo](https://github.com/golang/tools/tree/master/cmd/getgo).

## How it works

* `goup init` outputs a file (`$HOME/.go/env`) that exports goup's (`$HOME/.go/bin`) and Go's (`$HOME/.go/current/bin`) bin directory in yoru PATH environment variable.
* `goup update` downloads specified version of Go and symlinks downloaded version to `$HOME/.go/current`.
* `goup show` shows the current installed Go version.

## Installtion

```
curl https://raw.githubusercontent.com/jingweno/goup/master/install.sh | sh
```

## Quick Start

```
$ echo 'source $HOME/.go/env' > ~/.bashrc # Equivalent of adding export PATH="$HOME/.go/bin":"$HOME/.go/current/bin:$PATH" to ~/.bashrc
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
/home/owen/.go/current
```

## License

[Apache 2.0](https://github.com/jingweno/goup/blob/master/LICENSE)
