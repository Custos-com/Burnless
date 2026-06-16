# Contributing to Burnless

Thank you for your interest in contributing. This guide will get you from zero to your first PR.

## Development setup (5 commands)

```bash
git clone https://github.com/burnless/burnless
cd burnless
make setup      # installs golangci-lint, goimports, goreleaser
make build      # builds the CLI binary to ./bin/burnless
make test       # runs the test suite
```

## Finding something to work on

- [good first issue](https://github.com/burnless/burnless/labels/good%20first%20issue) — scoped, documented, ideal for first-time contributors
- [help wanted](https://github.com/burnless/burnless/labels/help%20wanted) — higher impact, more open-ended
- [area/docs](https://github.com/burnless/burnless/labels/area%2Fdocs) — no Go required

Ask in [Discord](https://discord.gg/burnless) `#contributors` before starting work on anything large.

## PR process

1. Fork the repo and create a branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run `make test lint` — both must pass
4. Commit using Conventional Commits format (see below)
5. Open a PR against `main`
6. A maintainer will review within 3 business days

## Commit message format

We use [Conventional Commits](https://www.conventionalcommits.org):

```
feat(slo): add latency SLO support
fix(runbook): handle missing steps gracefully
docs(readme): fix install instructions
test(toil): add cost calculation unit tests
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`

## Code style

- `gofmt` + `goimports` for formatting (enforced by CI)
- `golangci-lint` for linting (enforced by CI)
- Organise by domain, not by layer (see `internal/`)
- Unit tests live next to the code they test (`foo_test.go`)
- Integration tests live in `tests/integration/`
- Keep functions short — if it needs a comment to explain what it does, extract it

## Reporting security vulnerabilities

Do NOT open a public issue. See [SECURITY.md](SECURITY.md).
