package spec_test

import (
	"os"
	"strings"
	"testing"
)

func TestDistributionSpecDocumentsMakefile(t *testing.T) {
	content, err := os.ReadFile("../../specs/distribution.md")
	if err != nil {
		t.Fatalf("failed to read distribution.md: %v", err)
	}

	spec := string(content)

	t.Run("mentions Makefile in build process", func(t *testing.T) {
		if !strings.Contains(spec, "make build") {
			t.Error("distribution.md should mention 'make build' in build process")
		}
	})

	t.Run("includes make in build-time dependencies", func(t *testing.T) {
		if !strings.Contains(spec, "make") {
			t.Error("distribution.md should include 'make' in dependencies")
		}
	})

	t.Run("shows both make and direct go build examples", func(t *testing.T) {
		hasMakeBuild := strings.Contains(spec, "make build")
		hasGoBuild := strings.Contains(spec, "go build")
		
		if !hasMakeBuild || !hasGoBuild {
			t.Error("distribution.md should show both 'make build' and 'go build' examples")
		}
	})
}
