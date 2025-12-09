package checker

import (
	"context"
	"fmt"

	"gameserver-testing/internal/client"
)

// ConnectivityChecker verifies basic A2S_INFO response
type ConnectivityChecker struct{}

// Name returns the checker name
func (c *ConnectivityChecker) Name() string {
	return "connectivity"
}

// Check performs the connectivity check
func (c *ConnectivityChecker) Check(ctx context.Context, querier client.ServerQuerier) CheckResult {
	info, err := querier.QueryInfo()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Failed to query server: %v", err),
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: fmt.Sprintf("Server responded: %s", info.Name),
		Details: map[string]any{
			"server_name": info.Name,
			"game":        info.Game,
			"protocol":    info.Protocol,
		},
	}
}
