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

// LoadContextContent loads context content, checking if the value is a file path.
// Returns (content, isFile, error). If the file exists, reads and returns its content
// with isFile=true. Otherwise, treats the value as inline content with isFile=false.
func LoadContextContent(contextValue string) (string, bool, error) {
	// Check if file exists
	if _, err := os.Stat(contextValue); err == nil {
		// File exists - read content
		content, err := os.ReadFile(contextValue)
		if err != nil {
			return "", false, fmt.Errorf("failed to read context file %s: %v", contextValue, err)
		}
		return string(content), true, nil
	}
	
	// Not a file - treat as inline content
	return contextValue, false, nil
}
