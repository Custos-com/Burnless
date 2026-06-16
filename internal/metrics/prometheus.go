package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/common/model"
)

// PrometheusClient queries a Prometheus HTTP API endpoint.
type PrometheusClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewPrometheusClient creates a new Prometheus metrics client.
func NewPrometheusClient(baseURL string) *PrometheusClient {
	return &PrometheusClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Query executes an instant PromQL query.
func (p *PrometheusClient) Query(ctx context.Context, query string) (float64, error) {
	// TODO: implement Prometheus HTTP API /api/v1/query
	_ = model.Time(0) // ensure import used
	return 0, fmt.Errorf("Query not yet implemented — see internal/metrics/prometheus.go")
}

// QueryRange executes a range PromQL query.
func (p *PrometheusClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]DataPoint, error) {
	// TODO: implement Prometheus HTTP API /api/v1/query_range
	return nil, fmt.Errorf("QueryRange not yet implemented — see internal/metrics/prometheus.go")
}
