// Package toil implements the Burnless Toil Tracker.
// It logs manual SRE tasks, calculates their cost, and generates
// manager-facing reports showing the ROI of automation.
package toil

import "time"

// Event represents a single logged toil event.
type Event struct {
	ID           string    `yaml:"id"            json:"id"`
	Service      string    `yaml:"service"       json:"service"`
	Task         string    `yaml:"task"          json:"task"`
	Date         time.Time `yaml:"date"          json:"date"`
	DurationMins int       `yaml:"duration_minutes" json:"duration_minutes"`
	Trigger      string    `yaml:"trigger"       json:"trigger"`
	Automatable  bool      `yaml:"automatable"   json:"automatable"`
	Notes        string    `yaml:"notes,omitempty" json:"notes,omitempty"`
}

// Report is a monthly summary of toil events.
type Report struct {
	Month            string
	TotalHours       float64
	EstimatedCostUSD float64
	ToilPercentage   float64
	LinkedIncidents  int
	TotalIncidents   int
	TopToilTasks     []TaskSummary
}

// TaskSummary summarises toil by task type.
type TaskSummary struct {
	Task        string
	TotalHours  float64
	Automatable bool
	Priority    string // HIGH | MED | LOW
}
