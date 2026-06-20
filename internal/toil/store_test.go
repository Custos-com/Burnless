package toil

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStore_SaveAndRetrieve(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore() error: %v", err)
	}

	event := Event{
		ID:           "evt_001",
		Service:      "payments-api",
		Task:         "manual-rollback",
		Date:         time.Now(),
		DurationMins: 45,
		Trigger:      "bad-deployment",
		Automatable:  true,
		Notes:        "rolled back v2.3.1",
	}

	if err := store.Save(event); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// close the store to prove data survives being closed
	if err := store.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	// reopen and check the event is still there
	store2, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("re-opening store error: %v", err)
	}
	defer func() { _ = store2.Close() }()

	events, err := store2.All()
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Service != "payments-api" {
		t.Fatalf("Service = %q, want payments-api", events[0].Service)
	}
	if events[0].DurationMins != 45 {
		t.Fatalf("DurationMins = %d, want 45", events[0].DurationMins)
	}
}

func TestStore_MultipleEvents(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := OpenStore(dbPath)
	if err != nil {
		t.Fatalf("OpenStore() error: %v", err)
	}
	defer func() { _ = store.Close() }()

	for i := 0; i < 3; i++ {
		event := Event{
			ID:      filepath.Join("evt", string(rune('0'+i))),
			Service: "checkout-api",
			Task:    "service-restart",
		}
		if err := store.Save(event); err != nil {
			t.Fatalf("Save() error: %v", err)
		}
	}

	events, err := store.All()
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}
