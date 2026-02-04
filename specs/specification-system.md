# Specification System

## Jobs to be Done (JTBD)

A JTBD describes a high-level user outcome or need. It answers "What does the user want to accomplish?" without prescribing how.

**Examples:**
- Process payments securely
- Search products by category
- Export data in multiple formats

**Characteristics:**
- Outcome-focused, not mechanism-focused
- User-centric (what they want to achieve)
- Broad enough to encompass multiple activities
- Stable over time (implementation details change, JTBDs don't)

## Topics of Concern

A topic of concern is a distinct aspect or component within a JTBD. It represents a cohesive area of functionality that can be specified independently.

**Relationship:** 1 JTBD > multiple topics of concern

**Examples:**
- JTBD: "Process payments securely"
  - Topics: Payment validation, Transaction processing, Refund handling, Receipt generation

**Scope Test:** Can you describe it in one sentence without "and"?
- ✓ "Transaction processing executes payment authorization and captures funds"
- ✗ "User management handles authentication, profiles, and billing" > 3 topics

## Specification Documents

Each topic of concern has one specification document in `specs/*.md`.

**Structure:** See [TEMPLATE.md](TEMPLATE.md) for the complete specification template.

**Purpose:**
- Single source of truth for requirements
- Consumed by planning and building phases
- Defines acceptance criteria for backpressure (tests)
- Maps to implementation files

## README Structure

The README.md follows this organization:

1. **Header** - Visual identity and branding
2. **What is [Project]?** - Core value proposition and purpose
3. **Why [Project]?** - Problems solved and key features
4. **Installation** - Quick install, specific version, manual, verification, uninstall
5. **Quick Start** - Core concepts, setup steps, common commands
6. **Documentation** - Links to detailed docs organized by category

**Principles:**
- Progressive disclosure (quick start > detailed docs)
- Action-oriented (commands you can run immediately)
- Problem-solution framing (why before how)
- Links to detailed documentation (README stays concise)

## Specs Folder Structure

```
specs/
├── README.md              # Index of all JTBDs, topics, and specs
├── TEMPLATE.md            # Template for new specification documents
├── [topic-name].md        # One spec per topic of concern
└── ...
```

**specs/README.md** contains:
- **Jobs to be Done** - Grouped by category
- **Topics of Concern** - Grouped by category
- **Specification Documents** - Grouped by category with links to individual specs

**Naming Convention:**
- Lowercase with hyphens: `payment-validation.md`, `transaction-processing.md`
- Descriptive of the topic: name should clearly indicate what's specified
- One topic per file (no combining unrelated concerns)

**Organization:**
- specs/README.md provides the index and navigation
- Individual specs are flat (no subdirectories)
- TEMPLATE.md provides consistent structure for new specs
- Each spec is self-contained but references related specs in "Implementation Mapping"
