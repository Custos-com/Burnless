// Package slo implements SLO math: error budget calculation and
// burn rate computation using the multiwindow approach from the
// Google SRE book.
package slo

import "fmt"

// BurnRate represents the rate at which the error budget is being consumed.
// A burn rate of 1.0 means the budget will be exhausted exactly at the end
// of the SLO window. 14.4 means exhaustion in 2 days (for a 30d window).
type BurnRate struct {
	Window   string
	Rate     float64
	Severity string // critical | warning | info
}

// ErrorBudgetRemaining returns the percentage of error budget remaining
// given an SLO target and the observed success rate.
func ErrorBudgetRemaining(targetPct, observedPct float64) float64 {
	allowedErrors := 1 - (targetPct / 100)
	actualErrors  := 1 - (observedPct / 100)
	if allowedErrors == 0 {
		return 0
	}
	remaining := (allowedErrors - actualErrors) / allowedErrors * 100
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ComputeBurnRate computes the burn rate for a given window.
// burnRate = (actual error rate) / (allowed error rate)
func ComputeBurnRate(targetPct, observedPct float64) (float64, error) {
	allowedErrorRate := 1 - (targetPct / 100)
	if allowedErrorRate <= 0 {
		return 0, fmt.Errorf("SLO target of %.2f%% leaves no error budget", targetPct)
	}
	actualErrorRate := 1 - (observedPct / 100)
	return actualErrorRate / allowedErrorRate, nil
}

// IsCritical returns true if the burn rate will exhaust the budget in < 2 days
// for a 30-day window (14.4x = 30 days / 2 days = 15 * ~0.96 ≈ 14.4).
func IsCritical(burnRate float64) bool {
	return burnRate >= 14.4
}

// IsWarning returns true if the burn rate will exhaust the budget in < 5 days
// for a 30-day window (6x = 30 days / 5 days).
func IsWarning(burnRate float64) bool {
	return burnRate >= 6 && burnRate < 14.4
}
