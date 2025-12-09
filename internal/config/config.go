package config

import (
	"fmt"
	"net"
	"time"
)

// Default values
const (
	DefaultPort    = 27015
	DefaultTimeout = 5 * time.Second
)

// Config holds the application configuration
type Config struct {
	Host       string
	Port       int
	Timeout    time.Duration
	Checks     []string
	JSONOutput bool
	Verbose    bool
}

// AllChecks returns all available check names
func AllChecks() []string {
	return []string{"connectivity", "maploaded", "playerslots"}
}

// Validate ensures configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.Timeout < time.Millisecond*100 {
		return fmt.Errorf("timeout must be at least 100ms")
	}

	// Validate check names
	validChecks := make(map[string]bool)
	for _, name := range AllChecks() {
		validChecks[name] = true
	}
	for _, check := range c.Checks {
		if !validChecks[check] {
			return fmt.Errorf("unknown check: %s", check)
		}
	}
	return nil
}

// Address returns the server address in host:port format
func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}
