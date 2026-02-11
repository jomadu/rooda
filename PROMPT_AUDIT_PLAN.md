# Prompt Fragment Audit Plan

**Branch:** `audit-prompts`  
**Goal:** Review and improve all prompt fragments in `internal/prompt/fragments/`  
**Date Started:** 2026-02-11

## Audit Status Legend

- [ ] Not Started
- [~] In Progress
- [x] Complete

## Fragment Inventory (Sorted by Usage)

### High-Impact Fragments (10+ uses)

- [ ] `act/emit_success.md` - **16 procedures** (all)
- [ ] `observe/read_agents_md.md` - **16 procedures** (all)
- [ ] `observe/read_specs.md` - **10 procedures**
- [ ] `observe/read_task_input.md` - **10 procedures**
- [ ] `decide/check_if_blocked.md` - **10 procedures**
- [ ] `decide/break_down_into_tasks.md` - **10 procedures**
- [ ] `decide/prioritize_tasks.md` - **10 procedures**
- [ ] `act/write_draft_plan.md` - **10 procedures**

### Medium-Impact Fragments (3-9 uses)

- [ ] `observe/read_impl.md` - **9 procedures**
- [ ] `act/write_audit_report.md` - **3 procedures**

### Low-Impact Fragments (2 uses)

- [ ] `observe/scan_repo_structure.md` - **2 procedures**
- [ ] `observe/detect_build_system.md` - **2 procedures**
- [ ] `observe/query_work_tracking.md` - **2 procedures**
- [ ] `orient/identify_drift.md` - **2 procedures**
- [ ] `orient/evaluate_against_quality_criteria.md` - **2 procedures**
- [ ] `orient/understand_feature_requirements.md` - **2 procedures**
- [ ] `orient/understand_bug_root_cause.md` - **2 procedures**
- [ ] `orient/identify_maintenance_needs.md` - **2 procedures**
- [ ] `orient/identify_affected_code.md` - **2 procedures**
- [ ] `decide/identify_issues.md` - **2 procedures**
- [ ] `decide/prioritize_findings.md` - **2 procedures**
- [ ] `decide/prioritize_gaps_by_impact.md` - **2 procedures**
- [ ] `act/commit_changes.md` - **2 procedures**
- [ ] `act/write_gap_report.md` - **2 procedures**

### Single-Use Fragments (1 use)

- [ ] `observe/detect_work_tracking.md` - **1 procedure**
- [ ] `observe/read_task_details.md` - **1 procedure**
- [ ] `observe/read_draft_plan.md` - **1 procedure**
- [ ] `observe/run_tests.md` - **1 procedure**
- [ ] `observe/run_lints.md` - **1 procedure**
- [ ] `observe/verify_commands.md` - **1 procedure**
- [ ] `orient/compare_detected_vs_documented.md` - **1 procedure**
- [ ] `orient/understand_task_requirements.md` - **1 procedure**
- [ ] `orient/search_codebase.md` - **1 procedure**
- [ ] `orient/identify_affected_files.md` - **1 procedure**
- [ ] `orient/parse_plan_tasks.md` - **1 procedure**
- [ ] `orient/map_to_work_tracking_format.md` - **1 procedure**
- [ ] `orient/compare_documented_vs_actual.md` - **1 procedure**
- [ ] `orient/identify_specified_but_not_implemented.md` - **1 procedure**
- [ ] `orient/identify_implemented_but_not_specified.md` - **1 procedure**
- [ ] `orient/identify_affected_specs.md` - **1 procedure**
- [ ] `orient/identify_spec_deficiencies.md` - **1 procedure**
- [ ] `orient/identify_structural_issues.md` - **1 procedure**
- [ ] `orient/identify_duplication.md` - **1 procedure**
- [ ] `orient/identify_code_smells.md` - **1 procedure**
- [ ] `orient/identify_complexity_issues.md` - **1 procedure**
- [ ] `decide/determine_sections_to_update.md` - **1 procedure**
- [ ] `decide/pick_task.md` - **1 procedure**
- [ ] `decide/plan_implementation_approach.md` - **1 procedure**
- [ ] `decide/determine_import_strategy.md` - **1 procedure**
- [ ] `decide/categorize_drift_severity.md` - **1 procedure**
- [ ] `act/write_agents_md.md` - **1 procedure**
- [ ] `act/modify_files.md` - **1 procedure**
- [ ] `act/run_tests.md` - **1 procedure**
- [ ] `act/update_work_tracking.md` - **1 procedure**
- [ ] `act/create_work_items.md` - **1 procedure**
- [ ] `act/update_draft_plan_status.md` - **1 procedure**

### Unused Fragments (0 uses)

- [ ] `act/emit_failure.md` - **Not currently used**

**Total Fragments:** 57 files (56 unique)

---

## Procedure Breakdown

### 1. agents-sync
**Purpose:** Synchronize AGENTS.md with actual repository state

**Fragments:**
- Observe: `read_agents_md`, `scan_repo_structure`, `detect_build_system`, `detect_work_tracking`
- Orient: `compare_detected_vs_documented`, `identify_drift`
- Decide: `determine_sections_to_update`, `check_if_blocked`
- Act: `write_agents_md`, `commit_changes`, `emit_success`

**Audit Status:** [ ]

---

### 2. build
**Purpose:** Implement a task from work tracking

