#!/usr/bin/env bash

version() {
  ver_file="internal/commands/version.go"
  ver="$(git describe --tags HEAD 2>/dev/null || true)"

  if [ -z "$ver" ]; then
    ver="$(grep -w 'Version =' "$ver_file" | cut -d'"' -f2)"
    sha="$(git rev-parse --short HEAD 2>/dev/null || true)"
    [ -z "$sha" ] || ver="${ver}-g${sha}"
  fi

  echo "${ver#v}"
}

