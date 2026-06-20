package toil

import (
	"math"
	"testing"
)

// epsilon tolerates floating-point rounding when comparing dollar amounts.
const epsilon = 0.01

func TestDefaultHourlyRate(t *testing.T) {
	got := DefaultCostModel().HourlyRateUSD()
	const want = 86.54 // $180k / 2080h

	if math.Abs(got-want) > epsilon {
		t.Errorf("default HourlyRateUSD() = %.4f, want ~%.2f", got, want)
	}
}

func TestCostForHours(t *testing.T) {
	tests := []struct {
		name   string
		salary float64
		hours  float64
		want   float64
	}{
		{"default one hour", 0, 1, 86.54},
		{"default ten hours", 0, 10, 865.38},
		{"zero hours", 0, 0, 0},
		{"negative hours clamped to zero", 0, -5, 0},
		{"custom salary one hour", 208_000, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CostModel{AnnualSalaryUSD: tt.salary}
			got := m.CostForHours(tt.hours)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("CostForHours(%.2f) with salary %.0f = %.4f, want ~%.2f",
					tt.hours, tt.salary, got, tt.want)
			}
		})
	}
}

func TestCostForMinutes(t *testing.T) {
	m := DefaultCostModel()

	if got, want := m.CostForMinutes(30), 86.54/2; math.Abs(got-want) > epsilon {
		t.Errorf("CostForMinutes(30) = %.4f, want ~%.2f", got, want)
	}
	if got := m.CostForMinutes(0); got != 0 {
		t.Errorf("CostForMinutes(0) = %.4f, want 0", got)
	}
	if got := m.CostForMinutes(-15); got != 0 {
		t.Errorf("CostForMinutes(-15) = %.4f, want 0", got)
	}
}

func TestCostModelFallsBackToDefaultSalary(t *testing.T) {
	var m CostModel // zero value — no salary configured

	if got := m.HourlyRateUSD(); math.Abs(got-86.54) > epsilon {
		t.Errorf("zero-value HourlyRateUSD() = %.4f, want default ~86.54", got)
	}
}

func TestCostForEvent(t *testing.T) {
	m := CostModel{AnnualSalaryUSD: 208_000} // $100/hr
	e := Event{Service: "checkout-api", Task: "manual-restart", DurationMins: 120}

	if got := m.Cost(e); math.Abs(got-200) > epsilon {
		t.Errorf("Cost(event 120m) = %.4f, want ~200", got)
	}
}
