package prompt

import (
	"strings"
	"testing"
)

// TestEmitFailureFragmentFormat verifies that emit_failure.md contains
// the exact <promise>FAILURE</promise> signal format with clear instructions
func TestEmitFailureFragmentFormat(t *testing.T) {
	content, err := embeddedFS.ReadFile("fragments/act/emit_failure.md")
	if err != nil {
		t.Fatalf("failed to read emit_failure.md: %v", err)
	}

	contentStr := string(content)
	contentLower := strings.ToLower(contentStr)

	// Must contain exact signal format
	if !strings.Contains(contentStr, "<promise>FAILURE</promise>") {
		t.Error("emit_failure.md must contain exact signal format: <promise>FAILURE</promise>")
	}

	// Must clarify that signal is exact and explanation comes after
	if !strings.Contains(contentLower, "exact") || !strings.Contains(contentLower, "after") {
		t.Error("emit_failure.md must clarify that signal is exact and explanation comes after")
	}

	// Must show example with explanation after signal
	if !strings.Contains(contentStr, "Example:") {
		t.Error("emit_failure.md must include an example showing explanation after signal")
	}
}

// TestCheckIfBlockedFragmentFormat verifies that check_if_blocked.md contains
// the exact <promise>FAILURE</promise> signal format with clear instructions
func TestCheckIfBlockedFragmentFormat(t *testing.T) {
	content, err := embeddedFS.ReadFile("fragments/decide/check_if_blocked.md")
	if err != nil {
		t.Fatalf("failed to read check_if_blocked.md: %v", err)
	}

	contentStr := string(content)

	// Must contain exact signal format
	if !strings.Contains(contentStr, "<promise>FAILURE</promise>") {
		t.Error("check_if_blocked.md must contain exact signal format: <promise>FAILURE</promise>")
	}

	// Must clarify that signal is exact and explanation comes after
	if !strings.Contains(contentStr, "exact") && !strings.Contains(contentStr, "exactly") {
		t.Error("check_if_blocked.md must clarify that signal format is exact")
	}

	// Should not use vague language like "emit FAILURE with explanation"
	if strings.Contains(contentStr, "emit FAILURE with explanation") {
		t.Error("check_if_blocked.md should not use vague language 'emit FAILURE with explanation'")
	}
}
