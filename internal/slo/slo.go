package slo

import "fmt"

type BurnRate struct {
	Window   string
	Rate     float64
	Severity string
}

func ErrorBudgetRemaining(targetPct, observedPct float64) float64 {
	if targetPct == observedPct {
		return 100.0
	}
	allowedErrors := 100.0 - targetPct
	actualErrors := 100.0 - observedPct
	if allowedErrors <= 0 {
		return 0
	}
	remaining := (1.0 - (actualErrors / allowedErrors)) * 100.0
	if remaining < 0 {
		return 0
	}
	if remaining > 100 {
		return 100
	}
	return remaining
}

func ComputeBurnRate(targetPct, observedPct float64) (float64, error) {
	allowedErrorRate := 100.0 - targetPct
	if allowedErrorRate <= 0 {
		return 0, fmt.Errorf("SLO target of %.2f%% leaves no error budget", targetPct)
	}
	actualErrorRate := 100.0 - observedPct
	return actualErrorRate / allowedErrorRate, nil
}

func IsCritical(burnRate float64) bool {
	return burnRate >= 14.4
}

func IsWarning(burnRate float64) bool {
	return burnRate >= 6 && burnRate < 14.4
}
