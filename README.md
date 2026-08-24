# Pushman CLI

`pushman` sends push notifications to your own iPhone from terminals, scripts, servers, and CI jobs.

This repository is the first-release CLI scaffold. It contains the complete command tree, deterministic input validation, injectable command I/O, and a typed service boundary. The HTTP transport and credential persistence intentionally remain unimplemented until Pushman's shared OpenAPI contract is accepted.

## Requirements

- Go 1.27 or newer

## Develop

```sh
go mod download
go test ./...
go vet ./...
go run ./cmd/pushman help
```

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

Cobra commands depend on the `internal/cli.Service` interface. A later generated API adapter will implement that interface and own HTTPS requests, authentication, and error-code translation. This keeps the CLI's public UX independent from host or path layout decisions and prevents a provisional endpoint shape from becoming a compatibility promise.

Release metadata is injected with linker flags into `main.version`, `main.commit`, and `main.date`.

`api/openapi.yaml` is a bundled snapshot of the authoritative public contract in `WhiteKiwi/pushman`. After updating that snapshot, run `go generate ./internal/api` and commit the generated client. CI rejects stale generated code.

## Product contract

The source of truth for first-release behavior is the frozen specification in the adjacent private planning repository at `../pushman/docs/SPEC.md`. Public API documentation and generated clients will move to a shared public contract once its repository boundary is finalized.

## License

License selection is pending before publication.
