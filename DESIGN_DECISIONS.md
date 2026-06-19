# Design Decisions

This document explains *why* Burnless was built the way it was.
These are the decisions that shaped the architecture.

## Why YAML for config?

- Human readable — SRE teams already live in YAML (Kubernetes, Helm, GitHub Actions)
- Git friendly — diffs are meaningful, reviews are easy
- Declarative — describe what you want, not how to get it
- Tooling — editors, linters, and schema validators already exist

Alternatives considered: TOML (less familiar in SRE ecosystem), HCL (Terraform-specific), JSON (no comments).

## Why Go?

- Single binary — no runtime dependencies, easy to install and distribute
- Fast startup — CLI tools need to feel instant
- Same language as Kubernetes, Prometheus, and Terraform — contributors are familiar
- Strong standard library — JSON, YAML, HTTP, crypto all built in
- goroutines — the agent's polling loop is trivially concurrent

## Why bbolt for local storage?

The CLI needs to store toil events between runs without requiring a database server.
bbolt is an embedded key-value database — just a single file on disk.

- Zero setup — no server, no config, no migrations
- Single file — easy to back up, easy to delete
- ACID transactions — safe concurrent writes
- Pure Go — no CGo, works on all platforms

Alternative considered: SQLite — requires CGo, harder to cross-compile.

## Why interfaces for external tools?

```go
type MetricsSource interface {
    Query(ctx context.Context, query string) (float64, error)
}
```

The agent does not know or care whether metrics come from Prometheus, Datadog, or CloudWatch.
It calls `MetricsSource.Query()` and gets a float64 back.

This means:
- New integrations = new file, zero changes to core logic
- Testing is trivial — mock the interface
- Contributors can add integrations independently

## Why AGPLv3 for the core?

We want the CLI and agent to be free forever — including SSO and RBAC.
But we need to prevent cloud providers from offering Managed Burnless without contributing back.
AGPLv3 closes the SaaS loophole: if you run a modified version as a service, you must release your source.

See [licenses/](licenses/) for the full three-tier license strategy.

## Why Apache 2.0 for the SDK?

`pkg/types` and `pkg/sdk` are the public API for building tooling on top of Burnless.
We want zero friction for anyone building VS Code extensions, CI plugins, or third-party integrations.
Apache 2.0 is pre-approved at Google, Microsoft, and Amazon — no legal review needed.

## Why semi-auto mode as default?

Auto-remediation that runs without human approval will not be trusted in production.
The default is semi-auto: Burnless detects the issue and posts proposed steps to Slack for approval.
This builds trust before teams enable auto mode — which is opt-in per runbook.
