# Contributing

Thank you for helping improve this notifications engine. This document explains how we work in this repository and what we expect from contributions.

## Ground rules

- **Be respectful** — Assume good intent; keep review feedback constructive and specific.
- **Small, focused changes** — Prefer several small pull requests over one large unrelated diff.
- **Match existing style** — Follow patterns already used in `internal/` (naming, package layout, error handling, and Echo handler structure).

## Getting set up

1. Install **Go 1.24+** (see `go.mod` / toolchain).
2. Run **PostgreSQL**, **Kafka**, and **Redis** locally or via `docker compose up -d`.
3. Provide configuration the app can load (`go-config-library`, service name `github.com/viantonugroho11/go-notifications-engine`). For local work, use a `./config` directory or Consul as your team does in other environments.
4. Verify your setup:

```bash
go mod download
go test ./...
go run ./cmd/app
```

## Project architecture (where to put code)

We use a **clean architecture–style** layout:

| Layer | Responsibility |
|--------|----------------|
| `internal/entity` | Domain types and invariants |
| `internal/usecase` | Application services; depend on repository **interfaces** |
| `internal/repository` | Interfaces + GORM models + `postgres` implementations |
| `internal/transport/apis` | HTTP: Echo handlers, DTOs, routing |
| `internal/transport/event` | Async: Kafka handlers and routing |
| `internal/infrastructure` | DB, broker, cache — technical adapters |
| `internal/client` | Outbound HTTP or SDK integrations |

**Dependency direction:** handlers → usecases → repository interfaces ← infrastructure. Avoid importing transport or infrastructure from `entity` or `usecase` except through interfaces or small ports.

## How to contribute

1. **Open an issue or discuss** — For larger features, align on design first (API shape, events, schema).
2. **Branch** — Create a branch from the default branch (`feature/…`, `fix/…`, or your team’s convention).
3. **Implement** — Keep commits logical; avoid unrelated refactors in the same change as a bugfix unless necessary.
4. **Test** — Add or update tests when behavior changes (`go test ./...`). For JSON schema or pure functions, table-driven tests are welcome.
5. **Open a pull request** — Describe **what** changed and **why**, link issues, and note any config or migration steps for operators.

## Code review checklist (authors and reviewers)

- [ ] `go fmt ./...` and `go vet ./...` clean
- [ ] `go test ./...` passes
- [ ] No secrets committed (credentials, private keys, `.env` with real values)
- [ ] Config changes documented in README or deployment notes when behavior is user-visible
- [ ] API or Kafka contract changes called out explicitly in the PR description

## Commit messages

Use clear, imperative summaries (English is fine), for example:

- `Add notification inbox list filter by state`
- `Fix kafka consumer shutdown on context cancel`

Optionally add a body explaining tradeoffs or rollout steps.

## Security

If you discover a security vulnerability, **do not** open a public issue with exploit details. Contact the maintainers through the channel your organization uses for security reports.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT, unless the repository states otherwise).
