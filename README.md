<p align="center">
  <img src="docs/assets/pushman-icon.png" alt="Pushman" width="128" height="128">
</p>

<h1 align="center">Pushman CLI</h1>

<p align="center">
  Send push notifications to your iPhone from any terminal, script, server, or CI job.
</p>

<p align="center">
  <a href="https://github.com/WhiteKiwi/pushman-cli/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/WhiteKiwi/pushman-cli/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://go.dev/"><img alt="Go 1.27+" src="https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go&logoColor=white"></a>
  <a href="LICENSE"><img alt="MIT License" src="https://img.shields.io/badge/license-MIT-151515"></a>
  <a href="https://github.com/WhiteKiwi/pushman-cli/releases"><img alt="GitHub release" src="https://img.shields.io/github/v/release/WhiteKiwi/pushman-cli?display_name=tag&include_prereleases"></a>
</p>

```console
$ pushman push "Production deploy finished" --title "Acme API"
Accepted msg_01M0W2RDPVGEVX7D6ZWFK907B2 for 1 device
```

Pushman is a small, script-friendly companion for getting your own operational messages onto your iPhone. Pair once, then use the same command interactively or in automation. Credentials live in the operating-system keyring, not a plaintext config file.

> [!NOTE]
> Pushman for iPhone is currently in private beta. The App Store link will be added here when the public listing is available.

## Highlights

- One command for terminals, shell scripts, servers, and CI
- Rich notifications with title, subtitle, URL, image, sound, group, and update key
- Multiple receiving devices and explicit device targeting
- Seven-day history plus usage and delivery diagnostics
- Native credential storage through macOS Keychain, Windows Credential Manager, or Secret Service
- Stable JSON output and exit codes for automation
- Signed GitHub release provenance and published SHA-256 checksums

## Install

Install the latest tagged version with Go:

```sh
go install github.com/WhiteKiwi/pushman-cli/cmd/pushman@latest
```

This requires Go 1.27 or newer. Prebuilt binaries for macOS, Linux, and Windows are available on the [Releases](https://github.com/WhiteKiwi/pushman-cli/releases) page.

To verify a downloaded release asset:

```sh
release_asset=pushman_<version>_<platform>_<arch>.<archive>
awk -v name="$release_asset" '$2 == name' checksums.txt | shasum -a 256 -c -
gh attestation verify "$release_asset" -R WhiteKiwi/pushman-cli
```

## Quick start

Pair the CLI with Pushman on your iPhone:

```sh
pushman pair
```

Then send a notification:

```sh
pushman push "Database backup completed"
pushman push "Deploy completed" --title "Production" --url "https://example.com/runs/42"
printf '%s\n' "Build failed" | pushman push --title "CI"
```

Use `pushman help push` to see every notification field and output option.

## Commands

| Command | Purpose |
| --- | --- |
| `pushman pair` | Pair this CLI and assign its default nickname |
| `pushman push <body>` | Send a notification |
| `pushman devices` | List receiving devices |
| `pushman history` | Show recent messages |
| `pushman usage` | Show the current monthly allowance |
| `pushman status` | Show pairing and account state |
| `pushman rename <nickname>` | Rename this CLI |
| `pushman doctor` | Diagnose configuration and connectivity |
| `pushman logout` | Remove the local pairing credential |

## Automation

Pass reusable automation credentials only through `PUSHMAN_TOKEN`. Never place a token in command arguments, logs, or source control.

```sh
export PUSHMAN_TOKEN="..."
pushman push "Release $GITHUB_REF_NAME is live" --title "Deploy" --json
```

Successful command output is written to stdout; errors and diagnostics are written to stderr. Use `--json` for machine-readable output and `--quiet` when only the exit status matters.

## Service and future plans

The CLI source is MIT licensed. The hosted Pushman service, iPhone app, monthly allowance, and any future paid plan are separate from the source-code license. During beta, an account can submit up to 200 accepted push requests per month. The CLI discovers service capabilities at runtime so future plan and quota changes do not require embedding billing logic or secrets in this repository.

## Configuration

The production API is `https://api.pushman.whitekiwi.link/v1`. For local development, an explicit loopback override is supported:

```sh
PUSHMAN_API_URL=http://127.0.0.1:8080/v1 go run ./cmd/pushman status
```

Overrides must use HTTPS except for loopback development and must end in `/v1`. Redirects are rejected to prevent credentials from being forwarded to another origin.

## Development

```sh
go mod download
go generate ./internal/api
go test -race ./...
go vet ./...
go run ./cmd/pushman help
```

`api/openapi.yaml` is a bundled snapshot of the authoritative public contract. After updating it, run `go generate ./internal/api` and commit the generated client. CI rejects stale generated code.

See [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. For help, read [SUPPORT.md](SUPPORT.md). Please report vulnerabilities according to [SECURITY.md](SECURITY.md), never in a public issue.

## License

Pushman CLI is available under the [MIT License](LICENSE).
