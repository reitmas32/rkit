# AGENTS.md

Guidance for AI coding agents working **on** the rkit repository. (Consumers of
the library should start at [README.md](./README.md), [`llms.txt`](./llms.txt),
or https://pkg.go.dev/github.com/reitmas32/rkit.)

## Project

rkit is a general-purpose Go utilities library for microservices: structured
errors, a generic `Result[T]`, an error-accumulating context, generic
repositories, an event bus, and a structured logger with a Grafana Loki hook.

- Module: `github.com/reitmas32/rkit`
- Go: **1.25+**
- License: MIT

## Layout

- `core/` — pure contracts and value types (no heavy deps): `customctx`, `eventbus`, `http`, `kerrors`, `logger`, `result`, `types`.
- `infrastructure/` — concrete implementations of core contracts: `http` client, `eventbus/{inmemory,rabbit}`, `dtos`.
- `persistence/` — storage abstractions and backends: `contracts`, `criteria`, `pagination`, `models`, `inmemory`, `postgres`, `mongodb`.
- `observability/logger/loguru/` — logrus-based logger (`fields`, `hooks`).
- `examples/`, `mock/loki/` — runnable programs and a mock Loki server (each may be its own module; do not include them in the root module build).

## Commands

```bash
go build ./...        # build the module
go test ./...         # run unit tests AND runnable examples (Example_* with // Output:)
go vet ./...          # static checks (also typechecks tests)
gofmt -l .            # list unformatted files (must be empty)
```

Run a single package's tests: `go test ./core/result/...`.

`examples/` and `mock/loki/` have their own `go.mod`; build them from their own
directory, not via the root `./...`.

## Conventions

- **Errors:** return `*kerrors.KError` (code + metadata + cause), not bare
  `errors.New`. Across layers prefer returning `result.Result[T]` rather than
  `(T, error)`.
- **Context:** functions that can accumulate errors take `*customctx.CustomContext`.
- **Docs:** every package has a `doc.go` with a `// Package` comment. Add
  runnable `Example_*` functions in `*_example_test.go` / `example_test.go`;
  prefer a deterministic `// Output:` so `go test` validates them and pkg.go.dev
  renders them.
- **Public API docs are in English.** Some `docs/` guides are still Spanish and
  are being migrated; new and updated docs should be English.
- Keep `core/` dependency-light; heavy deps (GORM, AMQP, mongo-driver, logrus)
  belong in `infrastructure/`, `persistence/`, or `observability/`.

## Adding a package

1. Create the package directory with the implementation.
2. Add a `doc.go` with a `// Package <name>` overview and a short usage snippet.
3. Add an `example_test.go` with at least one runnable `Example`.
4. Add an entry to the package map in `README.md` and a section in `llms.txt`.
5. `go build ./... && go vet ./... && go test ./...` must pass.

## Releasing

Versions are git tags (semver). Current scheme: commit directly to `main`, then
create an annotated tag and push both:

```bash
git commit -m "feat(scope): ..."
git tag -a vX.Y.Z -m "vX.Y.Z: ..."
git push origin main
git push origin vX.Y.Z
```

Consumers then run `go get github.com/reitmas32/rkit@vX.Y.Z`.

## Security

- Do not log secrets. Loki/loguru fields become labels — keep them low
  cardinality and free of credentials/PII.
- `core/http` and `infrastructure/http` must not follow untrusted redirects
  without validation.
