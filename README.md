# goup

`goup` is a Go version manager based on https://github.com/golang/dl.
I have been annoyed with the long wait time for a package manager to update Go to the latest version, particularly on Linux.
I made `goup` as a simple solution to make Go version management easier without injecting junks into my shell. 

## How it works

`goup` downloads specified version of Go using `go get golang.org/dl/<version>` and symlinks downloaded version to `$HOME/sdk/current`.
`goup init` exports `$HOME/sdk/current/bin` to your `PATH`. That's it!

## Prerequisites

`goup` requires a version of Go (even an outdated version!) to function becuase it downloads Go using `go get golang.org/dl/<version>`.
You can install Go with a package manager. For example, `brew install go` or `apt-get install go`.

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

## License

[Apache 2.0](https://github.com/jingweno/goup/blob/master/LICENSE)
