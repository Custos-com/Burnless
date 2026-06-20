package toil

import (
	"testing"
	"time"
)

func TestFilterByMonth(t *testing.T) {
	events := []Event{
		{ID: "1", Date: time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)},
		{ID: "2", Date: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
		{ID: "3", Date: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)},  // outside June
		{ID: "4", Date: time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)}, // outside June
	}

	filtered, err := FilterByMonth(events, "2026-06")
	if err != nil {
		t.Fatalf("FilterByMonth() error: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 events in June, got %d", len(filtered))
	}
}

func TestFilterByMonth_InvalidFormat(t *testing.T) {
	_, err := FilterByMonth(nil, "not-a-month")
	if err == nil {
		t.Fatal("expected error for invalid month format, got nil")
	}
}

func TestTotalMinutes(t *testing.T) {
	events := []Event{
		{DurationMins: 30},
		{DurationMins: 45},
		{DurationMins: 15},
	}
	got := TotalMinutes(events)
	if got != 90 {
		t.Fatalf("TotalMinutes() = %d, want 90", got)
	}
}

func TestGroupByTask(t *testing.T) {
	events := []Event{
		{Task: "manual-rollback", DurationMins: 300, Automatable: true}, // 5h
		{Task: "manual-rollback", DurationMins: 300, Automatable: true}, // +5h = 10h total
		{Task: "service-restart", DurationMins: 60, Automatable: false}, // 1h
	}

	summaries := GroupByTask(events, DefaultCostModel())

	if len(summaries) != 2 {
		t.Fatalf("expected 2 task groups, got %d", len(summaries))
	}

	// should be sorted with manual-rollback first (10h > 1h)
	if summaries[0].Task != "manual-rollback" {
		t.Fatalf("expected manual-rollback first, got %s", summaries[0].Task)
	}
	if summaries[0].TotalHours != 10.0 {
		t.Fatalf("manual-rollback hours = %.2f, want 10.0", summaries[0].TotalHours)
	}
	if summaries[0].Priority != "HIGH" {
		t.Fatalf("manual-rollback priority = %s, want HIGH (10h+ and automatable)", summaries[0].Priority)
	}

	if summaries[1].Task != "service-restart" {
		t.Fatalf("expected service-restart second, got %s", summaries[1].Task)
	}
	if summaries[1].Priority != "LOW" {
		t.Fatalf("service-restart priority = %s, want LOW (not automatable)", summaries[1].Priority)
	}
}
