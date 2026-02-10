package specs_test

import (
	"os"
	"strings"
	"testing"
)

func TestPromptCompositionSpecHasEnhancedSectionMarkers(t *testing.T) {
	content, err := os.ReadFile("prompt-composition.md")
	if err != nil {
		t.Fatalf("Failed to read prompt-composition.md: %v", err)
	}

	spec := string(content)

	// Test that enhanced section marker format is documented
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "OBSERVE marker with description",
			expected: "Execute these observation tasks to gather information.",
		},
		{
			name:     "ORIENT marker with description",
			expected: "Analyze the information you gathered and form your understanding.",
		},
		{
			name:     "DECIDE marker with description",
			expected: "Make decisions about what actions to take.",
		},
		{
			name:     "ACT marker with description",
			expected: "Execute the actions you decided on. Modify files, run commands, commit changes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(spec, tt.expected) {
				t.Errorf("Spec does not contain expected phase description: %q", tt.expected)
			}
		})
	}

	// Test that enhanced format with double lines is documented
	enhancedFormatIndicators := []string{
		"═══════════════════════════════════════════════════════════════",
		"PHASE 1: OBSERVE",
		"PHASE 2: ORIENT",
		"PHASE 3: DECIDE",
		"PHASE 4: ACT",
	}

	for _, indicator := range enhancedFormatIndicators {
		if !strings.Contains(spec, indicator) {
			t.Errorf("Spec does not contain enhanced format indicator: %q", indicator)
		}
	}

	// Test that examples show enhanced markers
	exampleSectionStart := strings.Index(spec, "## Examples")
	if exampleSectionStart == -1 {
		t.Fatal("Spec does not contain Examples section")
	}

	examplesSection := spec[exampleSectionStart:]
	
	// At least one example should show the enhanced format
	if !strings.Contains(examplesSection, "═══════════════════════════════════════════════════════════════") {
		t.Error("Examples section does not show enhanced section marker format")
	}
}
