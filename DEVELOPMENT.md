# Development Guide

This guide is for contributors who want to build, test, and extend Burnless.

## Prerequisites

- Go 1.24+
- Make
- golangci-lint (installed by `make setup`)

## Setup

```bash
git clone https://github.com/Custos-com/Burnless.git
cd burnless
make setup       # install tools
go mod tidy      # download dependencies
make build       # build the binary
make test        # run all tests
make lint        # run linter
```

## Project structure
burnless/

├── cmd/

│   ├── burnless/     # main CLI entry point

│   ├── agent/        # agent daemon

│   └── controller/   # Kubernetes operator

├── internal/

│   ├── config/       # sre.yaml parser and validator

│   ├── slo/          # error budget math

│   ├── toil/         # toil tracker

│   ├── runbook/      # runbook executor

│   ├── metrics/      # metrics source interface

│   └── notify/       # notification interface

├── pkg/

│   ├── types/        # public Go types

│   └── sdk/          # public SDK

├── schema/           # sre.yaml JSON schema

├── examples/         # example sre.yaml files

├── docs/             # documentation

└── deploy/           # Helm charts and K8s manifests

## Workflow

```bash
# 1. Create a branch from dev
git checkout dev
git checkout -b feat/your-feature

# 2. Make changes and add tests
# 3. Run checks
make test
make lint

# 4. Open a PR against dev (not main)
```

## Adding a new integration

All external integrations are behind Go interfaces:

```go
// Add a new metrics source
type MetricsSource interface {
    Query(ctx context.Context, query string) (float64, error)
}

// Create internal/metrics/datadog.go
type DatadogSource struct { ... }
func (d *DatadogSource) Query(...) { ... }
```

No changes to core logic required.

## Running locally

```bash
# Build and run
make build
./bin/burnless init
./bin/burnless validate
./bin/burnless agent --dry-run
```
