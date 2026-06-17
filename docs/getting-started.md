# Getting Started with Burnless

## Prerequisites
- Go 1.24+
- A running Prometheus instance
- kubectl (optional, for Kubernetes)

## Install

```bash
git clone https://github.com/Custos-com/Burnless.git
cd burnless
make build
```

## Create your first sre.yaml

```yaml
service: my-api
team: platform-engineering
slos:
  - name: availability
    target: 99.9%
    window: 30d
error_budget:
  burn_rate_alerts:
    - rate: 14.4x
      severity: critical
      notify:
        slack: "#incidents"
```

## Run Burnless

```bash
burnless validate    # check your sre.yaml is valid
burnless status      # see current SLO status
burnless agent       # start the agent daemon
```
