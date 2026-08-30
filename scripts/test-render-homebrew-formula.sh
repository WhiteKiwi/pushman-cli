#!/bin/sh
set -eu

script_dir=$(CDPATH='' cd -P "$(dirname "$0")" && pwd)
renderer="$script_dir/render-homebrew-formula.sh"
fixture_dir=$(mktemp -d "${TMPDIR:-/tmp}/pushman-formula-test.XXXXXX")
trap 'find "$fixture_dir" -depth -delete' EXIT HUP INT TERM

formula="$fixture_dir/pushman.rb"
a_sha=$(printf '%064s' '' | tr ' ' a)
b_sha=$(printf '%064s' '' | tr ' ' b)
c_sha=$(printf '%064s' '' | tr ' ' c)
d_sha=$(printf '%064s' '' | tr ' ' d)

sh "$renderer" 1.2.3-beta.4 "$a_sha" "$b_sha" "$c_sha" "$d_sha" > "$formula"

grep -Fqx '  homepage "https://github.com/pushmanhq/pushman"' "$formula"
grep -Fqx '      url "https://github.com/pushmanhq/pushman-cli/releases/download/v1.2.3-beta.4/pushman_1.2.3-beta.4_macOS_arm64.tar.gz"' "$formula"
grep -Fqx '      url "https://github.com/pushmanhq/pushman-cli/releases/download/v1.2.3-beta.4/pushman_1.2.3-beta.4_macOS_x86_64.tar.gz"' "$formula"
grep -Fqx '      url "https://github.com/pushmanhq/pushman-cli/releases/download/v1.2.3-beta.4/pushman_1.2.3-beta.4_linux_arm64.tar.gz"' "$formula"
grep -Fqx '      url "https://github.com/pushmanhq/pushman-cli/releases/download/v1.2.3-beta.4/pushman_1.2.3-beta.4_linux_x86_64.tar.gz"' "$formula"
grep -Fqx "      sha256 \"$a_sha\"" "$formula"
grep -Fqx "      sha256 \"$b_sha\"" "$formula"
grep -Fqx "      sha256 \"$c_sha\"" "$formula"
grep -Fqx "      sha256 \"$d_sha\"" "$formula"
grep -Fqx '    assert_match "pushman #{version}", shell_output("#{bin}/pushman version")' "$formula"

if grep -Eq '@[A-Z0-9_]+@|[[:blank:]]$' "$formula"; then
  echo "rendered formula contains a template token or trailing whitespace" >&2
  exit 1
fi

if sh "$renderer" '1.2.3; false' "$a_sha" "$b_sha" "$c_sha" "$d_sha" >/dev/null 2>&1; then
  echo "renderer accepted an invalid version" >&2
  exit 1
fi
if sh "$renderer" 1.2.3 invalid "$b_sha" "$c_sha" "$d_sha" >/dev/null 2>&1; then
  echo "renderer accepted an invalid checksum" >&2
  exit 1
fi
