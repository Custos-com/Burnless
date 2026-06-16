package metrics

import (
	"context"
	"fmt"
	"time"
)

// PrometheusClient queries a Prometheus HTTP API endpoint.
type PrometheusClient struct {
	BaseURL string
}

// NewPrometheusClient creates a new Prometheus metrics client.
func NewPrometheusClient(baseURL string) *PrometheusClient {
	return &PrometheusClient{BaseURL: baseURL}
}

// Query executes an instant PromQL query.
func (p *PrometheusClient) Query(ctx context.Context, query string) (float64, error) {
	return 0, fmt.Errorf("not yet implemented")
}

// QueryRange executes a range PromQL query.
func (p *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]DataPoint, error) {
	return nil, fmt.Errorf("not yet implemented")
}
