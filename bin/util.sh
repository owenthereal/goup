#!/usr/bin/env bash

version() {
  ver="$(git describe --tags HEAD 2>/dev/null || true)"

  if [ -z "$ver" ]; then
    ver="$(grep -w 'Version =' "$version_file" | cut -d'"' -f2)"
    sha="$(git rev-parse --short HEAD 2>/dev/null || true)"
    [ -z "$sha" ] || ver="${ver}-g${sha}"
  fi

  echo "${ver#v}"
}

version_file="internal/commands/version.go"
