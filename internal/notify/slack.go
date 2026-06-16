package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackNotifier sends notifications to a Slack channel via webhook.
type SlackNotifier struct {
	WebhookURL string
	Channel    string
}

// Send posts a notification to Slack.
func (s *SlackNotifier) Send(ctx context.Context, n Notification) error {
	payload := map[string]string{
		"channel": s.Channel,
		"text":    fmt.Sprintf("*[%s] %s — %s*\n%s", n.Severity, n.Service, n.Title, n.Body),
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned %d", resp.StatusCode)
	}
	return nil
}
