# AGENTS.md

## Product contract

- Treat the frozen product specification in `../pushman/docs/SPEC.md` as the source of truth.
- Do not invent public API behavior in this repository. Keep transport work behind the `Service` interface until the shared OpenAPI contract is accepted.
- Preserve stdout for successful command output and stdout machine output. Send errors and diagnostics to stderr.

## Go

- Use Go 1.27 or newer and format changes with `gofmt`.
- Keep command I/O and terminal detection injectable so behavior remains testable.
- Validate user input before invoking the service.
- Run `go test ./...` and `go vet ./...` before committing.

## Git

- Commit messages must use `{type}: {message}`.
- Keep the `type` lowercase and concise.
- Write the `message` as an imperative, specific summary.
- Inspect staged and unstaged changes before every commit.
- Stage only files that belong to the current change.
