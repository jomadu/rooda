package prompt

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed fragments/*/*.md
var embeddedFS embed.FS

// LoadFragment loads a fragment from either embedded resources (builtin: prefix)
// or filesystem (relative to configDir).
func LoadFragment(path string, configDir string) (string, error) {
	// Check for builtin: prefix
	if strings.HasPrefix(path, "builtin:") {
		embeddedPath := strings.TrimPrefix(path, "builtin:")
		content, err := embeddedFS.ReadFile(embeddedPath)
		if err != nil {
			return "", fmt.Errorf("embedded fragment not found: %s", path)
		}
		return string(content), nil
	}

	// Filesystem path - resolve relative to config directory
	absPath := filepath.Join(configDir, path)
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("fragment file not found: %s (resolved to: %s)", path, absPath)
	}
	return string(content), nil
}
