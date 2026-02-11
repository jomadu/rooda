package agents

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseAgentsMD parses AGENTS.md content into structured data
func ParseAgentsMD(content string) (*AgentsMD, error) {
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("empty content")
	}

	agentsMD := &AgentsMD{
		RawContent: content,
	}

	// Parse Work Tracking System (required)
	if err := parseWorkTracking(content, agentsMD); err != nil {
		return nil, fmt.Errorf("failed to parse Work Tracking System: %w", err)
	}

	// Parse Build/Test/Lint Commands
	parseBuildTestLint(content, agentsMD)

	// Parse Specification Definition
	parseSpecDefinition(content, agentsMD)

	// Parse Implementation Definition
	parseImplDefinition(content, agentsMD)

	// Parse Task Input
	parseTaskInput(content, agentsMD)

	// Parse Planning System
	parsePlanningSystem(content, agentsMD)

	// Parse Audit Output
	parseAuditOutput(content, agentsMD)

	// Parse Quality Criteria
	parseQualityCriteria(content, agentsMD)

	return agentsMD, nil
}

func parseWorkTracking(content string, agentsMD *AgentsMD) error {
	section := extractSection(content, "Work Tracking System")
	if section == "" {
		return fmt.Errorf("Work Tracking System section not found")
	}

	// Parse system name
	systemRe := regexp.MustCompile(`\*\*System:\*\*\s+(.+)`)
	if matches := systemRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.WorkTracking.System = strings.TrimSpace(matches[1])
	}

	// Parse commands
	agentsMD.WorkTracking.QueryCommand = extractCommand(section, "Query ready work")
	agentsMD.WorkTracking.UpdateCommand = extractCommand(section, "Update status")
	agentsMD.WorkTracking.CloseCommand = extractCommand(section, "Close")
	agentsMD.WorkTracking.CreateCommand = extractCommand(section, "Create")

	return nil
}

func parseBuildTestLint(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Build/Test/Lint Commands")
	if section == "" {
		return
	}

	// Try new format first (Unified Interface via Makefile)
	// Look for the first code block after "Unified Interface"
	unifiedIdx := strings.Index(section, "**Unified Interface")
	if unifiedIdx >= 0 {
		afterUnified := section[unifiedIdx:]
		commands := extractAllCommands(afterUnified)
		if len(commands) > 0 {
			// Find "make test" command
			for _, cmd := range commands {
				if strings.Contains(cmd, "make test") {
					agentsMD.TestCommand = cmd
					break
				}
			}
			// Find "make build" command
			for _, cmd := range commands {
				if strings.Contains(cmd, "make build") {
					agentsMD.BuildCommand = cmd
					break
				}
			}
			// Find "make lint" command
			for _, cmd := range commands {
				if strings.Contains(cmd, "make lint") {
					agentsMD.LintCommands = append(agentsMD.LintCommands, cmd)
				}
			}
		}
	}

	// Try old format if new format didn't work
	if agentsMD.TestCommand == "" {
		agentsMD.TestCommand = extractCommand(section, "Test")
	}
	if agentsMD.BuildCommand == "" {
		agentsMD.BuildCommand = extractCommand(section, "Build")
	}
	if len(agentsMD.LintCommands) == 0 {
		lintSection := extractSubsection(section, "Lint")
		if lintSection != "" {
			commands := extractAllCommands(lintSection)
			agentsMD.LintCommands = commands
		}
	}
}

