# Integrations

## Prometheus
Burnless reads metrics from Prometheus every 60 seconds.

## Slack
Add your Slack webhook URL to sre.yaml under `notify.slack`.

## PagerDuty
Add your PagerDuty routing key under `notify.pagerduty`.

## Grafana
Burnless can generate Grafana dashboards from your sre.yaml.

## Kubernetes
Deploy Burnless as a Kubernetes operator using the manifests in deploy/.
