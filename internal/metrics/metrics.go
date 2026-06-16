// Package metrics provides adapters for querying metrics backends.
// The MetricsSource interface allows Burnless to support Prometheus,
// Datadog, and any OpenTelemetry-compatible backend interchangeably.
package metrics

import (
	"context"
	"time"
)

// DataPoint is a single metric sample.
type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

// MetricsSource is the interface all metrics backends must implement.
type MetricsSource interface {
	// Query executes a query and returns the current scalar result.
	Query(ctx context.Context, query string) (float64, error)

	// QueryRange executes a query over a time range.
	QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]DataPoint, error)
}
