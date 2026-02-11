package agents

// AgentsMD represents parsed AGENTS.md content
type AgentsMD struct {
	BuildCommand    string
	TestCommand     string
	LintCommands    []string
	SpecPaths       []string
	SpecExcludes    []string
	ImplPaths       []string
	ImplExcludes    []string
	DocsPaths       []string
	WorkTracking    WorkTrackingConfig
	TaskInput       TaskInputConfig
	PlanningSystem  PlanningConfig
	AuditOutput     AuditOutputConfig
	QualityCriteria []QualityCriterion
	RawContent      string
	FilePath        string
}

// WorkTrackingConfig represents work tracking system configuration
type WorkTrackingConfig struct {
	System        string
	QueryCommand  string
	UpdateCommand string
	CloseCommand  string
	CreateCommand string
}

// TaskInputConfig represents task input configuration
type TaskInputConfig struct {
	Location string
	Format   string
}

// PlanningConfig represents planning system configuration
type PlanningConfig struct {
	DraftPlanLocation string
	PublishingMethod  string
}

// AuditOutputConfig represents audit output configuration
type AuditOutputConfig struct {
	LocationPattern string
	Format          string
}

// QualityCriterion represents a single quality check
type QualityCriterion struct {
	Description string
	Command     string
	PassPattern string
	Category    string
}

// DriftDetection represents detected drift between AGENTS.md and reality
type DriftDetection struct {
	Field      string
	Expected   string
	Actual     string
	FixApplied string
	Rationale  string
}
