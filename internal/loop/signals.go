package loop

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// ScanOutputForSignals scans AI CLI output for promise signals.
// Returns (hasSuccess, hasFailure) booleans.
//
// Valid signals (case-sensitive, exact match):
//   - <promise>SUCCESS</promise>
//   - <promise>FAILURE</promise>
//
// Invalid signals (not recognized):
//   - <promise>success</promise> (lowercase)
//   - <promise> SUCCESS </promise> (extra spaces)
//   - <promise>FAILURE: reason</promise> (reason embedded in tag)
//   - <PROMISE>SUCCESS</PROMISE> (uppercase tags)
//   - <promise>SUCCESS (unclosed tag)
//
// Rationale: Strict format forces AI to follow exact specification,
// prevents ambiguity, and enables fast string matching.
//
// If explanations are needed, they should come AFTER the signal:
//   <promise>FAILURE</promise>
//   Reason: Missing API key configuration
func ScanOutputForSignals(output string) (hasSuccess bool, hasFailure bool) {
	hasSuccess = strings.Contains(output, "<promise>SUCCESS</promise>")
	hasFailure = strings.Contains(output, "<promise>FAILURE</promise>")
	return hasSuccess, hasFailure
}

// SetupSignalHandler sets up signal handling for SIGINT and SIGTERM.
// Returns a channel that will receive signals.
func SetupSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}