**Fragments:**
- Observe: `read_agents_md`, `query_work_tracking`, `read_specs`, `read_impl`, `read_task_details`
- Orient: `understand_task_requirements`, `search_codebase`, `identify_affected_files`
- Decide: `pick_task`, `plan_implementation_approach`, `check_if_blocked`
- Act: `modify_files`, `run_tests`, `update_work_tracking`, `commit_changes`, `emit_success`

**Audit Status:** [ ]

---

### 3. publish-plan
**Purpose:** Import draft plan into work tracking system

**Fragments:**
- Observe: `read_agents_md`, `read_draft_plan`, `query_work_tracking`
- Orient: `parse_plan_tasks`, `map_to_work_tracking_format`
- Decide: `determine_import_strategy`, `check_if_blocked`
- Act: `create_work_items`, `update_draft_plan_status`, `emit_success`

**Audit Status:** [ ]

---

### 4. audit-spec
**Purpose:** Audit specification files for quality issues

**Fragments:**
- Observe: `read_agents_md`, `read_specs`
- Orient: `evaluate_against_quality_criteria`
- Decide: `identify_issues`, `prioritize_findings`
- Act: `write_audit_report`, `emit_success`

**Audit Status:** [ ]

---

### 5. audit-impl
**Purpose:** Audit implementation files for quality issues

**Fragments:**
- Observe: `read_agents_md`, `read_impl`, `run_tests`, `run_lints`
- Orient: `evaluate_against_quality_criteria`
- Decide: `identify_issues`, `prioritize_findings`
- Act: `write_audit_report`, `emit_success`

**Audit Status:** [ ]

---

### 6. audit-agents
**Purpose:** Audit AGENTS.md for accuracy and completeness

**Fragments:**
- Observe: `read_agents_md`, `scan_repo_structure`, `detect_build_system`, `verify_commands`
- Orient: `compare_documented_vs_actual`, `identify_drift`
- Decide: `categorize_drift_severity`
- Act: `write_audit_report`, `emit_success`

**Audit Status:** [ ]

---

### 7. audit-spec-to-impl
**Purpose:** Find specifications not implemented in code

**Fragments:**
- Observe: `read_agents_md`, `read_specs`, `read_impl`
- Orient: `identify_specified_but_not_implemented`
- Decide: `prioritize_gaps_by_impact`
- Act: `write_gap_report`, `emit_success`

**Audit Status:** [ ]

---

### 8. audit-impl-to-spec
**Purpose:** Find implementation not covered by specifications

**Fragments:**
- Observe: `read_agents_md`, `read_impl`, `read_specs`
- Orient: `identify_implemented_but_not_specified`
- Decide: `prioritize_gaps_by_impact`
- Act: `write_gap_report`, `emit_success`

**Audit Status:** [ ]

---

### 9. draft-plan-spec-feat
**Purpose:** Create plan for new specification feature

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`, `read_impl`
- Orient: `understand_feature_requirements`, `identify_affected_specs`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 10. draft-plan-spec-fix
**Purpose:** Create plan for specification bug fix

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`, `read_impl`
- Orient: `understand_bug_root_cause`, `identify_spec_deficiencies`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 11. draft-plan-spec-refactor
**Purpose:** Create plan for specification refactoring

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`
- Orient: `identify_structural_issues`, `identify_duplication`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 12. draft-plan-spec-chore
**Purpose:** Create plan for specification maintenance tasks

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`
- Orient: `identify_maintenance_needs`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 13. draft-plan-impl-feat
**Purpose:** Create plan for new implementation feature

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`, `read_impl`
- Orient: `understand_feature_requirements`, `identify_affected_code`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 14. draft-plan-impl-fix
**Purpose:** Create plan for implementation bug fix

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_specs`, `read_impl`
- Orient: `understand_bug_root_cause`, `identify_affected_code`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 15. draft-plan-impl-refactor
**Purpose:** Create plan for code refactoring

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_impl`
- Orient: `identify_code_smells`, `identify_complexity_issues`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

### 16. draft-plan-impl-chore
**Purpose:** Create plan for code maintenance tasks

**Fragments:**
- Observe: `read_agents_md`, `read_task_input`, `read_impl`
- Orient: `identify_maintenance_needs`
- Decide: `break_down_into_tasks`, `prioritize_tasks`, `check_if_blocked`
- Act: `write_draft_plan`, `emit_success`

**Audit Status:** [ ]

---

## Audit Criteria

For each fragment, evaluate:

1. **Clarity** - Is the instruction clear and unambiguous?
2. **Completeness** - Does it cover all necessary aspects?
3. **Consistency** - Does it align with other fragments and OODA principles?
4. **Actionability** - Can an AI agent execute this effectively?
5. **Specificity** - Is it specific enough without being overly prescriptive?
6. **Error Handling** - Does it address failure cases appropriately?

## Audit Process

1. Review fragment content
2. Check usage across procedures
3. Identify improvement opportunities
4. Document findings
5. Implement changes
6. Update status checkboxes

## Notes

- `emit_failure.md` exists but is not used by any procedure - investigate if needed
- `emit_success.md` is used by all 16 procedures - critical fragment
- `check_if_blocked.md` is used by 10 procedures - high-impact fragment
- `read_agents_md.md` is used by all 16 procedures - foundational fragment
