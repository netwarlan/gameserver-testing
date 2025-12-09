package checker

import (
	"context"
	"fmt"

	"gameserver-testing/internal/client"
)

// CheckResult represents the outcome of a single check
type CheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Checker defines the interface for all health checks
type Checker interface {
	Name() string
	Check(ctx context.Context, querier client.ServerQuerier) CheckResult
}

// Registry holds all available checkers
type Registry struct {
	checkers map[string]Checker
}

// NewRegistry creates a registry with all built-in checkers
func NewRegistry() *Registry {
	r := &Registry{
		checkers: make(map[string]Checker),
	}

	// Register built-in checkers
	r.Register(&ConnectivityChecker{})
	r.Register(&MapLoadedChecker{})
	r.Register(&PlayerSlotsChecker{})

	return r
}

// Register adds a checker to the registry
func (r *Registry) Register(c Checker) {
	r.checkers[c.Name()] = c
}

// Get retrieves a checker by name
func (r *Registry) Get(name string) (Checker, bool) {
	c, ok := r.checkers[name]
	return c, ok
}

// GetAll retrieves multiple checkers by name
func (r *Registry) GetAll(names []string) ([]Checker, error) {
	result := make([]Checker, 0, len(names))
	for _, name := range names {
		c, ok := r.Get(name)
		if !ok {
			return nil, fmt.Errorf("unknown checker: %s", name)
		}
		result = append(result, c)
	}
	return result, nil
}
