package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/rumblefrog/go-a2s"

	"gameserver-testing/internal/checker"
)

// TestReport is the complete JSON output structure
type TestReport struct {
	Timestamp time.Time              `json:"timestamp"`
	Server    ServerInfo             `json:"server"`
	Results   []checker.CheckResult  `json:"results"`
	AllPassed bool                   `json:"all_passed"`
	ExitCode  int                    `json:"exit_code"`
}

// ServerInfo contains the queried server information
type ServerInfo struct {
	Address    string `json:"address"`
	Name       string `json:"name,omitempty"`
	Map        string `json:"map,omitempty"`
	Game       string `json:"game,omitempty"`
	Players    uint8  `json:"players,omitempty"`
	MaxPlayers uint8  `json:"max_players,omitempty"`
	Bots       uint8  `json:"bots,omitempty"`
	ServerType string `json:"server_type,omitempty"`
	VAC        bool   `json:"vac,omitempty"`
}

// FromA2SInfo converts a2s.ServerInfo to our ServerInfo
func FromA2SInfo(info *a2s.ServerInfo, address string) ServerInfo {
	if info == nil {
		return ServerInfo{Address: address}
	}
	return ServerInfo{
		Address:    address,
		Name:       info.Name,
		Map:        info.Map,
		Game:       info.Game,
		Players:    info.Players,
		MaxPlayers: info.MaxPlayers,
		Bots:       info.Bots,
		ServerType: info.ServerType.String(),
		VAC:        info.VAC,
	}
}

// Writer handles output formatting
type Writer struct {
	out      io.Writer
	jsonMode bool
	verbose  bool
}

// NewWriter creates a new output writer
func NewWriter(out io.Writer, jsonMode, verbose bool) *Writer {
	return &Writer{
		out:      out,
		jsonMode: jsonMode,
		verbose:  verbose,
	}
}

// WriteReport writes the test report in the configured format
func (w *Writer) WriteReport(report TestReport) error {
	if w.jsonMode {
		return w.writeJSON(report)
	}
	return w.writeText(report)
}

func (w *Writer) writeJSON(report TestReport) error {
	encoder := json.NewEncoder(w.out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func (w *Writer) writeText(report TestReport) error {
	// Header
	fmt.Fprintf(w.out, "Testing server: %s\n", report.Server.Address)
	if report.Server.Name != "" {
		fmt.Fprintf(w.out, "Server name: %s\n", report.Server.Name)
	}
	fmt.Fprintf(w.out, "\n")

	// Results
	for _, result := range report.Results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(w.out, "[%s] %s: %s\n", status, result.Name, result.Message)
	}

	// Summary
	fmt.Fprintf(w.out, "\n")
	if report.AllPassed {
		fmt.Fprintf(w.out, "All checks passed.\n")
	} else {
		fmt.Fprintf(w.out, "Some checks failed.\n")
	}

	return nil
}
