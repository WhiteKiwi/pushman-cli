# Pushman CLI

`pushman` sends push notifications to your own iPhone from terminals, scripts, servers, and CI jobs.

This repository contains the first-release command tree, deterministic input validation, a generated OpenAPI client, interactive pairing polling, operating-system credential storage, and HTTP adapters for push, status, rename, logout, devices, history list, usage, and diagnostics. History detail remains a follow-up until its API operation is finalized.

## Requirements

- Go 1.27 or newer

## Develop

```sh
go mod download
go test ./...
go vet ./...
go run ./cmd/pushman help
```

The compiled default is `https://api.pushman.whitekiwi.link/v1`. For a local server use an explicit loopback override:

```sh
PUSHMAN_API_URL=http://127.0.0.1:8080/v1 go run ./cmd/pushman status
```

Overrides must use HTTPS except for loopback development and must end at `/v1`. Redirects are rejected so credentials cannot be forwarded to another origin.

## Command surface

```text
pushman pair
pushman status
pushman rename <nickname>
pushman logout
pushman push <body>
pushman devices
pushman history
pushman history show <message-id>
pushman usage
pushman doctor
pushman version
pushman help [command]
```

The `push` command accepts one positional body or standard input:

```sh
pushman push "Deployment finished"
printf '%s\n' "Deployment finished" | pushman push
pushman push - < message.txt
```

Frequently used options include `-t, --title`, `-d, --device`, `--json`, and `--quiet`. Run `pushman help push` for the full set. Automation credentials will be read only from `PUSHMAN_TOKEN`; never place reusable tokens in command arguments.

## Architecture boundary

Cobra commands depend on the `internal/cli.Service` interface. The generated-client adapter owns HTTPS requests, authentication, pairing polling, and stable error translation. This keeps the CLI's public UX independent from transport implementation details.

Paired credentials are stored through the operating-system keyring: macOS Keychain, Windows Credential Manager, or Secret Service on Linux. Storage is isolated by an API-base digest so development and production credentials cannot overwrite one another. `PUSHMAN_TOKEN` is process-only, is never persisted, and takes precedence only for `push`.

Release metadata is injected with linker flags into `main.version`, `main.commit`, and `main.date`.

`api/openapi.yaml` is a bundled snapshot of the authoritative public contract in `WhiteKiwi/pushman`. After updating that snapshot, run `go generate ./internal/api` and commit the generated client. CI rejects stale generated code.

## Product contract

The source of truth for first-release behavior is the frozen specification in the adjacent private planning repository at `../pushman/docs/SPEC.md`. Public API documentation and generated clients will move to a shared public contract once its repository boundary is finalized.

## License

License selection is pending before publication.
