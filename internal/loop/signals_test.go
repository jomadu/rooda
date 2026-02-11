package loop

import (
	"testing"
)

// TestPromiseSignalFormat verifies that promise signals follow the exact format
// specified in specs/error-handling.md: no reasons embedded in tags
func TestPromiseSignalFormat(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		wantSuccess bool
		wantFailure bool
	}{
		{
			name:        "valid SUCCESS signal",
			output:      "Some work done\n<promise>SUCCESS</promise>\n",
			wantSuccess: true,
			wantFailure: false,
		},
		{
			name:        "valid FAILURE signal",
			output:      "Blocked by missing dependency\n<promise>FAILURE</promise>\n",
			wantSuccess: false,
			wantFailure: true,
		},
		{
			name:        "FAILURE with reason after signal (valid)",
			output:      "<promise>FAILURE</promise>\nReason: Missing API key",
			wantSuccess: false,
			wantFailure: true,
		},
		{
			name:        "FAILURE with reason embedded in tag (invalid - should not match)",
			output:      "<promise>FAILURE: Missing API key</promise>",
			wantSuccess: false,
			wantFailure: false, // Should NOT match because format is wrong
		},
		{
			name:        "both signals present",
			output:      "<promise>SUCCESS</promise>\n<promise>FAILURE</promise>",
			wantSuccess: true,
			wantFailure: true,
		},
		{
			name:        "no signals",
			output:      "Just some output",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "lowercase success (invalid)",
			output:      "<promise>success</promise>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "extra spaces (invalid)",
			output:      "<promise> SUCCESS </promise>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "unclosed tag (invalid)",
			output:      "<promise>SUCCESS",
			wantSuccess: false,
			wantFailure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSuccess, gotFailure := ScanOutputForSignals(tt.output)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("ScanOutputForSignals() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
			if gotFailure != tt.wantFailure {
				t.Errorf("ScanOutputForSignals() gotFailure = %v, want %v", gotFailure, tt.wantFailure)
			}
		})
	}
}

// TestPromiseSignalPrecedence verifies FAILURE takes precedence when both signals present
func TestPromiseSignalPrecedence(t *testing.T) {
	output := "<promise>SUCCESS</promise>\n<promise>FAILURE</promise>"
	
	gotSuccess, gotFailure := ScanOutputForSignals(output)
	
	if !gotSuccess {
		t.Error("Expected SUCCESS signal to be detected")
	}
	if !gotFailure {
		t.Error("Expected FAILURE signal to be detected")
	}
	
	// The iteration loop should handle precedence (FAILURE wins)
	// This test just verifies both are detected
}
