# Good First Issues

Welcome! These are well-defined tasks for new contributors.
Each one is scoped to under 2 hours and has clear acceptance criteria.

## Beginner — under 1 hour

### Add unit tests for internal/toil
**Files:** `internal/toil/toil_test.go`
**What:** Write unit tests for toil event storage and reporting.
**Acceptance criteria:**
- Test that events can be saved and retrieved
- Test that monthly reports aggregate correctly
- All tests pass with `make test`

### Add unit tests for internal/runbook
**Files:** `internal/runbook/runbook_test.go`
**What:** Write unit tests for runbook execution.
**Acceptance criteria:**
- Test dry-run mode (no actual execution)
- Test that steps are executed in order
- All tests pass with `make test`

### Add a sre.yaml example for a background worker
**Files:** `examples/background-worker/sre.yaml`
**What:** Write a realistic sre.yaml for a background job service.
**Acceptance criteria:**
- Valid sre.yaml (passes `burnless validate`)
- Covers availability and throughput SLOs
- Includes a runbook for when throughput drops

## Intermediate — 1-2 hours

### Implement burnless status command
**Files:** `internal/cli/status.go`
**What:** Show current SLO status in a table format.
**Output should look like:**
SERVICE          SLO              STATUS    BUDGET

payments-api     availability     ✓ OK      77% remaining

checkout-api     latency          ✗ WARN    12% remaining

### Add PagerDuty notification support
**Files:** `internal/notify/pagerduty.go`
**What:** Implement the Notifier interface for PagerDuty.
**Acceptance criteria:**
- Sends alert to PagerDuty Events API v2
- Includes service name, severity, and burn rate in payload
- Has unit tests with mocked HTTP

### Add Prometheus metrics scraping
**Files:** `internal/metrics/prometheus.go`
**What:** Implement MetricsSource interface for Prometheus HTTP API.
**Acceptance criteria:**
- Queries Prometheus HTTP API with PromQL
- Returns float64 metric value
- Has unit tests with mocked HTTP server

## Browse all issues

https://github.com/Custos-com/Burnless/issues?q=label%3A%22good+first+issue%22
