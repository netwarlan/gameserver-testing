package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"gameserver-testing/internal/checker"
	"gameserver-testing/internal/client"
	"gameserver-testing/internal/config"
	"gameserver-testing/internal/output"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "gstest [host]",
	Short: "Test Source engine game server health",
	Long: `gstest is a testing tool for Source engine game servers.
It uses the A2S protocol to verify server health and readiness.

Example:
  gstest 192.168.1.100 --port 27015 --checks connectivity,maploaded`,
	Args: cobra.ExactArgs(1),
	RunE: runTests,
}

func init() {
	rootCmd.Flags().IntVarP(&cfg.Port, "port", "p", config.DefaultPort,
		"Server port")
	rootCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", config.DefaultTimeout,
		"Query timeout")
	rootCmd.Flags().StringSliceVarP(&cfg.Checks, "checks", "c", config.AllChecks(),
		"Checks to run (comma-separated)")
	rootCmd.Flags().BoolVarP(&cfg.JSONOutput, "json", "j", false,
		"Output results as JSON")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false,
		"Verbose output")

	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(output.ExitConfigError)
	}
}

func runTests(cmd *cobra.Command, args []string) error {
	cfg.Host = args[0]

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(output.ExitConfigError)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout*2)
	defer cancel()

	// Create A2S client
	a2sClient, err := client.NewClient(cfg.Address(), cfg.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client: %v\n", err)
		os.Exit(output.ExitConnectivityFailed)
		return nil
	}
	defer a2sClient.Close()

	// Get checkers
	registry := checker.NewRegistry()
	checkers, err := registry.GetAll(cfg.Checks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid checks: %v\n", err)
		os.Exit(output.ExitConfigError)
		return nil
	}

	// Run checks
	results := make([]checker.CheckResult, 0, len(checkers))
	var failedChecks []string

	for _, c := range checkers {
		result := c.Check(ctx, a2sClient)
		results = append(results, result)
		if !result.Passed {
			failedChecks = append(failedChecks, c.Name())
		}
	}

	// Get server info for report (best effort)
	serverInfo, _ := a2sClient.QueryInfo()

	// Build report
	report := output.TestReport{
		Timestamp: time.Now().UTC(),
		Server:    output.FromA2SInfo(serverInfo, cfg.Address()),
		Results:   results,
		AllPassed: len(failedChecks) == 0,
		ExitCode:  output.GetExitCode(failedChecks),
	}

	// Write output
	writer := output.NewWriter(os.Stdout, cfg.JSONOutput, cfg.Verbose)
	if err := writer.WriteReport(report); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
		os.Exit(output.ExitUnknownError)
		return nil
	}

	os.Exit(report.ExitCode)
	return nil
}