func parseSpecDefinition(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Specification Definition")
	if section == "" {
		return
	}

	// Parse location
	locationRe := regexp.MustCompile(`\*\*Location:\*\*\s+` + "`([^`]+)`")
	if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
		paths := strings.Split(matches[1], ",")
		for _, p := range paths {
			agentsMD.SpecPaths = append(agentsMD.SpecPaths, strings.TrimSpace(p))
		}
	} else {
		// Try without backticks
		locationRe = regexp.MustCompile(`\*\*Location:\*\*\s+(.+)`)
		if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
			path := strings.TrimSpace(matches[1])
			agentsMD.SpecPaths = append(agentsMD.SpecPaths, path)
		}
	}

	// Parse excludes
	excludeRe := regexp.MustCompile(`\*\*Exclude:\*\*\s+` + "`([^`]+)`")
	if matches := excludeRe.FindStringSubmatch(section); len(matches) > 1 {
		excludes := strings.Split(matches[1], ",")
		for _, e := range excludes {
			agentsMD.SpecExcludes = append(agentsMD.SpecExcludes, strings.TrimSpace(e))
		}
	} else {
		// Try without backticks
		excludeRe = regexp.MustCompile(`\*\*Exclude:\*\*\s+(.+)`)
		if matches := excludeRe.FindStringSubmatch(section); len(matches) > 1 {
			exclude := strings.TrimSpace(matches[1])
			agentsMD.SpecExcludes = append(agentsMD.SpecExcludes, exclude)
		}
	}
}

func parseImplDefinition(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Implementation Definition")
	if section == "" {
		return
	}

	// Parse location
	locationRe := regexp.MustCompile(`\*\*Location:\*\*\s+` + "`([^`]+)`")
	if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
		paths := strings.Split(matches[1], ",")
		for _, p := range paths {
			agentsMD.ImplPaths = append(agentsMD.ImplPaths, strings.TrimSpace(p))
		}
	} else {
		// Try without backticks
		locationRe = regexp.MustCompile(`\*\*Location:\*\*\s+(.+)`)
		if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
			path := strings.TrimSpace(matches[1])
			agentsMD.ImplPaths = append(agentsMD.ImplPaths, path)
		}
	}

	// Parse patterns (list items starting with -)
	patternRe := regexp.MustCompile(`-\s+` + "`([^`]+)`")
	for _, matches := range patternRe.FindAllStringSubmatch(section, -1) {
		if len(matches) > 1 {
			agentsMD.ImplPaths = append(agentsMD.ImplPaths, strings.TrimSpace(matches[1]))
		}
	}

	// Parse excludes
	excludeRe := regexp.MustCompile(`\*\*Exclude:\*\*`)
	if excludeRe.MatchString(section) {
		excludeSection := section[excludeRe.FindStringIndex(section)[0]:]
		excludeListRe := regexp.MustCompile(`-\s+` + "`([^`]+)`")
		for _, matches := range excludeListRe.FindAllStringSubmatch(excludeSection, -1) {
			if len(matches) > 1 {
				agentsMD.ImplExcludes = append(agentsMD.ImplExcludes, strings.TrimSpace(matches[1]))
			}
		}
	}
}

func parseTaskInput(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Task Input")
	if section == "" {
		section = extractSection(content, "Story/Bug Input")
	}
	if section == "" {
		return
	}

	// Parse location
	locationRe := regexp.MustCompile(`\*\*Location:\*\*\s+` + "`([^`]+)`")
	if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.TaskInput.Location = strings.TrimSpace(matches[1])
	} else {
		locationRe = regexp.MustCompile(`\*\*Location:\*\*\s+(.+)`)
		if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
			agentsMD.TaskInput.Location = strings.TrimSpace(matches[1])
		}
	}
}

func parsePlanningSystem(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Planning System")
	if section == "" {
		return
	}

	// Parse draft plan location
	locationRe := regexp.MustCompile(`\*\*Draft plan location:\*\*\s+` + "`([^`]+)`")
	if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.PlanningSystem.DraftPlanLocation = strings.TrimSpace(matches[1])
	} else {
		locationRe = regexp.MustCompile(`\*\*Draft plan location:\*\*\s+(.+)`)
		if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
			agentsMD.PlanningSystem.DraftPlanLocation = strings.TrimSpace(matches[1])
		}
	}

	// Parse publishing method
	methodRe := regexp.MustCompile(`\*\*Publishing mechanism:\*\*\s+(.+)`)
	if matches := methodRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.PlanningSystem.PublishingMethod = strings.TrimSpace(matches[1])
	}
}

