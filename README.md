<div align="center">

<br/>
**Stop burning your error budget. Stop burning out your team.**

Full SRE config as code — SLOs, error budgets, runbooks, on-call, and dashboards in one repo.
Auto-remediates before alerts even fire.

[![CI](https://img.shields.io/github/actions/workflow/status/burnless/burnless/ci.yml?label=CI&style=flat-square)](https://github.com/burnless/burnless/actions)
[![Release](https://img.shields.io/github/v/release/burnless/burnless?style=flat-square)](https://github.com/burnless/burnless/releases)
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

- [Getting started](https://docs.burnless.dev/getting-started)
- [sre.yaml reference](https://docs.burnless.dev/reference/sre-yaml)
- [CLI reference](https://docs.burnless.dev/reference/cli)
- [Integrations](https://docs.burnless.dev/integrations)
- [Contributing](CONTRIBUTING.md)

## Community

- [Discord](https://discord.gg/burnless)
- [Twitter / X](https://twitter.com/burnlessdev)
- [Blog](https://burnless.dev/blog)

## License

Apache 2.0 — see [LICENSE](LICENSE).
