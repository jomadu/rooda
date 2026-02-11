package prompt

import (
	"strings"
	"testing"
)

// TestTop5FragmentsImperativeVoice verifies that the 5 most-used fragments
// use imperative voice with clear commands and tool usage hints.
func TestTop5FragmentsImperativeVoice(t *testing.T) {
	tests := []struct {
		name                string
		path                string
		mustContainPhrases  []string
		mustNotStartWith    []string
	}{
		{
			name: "read_agents_md uses imperative voice",
			path: "fragments/observe/read_agents_md.md",
			mustContainPhrases: []string{
				"You must",
				"Use the",
				"tool",
			},
			mustNotStartWith: []string{
				"Load and parse", // passive start
			},
		},
		{
			name: "read_specs uses imperative voice",
			path: "fragments/observe/read_specs.md",
			mustContainPhrases: []string{
				"You must",
				"Use the",
				"tool",
			},
			mustNotStartWith: []string{
				"Load specification", // passive start
			},
		},
		{
			name: "read_impl uses imperative voice",
			path: "fragments/observe/read_impl.md",
			mustContainPhrases: []string{
				"You must",
				"Use the",
				"tool",
			},
			mustNotStartWith: []string{
				"Load implementation", // passive start
			},
		},
		{
			name: "understand_task_requirements uses imperative voice",
			path: "fragments/orient/understand_task_requirements.md",
			mustContainPhrases: []string{
				"Your task is to",
			},
			mustNotStartWith: []string{
				"Parse task description", // passive start
			},
		},
		{
			name: "plan_implementation_approach uses imperative voice",
			path: "fragments/decide/plan_implementation_approach.md",
			mustContainPhrases: []string{
				"Your task is to",
			},
			mustNotStartWith: []string{
				"Determine how to", // passive start
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := LoadFragment(tt.path, "")
			if err != nil {
				t.Fatalf("failed to load fragment: %v", err)
			}

			contentLower := strings.ToLower(content)

			// Check for required imperative phrases
			for _, phrase := range tt.mustContainPhrases {
				if !strings.Contains(contentLower, strings.ToLower(phrase)) {
					t.Errorf("fragment must contain phrase %q but doesn't", phrase)
				}
			}

			// Check that passive starts are not present
			// Skip the title line (starts with #)
			lines := strings.Split(content, "\n")
			var bodyStart int
			for i, line := range lines {
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "#") {
					bodyStart = i
					break
				}
			}
			
			if bodyStart < len(lines) {
				firstBodyLine := strings.TrimSpace(lines[bodyStart])
				for _, phrase := range tt.mustNotStartWith {
					if strings.HasPrefix(firstBodyLine, phrase) {
						t.Errorf("fragment body must not start with passive phrase %q but does", phrase)
					}
				}
			}
		})
	}
}
