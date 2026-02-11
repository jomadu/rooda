package procedures

import (
	"testing"

	"github.com/jomadu/rooda/internal/config"
)

func TestBuiltInProcedures(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that all 16 built-in procedures are defined
	expectedProcedures := []string{
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

	if len(procedures) != len(expectedProcedures) {
		t.Errorf("expected %d procedures, got %d", len(expectedProcedures), len(procedures))
	}

	for _, name := range expectedProcedures {
		proc, exists := procedures[name]
		if !exists {
			t.Errorf("missing built-in procedure: %s", name)
			continue
		}

		// Verify required fields
		if proc.Display == "" {
			t.Errorf("procedure %s: missing Display field", name)
		}
		if proc.Summary == "" {
			t.Errorf("procedure %s: missing Summary field", name)
		}
		if proc.Description == "" {
			t.Errorf("procedure %s: missing Description field", name)
		}

		// Verify all OODA phases have at least one fragment
		if len(proc.Observe) == 0 {
			t.Errorf("procedure %s: Observe phase is empty", name)
		}
		if len(proc.Orient) == 0 {
			t.Errorf("procedure %s: Orient phase is empty", name)
		}
		if len(proc.Decide) == 0 {
			t.Errorf("procedure %s: Decide phase is empty", name)
		}
		if len(proc.Act) == 0 {
			t.Errorf("procedure %s: Act phase is empty", name)
		}

		// Verify all fragments use builtin: prefix
		allFragments := append(append(append(proc.Observe, proc.Orient...), proc.Decide...), proc.Act...)
		for i, fragment := range allFragments {
			if fragment.Path == "" && fragment.Content == "" {
				t.Errorf("procedure %s: fragment %d has neither path nor content", name, i)
			}
			if fragment.Path != "" && fragment.Content != "" {
				t.Errorf("procedure %s: fragment %d has both path and content", name, i)
			}
			if fragment.Path != "" && fragment.Path[:8] != "builtin:" {
				t.Errorf("procedure %s: fragment %d path does not start with 'builtin:': %s", name, i, fragment.Path)
			}
		}
	}
}

func TestBuiltInProcedureFragmentPaths(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test specific procedures have expected fragment structure
	testCases := []struct {
		name          string
		observeCount  int
		orientCount   int
		decideCount   int
		actCount      int
	}{
		{"agents-sync", 4, 2, 2, 3},
		{"build", 5, 3, 3, 5},
		{"publish-plan", 3, 2, 2, 3},
		{"audit-spec", 2, 1, 2, 2},
		{"audit-impl", 4, 1, 2, 2},
		{"audit-agents", 4, 2, 1, 2},
		{"audit-spec-to-impl", 3, 1, 1, 2},
		{"audit-impl-to-spec", 3, 1, 1, 2},
	}

	for _, tc := range testCases {
		proc, exists := procedures[tc.name]
		if !exists {
			t.Errorf("procedure %s not found", tc.name)
			continue
		}

		if len(proc.Observe) != tc.observeCount {
			t.Errorf("procedure %s: expected %d observe fragments, got %d", tc.name, tc.observeCount, len(proc.Observe))
		}
		if len(proc.Orient) != tc.orientCount {
			t.Errorf("procedure %s: expected %d orient fragments, got %d", tc.name, tc.orientCount, len(proc.Orient))
		}
		if len(proc.Decide) != tc.decideCount {
			t.Errorf("procedure %s: expected %d decide fragments, got %d", tc.name, tc.decideCount, len(proc.Decide))
		}
		if len(proc.Act) != tc.actCount {
			t.Errorf("procedure %s: expected %d act fragments, got %d", tc.name, tc.actCount, len(proc.Act))
		}
	}
}

func TestBuiltInProcedureMetadata(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that agents-sync has correct metadata
	agentsSync := procedures["agents-sync"]
	if agentsSync.Display != "Agents Sync" {
		t.Errorf("agents-sync: expected Display 'Agents Sync', got '%s'", agentsSync.Display)
	}
	if agentsSync.Summary != "Synchronize AGENTS.md with actual repository state" {
		t.Errorf("agents-sync: unexpected Summary: %s", agentsSync.Summary)
	}

	// Test that build has correct metadata
	build := procedures["build"]
	if build.Display != "Build" {
		t.Errorf("build: expected Display 'Build', got '%s'", build.Display)
	}
	if build.Summary != "Implement a task from work tracking" {
		t.Errorf("build: unexpected Summary: %s", build.Summary)
	}
}

func TestBuiltInProcedureFragmentContent(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that agents-sync uses expected fragments
	agentsSync := procedures["agents-sync"]
	
	expectedObserveFragments := []string{
		"builtin:fragments/observe/read_agents_md.md",
		"builtin:fragments/observe/scan_repo_structure.md",
		"builtin:fragments/observe/detect_build_system.md",
		"builtin:fragments/observe/detect_work_tracking.md",
	}
	
	for i, expected := range expectedObserveFragments {
		if i >= len(agentsSync.Observe) {
			t.Errorf("agents-sync: missing observe fragment %d", i)
			continue
		}
		if agentsSync.Observe[i].Path != expected {
			t.Errorf("agents-sync: observe fragment %d: expected %s, got %s", i, expected, agentsSync.Observe[i].Path)
		}
	}

	// Test that build uses expected fragments
	build := procedures["build"]
	
	expectedBuildObserve := []string{
		"builtin:fragments/observe/read_agents_md.md",
		"builtin:fragments/observe/query_work_tracking.md",
		"builtin:fragments/observe/read_specs.md",
		"builtin:fragments/observe/read_impl.md",
		"builtin:fragments/observe/read_task_details.md",
	}
	
	for i, expected := range expectedBuildObserve {
		if i >= len(build.Observe) {
			t.Errorf("build: missing observe fragment %d", i)
			continue
		}
		if build.Observe[i].Path != expected {
			t.Errorf("build: observe fragment %d: expected %s, got %s", i, expected, build.Observe[i].Path)
		}
	}
}

func TestBuiltInProceduresNoParameters(t *testing.T) {
	procedures := BuiltInProcedures()

	// Built-in procedures should not have template parameters
	for name, proc := range procedures {
		allFragments := append(append(append(proc.Observe, proc.Orient...), proc.Decide...), proc.Act...)
		for i, fragment := range allFragments {
			if len(fragment.Parameters) > 0 {
				t.Errorf("procedure %s: fragment %d has parameters (built-ins should not use parameters)", name, i)
			}
		}
	}
}

func TestBuiltInProceduresNoInlineContent(t *testing.T) {
	procedures := BuiltInProcedures()

	// Built-in procedures should use fragment files, not inline content
	for name, proc := range procedures {
		allFragments := append(append(append(proc.Observe, proc.Orient...), proc.Decide...), proc.Act...)
		for i, fragment := range allFragments {
			if fragment.Content != "" {
				t.Errorf("procedure %s: fragment %d uses inline content (built-ins should use fragment files)", name, i)
			}
		}
	}
}

func TestBuiltInProceduresIterationSettings(t *testing.T) {
	procedures := BuiltInProcedures()

	// Built-in procedures should not override iteration settings
	// (they inherit from loop config)
	for name, proc := range procedures {
		if proc.IterationMode != "" {
			t.Errorf("procedure %s: has iteration_mode set (should inherit from loop)", name)
		}
		if proc.DefaultMaxIterations != nil {
			t.Errorf("procedure %s: has default_max_iterations set (should inherit from loop)", name)
		}
		if proc.IterationTimeout != nil {
			t.Errorf("procedure %s: has iteration_timeout set (should inherit from loop)", name)
		}
		if proc.MaxOutputBuffer != nil {
			t.Errorf("procedure %s: has max_output_buffer set (should inherit from loop)", name)
		}
	}
}

func TestBuiltInProceduresAICommandSettings(t *testing.T) {
	procedures := BuiltInProcedures()

	// Built-in procedures should not override AI command settings
	// (they inherit from loop config)
	for name, proc := range procedures {
		if proc.AICmd != "" {
			t.Errorf("procedure %s: has ai_cmd set (should inherit from loop)", name)
		}
		if proc.AICmdAlias != "" {
			t.Errorf("procedure %s: has ai_cmd_alias set (should inherit from loop)", name)
		}
	}
}

func TestBuiltInProceduresFragmentOrdering(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that fragments are in logical order
	// For agents-sync, observe should read AGENTS.md first
	agentsSync := procedures["agents-sync"]
	if len(agentsSync.Observe) > 0 && agentsSync.Observe[0].Path != "builtin:fragments/observe/read_agents_md.md" {
		t.Errorf("agents-sync: first observe fragment should be read_agents_md.md, got %s", agentsSync.Observe[0].Path)
	}

	// For build, observe should read AGENTS.md first
	build := procedures["build"]
	if len(build.Observe) > 0 && build.Observe[0].Path != "builtin:fragments/observe/read_agents_md.md" {
		t.Errorf("build: first observe fragment should be read_agents_md.md, got %s", build.Observe[0].Path)
	}

	// For build, act should emit success last
	if len(build.Act) > 0 && build.Act[len(build.Act)-1].Path != "builtin:fragments/act/emit_success.md" {
		t.Errorf("build: last act fragment should be emit_success.md, got %s", build.Act[len(build.Act)-1].Path)
	}
}

func TestBuiltInProceduresCommonFragments(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that common fragments are reused across procedures
	readAgentsMdCount := 0
	emitSuccessCount := 0

	for _, proc := range procedures {
		// Count read_agents_md usage
		for _, fragment := range proc.Observe {
			if fragment.Path == "builtin:fragments/observe/read_agents_md.md" {
				readAgentsMdCount++
				break
			}
		}

		// Count emit_success usage
		for _, fragment := range proc.Act {
			if fragment.Path == "builtin:fragments/act/emit_success.md" {
				emitSuccessCount++
				break
			}
		}
	}

	// Most procedures should read AGENTS.md first
	if readAgentsMdCount < 10 {
		t.Errorf("expected at least 10 procedures to read AGENTS.md, got %d", readAgentsMdCount)
	}

	// All procedures should emit success
	if emitSuccessCount != 16 {
		t.Errorf("expected all 16 procedures to emit success, got %d", emitSuccessCount)
	}
}

func TestBuiltInProceduresReturnsCopy(t *testing.T) {
	// Test that BuiltInProcedures returns a new map each time
	// (not a shared reference that could be mutated)
	procedures1 := BuiltInProcedures()
	procedures2 := BuiltInProcedures()

	// Modify procedures1
	delete(procedures1, "build")

	// Verify procedures2 is unaffected
	if _, exists := procedures2["build"]; !exists {
		t.Error("modifying returned map affected subsequent calls (not a copy)")
	}

	// Verify we can get a fresh copy
	procedures3 := BuiltInProcedures()
	if len(procedures3) != 16 {
		t.Errorf("expected 16 procedures in fresh copy, got %d", len(procedures3))
	}
}

func TestBuiltInProceduresFragmentActionStructure(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that all FragmentActions have valid structure
	for name, proc := range procedures {
		allFragments := append(append(append(proc.Observe, proc.Orient...), proc.Decide...), proc.Act...)
		for i, fragment := range allFragments {
			// Must have exactly one of path or content
			hasPath := fragment.Path != ""
			hasContent := fragment.Content != ""

			if !hasPath && !hasContent {
				t.Errorf("procedure %s: fragment %d has neither path nor content", name, i)
			}
			if hasPath && hasContent {
				t.Errorf("procedure %s: fragment %d has both path and content", name, i)
			}

			// If path is set, must start with builtin:
			if hasPath && fragment.Path[:8] != "builtin:" {
				t.Errorf("procedure %s: fragment %d path must start with 'builtin:', got %s", name, i, fragment.Path)
			}

			// Parameters should be nil or empty for built-ins
			if fragment.Parameters != nil && len(fragment.Parameters) > 0 {
				t.Errorf("procedure %s: fragment %d has parameters (built-ins should not use parameters)", name, i)
			}
		}
	}
}

func TestBuiltInProceduresPhaseBalance(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that procedures have reasonable phase balance
	// (no phase should be dramatically larger than others)
	for name, proc := range procedures {
		observeCount := len(proc.Observe)
		orientCount := len(proc.Orient)
		decideCount := len(proc.Decide)
		actCount := len(proc.Act)

		total := observeCount + orientCount + decideCount + actCount

		// Each phase should have at least 1 fragment
		if observeCount == 0 || orientCount == 0 || decideCount == 0 || actCount == 0 {
			t.Errorf("procedure %s: has empty phase (observe=%d, orient=%d, decide=%d, act=%d)",
				name, observeCount, orientCount, decideCount, actCount)
		}

		// No phase should be more than 70% of total
		maxPhase := max(observeCount, orientCount, decideCount, actCount)
		if float64(maxPhase)/float64(total) > 0.7 {
			t.Errorf("procedure %s: phase imbalance detected (largest phase is %d of %d total)",
				name, maxPhase, total)
		}
	}
}

func max(a, b, c, d int) int {
	result := a
	if b > result {
		result = b
	}
	if c > result {
		result = c
	}
	if d > result {
		result = d
	}
	return result
}

func TestBuiltInProceduresFragmentPathFormat(t *testing.T) {
	procedures := BuiltInProcedures()

	// Test that all fragment paths follow the expected format:
	// builtin:fragments/<phase>/<name>.md

	for procName, proc := range procedures {
		phases := []struct {
			name      string
			fragments []config.FragmentAction
		}{
			{"observe", proc.Observe},
			{"orient", proc.Orient},
			{"decide", proc.Decide},
			{"act", proc.Act},
		}

		for _, phase := range phases {
			for i, fragment := range phase.fragments {
				if fragment.Path == "" {
					continue // inline content
				}

				// Check format: builtin:fragments/<phase>/<name>.md
				expectedPrefix := "builtin:fragments/" + phase.name + "/"
				if len(fragment.Path) < len(expectedPrefix) || fragment.Path[:len(expectedPrefix)] != expectedPrefix {
					t.Errorf("procedure %s: %s phase fragment %d: path should start with '%s', got %s",
						procName, phase.name, i, expectedPrefix, fragment.Path)
				}

				// Check .md extension
				if len(fragment.Path) < 3 || fragment.Path[len(fragment.Path)-3:] != ".md" {
					t.Errorf("procedure %s: %s phase fragment %d: path should end with '.md', got %s",
						procName, phase.name, i, fragment.Path)
				}
			}
		}
	}
}
