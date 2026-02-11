package agents

import (
	"os"
	"testing"
)

func TestParseAgentsMD_RealAGENTSMD(t *testing.T) {
	content, err := os.ReadFile("../../AGENTS.md")
	if err != nil {
		t.Fatalf("Failed to read AGENTS.md: %v", err)
	}

	agentsMD, err := ParseAgentsMD(string(content))
	if err != nil {
		t.Fatalf("ParseAgentsMD failed: %v", err)
	}

	// Verify work tracking
	if agentsMD.WorkTracking.System != "beads (bd CLI)" {
		t.Errorf("Expected system 'beads (bd CLI)', got '%s'", agentsMD.WorkTracking.System)
	}
	if agentsMD.WorkTracking.QueryCommand != "bd ready --json" {
		t.Errorf("Expected query command 'bd ready --json', got '%s'", agentsMD.WorkTracking.QueryCommand)
	}

	// Verify test command (should be "make test" in new format)
	if agentsMD.TestCommand == "" {
		t.Errorf("Expected test command to be parsed, got empty string")
		t.Logf("Debug: TestCommand='%s', BuildCommand='%s', LintCommands=%v", 
			agentsMD.TestCommand, agentsMD.BuildCommand, agentsMD.LintCommands)
	}

	// Verify spec paths
	if len(agentsMD.SpecPaths) == 0 {
		t.Error("Expected spec paths to be parsed")
	}

	// Verify impl paths
	if len(agentsMD.ImplPaths) == 0 {
		t.Error("Expected impl paths to be parsed")
	}

	// Verify quality criteria
	if len(agentsMD.QualityCriteria) == 0 {
		t.Error("Expected quality criteria to be parsed")
	}
}

func TestParseAgentsMD_MinimalValid(t *testing.T) {
	content := `# Agent Instructions

## Work Tracking System

**System:** beads (bd CLI)

**Query ready work:**
` + "```bash\nbd ready --json\n```" + `

## Build/Test/Lint Commands

**Test:**
` + "```bash\ngo test ./...\n```" + `

## Specification Definition

**Location:** ` + "`specs/*.md`" + `

## Implementation Definition

**Location:** ` + "`internal/`" + `

## Quality Criteria

**For specifications:**
- All specs have JTBD section (PASS/FAIL)
`

	agentsMD, err := ParseAgentsMD(content)
	if err != nil {
		t.Fatalf("ParseAgentsMD failed: %v", err)
	}

	if agentsMD.WorkTracking.System != "beads (bd CLI)" {
		t.Errorf("Expected system 'beads (bd CLI)', got '%s'", agentsMD.WorkTracking.System)
	}
	if agentsMD.TestCommand != "go test ./..." {
		t.Errorf("Expected test command 'go test ./...', got '%s'", agentsMD.TestCommand)
	}
	if len(agentsMD.SpecPaths) != 1 || agentsMD.SpecPaths[0] != "specs/*.md" {
		t.Errorf("Expected spec paths ['specs/*.md'], got %v", agentsMD.SpecPaths)
	}
	if len(agentsMD.ImplPaths) != 1 || agentsMD.ImplPaths[0] != "internal/" {
		t.Errorf("Expected impl paths ['internal/'], got %v", agentsMD.ImplPaths)
	}
	if len(agentsMD.QualityCriteria) != 1 {
		t.Errorf("Expected 1 quality criterion, got %d", len(agentsMD.QualityCriteria))
	}
}

func TestParseAgentsMD_MissingOptionalSections(t *testing.T) {
	content := `# Agent Instructions

## Work Tracking System

**System:** beads

**Query ready work:**
` + "```bash\nbd ready\n```" + `

## Build/Test/Lint Commands

**Test:**
` + "```bash\ngo test\n```" + `

## Specification Definition

**Location:** specs/

## Implementation Definition

**Location:** src/
`

	agentsMD, err := ParseAgentsMD(content)
	if err != nil {
		t.Fatalf("ParseAgentsMD should succeed with missing optional sections: %v", err)
	}

	if agentsMD.WorkTracking.System != "beads" {
		t.Errorf("Expected system 'beads', got '%s'", agentsMD.WorkTracking.System)
	}
	if len(agentsMD.QualityCriteria) != 0 {
		t.Errorf("Expected 0 quality criteria, got %d", len(agentsMD.QualityCriteria))
	}
}

func TestParseAgentsMD_EmptyContent(t *testing.T) {
	_, err := ParseAgentsMD("")
	if err == nil {
		t.Error("Expected error for empty content")
	}
}

func TestParseAgentsMD_MissingRequiredSection(t *testing.T) {
	content := `# Agent Instructions

## Build/Test/Lint Commands

**Test:**
` + "```bash\ngo test\n```" + `
`

	_, err := ParseAgentsMD(content)
	if err == nil {
		t.Error("Expected error for missing Work Tracking System section")
	}
}
