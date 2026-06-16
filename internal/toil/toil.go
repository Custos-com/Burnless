// Package toil implements the Burnless Toil Tracker.
// It logs manual SRE tasks, calculates their cost, and generates
// manager-facing reports showing the ROI of automation.
package toil

import "time"

// Event represents a single logged toil event.
type Event struct {
	ID           string    `yaml:"id"`
	Service      string    `yaml:"service"`
	Task         string    `yaml:"task"`
	Date         time.Time `yaml:"date"`
	DurationMins int       `yaml:"duration_minutes"`
	Trigger      string    `yaml:"trigger"`
	Automatable  bool      `yaml:"automatable"`
	Notes        string    `yaml:"notes,omitempty"`
}

// Report is a monthly summary of toil events.
type Report struct {
	Month              string
	TotalHours         float64
	EstimatedCostUSD   float64
	ToilPercentage     float64
	LinkedIncidents    int
	TotalIncidents     int
	TopToilTasks       []TaskSummary
}

// TaskSummary summarises toil by task type.
type TaskSummary struct {
	Task         string
	TotalHours   float64
	Automatable  bool
	Priority     string // HIGH | MED | LOW
}

// Log stores a toil event.
func Log(event Event) error {
	// TODO: persist to bbolt local store
	return nil
}

// GenerateReport produces a monthly toil report.
func GenerateReport(month string) (*Report, error) {
	// TODO: query bbolt, aggregate by task, compute cost
	// Assume average SRE salary $180k/year = $86.54/hr
	return &Report{Month: month}, nil
}
