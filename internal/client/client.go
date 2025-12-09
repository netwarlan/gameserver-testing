package client

import (
	"fmt"
	"time"

	"github.com/rumblefrog/go-a2s"
)

// ServerQuerier defines the interface for querying game servers.
// This interface enables mocking for unit tests.
type ServerQuerier interface {
	QueryInfo() (*a2s.ServerInfo, error)
	QueryPlayer() (*a2s.PlayerInfo, error)
	QueryRules() (*a2s.RulesInfo, error)
	Close() error
}

// Client wraps the a2s.Client with our configuration
type Client struct {
	inner   *a2s.Client
	address string
}

// NewClient creates a new A2S client
func NewClient(address string, timeout time.Duration) (*Client, error) {
	inner, err := a2s.NewClient(
		address,
		a2s.TimeoutOption(timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create A2S client: %w", err)
	}

	return &Client{
		inner:   inner,
		address: address,
	}, nil
}

// QueryInfo queries the server for basic information
func (c *Client) QueryInfo() (*a2s.ServerInfo, error) {
	return c.inner.QueryInfo()
}

// QueryPlayer queries the server for player information
func (c *Client) QueryPlayer() (*a2s.PlayerInfo, error) {
	return c.inner.QueryPlayer()
}

// QueryRules queries the server for rules/cvars
func (c *Client) QueryRules() (*a2s.RulesInfo, error) {
	return c.inner.QueryRules()
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.inner.Close()
}

// Address returns the server address
func (c *Client) Address() string {
	return c.address
}
