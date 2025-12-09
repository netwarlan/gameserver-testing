package output

// Exit codes for different failure scenarios
const (
	ExitSuccess            = 0
	ExitConnectivityFailed = 1
	ExitMapNotLoaded       = 2
	ExitServerFull         = 3
	ExitConfigError        = 10
	ExitUnknownError       = 99
)

// CheckNameToExitCode maps check names to their failure exit codes
var CheckNameToExitCode = map[string]int{
	"connectivity": ExitConnectivityFailed,
	"maploaded":    ExitMapNotLoaded,
	"playerslots":  ExitServerFull,
}

// GetExitCode returns the appropriate exit code for failed checks
func GetExitCode(failedChecks []string) int {
	if len(failedChecks) == 0 {
		return ExitSuccess
	}

	// Return the exit code for the first failed check
	// (maintains predictable behavior)
	if code, ok := CheckNameToExitCode[failedChecks[0]]; ok {
		return code
	}
	return ExitUnknownError
}
