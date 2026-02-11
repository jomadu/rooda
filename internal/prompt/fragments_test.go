package prompt

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestFragmentsUseImperativeVoice validates that all fragments use imperative voice
// and follow the established pattern for executable procedures.
func TestFragmentsUseImperativeVoice(t *testing.T) {
	fragmentsDir := "fragments"
	
	// Strong imperative patterns that indicate direct commands
	strongImperativePatterns := []string{
		"you must",
		"your task is",
		"use the",
		"execute the",
		"load the",
		"modify the",
		"create the",
		"update the",
	}
	
	// Passive patterns that indicate documentation rather than instructions
	passivePatterns := []string{
		"this phase",
		"this section",
		"the agent should",
		"the system will",
		"is used to",
		"can be used",
		"allows you to",
	}
	
	phases := []string{"observe", "orient", "decide", "act"}
	
	for _, phase := range phases {
		phaseDir := filepath.Join(fragmentsDir, phase)
		entries, err := os.ReadDir(phaseDir)
		if err != nil {
			t.Fatalf("Failed to read %s directory: %v", phaseDir, err)
		}
		
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			
			fragmentPath := filepath.Join(phaseDir, entry.Name())
			content, err := os.ReadFile(fragmentPath)
			if err != nil {
				t.Errorf("Failed to read fragment %s: %v", fragmentPath, err)
				continue
			}
			
			contentLower := strings.ToLower(string(content))
			
			// Check for strong imperative patterns
			hasStrongImperative := false
			for _, pattern := range strongImperativePatterns {
				if strings.Contains(contentLower, pattern) {
					hasStrongImperative = true
					break
				}
			}
			
			if !hasStrongImperative {
				t.Errorf("Fragment %s does not use strong imperative voice (should start with 'You must...' or 'Your task is...')", fragmentPath)
			}
			
			// Check for passive patterns (should not be present)
			for _, pattern := range passivePatterns {
				if strings.Contains(contentLower, pattern) {
					t.Errorf("Fragment %s uses passive documentation language: contains '%s'", fragmentPath, pattern)
				}
			}
			
			// Check that content is substantial (not just a title)
			lines := strings.Split(string(content), "\n")
			nonEmptyLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "#") {
					nonEmptyLines++
				}
			}
			
			if nonEmptyLines < 2 {
				t.Errorf("Fragment %s appears to be incomplete (less than 2 non-empty, non-header lines)", fragmentPath)
			}
		}
	}
}


// TestPromiseSignalFormat validates that fragments use exact promise signal format
// with no variations and explanations come after the signal.
func TestPromiseSignalFormat(t *testing.T) {
	fragmentsDir := "fragments"
	
	// Exact formats required
	successSignal := "<promise>SUCCESS</promise>"
	failureSignal := "<promise>FAILURE</promise>"
	
	// Invalid patterns that should NOT appear
	invalidPatterns := []*regexp.Regexp{
		regexp.MustCompile(`<promise>SUCCESS:.*?</promise>`),  // SUCCESS with reason inside
		regexp.MustCompile(`<promise>FAILURE:.*?</promise>`),  // FAILURE with reason inside
		regexp.MustCompile(`<promise>success</promise>`),      // lowercase
		regexp.MustCompile(`<promise>failure</promise>`),      // lowercase
		regexp.MustCompile(`<Promise>.*?</Promise>`),          // wrong case
	}
	
	phases := []string{"observe", "orient", "decide", "act"}
	
	for _, phase := range phases {
		phaseDir := filepath.Join(fragmentsDir, phase)
		entries, err := os.ReadDir(phaseDir)
		if err != nil {
			t.Fatalf("Failed to read %s directory: %v", phaseDir, err)
		}
		
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			
			fragmentPath := filepath.Join(phaseDir, entry.Name())
			content, err := os.ReadFile(fragmentPath)
			if err != nil {
				t.Errorf("Failed to read fragment %s: %v", fragmentPath, err)
				continue
			}
			
			contentStr := string(content)
			
			// Check for invalid patterns
			for _, pattern := range invalidPatterns {
				if pattern.MatchString(contentStr) {
					t.Errorf("Fragment %s contains invalid promise signal format: %s", 
						fragmentPath, pattern.String())
				}
			}
			
			// If fragment mentions SUCCESS or FAILURE signals, verify exact format
			if strings.Contains(contentStr, "SUCCESS") || strings.Contains(contentStr, "FAILURE") {
				// Check emit_success.md specifically
				if entry.Name() == "emit_success.md" {
					if !strings.Contains(contentStr, successSignal) {
						t.Errorf("Fragment %s must show exact SUCCESS signal format: %s", 
							fragmentPath, successSignal)
					}
					
					// Must have examples showing signal with explanation after
					if !strings.Contains(contentStr, "Example") {
						t.Errorf("Fragment %s must include examples showing signal format", 
							fragmentPath)
					}
					
					// Must clarify SUCCESS means procedure goal achieved
					if !strings.Contains(contentStr, "procedure") && !strings.Contains(contentStr, "goal") {
						t.Errorf("Fragment %s must clarify SUCCESS means procedure goal achieved", 
							fragmentPath)
					}
				}
				
				// Check emit_failure.md specifically
				if entry.Name() == "emit_failure.md" {
					if !strings.Contains(contentStr, failureSignal) {
						t.Errorf("Fragment %s must show exact FAILURE signal format: %s", 
							fragmentPath, failureSignal)
					}
					
					// Must have examples showing signal with explanation after
					if !strings.Contains(contentStr, "Example") {
						t.Errorf("Fragment %s must include examples showing signal format", 
							fragmentPath)
					}
				}
				
				// Check check_if_blocked.md specifically
				if entry.Name() == "check_if_blocked.md" {
					if strings.Contains(contentStr, "FAILURE") {
						if !strings.Contains(contentStr, failureSignal) {
							t.Errorf("Fragment %s must show exact FAILURE signal format: %s", 
								fragmentPath, failureSignal)
						}
						
						// Should not say "emit FAILURE with explanation" - should separate signal and explanation
						if strings.Contains(strings.ToLower(contentStr), "with explanation") {
							t.Errorf("Fragment %s should separate signal from explanation, not combine them", 
								fragmentPath)
						}
					}
				}
			}
		}
	}
}
