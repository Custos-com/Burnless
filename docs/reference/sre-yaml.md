# sre.yaml Reference

## service
Name of your service.

## team
Team responsible for this service.

## slos
List of SLO definitions.
- `name` — name of the SLO
- `target` — target percentage (e.g. 99.9%)
- `window` — rolling window (e.g. 30d)

## error_budget
- `burn_rate_alerts` — list of burn rate thresholds
  - `rate` — burn rate multiplier (e.g. 14.4x = critical)
  - `severity` — critical or warning
  - `notify` — slack channel or pagerduty routing key
