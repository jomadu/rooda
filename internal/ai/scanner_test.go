package ai

import (
	"testing"
)

func TestScanOutputForSignals(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		wantSuccess bool
		wantFailure bool
	}{
		{
			name:        "success signal on own line",
			output:      "Some output\n<promise>SUCCESS</promise>\nMore output",
			wantSuccess: true,
			wantFailure: false,
		},
		{
			name:        "failure signal on own line",
			output:      "Some output\n<promise>FAILURE</promise>\nMore output",
			wantSuccess: false,
			wantFailure: true,
		},
		{
			name:        "both signals present",
			output:      "Some output\n<promise>SUCCESS</promise>\n<promise>FAILURE</promise>\n",
			wantSuccess: true,
			wantFailure: true,
		},
		{
			name:        "no signals",
			output:      "No signals here",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "success signal with surrounding text on same line",
			output:      "Task complete <promise>SUCCESS</promise> - all tests passing",
			wantSuccess: true,
			wantFailure: false,
		},
		{
			name:        "case sensitive - lowercase not recognized",
			output:      "<promise>success</promise>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "extra spaces not recognized",
			output:      "<promise> SUCCESS </promise>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "uppercase tags not recognized",
			output:      "<PROMISE>SUCCESS</PROMISE>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "unclosed tag not recognized",
			output:      "<promise>SUCCESS",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "mixed case not recognized",
			output:      "<Promise>Success</Promise>",
			wantSuccess: false,
			wantFailure: false,
		},
		{
			name:        "signal at end of output",
			output:      "Some work done\n<promise>SUCCESS</promise>",
			wantSuccess: true,
			wantFailure: false,
		},
		{
			name:        "signal at beginning of output",
			output:      "<promise>SUCCESS</promise>\nSome work done",
			wantSuccess: true,
			wantFailure: false,
		},
		{
			name:        "multiple success signals",
			output:      "<promise>SUCCESS</promise>\nMore work\n<promise>SUCCESS</promise>",
			wantSuccess: true,
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
