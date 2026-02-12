package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHighUsageFragmentsWordCount(t *testing.T) {
	tests := []struct {
		path     string
		maxWords int
	}{
		{"observe/study_specs.md", 75},
		{"observe/study_impl.md", 75},
		{"act/write_draft_plan.md", 75},
		{"observe/study_task_input.md", 75},
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
