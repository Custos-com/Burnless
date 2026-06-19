# Burnless Roadmap

> Last updated: June 2026

## v0.1.0 — Foundation (current)
- [x] sre.yaml schema v0.1
- [x] SLO error budget math engine
- [x] Burn rate detection (14.4x critical, 6x warning)
- [x] Config parser with validation
- [x] burnless init command
- [x] burnless validate command
- [x] Slack notification support
- [x] Agent skeleton
- [x] Kubernetes operator skeleton
- [x] Full CI/CD pipeline
- [x] Community files

## v0.2.0 — Agent & Integrations
- [ ] burnless agent main loop
- [ ] Prometheus metrics scraping
- [ ] Burn rate auto-detection from live metrics
- [ ] PagerDuty on-call routing
- [ ] Executable runbooks (auto mode)
- [ ] Semi-auto runbook approval via Slack
- [ ] burnless status command
- [ ] bbolt local storage for toil events
- [ ] burnless toil log command
- [ ] burnless toil report command

## v0.3.0 — Kubernetes Native
- [ ] Kubernetes CRD for SLO objects
- [ ] Controller reconciliation loop
- [ ] kubectl plugin
- [ ] Helm chart for agent deployment
- [ ] Grafana dashboard auto-generation

## v1.0.0 — Production Ready
- [ ] Multi-service support
- [ ] Multi-cloud (AWS CloudWatch, GCP Operations, Azure Monitor)
- [ ] Datadog metrics source
- [ ] OpenTelemetry support
- [ ] Full test coverage
- [ ] SOC 2 readiness

## SaaS (Burnless Cloud)
- [ ] Web dashboard (burnless.dev)
- [ ] Multi-tenant architecture
- [ ] SSO / SAML
- [ ] Audit log UI
- [ ] Stripe billing
- [ ] SOC 2 Type II certification

---

Track progress: https://github.com/orgs/Custos-com/projects/2
