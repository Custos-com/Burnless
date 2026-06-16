// Package sdk provides a high-level Go client for interacting with Burnless
// programmatically. Teams can use this to build CI/CD integrations and
// custom tooling on top of Burnless.
package sdk

// Client is the Burnless SDK client.
type Client struct {
	configPath string
}

// New creates a new SDK client pointing at the given sre.yaml path.
func New(configPath string) *Client {
	return &Client{configPath: configPath}
}
