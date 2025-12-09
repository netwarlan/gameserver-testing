package checker

import (
	"context"
	"fmt"

	"gameserver-testing/internal/client"
)

// PlayerSlotsChecker verifies players can join the server
type PlayerSlotsChecker struct{}

// Name returns the checker name
func (c *PlayerSlotsChecker) Name() string {
	return "playerslots"
}

// Check performs the player slots check
func (c *PlayerSlotsChecker) Check(ctx context.Context, querier client.ServerQuerier) CheckResult {
	info, err := querier.QueryInfo()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Failed to query server: %v", err),
		}
	}

	details := map[string]any{
		"players":     info.Players,
		"max_players": info.MaxPlayers,
		"bots":        info.Bots,
	}

	// Check max_players > 0
	if info.MaxPlayers == 0 {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Server has no player slots configured (max_players=0)",
			Details: details,
		}
	}

	// Check server isn't full
	if info.Players >= info.MaxPlayers {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Server is full: %d/%d players", info.Players, info.MaxPlayers),
			Details: details,
		}
	}

	available := int(info.MaxPlayers) - int(info.Players)
	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: fmt.Sprintf("Player slots available: %d/%d (free: %d)", info.Players, info.MaxPlayers, available),
		Details: details,
	}
}
