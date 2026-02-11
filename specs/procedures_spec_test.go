package specs

import (
	"os"
	"strings"
	"testing"
)

// TestProceduresSpecSignalFormat verifies all signal references use exact format
func TestProceduresSpecSignalFormat(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// Verify SUCCESS signal format
	if !strings.Contains(spec, "<promise>SUCCESS</promise>") {
		t.Error("Spec must contain exact SUCCESS signal format: <promise>SUCCESS</promise>")
	}

	// Verify FAILURE signal format
	if !strings.Contains(spec, "<promise>FAILURE</promise>") {
		t.Error("Spec must contain exact FAILURE signal format: <promise>FAILURE</promise>")
	}

	// Check for incorrect variations
	incorrectVariations := []string{
		"SUCCESS promise",
		"promise SUCCESS",
		"<SUCCESS>",
		"<FAILURE>",
		"promise: SUCCESS",
		"promise: FAILURE",
	}

	for _, variation := range incorrectVariations {
		if strings.Contains(spec, variation) {
			t.Errorf("Spec contains incorrect signal variation: %s", variation)
		}
	}
}

// TestProceduresSpecHasSUCCESSCriteria verifies SUCCESS criteria section exists
func TestProceduresSpecHasSUCCESSCriteria(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// Verify SUCCESS Criteria section exists
	if !strings.Contains(spec, "## SUCCESS Criteria by Procedure Type") {
		t.Error("Spec must contain '## SUCCESS Criteria by Procedure Type' section")
	}

	// Verify it documents all three procedure categories
	requiredSections := []string{
		"### Direct Action Procedures",
		"### Audit Procedures",
		"### Planning Procedures",
	}

	for _, section := range requiredSections {
		if !strings.Contains(spec, section) {
			t.Errorf("SUCCESS Criteria section must contain: %s", section)
		}
	}

	// Verify it documents when to emit SUCCESS for each type
	requiredProcedures := []string{
		"**agents-sync:**",
		"**build:**",
		"**publish-plan:**",
		"**audit-spec, audit-impl, audit-agents:**",
		"**audit-spec-to-impl, audit-impl-to-spec:**",
		"**draft-plan-* (all 8 variants):**",
	}

	for _, proc := range requiredProcedures {
		if !strings.Contains(spec, proc) {
			t.Errorf("SUCCESS Criteria section must document: %s", proc)
		}
	}
}

// TestProceduresSpecHasExampleOutputs verifies Example Outputs section exists
func TestProceduresSpecHasExampleOutputs(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// Verify Example Outputs section exists
	if !strings.Contains(spec, "## Example Outputs") {
		t.Error("Spec must contain '## Example Outputs' section")
	}

	// Verify it has examples for key procedures
	requiredExamples := []string{
		"### agents-sync Procedure",
		"### build Procedure",
		"### audit-spec Procedure",
		"### draft-plan-impl-feat Procedure",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(spec, example) {
			t.Errorf("Example Outputs section must contain: %s", example)
		}
	}

	// Verify examples show signals on their own line
	// Count occurrences of signal in code blocks
	successCount := strings.Count(spec, "<promise>SUCCESS</promise>")
	failureCount := strings.Count(spec, "<promise>FAILURE</promise>")

	if successCount < 5 {
		t.Errorf("Expected at least 5 SUCCESS signal examples, found %d", successCount)
	}

	if failureCount < 2 {
		t.Errorf("Expected at least 2 FAILURE signal examples, found %d", failureCount)
	}
}

// TestProceduresSpecAll16ProceduresDefined verifies all 16 procedures are defined
func TestProceduresSpecAll16ProceduresDefined(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// All 16 built-in procedures
	procedures := []string{
		"agents-sync",
		"build",
		"publish-plan",
		"audit-spec",
		"audit-impl",
		"audit-agents",
		"audit-spec-to-impl",
		"audit-impl-to-spec",
		"draft-plan-spec-feat",
		"draft-plan-spec-fix",
		"draft-plan-spec-refactor",
		"draft-plan-spec-chore",
		"draft-plan-impl-feat",
		"draft-plan-impl-fix",
		"draft-plan-impl-refactor",
		"draft-plan-impl-chore",
	}

	for _, proc := range procedures {
		// Check for procedure definition in YAML config
		if !strings.Contains(spec, proc+":") {
			t.Errorf("Procedure %s must be defined in Built-in Procedures Configuration", proc)
		}
	}
}

// TestProceduresSpecFragmentReferences verifies all procedures reference emit_success
func TestProceduresSpecFragmentReferences(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// All procedures should reference emit_success.md in their act phase
	emitSuccessCount := strings.Count(spec, "builtin:fragments/act/emit_success.md")

	// Should be 16 (one per procedure)
	if emitSuccessCount != 16 {
		t.Errorf("Expected 16 references to emit_success.md (one per procedure), found %d", emitSuccessCount)
	}

	// Some procedures should reference emit_failure.md in decide phase (check_if_blocked)
	// Not all procedures have this, but at least build should
	if !strings.Contains(spec, "builtin:fragments/decide/check_if_blocked.md") {
		t.Error("Expected at least one procedure to reference check_if_blocked.md")
	}
}

// TestProceduresSpecConsistentTerminology verifies consistent terminology
func TestProceduresSpecConsistentTerminology(t *testing.T) {
	content, err := os.ReadFile("procedures.md")
	if err != nil {
		t.Fatalf("Failed to read procedures.md: %v", err)
	}

	spec := string(content)

	// Verify consistent use of "emit" terminology
	if !strings.Contains(spec, "Emit `<promise>SUCCESS</promise>` when:") {
		t.Error("SUCCESS criteria should use 'Emit' terminology")
	}

	if !strings.Contains(spec, "Emit `<promise>FAILURE</promise>` when:") {
		t.Error("FAILURE criteria should use 'Emit' terminology")
	}

	// Verify consistent use of "Continue iterating when:"
	continueCount := strings.Count(spec, "Continue iterating when:")
	if continueCount < 6 {
		t.Errorf("Expected at least 6 'Continue iterating when:' statements, found %d", continueCount)
	}
}
