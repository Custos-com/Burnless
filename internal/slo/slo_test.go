package slo

import (
	"math"
	"testing"
)

func TestErrorBudgetRemaining(t *testing.T) {
	tests := []struct {
		name        string
		target      float64
		observed    float64
		wantApprox  float64
	}{
		{"full budget", 99.9, 99.9, 100.0},
		{"half budget gone", 99.9, 99.85, 50.0},
		{"no budget left", 99.9, 99.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorBudgetRemaining(tt.target, tt.observed)
			if math.Abs(got-tt.wantApprox) > 1.0 {
				t.Errorf("ErrorBudgetRemaining(%.2f, %.2f) = %.2f, want ~%.2f",
					tt.target, tt.observed, got, tt.wantApprox)
			}
		})
	}
}

func TestComputeBurnRate(t *testing.T) {
	rate, err := ComputeBurnRate(99.9, 99.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// observed error rate 1%, allowed 0.1% → burn rate ~10
	if rate < 9 || rate > 11 {
		t.Errorf("ComputeBurnRate(99.9, 99.0) = %.2f, want ~10.0", rate)
	}
}

func TestIsCritical(t *testing.T) {
	if !IsCritical(14.4) { t.Error("14.4x should be critical") }
	if !IsCritical(20.0) { t.Error("20.0x should be critical") }
	if  IsCritical(6.0)  { t.Error("6.0x should not be critical") }
}