func parseAuditOutput(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Audit Output")
	if section == "" {
		return
	}

	// Parse location pattern
	locationRe := regexp.MustCompile(`\*\*Location:\*\*\s+` + "`([^`]+)`")
	if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.AuditOutput.LocationPattern = strings.TrimSpace(matches[1])
	} else {
		locationRe = regexp.MustCompile(`\*\*Location:\*\*\s+(.+)`)
		if matches := locationRe.FindStringSubmatch(section); len(matches) > 1 {
			agentsMD.AuditOutput.LocationPattern = strings.TrimSpace(matches[1])
		}
	}

	// Parse format
	formatRe := regexp.MustCompile(`\*\*Format:\*\*\s+(.+)`)
	if matches := formatRe.FindStringSubmatch(section); len(matches) > 1 {
		agentsMD.AuditOutput.Format = strings.TrimSpace(matches[1])
	}
}

func parseQualityCriteria(content string, agentsMD *AgentsMD) {
	section := extractSection(content, "Quality Criteria")
	if section == "" {
		return
	}

	// Parse criteria by category
	categories := []string{"For specifications", "For implementation"}
	for _, category := range categories {
		categorySection := extractSubsection(section, category)
		if categorySection == "" {
			continue
		}

		// Parse list items
		criteriaRe := regexp.MustCompile(`-\s+(.+?)\s+\(PASS/FAIL\)`)
		for _, matches := range criteriaRe.FindAllStringSubmatch(categorySection, -1) {
			if len(matches) > 1 {
				criterion := QualityCriterion{
					Description: strings.TrimSpace(matches[1]),
					Category:    category,
				}
				agentsMD.QualityCriteria = append(agentsMD.QualityCriteria, criterion)
			}
		}
	}
}

// Helper functions

func extractSection(content, sectionName string) string {
	// Match ## Section Name
	sectionRe := regexp.MustCompile(`(?m)^##\s+` + regexp.QuoteMeta(sectionName) + `\s*$`)
	loc := sectionRe.FindStringIndex(content)
	if loc == nil {
		return ""
	}

	start := loc[1]
	
	// Find next ## section
	nextSectionRe := regexp.MustCompile(`(?m)^##\s+`)
	nextLoc := nextSectionRe.FindStringIndex(content[start:])
	
	var end int
	if nextLoc != nil {
		end = start + nextLoc[0]
	} else {
		end = len(content)
	}

	return content[start:end]
}

func extractSubsection(content, subsectionName string) string {
	// Match **Subsection Name:**
	subsectionRe := regexp.MustCompile(`\*\*` + regexp.QuoteMeta(subsectionName) + `[:\*]`)
	loc := subsectionRe.FindStringIndex(content)
	if loc == nil {
		return ""
	}

	start := loc[0]
	
	// Find next ** or ##
	nextRe := regexp.MustCompile(`(?m)(^##\s+|\*\*[A-Z])`)
	nextLoc := nextRe.FindStringIndex(content[start+len(subsectionName)+4:])
	
	var end int
	if nextLoc != nil {
		end = start + len(subsectionName) + 4 + nextLoc[0]
	} else {
		end = len(content)
	}

	return content[start:end]
}

func extractCommand(section, label string) string {
	// Find label followed by code block
	labelRe := regexp.MustCompile(`\*\*` + regexp.QuoteMeta(label) + `[:\*]`)
	loc := labelRe.FindStringIndex(section)
	if loc == nil {
		return ""
	}

	afterLabel := section[loc[1]:]
	
	// Extract code block
	codeBlockRe := regexp.MustCompile("```(?:bash)?\\s*\n([^`]+)```")
	if matches := codeBlockRe.FindStringSubmatch(afterLabel); len(matches) > 1 {
		lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
		// Return first non-comment line
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				return line
			}
		}
	}

	return ""
}

func extractAllCommands(section string) []string {
	var commands []string
	
	codeBlockRe := regexp.MustCompile("```(?:bash)?\\s*\n([^`]+)```")
	for _, matches := range codeBlockRe.FindAllStringSubmatch(section, -1) {
		if len(matches) > 1 {
			lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					commands = append(commands, line)
				}
			}
		}
	}

	return commands
}
