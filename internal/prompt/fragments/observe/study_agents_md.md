# Study AGENTS.md

You must study and parse the AGENTS.md configuration file from the repository root. Use the file reading tool to access this file.

**If AGENTS.md doesn't exist:** Note this fact - the repository needs bootstrap procedure first.

**Extract information for all 10 required topics** (search flexibly - exact section names/formats may vary):

1. **Work Tracking System** - system name, onboarding, query/update/close/create commands
2. **Quick Reference** - cross-section command summary
3. **Task Input** - location and format of task documentation
4. **Planning System** - draft plan location, publishing mechanism
5. **Build/Test/Lint Commands** - dependencies, test/build/lint commands
6. **Specification Definition** - location, format, exclude patterns, current state
7. **Implementation Definition** - location, patterns, exclude paths
8. **Audit Output** - location pattern, format
9. **Quality Criteria** - for specs, for implementation, refactoring triggers
10. **Operational Learnings** - last verification date, verified working/not working, rationale

**Parsing approach:**
- Search for topic names in headers (any level), bold text, or plain text
- Extract commands from code blocks or inline code
- Extract paths/patterns from any format (bold labels, plain text, lists)
- Infer missing information from context when possible

**Validate:**
- All 10 required topics covered (if missing, note which for error reporting)
- Commands are identifiable and extractable
- File paths/patterns are present

**Remember extracted information** for use in Orient/Decide/Act phases - you'll reference these commands, paths, and criteria throughout the loop iteration.
