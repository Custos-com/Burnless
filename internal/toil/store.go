package toil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "toil_events"

// DefaultDBPath returns the standard location for the local toil database:
// ~/.burnless/data.db
func DefaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	dir := filepath.Join(home, ".burnless")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating %s: %w", dir, err)
	}
	return filepath.Join(dir, "data.db"), nil
}

// Store wraps a bbolt database for storing toil events.
type Store struct {
	db *bolt.DB
}

// OpenStore opens (or creates) the bbolt database at path.
func OpenStore(path string) (*Store, error) {
	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening database at %s: %w", path, err)
	}

	// make sure the bucket exists before we try to use it
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("creating bucket: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the underlying database file.
func (s *Store) Close() error {
	return s.db.Close()
}

// Save writes one toil event to the database, keyed by its ID.
func (s *Store) Save(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encoding event: %w", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Put([]byte(event.ID), data)
	})
}

// All returns every toil event currently stored.
func (s *Store) All() ([]Event, error) {
	var events []Event

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.ForEach(func(k, v []byte) error {
			var event Event
			if err := json.Unmarshal(v, &event); err != nil {
				return fmt.Errorf("decoding event %s: %w", k, err)
			}
			events = append(events, event)
			return nil
		})
	})

	return events, err
}
