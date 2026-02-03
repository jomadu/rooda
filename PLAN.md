# Draft Plan: Spec to Implementation Gap Analysis

## Priority 1: Core Framework Behavior (Missing Specs)

### 1. CLI Interface Specification
**Gap:** rooda.sh has complex argument parsing (procedure names, explicit OODA flags, config override, max-iterations) but no specification defining the CLI interface design.

**Tasks:**
- Document CLI invocation patterns (procedure-based vs explicit flags)
- Specify argument precedence (explicit flags override config)
- Define error handling for invalid arguments
- Specify config file resolution (relative to script location)

**Acceptance Criteria:**
- [ ] All CLI invocation patterns documented
- [ ] Argument parsing behavior specified
- [ ] Error messages defined for invalid usage

### 2. Prompt Composition Specification
**Gap:** create_prompt() function concatenates 4 OODA phase files into single prompt, but no specification for how prompt assembly works.

**Tasks:**
- Document prompt composition algorithm (concatenate observe + orient + decide + act)
- Specify file reading and validation
- Define error handling for missing prompt files
- Specify output format (single combined prompt to stdout)

**Acceptance Criteria:**
- [ ] Prompt assembly algorithm documented
- [ ] File validation behavior specified
- [ ] Error handling for missing files defined

### 3. Iteration Loop Control Specification
**Gap:** Script supports --max-iterations, default iterations per procedure, but no specification for iteration behavior.

**Tasks:**
- Document iteration loop semantics (exit and restart per iteration)
- Specify max-iterations override behavior
- Define default iteration fallback (from config or 1)
- Specify termination conditions (max reached, Ctrl+C)

**Acceptance Criteria:**
- [ ] Iteration loop behavior documented
- [ ] Termination conditions specified
- [ ] Default iteration logic defined

### 4. Configuration Schema Specification
**Gap:** rooda-config.yml has specific structure (procedures, display, summary, description, OODA paths, default_iterations) but no specification.

**Tasks:**
- Document YAML configuration schema
- Specify required vs optional fields per procedure
- Define procedure lookup mechanism (yq queries)
- Specify validation for missing procedures

**Acceptance Criteria:**
- [ ] Configuration schema documented
- [ ] Field requirements specified
- [ ] Lookup and validation behavior defined

## Priority 2: External Dependencies (Missing Specs)

### 5. AI CLI Integration Specification
**Gap:** Script pipes prompt to kiro-cli with specific flags (--no-interactive --trust-all-tools) but no specification for AI CLI integration.

**Tasks:**
- Document AI CLI requirements (kiro-cli or compatible)
- Specify required flags and their purpose
- Define expected input/output format
- Specify error handling for CLI failures

**Acceptance Criteria:**
- [ ] AI CLI integration requirements documented
- [ ] Flag usage and rationale specified
- [ ] Error handling defined

### 6. External Dependencies Specification
**Gap:** Script requires yq for YAML parsing but no specification for external dependencies.

**Tasks:**
- Document all external dependencies (yq, kiro-cli, bd)
- Specify version requirements if applicable
- Define dependency checking mechanism
- Specify installation instructions per platform

**Acceptance Criteria:**
- [ ] All dependencies documented
- [ ] Version requirements specified
- [ ] Installation instructions provided

## Priority 3: Quality Enforcement (Missing Implementation)

### 7. AGENTS.md Section Validation
**Gap:** agents-md-format.md defines required sections but bootstrap doesn't validate completeness.

**Tasks:**
- Add validation in bootstrap to check for required sections
- Warn if sections are missing or incomplete
- Provide guidance on what to add
- Update AGENTS.md with validation results

**Acceptance Criteria:**
- [ ] Bootstrap validates AGENTS.md structure
- [ ] Missing sections trigger warnings
- [ ] Guidance provided for incomplete sections

### 8. Quality Criteria Boolean Enforcement
**Gap:** Specs emphasize boolean criteria but no validation that AGENTS.md quality criteria are actually boolean.

**Tasks:**
- Add validation in quality assessment procedures
- Check that criteria are PASS/FAIL (not subjective scores)
- Warn if criteria are ambiguous or subjective
- Provide examples of boolean vs non-boolean criteria

**Acceptance Criteria:**
- [ ] Quality procedures validate boolean criteria
- [ ] Non-boolean criteria trigger warnings
- [ ] Examples provided for correction

## Priority 4: Spec System Automation (Missing Implementation)

### 9. Spec Template Structure Validation
**Gap:** TEMPLATE.md defines structure but no validation or enforcement in rooda.sh.

**Tasks:**
- Add validation in build procedure when creating/updating specs
- Check for required sections (JTBD, Activities, Acceptance Criteria, etc.)
- Warn if sections are missing
- Suggest using TEMPLATE.md for new specs

**Acceptance Criteria:**
- [ ] Build procedure validates spec structure
- [ ] Missing sections trigger warnings
- [ ] Template usage suggested for new specs

### 10. Spec Naming Convention Validation
**Gap:** specification-system.md specifies lowercase-with-hyphens naming but no validation.

**Tasks:**
- Add validation in build procedure when creating specs
- Check filename matches convention (lowercase, hyphens, .md extension)
- Warn if naming doesn't match convention
- Suggest correct naming

**Acceptance Criteria:**
- [ ] Build procedure validates spec naming
- [ ] Non-conforming names trigger warnings
- [ ] Correct naming suggested

### 11. Spec README Index Generation
**Gap:** specification-system.md describes specs/README.md index but no automation to generate/maintain.

**Tasks:**
- Create procedure to generate specs/README.md from existing specs
- Extract JTBD from each spec file
- Group by category (if categorization exists)
- Link to individual spec files

**Acceptance Criteria:**
- [ ] Automated index generation implemented
- [ ] JTBDs extracted from specs
- [ ] Links to spec files included

## Priority 5: Documentation Alignment (Drift)

### 12. Component Path Convention Documentation
**Gap:** rooda-config.yml uses src/components/*.md but agents-md-format.md examples show prompts/*.md.

**Tasks:**
- Clarify in agents-md-format.md that examples use consumer convention (prompts/)
- Document that framework internal structure uses src/components/
- Explain that consumers copy to project root (flat structure)
- Update examples to show both perspectives

**Acceptance Criteria:**
- [ ] Path convention clarified in agents-md-format.md
- [ ] Consumer vs framework structure explained
- [ ] Examples updated for clarity

---

## Implementation Notes

**Rationale for prioritization:**
- Priority 1: Core behavior specs enable understanding how the framework works
- Priority 2: Dependency specs enable proper installation and troubleshooting
- Priority 3: Quality enforcement improves AGENTS.md and spec quality
- Priority 4: Spec automation reduces manual maintenance burden
- Priority 5: Documentation alignment prevents confusion

**Dependencies:**
- Tasks 1-6 are independent (can be parallelized)
- Tasks 7-8 depend on understanding bootstrap and quality procedures (tasks 1-4)
- Tasks 9-11 depend on understanding build procedure (task 1)
- Task 12 is independent (documentation only)

**Estimated scope:**
- Each spec task: 1 iteration (write spec following TEMPLATE.md)
- Each validation task: 1-2 iterations (add validation logic, test, update AGENTS.md)
- Index generation: 2-3 iterations (implement, test, integrate with build)
