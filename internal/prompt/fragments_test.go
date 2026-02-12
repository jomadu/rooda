package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFragmentWordCount(t *testing.T) {
	tests := []struct {
		path     string
		maxWords int
	}{
		{"observe/study_agents_md.md", 100},
		{"decide/decide_signal.md", 100},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			fullPath := filepath.Join("fragments", tt.path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("failed to read fragment %s: %v", tt.path, err)
			}

			wordCount := countWords(string(content))
			if wordCount > tt.maxWords {
				t.Errorf("fragment %s has %d words, expected ≤%d", tt.path, wordCount, tt.maxWords)
			}
		})
	}
}

func TestStudyAgentsMdPreservesTopics(t *testing.T) {
	content, err := os.ReadFile("fragments/observe/study_agents_md.md")
	if err != nil {
		t.Fatalf("failed to read study_agents_md.md: %v", err)
	}

	text := string(content)
	requiredTopics := []string{
		"Work Tracking",
		"Quick Reference",
		"Task Input",
		"Planning System",
		"Build",
		"Test",
		"Lint",
		"Specification",
		"Implementation",
		"Audit",
		"Quality Criteria",
		"Operational Learnings",
	}

	for _, topic := range requiredTopics {
		if !strings.Contains(text, topic) {
			t.Errorf("fragment missing required topic: %s", topic)
		}
	}
}

func countWords(s string) int {
	return len(strings.Fields(s))
}
