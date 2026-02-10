package prompt

import (
	"os"
	"strings"
	"testing"

	"github.com/jomadu/rooda/internal/config"
)

func TestComposePhasePrompt_SingleFragment(t *testing.T) {
	fragments := []config.FragmentAction{
		{Path: "builtin:fragments/observe/read_agents_md.md"},
	}

	result, err := ComposePhasePrompt(fragments, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !strings.Contains(result, "# Read AGENTS.md") {
		t.Errorf("expected fragment content, got: %s", result)
	}
}

func TestComposePhasePrompt_MultipleFragments(t *testing.T) {
	fragments := []config.FragmentAction{
		{Path: "builtin:fragments/observe/read_agents_md.md"},
		{Path: "builtin:fragments/observe/read_specs.md"},
	}

	result, err := ComposePhasePrompt(fragments, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Should contain content from both fragments separated by double newlines
	if !strings.Contains(result, "# Read AGENTS.md") {
		t.Errorf("expected first fragment content")
	}
	if !strings.Contains(result, "# Read Specifications") {
		t.Errorf("expected second fragment content")
	}
	if !strings.Contains(result, "\n\n") {
		t.Errorf("expected double newlines between fragments")
	}
}

func TestComposePhasePrompt_InlineContent(t *testing.T) {
	fragments := []config.FragmentAction{
		{Content: "This is inline content."},
	}

	result, err := ComposePhasePrompt(fragments, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != "This is inline content." {
		t.Errorf("expected inline content, got: %s", result)
	}
}

func TestComposePhasePrompt_WithTemplate(t *testing.T) {
	fragments := []config.FragmentAction{
		{
			Content:    "Hello {{.name}}",
			Parameters: map[string]interface{}{"name": "World"},
		},
	}

	result, err := ComposePhasePrompt(fragments, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got: %s", result)
	}
}

func TestComposePhasePrompt_EmptyFragments(t *testing.T) {
	fragments := []config.FragmentAction{}

	result, err := ComposePhasePrompt(fragments, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got: %s", result)
	}
}

func TestAssemblePrompt_AllPhases(t *testing.T) {
	procedure := config.Procedure{
		Observe: []config.FragmentAction{
			{Content: "Observe content"},
		},
		Orient: []config.FragmentAction{
			{Content: "Orient content"},
		},
		Decide: []config.FragmentAction{
			{Content: "Decide content"},
		},
		Act: []config.FragmentAction{
			{Content: "Act content"},
		},
	}

	result, err := AssemblePrompt(procedure, "", "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check enhanced section markers with double lines and phase descriptions
	expectedMarkers := []string{
		"═══════════════════════════════════════════════════════════════\nPHASE 1: OBSERVE\nExecute these observation tasks to gather information.\n═══════════════════════════════════════════════════════════════",
		"═══════════════════════════════════════════════════════════════\nPHASE 2: ORIENT\nAnalyze the information you gathered and form your understanding.\n═══════════════════════════════════════════════════════════════",
		"═══════════════════════════════════════════════════════════════\nPHASE 3: DECIDE\nMake decisions about what actions to take.\n═══════════════════════════════════════════════════════════════",
		"═══════════════════════════════════════════════════════════════\nPHASE 4: ACT\nExecute the actions you decided on. Modify files, run commands, commit changes.\n═══════════════════════════════════════════════════════════════",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(result, marker) {
			t.Errorf("expected enhanced section marker:\n%s\n\nGot result:\n%s", marker, result)
		}
	}

	// Check content
	if !strings.Contains(result, "Observe content") {
		t.Errorf("expected observe content")
	}
	if !strings.Contains(result, "Orient content") {
		t.Errorf("expected orient content")
	}
	if !strings.Contains(result, "Decide content") {
		t.Errorf("expected decide content")
	}
	if !strings.Contains(result, "Act content") {
		t.Errorf("expected act content")
	}
}

func TestAssemblePrompt_WithUserContext(t *testing.T) {
	procedure := config.Procedure{
		Display: "Test Procedure",
		Observe: []config.FragmentAction{
			{Content: "Observe content"},
		},
		Orient: []config.FragmentAction{
			{Content: "Orient content"},
		},
		Decide: []config.FragmentAction{
			{Content: "Decide content"},
		},
		Act: []config.FragmentAction{
			{Content: "Act content"},
		},
	}

	userContext := "Focus on authentication module"
	result, err := AssemblePrompt(procedure, userContext, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check preamble appears first
	if !strings.Contains(result, "ROODA PROCEDURE EXECUTION") {
		t.Errorf("expected preamble at start")
	}

	// Check context appears after preamble with marker
	if !strings.Contains(result, "=== CONTEXT ===") {
		t.Errorf("expected CONTEXT section marker")
	}
	if !strings.Contains(result, "Focus on authentication module") {
		t.Errorf("expected user context content")
	}

	// Preamble should appear before context
	preambleIdx := strings.Index(result, "ROODA PROCEDURE EXECUTION")
	contextIdx := strings.Index(result, "=== CONTEXT ===")
	if preambleIdx == -1 || contextIdx == -1 || preambleIdx >= contextIdx {
		t.Errorf("expected preamble before CONTEXT")
	}

	// Context should appear before OBSERVE
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")
	if contextIdx == -1 || observeIdx == -1 || contextIdx >= observeIdx {
		t.Errorf("expected CONTEXT before OBSERVE")
	}
}

func TestAssemblePrompt_EmptyPhases(t *testing.T) {
	procedure := config.Procedure{
		Observe: []config.FragmentAction{},
		Orient:  []config.FragmentAction{},
		Decide:  []config.FragmentAction{},
		Act:     []config.FragmentAction{},
	}

	result, err := AssemblePrompt(procedure, "", "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Empty phases should not produce section markers
	if strings.Contains(result, "PHASE 1: OBSERVE") {
		t.Errorf("expected no OBSERVE marker for empty phase")
	}
	// Preamble should still be present
	if !strings.Contains(result, "ROODA PROCEDURE EXECUTION") {
		t.Errorf("expected preamble even with empty phases")
	}
}

func TestAssemblePrompt_WithContextFile(t *testing.T) {
	// Create temp file for context
	tmpFile, err := os.CreateTemp("", "context-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	contextContent := "this repository should use make"
	if _, err := tmpFile.WriteString(contextContent); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	procedure := config.Procedure{
		Observe: []config.FragmentAction{{Content: "Observe content"}},
		Orient:  []config.FragmentAction{{Content: "Orient content"}},
		Decide:  []config.FragmentAction{{Content: "Decide content"}},
		Act:     []config.FragmentAction{{Content: "Act content"}},
	}
	
	result, err := AssemblePrompt(procedure, tmpFile.Name(), "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	
	// Check for Source: line
	if !strings.Contains(result, "Source: "+tmpFile.Name()) {
		t.Errorf("expected Source line with file path")
	}
	
	// Check for file content
	if !strings.Contains(result, contextContent) {
		t.Errorf("expected file content in prompt")
	}
	
	// Verify Source line comes before content
	sourceIdx := strings.Index(result, "Source:")
	contentIdx := strings.Index(result, contextContent)
	if sourceIdx == -1 || contentIdx == -1 || sourceIdx >= contentIdx {
		t.Errorf("expected Source line before content")
	}
}

func TestAssemblePrompt_WithInlineContext(t *testing.T) {
	procedure := config.Procedure{
		Observe: []config.FragmentAction{{Content: "Observe content"}},
		Orient:  []config.FragmentAction{{Content: "Orient content"}},
		Decide:  []config.FragmentAction{{Content: "Decide content"}},
		Act:     []config.FragmentAction{{Content: "Act content"}},
	}
	
	inlineContext := "focus on auth module"
	result, err := AssemblePrompt(procedure, inlineContext, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	
	// Should NOT have Source: line for inline content
	if strings.Contains(result, "Source:") {
		t.Errorf("expected no Source line for inline context")
	}
	
	// Should have inline content directly
	if !strings.Contains(result, inlineContext) {
		t.Errorf("expected inline context in prompt")
	}
}

func TestLoadContextContent_File(t *testing.T) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "context-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	expected := "file content here"
	if _, err := tmpFile.WriteString(expected); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	content, isFile, err := LoadContextContent(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	
	if !isFile {
		t.Errorf("expected isFile=true")
	}
	
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestLoadContextContent_Inline(t *testing.T) {
	inline := "inline context text"
	
	content, isFile, err := LoadContextContent(inline)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	
	if isFile {
		t.Errorf("expected isFile=false")
	}
	
	if content != inline {
		t.Errorf("expected %q, got %q", inline, content)
	}
}

// Integration Tests

func TestAssemblePrompt_Integration_SeparatorFormat(t *testing.T) {
	procedure := config.Procedure{
		Observe: []config.FragmentAction{
			{Path: "builtin:fragments/observe/read_agents_md.md"},
			{Path: "builtin:fragments/observe/read_specs.md"},
		},
		Orient: []config.FragmentAction{
			{Path: "builtin:fragments/orient/understand_task_requirements.md"},
		},
		Decide: []config.FragmentAction{
			{Path: "builtin:fragments/decide/plan_implementation_approach.md"},
		},
		Act: []config.FragmentAction{
			{Path: "builtin:fragments/act/modify_files.md"},
			{Path: "builtin:fragments/act/run_tests.md"},
		},
	}

	result, err := AssemblePrompt(procedure, "", "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify enhanced separator format for all phases
	expectedSeparators := []string{
		"PHASE 1: OBSERVE",
		"PHASE 2: ORIENT",
		"PHASE 3: DECIDE",
		"PHASE 4: ACT",
	}

	for _, sep := range expectedSeparators {
		if !strings.Contains(result, sep) {
			t.Errorf("expected separator %q in assembled prompt", sep)
		}
	}

	// Verify phase order
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")
	orientIdx := strings.Index(result, "PHASE 2: ORIENT")
	decideIdx := strings.Index(result, "PHASE 3: DECIDE")
	actIdx := strings.Index(result, "PHASE 4: ACT")

	if observeIdx >= orientIdx || orientIdx >= decideIdx || decideIdx >= actIdx {
		t.Errorf("phases not in correct order: OBSERVE(%d) ORIENT(%d) DECIDE(%d) ACT(%d)",
			observeIdx, orientIdx, decideIdx, actIdx)
	}

	// Verify fragments within phases are separated by double newlines
	if !strings.Contains(result, "\n\n") {
		t.Errorf("expected double newlines between fragments")
	}
}

func TestAssemblePrompt_Integration_ContextFileWithSource(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "context-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	contextContent := "Focus on authentication and authorization modules.\nEnsure backward compatibility."
	if _, err := tmpFile.WriteString(contextContent); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	procedure := config.Procedure{
		Observe: []config.FragmentAction{{Content: "Observe phase"}},
		Orient:  []config.FragmentAction{{Content: "Orient phase"}},
		Decide:  []config.FragmentAction{{Content: "Decide phase"}},
		Act:     []config.FragmentAction{{Content: "Act phase"}},
	}

	result, err := AssemblePrompt(procedure, tmpFile.Name(), "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify CONTEXT section exists
	if !strings.Contains(result, "=== CONTEXT ===") {
		t.Errorf("expected CONTEXT section marker")
	}

	// Verify Source: line with file path
	expectedSource := "Source: " + tmpFile.Name()
	if !strings.Contains(result, expectedSource) {
		t.Errorf("expected Source line %q", expectedSource)
	}

	// Verify file content is included
	if !strings.Contains(result, contextContent) {
		t.Errorf("expected file content in prompt")
	}

	// Verify Source line comes before content
	sourceIdx := strings.Index(result, "Source:")
	contentIdx := strings.Index(result, contextContent)
	if sourceIdx >= contentIdx {
		t.Errorf("expected Source line before content")
	}

	// Verify CONTEXT comes before OBSERVE
	contextIdx := strings.Index(result, "=== CONTEXT ===")
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")
	if contextIdx >= observeIdx {
		t.Errorf("expected CONTEXT before OBSERVE")
	}
}

func TestAssemblePrompt_Integration_InlineContextNoSource(t *testing.T) {
	procedure := config.Procedure{
		Observe: []config.FragmentAction{{Content: "Observe phase"}},
		Orient:  []config.FragmentAction{{Content: "Orient phase"}},
		Decide:  []config.FragmentAction{{Content: "Decide phase"}},
		Act:     []config.FragmentAction{{Content: "Act phase"}},
	}

	inlineContext := "Focus on performance optimization and reduce memory usage"
	result, err := AssemblePrompt(procedure, inlineContext, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify CONTEXT section exists
	if !strings.Contains(result, "=== CONTEXT ===") {
		t.Errorf("expected CONTEXT section marker")
	}

	// Verify NO Source: line for inline content
	if strings.Contains(result, "Source:") {
		t.Errorf("expected no Source line for inline context")
	}

	// Verify inline content is included directly
	if !strings.Contains(result, inlineContext) {
		t.Errorf("expected inline context in prompt")
	}

	// Verify CONTEXT comes before OBSERVE
	contextIdx := strings.Index(result, "=== CONTEXT ===")
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")
	if contextIdx >= observeIdx {
		t.Errorf("expected CONTEXT before OBSERVE")
	}
}

func TestAssemblePrompt_Integration_MixedFileAndInlineContexts(t *testing.T) {
	tmpFile1, err := os.CreateTemp("", "context1-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file 1: %v", err)
	}
	defer os.Remove(tmpFile1.Name())

	content1 := "File context 1: API requirements"
	if _, err := tmpFile1.WriteString(content1); err != nil {
		t.Fatalf("failed to write temp file 1: %v", err)
	}
	tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "context2-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file 2: %v", err)
	}
	defer os.Remove(tmpFile2.Name())

	content2 := "File context 2: Database schema"
	if _, err := tmpFile2.WriteString(content2); err != nil {
		t.Fatalf("failed to write temp file 2: %v", err)
	}
	tmpFile2.Close()

	procedure := config.Procedure{
		Observe: []config.FragmentAction{{Content: "Observe phase"}},
		Orient:  []config.FragmentAction{{Content: "Orient phase"}},
		Decide:  []config.FragmentAction{{Content: "Decide phase"}},
		Act:     []config.FragmentAction{{Content: "Act phase"}},
	}

	// Simulate multiple --context flags: file, inline, file
	mixedContext := tmpFile1.Name() + "\n\n" + "Inline: ensure backward compatibility" + "\n\n" + tmpFile2.Name()
	result, err := AssemblePrompt(procedure, mixedContext, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify CONTEXT section exists
	if !strings.Contains(result, "=== CONTEXT ===") {
		t.Errorf("expected CONTEXT section marker")
	}

	// Verify first file has Source line
	expectedSource1 := "Source: " + tmpFile1.Name()
	if !strings.Contains(result, expectedSource1) {
		t.Errorf("expected Source line for first file")
	}

	// Verify first file content
	if !strings.Contains(result, content1) {
		t.Errorf("expected first file content")
	}

	// Verify inline content (no Source line for this part)
	if !strings.Contains(result, "Inline: ensure backward compatibility") {
		t.Errorf("expected inline context")
	}

	// Verify second file has Source line
	expectedSource2 := "Source: " + tmpFile2.Name()
	if !strings.Contains(result, expectedSource2) {
		t.Errorf("expected Source line for second file")
	}

	// Verify second file content
	if !strings.Contains(result, content2) {
		t.Errorf("expected second file content")
	}

	// Verify order: file1, inline, file2
	source1Idx := strings.Index(result, expectedSource1)
	inlineIdx := strings.Index(result, "Inline: ensure backward compatibility")
	source2Idx := strings.Index(result, expectedSource2)

	if source1Idx >= inlineIdx || inlineIdx >= source2Idx {
		t.Errorf("contexts not in correct order")
	}
}

func TestAssemblePrompt_PreambleStructure(t *testing.T) {
	procedure := config.Procedure{
		Display: "Build Procedure",
		Observe: []config.FragmentAction{{Content: "Observe"}},
		Orient:  []config.FragmentAction{{Content: "Orient"}},
		Decide:  []config.FragmentAction{{Content: "Decide"}},
		Act:     []config.FragmentAction{{Content: "Act"}},
	}

	result, err := AssemblePrompt(procedure, "", "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check preamble header
	if !strings.Contains(result, "ROODA PROCEDURE EXECUTION") {
		t.Errorf("expected preamble header")
	}

	// Check procedure name
	if !strings.Contains(result, "Procedure: Build Procedure") {
		t.Errorf("expected procedure name")
	}

	// Check role section
	if !strings.Contains(result, "Your Role") {
		t.Errorf("expected role section")
	}

	// Check success signaling
	if !strings.Contains(result, "SUCCESS") {
		t.Errorf("expected success signal instruction")
	}
	if !strings.Contains(result, "FAILURE") {
		t.Errorf("expected failure signal instruction")
	}

	// Verify preamble comes before OODA phases
	preambleIdx := strings.Index(result, "ROODA PROCEDURE EXECUTION")
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")
	if preambleIdx == -1 || observeIdx == -1 || preambleIdx >= observeIdx {
		t.Errorf("expected preamble before OBSERVE phase")
	}
}

func TestAssemblePrompt_PreambleOrder(t *testing.T) {
	procedure := config.Procedure{
		Display: "Test",
		Observe: []config.FragmentAction{{Content: "Observe"}},
		Orient:  []config.FragmentAction{{Content: "Orient"}},
		Decide:  []config.FragmentAction{{Content: "Decide"}},
		Act:     []config.FragmentAction{{Content: "Act"}},
	}

	userContext := "Focus on testing"
	result, err := AssemblePrompt(procedure, userContext, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify order: preamble -> context -> OODA phases
	preambleIdx := strings.Index(result, "ROODA PROCEDURE EXECUTION")
	contextIdx := strings.Index(result, "=== CONTEXT ===")
	observeIdx := strings.Index(result, "PHASE 1: OBSERVE")

	if preambleIdx >= contextIdx {
		t.Errorf("expected preamble before context")
	}
	if contextIdx >= observeIdx {
		t.Errorf("expected context before OBSERVE")
	}
}
