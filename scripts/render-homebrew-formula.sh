#!/bin/sh
set -eu

if [ "$#" -ne 5 ]; then
  echo "usage: $0 VERSION ARM64_MACOS_SHA X86_MACOS_SHA ARM64_LINUX_SHA X86_LINUX_SHA" >&2
  exit 2
fi

version=$1
arm64_macos_sha=$2
x86_macos_sha=$3
arm64_linux_sha=$4
x86_linux_sha=$5

if ! printf '%s\n' "$version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z][0-9A-Za-z.-]*)?$'; then
  echo "invalid release version: $version" >&2
  exit 2
fi

for checksum in "$arm64_macos_sha" "$x86_macos_sha" "$arm64_linux_sha" "$x86_linux_sha"; do
  if ! printf '%s\n' "$checksum" | grep -Eq '^[0-9a-f]{64}$'; then
    echo "invalid SHA-256: $checksum" >&2
    exit 2
  fi
done

script_dir=$(CDPATH='' cd -P "$(dirname "$0")" && pwd)
template="$script_dir/../packaging/homebrew/pushman.rb.in"

awk \
  -v version="$version" \
  -v arm64_macos_sha="$arm64_macos_sha" \
  -v x86_macos_sha="$x86_macos_sha" \
  -v arm64_linux_sha="$arm64_linux_sha" \
  -v x86_linux_sha="$x86_linux_sha" \
  '
    {
      version_count += gsub(/@VERSION@/, version)
      arm64_macos_count += gsub(/@ARM64_MACOS_SHA@/, arm64_macos_sha)
      x86_macos_count += gsub(/@X86_MACOS_SHA@/, x86_macos_sha)
      arm64_linux_count += gsub(/@ARM64_LINUX_SHA@/, arm64_linux_sha)
      x86_linux_count += gsub(/@X86_LINUX_SHA@/, x86_linux_sha)
      if ($0 ~ /@[A-Z0-9_]+@/) {
        unresolved_token = 1
      }
      print
    }
    END {
      if (version_count != 8 || arm64_macos_count != 1 || x86_macos_count != 1 ||
          arm64_linux_count != 1 || x86_linux_count != 1 || unresolved_token) {
        print "unexpected Homebrew formula template token counts" > "/dev/stderr"
        exit 2
      }
    }
  ' "$template"
