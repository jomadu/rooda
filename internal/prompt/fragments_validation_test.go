package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestOutputFragmentsHaveValidation validates that fragments which produce output
// (write_audit_report, write_draft_plan, write_gap_report) include validation steps
// before proceeding to emit_success.
func TestOutputFragmentsHaveValidation(t *testing.T) {
	outputFragments := []string{
		"act/write_audit_report.md",
		"act/write_draft_plan.md",
		"act/write_gap_report.md",
	}
	
	requiredValidationQuestions := []string{
		"minimal",
		"complete",
		"accurate",
	}
	
	requiredInstruction := "only proceed to emit_success if validation passes"
	
	for _, fragmentPath := range outputFragments {
		fullPath := filepath.Join("fragments", fragmentPath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read fragment %s: %v", fullPath, err)
		}
		
		contentLower := strings.ToLower(string(content))
		
		// Check for validation questions
		for _, question := range requiredValidationQuestions {
			if !strings.Contains(contentLower, question) {
				t.Errorf("Fragment %s missing validation question: %s", fragmentPath, question)
			}
		}
		
		// Check for instruction to only proceed if validation passes
		if !strings.Contains(contentLower, requiredInstruction) {
			t.Errorf("Fragment %s missing instruction: %s", fragmentPath, requiredInstruction)
		}
	}
}
