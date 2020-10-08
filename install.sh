#!/bin/sh
# shellcheck shell=dash

# This is just a little script that can be downloaded from the internet to
# install goup. It just does platform detection, downloads the installer
# and runs it.

# It runs on Unix shells like {a,ba,da,k,z}sh. It uses the common `local`
# extension. Note: Most shells limit `local` to 1 var per line, contra bash.

set -u

GOUP_GH_RELEASE_API="https://api.github.com/repos/owenthereal/goup/releases/latest"
GOUP_GH_RELEASE_PREFIX="https://github.com/owenthereal/goup/releases/download"


main() {
  downloader --check
  need_cmd uname
  need_cmd mkdir
  need_cmd chmod

  get_architecture || return 1
  local _arch="$RETVAL"

  local _ext=""
  case "$_arch" in
    *windows*)
      _ext=".exe"
      ;;
  esac

  local _latest_tag
  _latest_tag="$(ensure downloader "$GOUP_GH_RELEASE_API" "" | grep -oP '"tag_name": "\K(.*)(?=")')"

  if [ -z "$_latest_tag" ]; then
    say "latest release tag not found"
    return 1
  fi

  local _url="${GOUP_GH_RELEASE_PREFIX}/${_latest_tag}/${_arch}-${_latest_tag}${_ext}"
  local _dir="$HOME/.go/bin"
  local _file="${_dir}/goup${_ext}"

  ensure mkdir -p "$_dir"
  ensure downloader "$_url" "$_file"
  ensure chmod u+x "$_file"
  if [ ! -x "$_file" ]; then
    printf '%s\n' "Cannot execute $_file." 1>&2
    printf '%s\n' "Please copy the file to a location where you can execute binaries and run ./goup${_ext} init --install." 1>&2
    exit 1
  fi

  ignore "$_file" init --install < /dev/tty

  local _retval=$?
  return "$_retval"
}

# This is just for indicating that commands' results are being
# intentionally ignored. Usually, because it's being executed
# as part of error handling.
ignore() {
  "$@"
}

need_cmd() {
  if ! check_cmd "$1"; then
    err "need '$1' (command not found)"
  fi
}

check_cmd() {
  command -v "$1" > /dev/null 2>&1
}

# Run a command that should never fail. If the command fails execution
# will immediately terminate with an error showing the failing
# command.
ensure() {
  if ! "$@"; then err "command failed: $*"; fi
}

say() {
  printf 'goup: %s\n' "$1"
}

err() {
  say "$1" >&2
  exit 1
}

# This wraps curl or wget. Try curl first, if not installed,
# use wget instead.
downloader() {
  local _dld
  if check_cmd curl; then
    _dld=curl
  elif check_cmd wget; then
    _dld=wget
  else
    _dld='curl or wget' # to be used in error message of need_cmd
  fi

  if [ "$1" = --check ]; then
    need_cmd "$_dld"
  elif [ "$_dld" = curl ]; then
    if [ -z "$2" ]; then
      curl --silent --show-error --fail --location "$1"
    else
      curl --silent --show-error --fail --location "$1" --output "$2"
    fi
  elif [ "$_dld" = wget ]; then
    if [ -z "$2" ]; then
      wget "$1"
    else
      wget "$1" -O "$2"
    fi
  else
    err "Unknown downloader"   # should not reach here
  fi
}

get_architecture() {
  local _ostype _cputype _arch
  _ostype="$(uname -s)"
  _cputype="$(uname -m)"

  if [ "$_ostype" = Darwin ] && [ "$_cputype" = i386 ]; then
    # Darwin `uname -m` lies
    if sysctl hw.optional.x86_64 | grep -q ': 1'; then
      _cputype=amd64
    fi
  fi

  case "$_ostype" in
    Linux)
      _ostype=linux
      ;;

    FreeBSD)
      _ostype=freebsd
      ;;

    Darwin)
      _ostype=darwin
      ;;

    MINGW* | MSYS* | CYGWIN*)
      _ostype=windows
      ;;

    *)
      err "unrecognized OS type: $_ostype"
      ;;
  esac

  case "$_cputype" in
    i386 | i486 | i686 | i786 | x86)
      _cputype=386
      ;;

    xscale | arm | armv6l | armv7l |armv8l)
      _cputype=arm
      ;;

    aarch64 | x86_64 | x86-64 | x64 | amd64)
      _cputype=amd64
      ;;

    *)
      err "unknown CPU type: $_cputype"
  esac

  _arch="${_ostype}-${_cputype}"

  RETVAL="$_arch"
}

main "$@" || exit 1
