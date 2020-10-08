#!/bin/sh
# shellcheck shell=dash

# This is just a little script that can be downloaded from the internet to
# install goup. It just does platform detection, downloads the installer
# and runs it.

# It runs on Unix shells like {a,ba,da,k,z}sh. It uses the common `local`
# extension. Note: Most shells limit `local` to 1 var per line, contra bash.

set -u

GOUP_GH_RELEASE_API="${GOUP_GH_RELEASE_API:-https://api.github.com/repos/owenthereal/goup/releases/latest}"


main() {
  downloader --check
  need_cmd uname
  need_cmd mkdir
  need_cmd chmod

  get_architecture || return 1
  local _arch="$RETVAL"
  assert_nz "$_arch" "arch"

  local _ext=""
  case "$_arch" in
    *windows*)
      _ext=".exe"
      ;;
  esac

  local _latest_tag
  _latest_tag="$(ensure downloader "$GOUP_GH_RELEASE_API" "" | grep -oP '"tag_name": "\K(.*)(?=")')"

  local _url="https://github.com/owenthereal/upterm/releases/download/v${_latest_tag}/${_arch}${_ext}"
  local _dir="$HOME/.go/bin"
  local _file="${_dir}/goup${_ext}"

  ensure mkdir -p "$_dir"
  ensure downloader "$_url" "$_file"
  ensure chmod u+x "$_file"
  if [ ! -x "$_file" ]; then
    printf '%s\n' "Cannot execute $_file (likely because of mounting /tmp as noexec)." 1>&2
    printf '%s\n' "Please copy the file to a location where you can execute binaries and run ./goup${_ext} init --update." 1>&2
    exit 1
  fi

  ignore "$_file" init --update < /dev/tty

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

assert_nz() {
  if [ -z "$1" ]; then err "assert_nz $2"; fi
}

# Run a command that should never fail. If the command fails execution
# will immediately terminate with an error showing the failing
# command.
ensure() {
  if ! "$@"; then err "command failed: $*"; fi
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

# Return strong TLS 1.2-1.3 cipher suites in OpenSSL or GnuTLS syntax. TLS 1.2 
# excludes non-ECDHE and non-AEAD cipher suites. DHE is excluded due to bad 
# DH params often found on servers (see RFC 7919). Sequence matches or is
# similar to Firefox 68 ESR with weak cipher suites disabled via about:config.  
# $1 must be openssl or gnutls.
get_strong_ciphersuites_for() {
  if [ "$1" = "openssl" ]; then
    # OpenSSL is forgiving of unknown values, no problems with TLS 1.3 values on versions that don't support it yet.
    echo "TLS_AES_128_GCM_SHA256:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_256_GCM_SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384"
  elif [ "$1" = "gnutls" ]; then
    # GnuTLS isn't forgiving of unknown values, so this may require a GnuTLS version that supports TLS 1.3 even if wget doesn't.
    # Begin with SECURE128 (and higher) then remove/add to build cipher suites. Produces same 9 cipher suites as OpenSSL but in slightly different order.
    echo "SECURE128:-VERS-SSL3.0:-VERS-TLS1.0:-VERS-TLS1.1:-VERS-DTLS-ALL:-CIPHER-ALL:-MAC-ALL:-KX-ALL:+AEAD:+ECDHE-ECDSA:+ECDHE-RSA:+AES-128-GCM:+CHACHA20-POLY1305:+AES-256-GCM"
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
