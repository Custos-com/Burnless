package toil

import (
	"fmt"
	"sort"
	"time"
)

// FilterByMonth returns only the events whose Date falls within the given
// month. month must be in "YYYY-MM" format, e.g. "2026-06".
func FilterByMonth(events []Event, month string) ([]Event, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month %q, expected format YYYY-MM: %w", month, err)
	}
	end := start.AddDate(0, 1, 0) // first day of the NEXT month

	var filtered []Event
	for _, e := range events {
		if !e.Date.Before(start) && e.Date.Before(end) {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// TotalMinutes sums the DurationMins of all events.
func TotalMinutes(events []Event) int {
	total := 0
	for _, e := range events {
		total += e.DurationMins
	}
	return total
}

// GroupByTask aggregates events into TaskSummary entries, sorted by total
// time spent (descending) so the biggest sources of toil appear first.
func GroupByTask(events []Event, model CostModel) []TaskSummary {
	totals := map[string]int{}       // task name -> total minutes
	automatable := map[string]bool{} // task name -> is it automatable

	for _, e := range events {
		totals[e.Task] += e.DurationMins
		if e.Automatable {
			automatable[e.Task] = true
		}
	}

	summaries := make([]TaskSummary, 0, len(totals))
	for task, minutes := range totals {
		hours := float64(minutes) / 60.0
		summaries = append(summaries, TaskSummary{
			Task:        task,
			TotalHours:  hours,
			Automatable: automatable[task],
			Priority:    priorityFor(hours, automatable[task]),
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalHours > summaries[j].TotalHours
	})

	return summaries
}

// priorityFor assigns a HIGH/MED/LOW automation priority based on how much
// time a task costs and whether it is automatable at all.
func priorityFor(hours float64, automatable bool) string {
	if !automatable {
		return "LOW"
	}
	switch {
	case hours >= 10:
		return "HIGH"
	case hours >= 3:
		return "MED"
	default:
		return "LOW"
	}
}
