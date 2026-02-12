package prompt

import (
	"fmt"
	"strings"

	"github.com/jomadu/rooda/internal/config"
)

// IterationContext contains iteration state information for prompt assembly.
// This is a minimal struct to avoid circular dependencies with the loop package.
type IterationContext struct {
	CurrentIteration int  // 0-indexed current iteration number
	MaxIterations    *int // nil for unlimited mode
}

// AssemblePrompt assembles a complete prompt from a procedure definition.
// It concatenates fragments from each OODA phase with section markers and
// optionally injects user context at the top. Context string may contain
// multiple values separated by \n\n. Each value is checked for file existence -
// if a file exists, its content is read and prefixed with "Source: <path>".
// Otherwise, the value is treated as inline content.
// If iterCtx is provided, iteration context is included in the preamble.
func AssemblePrompt(procedure config.Procedure, userContext string, configDir string, iterCtx *IterationContext) (string, error) {
	var prompt strings.Builder

	// Inject preamble first
	preamble := generatePreamble(procedure, iterCtx)
	prompt.WriteString(preamble)
	prompt.WriteString("\n\n")

	// Inject user context first if provided
	if userContext != "" {
		prompt.WriteString("=== CONTEXT ===\n")
		
		// Split on double newlines (multiple --context flags)
		contextValues := strings.Split(userContext, "\n\n")
		
		for i, contextValue := range contextValues {
			if strings.TrimSpace(contextValue) == "" {
				continue
			}
			
			// Check if context is a file path (file existence heuristic)
			contextContent, isFile, err := LoadContextContent(contextValue)
			if err != nil {
				return "", fmt.Errorf("failed to load context: %v", err)
			}
			
			if isFile {
				// Add source path for file-based context
				prompt.WriteString("Source: ")
				prompt.WriteString(contextValue)
				prompt.WriteString("\n\n")
			}
			
			prompt.WriteString(contextContent)
			
			// Add separator between multiple contexts
			if i < len(contextValues)-1 {
				prompt.WriteString("\n\n")
			}
		}
		
		prompt.WriteString("\n\n")
	}

	// Process each OODA phase in order
	phases := []struct {
		name        string
		number      int
		description string
		fragments   []config.FragmentAction
	}{
		{"OBSERVE", 1, "Execute these observation tasks to gather information.", procedure.Observe},
		{"ORIENT", 2, "Analyze the information you gathered and form your understanding.", procedure.Orient},
		{"DECIDE", 3, "Make decisions about what actions to take.", procedure.Decide},
		{"ACT", 4, "Execute the actions you decided on. Modify files, run commands, commit changes.", procedure.Act},
	}

	for _, phase := range phases {
		phaseContent, err := ComposePhasePrompt(phase.fragments, configDir)
		if err != nil {
			return "", fmt.Errorf("failed to compose %s phase: %v", phase.name, err)
		}

		// Add section marker and content if phase has content
		trimmed := strings.TrimSpace(phaseContent)
		if trimmed != "" {
			// Enhanced section marker with double lines and phase description
			prompt.WriteString("═══════════════════════════════════════════════════════════════\n")
			prompt.WriteString(fmt.Sprintf("PHASE %d: %s\n", phase.number, phase.name))
			prompt.WriteString(phase.description)
			prompt.WriteString("\n═══════════════════════════════════════════════════════════════\n")
			prompt.WriteString(trimmed)
			prompt.WriteString("\n\n")
		}
	}

	return prompt.String(), nil
}

// generatePreamble creates the procedure execution preamble with agent role and success signaling instructions.
// If iterCtx is provided, includes iteration context (current iteration and max iterations or unlimited).
func generatePreamble(procedure config.Procedure, iterCtx *IterationContext) string {
	var preamble strings.Builder

	preamble.WriteString("═══════════════════════════════════════════════════════════════\n")
	preamble.WriteString("ROODA PROCEDURE EXECUTION\n")
	preamble.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	if procedure.Display != "" {
		preamble.WriteString("Procedure: ")
		preamble.WriteString(procedure.Display)
		preamble.WriteString("\n\n")
	}

	// Add iteration context if provided
	if iterCtx != nil {
		preamble.WriteString("Iteration: ")
		preamble.WriteString(fmt.Sprintf("%d", iterCtx.CurrentIteration+1)) // Convert 0-indexed to 1-indexed
		if iterCtx.MaxIterations != nil {
			preamble.WriteString(fmt.Sprintf(" of %d", *iterCtx.MaxIterations))
		} else {
			preamble.WriteString(" (unlimited)")
		}
		preamble.WriteString("\n\n")
	}

	preamble.WriteString("Your Role:\n")
	preamble.WriteString("You are an AI coding agent executing a structured OODA loop procedure.\n")
	preamble.WriteString("This is NOT a template or example - this is an EXECUTABLE PROCEDURE.\n")
	preamble.WriteString("You must complete all phases and produce concrete outputs.\n\n")

	preamble.WriteString("Success Signaling:\n")
	preamble.WriteString("- When you complete all tasks successfully, output: <promise>SUCCESS</promise>\n")
	preamble.WriteString("- If you cannot proceed due to blockers, output: <promise>FAILURE</promise>\n")
	preamble.WriteString("- Explanations should come AFTER the signal, not embedded in the tag\n")
	preamble.WriteString("- The loop orchestrator uses these signals to determine iteration outcome.\n")

	return preamble.String()
}

// ComposePhasePrompt composes a single phase prompt from an array of fragment actions.
// It loads fragments, processes templates if parameters are provided, and concatenates
// with double newlines.
func ComposePhasePrompt(fragments []config.FragmentAction, configDir string) (string, error) {
	if len(fragments) == 0 {
		return "", nil
	}

	var parts []string

	for _, fragment := range fragments {
		var content string

		// Determine content source (inline content or file path)
		if fragment.Content != "" {
			content = fragment.Content
		} else if fragment.Path != "" {
			var err error
			content, err = LoadFragment(fragment.Path, configDir)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("fragment must specify either content or path")
		}

		// Process template if parameters provided
		if len(fragment.Parameters) > 0 {
			var err error
			content, err = ProcessTemplate(content, fragment.Parameters)
			if err != nil {
				return "", err
			}
		}

		// Append to parts
		parts = append(parts, strings.TrimSpace(content))
	}

	// Concatenate with double newlines
	return strings.Join(parts, "\n\n"), nil
}
