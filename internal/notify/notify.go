// Package notify provides notification adapters for Slack, PagerDuty, and email.
package notify

import "context"

// Notification is a message to send to a notification target.
type Notification struct {
	Title   string
	Body    string
	Service string
	SLO     string
	Severity string // critical | warning | info
}

// Notifier is the interface all notification backends must implement.
type Notifier interface {
	Send(ctx context.Context, n Notification) error
}

// Multi sends a notification to multiple backends.
type Multi struct {
	backends []Notifier
}

// NewMulti creates a multi-backend notifier.
func NewMulti(backends ...Notifier) *Multi {
	return &Multi{backends: backends}
}

// Send sends the notification to all backends. Returns the first error encountered.
func (m *Multi) Send(ctx context.Context, n Notification) error {
	for _, b := range m.backends {
		if err := b.Send(ctx, n); err != nil {
			return err
		}
	}
	return nil
}
