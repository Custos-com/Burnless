<div align="center">

```
██████╗ ██╗   ██╗██████╗ ███╗   ██╗██╗     ███████╗███████╗███████╗
██╔══██╗██║   ██║██╔══██╗████╗  ██║██║     ██╔════╝██╔════╝██╔════╝
██████╔╝██║   ██║██████╔╝██╔██╗ ██║██║     █████╗  ███████╗███████╗
██╔══██╗██║   ██║██╔══██╗██║╚██╗██║██║     ██╔══╝  ╚════██║╚════██║
██████╔╝╚██████╔╝██║  ██║██║ ╚████║███████╗███████╗███████║███████║
╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝╚══════╝
```

**Stop burning your error budget. Stop burning out your team.**

Full SRE config as code — SLOs, error budgets, runbooks, on-call, and dashboards in one repo.
Auto-remediates before alerts even fire.

[![CI](https://img.shields.io/github/actions/workflow/status/Custos-com/Burnless/ci.yml?label=CI&style=flat-square)](https://github.com/Custos-com/Burnless/actions)
[![Release](https://img.shields.io/github/v/release/Custos-com/Burnless?style=flat-square)](https://github.com/Custos-com/Burnless/releases)
[![Go Version](https://img.shields.io/badge/go-1.24-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-community-5865F2?style=flat-square&logo=discord)](https://discord.gg/burnless)

</div>

---

## The problem

Your SLOs live in a Datadog dashboard. Your runbooks live in Confluence — stale since last quarter.
Your alert thresholds are configured manually in Grafana, differently in every environment.

None of it is versioned. None of it is reviewable. None of it executes itself at 3am.

**Burnless changes that.**

---

## Installation

### macOS — Homebrew (recommended)

```bash
brew install --cask Custos-com/tap/burnless
```

### macOS + Linux — one-liner

```bash
curl -fsSL https://raw.githubusercontent.com/Custos-com/Burnless/main/scripts/install.sh | sh
```

### Go developers

```bash
go install github.com/Custos-com/Burnless/cmd/burnless@latest
```

### Download binary directly

Download the latest binary for your platform from
[GitHub Releases](https://github.com/Custos-com/Burnless/releases/latest).

| Platform | File |
|---|---|
| macOS Apple Silicon | `burnless_darwin_arm64.tar.gz` |
| macOS Intel | `burnless_darwin_amd64.tar.gz` |
| Linux amd64 | `burnless_linux_amd64.tar.gz` |
| Linux arm64 | `burnless_linux_arm64.tar.gz` |
| Windows | `burnless_windows_amd64.zip` |

### Verify installation

```bash
burnless version
# burnless 0.2.0
```

---

## Quick start

```bash
# Install
curl -fsSL https://burnless.dev/install.sh | sh

# Scaffold a new sre.yaml
burnless init

# Validate
burnless validate sre.yaml

# Deploy to Prometheus + Grafana + PagerDuty
burnless apply sre.yaml

# Watch live burn rate
burnless status payments-api

# Start the auto-remediation agent
burnless agent start
```

---

## The sre.yaml

> **VS Code tip:** Open any file named `sre.yaml` and get autocomplete automatically.
> Schema is published on [SchemaStore.org](https://www.schemastore.org/json/).

```yaml
service: payments-api
team: platform-engineering

slos:
  - name: availability
    target: 99.9%
    window: 30d
    indicator:
      metric: http_requests_total
      good_filter: 'status!~"5.."'

error_budget:
  burn_rate_alerts:
    - severity: critical
      rate: 14.4x
      window: 1h
      remediate: scale-up

    - severity: warning
      rate: 6x
      window: 6h
      remediate: restart-pods

runbooks:
  scale-up:
    mode: auto
    steps:
      - kubectl scale deploy/payments --replicas=+2
      - wait 60s
      - assert slo.availability > 99.5%

  restart-pods:
    mode: semi-auto
    steps:
      - kubectl rollout restart deploy/payments

oncall:
  provider: pagerduty
  escalation_minutes: 10
  notify_slack: "#sre-incidents"

dashboards:
  provider: grafana
  auto_generate: true
  panels:
    - error_budget_remaining
    - burn_rate_1h
    - burn_rate_6h
```

---

## Documentation

- [Getting started](docs/getting-started.md)
- [sre.yaml reference](docs/reference/sre-yaml.md)
- [CLI reference](docs/reference/cli.md)
- [Integrations](docs/integrations.md)
- [Security & Data Handling](docs/security.md)
- [Data Flow Diagram](docs/data-flow.html)
- [Contributing](CONTRIBUTING.md)

## Community

- [Discord](https://discord.gg/burnless)
- [Twitter / X](https://twitter.com/burnlessdev)
- [Blog](https://burnless.dev/blog)

## License

Apache 2.0 — see [LICENSE](LICENSE).

## License

Burnless uses a three-tier license strategy:

| Layer | Files | License |
|-------|-------|---------|
| SDK & Schema | `pkg/` `schema/` `examples/` | Apache 2.0 — maximum ecosystem reach |
| Core CLI & Agent | `cmd/` `internal/` `deploy/` | AGPLv3 — free forever including SSO |
| SaaS Dashboard | `saas/` `dashboard/` | BSL 1.1 — free for dev, paid for production SaaS |

**TL;DR:** If you're an SRE engineer using the CLI or self-hosting the agent — it's free, forever.
If you want to offer a managed Burnless service — [contact us](mailto:licensing@burnless.dev).

See [licenses/](licenses/) for full license texts.
