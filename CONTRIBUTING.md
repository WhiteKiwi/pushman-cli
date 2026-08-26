# Contributing to Pushman CLI

Thanks for helping improve Pushman CLI. Bug reports, focused fixes, tests, documentation, and usability feedback are welcome.

## Before you start

- Search existing issues before opening a new one.
- Use GitHub Discussions for support and open an issue for reproducible bugs or concrete proposals.
- Never include credentials, tokens, private notification content, or unredacted logs.
- For vulnerabilities, follow [SECURITY.md](SECURITY.md) instead of opening a public issue.

For a substantial feature or public command change, open a proposal first. The hosted API and product behavior are maintained outside this repository, so a CLI change may require an accepted contract change before implementation.

## Development

Pushman CLI requires Go 1.27 or newer.

```sh
git clone https://github.com/pushmanhq/pushman-cli.git
cd pushman-cli
go mod download
go generate ./internal/api
go test -race ./...
go vet ./...
go run ./cmd/pushman help
```

Before submitting a pull request:

```sh
test -z "$(gofmt -l .)"
go generate ./internal/api
git diff --exit-code -- internal/api/api.gen.go
go test -race ./...
go vet ./...
```

Keep changes focused. Add or update tests for user-visible behavior, preserve the stdout/stderr contract, and do not hand-edit the generated API client. Commit messages use `type: imperative summary`, for example `fix: reject an empty notification body`.

By contributing, you agree that your contributions are licensed under the repository's [MIT License](LICENSE). All participants must follow the [Code of Conduct](CODE_OF_CONDUCT.md).
