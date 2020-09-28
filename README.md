# goup

`goup` is a simple Go version manager based on https://github.com/golang/dl.
It is notorious that an operating system's package manager takes a long time to update Go to the latest version, particularly on Linux distros.
`goup` makes Go version management easier with less than 100 lines of bash.
And it does not inject junks into your shell, like other version managers.

## How it works

`goup update` downloads specified version of Go using `go get golang.org/dl/<version>` and symlinks downloaded version to `$HOME/sdk/current`.
`goup init` exports `$HOME/sdk/current/bin` to your `PATH`. That's it!

## Prerequisites

`goup` requires a version of Go (even an outdated one!) to compile the [downloader](https://github.com/golang/dl) with `go get golang.org/dl/<version>`.
Go should be available in most package managers, even if the version is not always up-to-date.
For example, `brew install go` on Mac or `apt-get install go` on Debian.

## Installtion

Just drop `goup` to your `PATH`. For example,

```
curl https://raw.githubusercontent.com/jingweno/goup/master/goup > /usr/local/bin/goup
```

## Quick Start

```
$ echo 'eval "$(goup init -)"' > ~/.bashrc # Equivalent of adding export PATH="$HOME/sdk/current/bin:$PATH" to ~/.bashrc
$ goup update
Downloading go1.15.2...
go: found golang.org/dl/go1.15.2 in golang.org/dl v0.0.0-20200909201834-1fb66e01de4d
Downloaded   0.0% (    14448 / 122469402 bytes) ...
Downloaded   1.0% (  1259632 / 122469402 bytes) ...
Downloaded   5.9% (  7190640 / 122469402 bytes) ...
Downloaded  10.1% ( 12400752 / 122469402 bytes) ...
Downloaded  28.5% ( 34879600 / 122469402 bytes) ...
Downloaded  46.1% ( 56506480 / 122469402 bytes) ...
Downloaded  64.3% ( 78755952 / 122469402 bytes) ...
Downloaded  81.6% ( 99956848 / 122469402 bytes) ...
Downloaded  99.7% (122140784 / 122469402 bytes) ...
Downloaded 100.0% (122469402 / 122469402 bytes)
Unpacking /Users/jou/sdk/go1.15.2/go1.15.2.darwin-amd64.tar.gz ...
Success. You may now run 'go1.15.2'
Activated go1.15.2
$ goup show
go1.15.2
$ go env GOROOT
/Users/owen/sdk/current
```

## Wish List

The Go team is working on a `getgo` CLI to provide a one-liner installation of the latest Go, similar to rustup: https://github.com/golang/go/issues/23381.
A prototype was done, but it is never productionized.
I hope that `getgo` is ready soon.

## License

[Apache 2.0](https://github.com/jingweno/goup/blob/master/LICENSE)
