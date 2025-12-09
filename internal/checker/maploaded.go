package checker

import (
	"context"
	"fmt"
	"strings"

	"gameserver-testing/internal/client"
)

// MapLoadedChecker verifies a map is loaded and not in loading state
type MapLoadedChecker struct{}

// Name returns the checker name
func (c *MapLoadedChecker) Name() string {
	return "maploaded"
}

// loadingIndicators are common states that indicate a map is not ready
var loadingIndicators = []string{
	"loading",
	"changelevel",
}

// Check performs the map loaded check
func (c *MapLoadedChecker) Check(ctx context.Context, querier client.ServerQuerier) CheckResult {
	info, err := querier.QueryInfo()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Failed to query server: %v", err),
		}
	}

	mapName := strings.TrimSpace(info.Map)
	mapLower := strings.ToLower(mapName)

	// Check for empty map name
	if mapName == "" {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Map not loaded: map name is empty",
			Details: map[string]any{
				"map": mapName,
			},
		}
	}

	// Check for loading states
	for _, indicator := range loadingIndicators {
		if strings.Contains(mapLower, indicator) {
			return CheckResult{
				Name:    c.Name(),
				Passed:  false,
				Message: fmt.Sprintf("Map not loaded or in loading state: '%s'", mapName),
				Details: map[string]any{
					"map": mapName,
				},
			}
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: fmt.Sprintf("Map loaded: %s", mapName),
		Details: map[string]any{
			"map": mapName,
		},
	}
}
