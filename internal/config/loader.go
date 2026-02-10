package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"gopkg.in/yaml.v3"
)

// CLIFlags represents command-line flag overrides
type CLIFlags struct {
	MaxIterations *int
	AICmd         string
	AICmdAlias    string
	ConfigPath    string
}

// LoadConfig loads and merges configuration from all tiers
func LoadConfig(cliFlags CLIFlags) (*Config, error) {
	// 1. Start with built-in defaults
	config := builtInDefaults()
	provenance := initProvenance(config)

	// 2. Resolve global config directory and load config
	globalDir := resolveGlobalConfigDir()
	globalPath := filepath.Join(globalDir, "rooda-config.yml")
	if fileExists(globalPath) {
		globalConfig, err := parseYAML(globalPath)
		if err != nil {
			return nil, fmt.Errorf("global config %s: %w", globalPath, err)
		}
		mergeConfig(config, globalConfig, provenance, TierGlobal, globalPath, globalDir)
	}

	// 3. Load workspace config
	workspacePath := "./rooda-config.yml"
	if cliFlags.ConfigPath != "" {
		workspacePath = cliFlags.ConfigPath
	}
	if fileExists(workspacePath) {
		workspaceConfig, err := parseYAML(workspacePath)
		if err != nil {
			return nil, fmt.Errorf("workspace config %s: %w", workspacePath, err)
		}
		workspaceDir := filepath.Dir(workspacePath)
		if workspaceDir == "" {
			workspaceDir = "."
		}
		mergeConfig(config, workspaceConfig, provenance, TierWorkspace, workspacePath, workspaceDir)
	}

	// 4. Apply environment variables
	applyEnvVars(config, provenance)

	// 5. Apply CLI flags
	if cliFlags.MaxIterations != nil {
		config.Loop.DefaultMaxIterations = cliFlags.MaxIterations
		provenance["loop.default_max_iterations"] = ConfigSource{TierCLIFlag, "", *cliFlags.MaxIterations}
	}
	if cliFlags.AICmd != "" {
		config.Loop.AICmd = cliFlags.AICmd
		provenance["loop.ai_cmd"] = ConfigSource{TierCLIFlag, "", cliFlags.AICmd}
	}
	if cliFlags.AICmdAlias != "" {
		config.Loop.AICmdAlias = cliFlags.AICmdAlias
		provenance["loop.ai_cmd_alias"] = ConfigSource{TierCLIFlag, "", cliFlags.AICmdAlias}
	}

	// Assign provenance to config
	config.Provenance = provenance

	return config, nil
}

// resolveGlobalConfigDir resolves the global config directory
func resolveGlobalConfigDir() string {
	// 1. ROODA_CONFIG_HOME env var
	if dir := os.Getenv("ROODA_CONFIG_HOME"); dir != "" {
		return dir
	}

	// 2. XDG_CONFIG_HOME
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "rooda")
	}

	// 3. Platform-specific defaults
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "rooda")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "rooda")
}

// builtInDefaults returns the built-in default configuration
func builtInDefaults() *Config {
	maxIter := DefaultMaxIterations
	return &Config{
		Loop: LoopConfig{
			IterationMode:        DefaultIterationMode,
			DefaultMaxIterations: &maxIter,
			MaxOutputBuffer:      DefaultMaxOutputBuffer,
			FailureThreshold:     DefaultFailureThreshold,
			LogLevel:             DefaultLogLevel,
			LogTimestampFormat:   DefaultTimestampFormat,
			ShowAIOutput:         DefaultShowAIOutput,
		},
		Procedures:   make(map[string]Procedure),
		AICmdAliases: builtInAliases(),
		Provenance:   make(map[string]ConfigSource),
	}
}

// builtInAliases returns the built-in AI command aliases
func builtInAliases() map[string]string {
	return map[string]string{
		"kiro-cli":     "kiro-cli chat --no-interactive --trust-all-tools",
		"claude":       "claude -p --dangerously-skip-permissions",
		"copilot":      "copilot --yolo",
		"cursor-agent": "cursor-wrapper.sh",
	}
}

// initProvenance initializes provenance tracking for built-in defaults
func initProvenance(config *Config) map[string]ConfigSource {
	p := make(map[string]ConfigSource)
	p["loop.iteration_mode"] = ConfigSource{TierBuiltIn, "", config.Loop.IterationMode}
	p["loop.default_max_iterations"] = ConfigSource{TierBuiltIn, "", *config.Loop.DefaultMaxIterations}
	p["loop.max_output_buffer"] = ConfigSource{TierBuiltIn, "", config.Loop.MaxOutputBuffer}
	p["loop.failure_threshold"] = ConfigSource{TierBuiltIn, "", config.Loop.FailureThreshold}
	p["loop.log_level"] = ConfigSource{TierBuiltIn, "", config.Loop.LogLevel}
	p["loop.log_timestamp_format"] = ConfigSource{TierBuiltIn, "", config.Loop.LogTimestampFormat}
	p["loop.show_ai_output"] = ConfigSource{TierBuiltIn, "", config.Loop.ShowAIOutput}
	for name, cmd := range config.AICmdAliases {
		p["ai_cmd_aliases."+name] = ConfigSource{TierBuiltIn, "", cmd}
	}
	return p
}

// configFile represents the YAML config file structure
type configFile struct {
	Loop struct {
		IterationMode        string `yaml:"iteration_mode"`
		DefaultMaxIterations *int   `yaml:"default_max_iterations"`
		IterationTimeout     *int   `yaml:"iteration_timeout"`
		MaxOutputBuffer      int    `yaml:"max_output_buffer"`
		FailureThreshold     int    `yaml:"failure_threshold"`
		LogLevel             string `yaml:"log_level"`
		LogTimestampFormat   string `yaml:"log_timestamp_format"`
		ShowAIOutput         bool   `yaml:"show_ai_output"`
		AICmd                string `yaml:"ai_cmd"`
		AICmdAlias           string `yaml:"ai_cmd_alias"`
	} `yaml:"loop"`
	AICmdAliases map[string]string            `yaml:"ai_cmd_aliases"`
	Procedures   map[string]procedureYAML     `yaml:"procedures"`
}

type procedureYAML struct {
	Display              string                   `yaml:"display"`
	Summary              string                   `yaml:"summary"`
	Description          string                   `yaml:"description"`
	Observe              []fragmentActionYAML     `yaml:"observe"`
	Orient               []fragmentActionYAML     `yaml:"orient"`
	Decide               []fragmentActionYAML     `yaml:"decide"`
	Act                  []fragmentActionYAML     `yaml:"act"`
	IterationMode        string                   `yaml:"iteration_mode"`
	DefaultMaxIterations *int                     `yaml:"default_max_iterations"`
	IterationTimeout     *int                     `yaml:"iteration_timeout"`
	MaxOutputBuffer      *int                     `yaml:"max_output_buffer"`
	AICmd                string                   `yaml:"ai_cmd"`
	AICmdAlias           string                   `yaml:"ai_cmd_alias"`
}

type fragmentActionYAML struct {
	Content    string                 `yaml:"content"`
	Path       string                 `yaml:"path"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

// parseYAML parses a YAML config file
func parseYAML(path string) (*configFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// mergeConfig merges overlay config into base config
func mergeConfig(base *Config, overlay *configFile, provenance map[string]ConfigSource, tier ConfigTier, filePath string, configDir string) {
	// Merge loop settings
	if overlay.Loop.IterationMode != "" {
		base.Loop.IterationMode = IterationMode(overlay.Loop.IterationMode)
		provenance["loop.iteration_mode"] = ConfigSource{tier, filePath, overlay.Loop.IterationMode}
	}
	if overlay.Loop.DefaultMaxIterations != nil {
		base.Loop.DefaultMaxIterations = overlay.Loop.DefaultMaxIterations
		provenance["loop.default_max_iterations"] = ConfigSource{tier, filePath, *overlay.Loop.DefaultMaxIterations}
	}
	if overlay.Loop.IterationTimeout != nil {
		base.Loop.IterationTimeout = overlay.Loop.IterationTimeout
		provenance["loop.iteration_timeout"] = ConfigSource{tier, filePath, *overlay.Loop.IterationTimeout}
	}
	if overlay.Loop.MaxOutputBuffer != 0 {
		base.Loop.MaxOutputBuffer = overlay.Loop.MaxOutputBuffer
		provenance["loop.max_output_buffer"] = ConfigSource{tier, filePath, overlay.Loop.MaxOutputBuffer}
	}
	if overlay.Loop.FailureThreshold != 0 {
		base.Loop.FailureThreshold = overlay.Loop.FailureThreshold
		provenance["loop.failure_threshold"] = ConfigSource{tier, filePath, overlay.Loop.FailureThreshold}
	}
	if overlay.Loop.LogLevel != "" {
		base.Loop.LogLevel = LogLevel(overlay.Loop.LogLevel)
		provenance["loop.log_level"] = ConfigSource{tier, filePath, overlay.Loop.LogLevel}
	}
	if overlay.Loop.LogTimestampFormat != "" {
		base.Loop.LogTimestampFormat = TimestampFormat(overlay.Loop.LogTimestampFormat)
		provenance["loop.log_timestamp_format"] = ConfigSource{tier, filePath, overlay.Loop.LogTimestampFormat}
	}
	if overlay.Loop.ShowAIOutput {
		base.Loop.ShowAIOutput = overlay.Loop.ShowAIOutput
		provenance["loop.show_ai_output"] = ConfigSource{tier, filePath, overlay.Loop.ShowAIOutput}
	}
	if overlay.Loop.AICmd != "" {
		base.Loop.AICmd = overlay.Loop.AICmd
		provenance["loop.ai_cmd"] = ConfigSource{tier, filePath, overlay.Loop.AICmd}
	}
	if overlay.Loop.AICmdAlias != "" {
		base.Loop.AICmdAlias = overlay.Loop.AICmdAlias
		provenance["loop.ai_cmd_alias"] = ConfigSource{tier, filePath, overlay.Loop.AICmdAlias}
	}

	// Merge AI command aliases
	for name, command := range overlay.AICmdAliases {
		base.AICmdAliases[name] = command
		provenance["ai_cmd_aliases."+name] = ConfigSource{tier, filePath, command}
	}

	// Merge procedures
	for name, proc := range overlay.Procedures {
		baseProcedure, exists := base.Procedures[name]
		if !exists {
			baseProcedure = Procedure{}
		}

		// Merge fields
		if proc.Display != "" {
			baseProcedure.Display = proc.Display
		}
		if proc.Summary != "" {
			baseProcedure.Summary = proc.Summary
		}
		if proc.Description != "" {
			baseProcedure.Description = proc.Description
		}
		if len(proc.Observe) > 0 {
			baseProcedure.Observe = resolveFragmentPaths(configDir, proc.Observe)
		}
		if len(proc.Orient) > 0 {
			baseProcedure.Orient = resolveFragmentPaths(configDir, proc.Orient)
		}
		if len(proc.Decide) > 0 {
			baseProcedure.Decide = resolveFragmentPaths(configDir, proc.Decide)
		}
		if len(proc.Act) > 0 {
			baseProcedure.Act = resolveFragmentPaths(configDir, proc.Act)
		}
		if proc.IterationMode != "" {
			baseProcedure.IterationMode = IterationMode(proc.IterationMode)
		}
		if proc.DefaultMaxIterations != nil {
			baseProcedure.DefaultMaxIterations = proc.DefaultMaxIterations
		}
		if proc.IterationTimeout != nil {
			baseProcedure.IterationTimeout = proc.IterationTimeout
		}
		if proc.MaxOutputBuffer != nil {
			baseProcedure.MaxOutputBuffer = proc.MaxOutputBuffer
		}
		if proc.AICmd != "" {
			baseProcedure.AICmd = proc.AICmd
		}
		if proc.AICmdAlias != "" {
			baseProcedure.AICmdAlias = proc.AICmdAlias
		}

		base.Procedures[name] = baseProcedure
		provenance["procedures."+name] = ConfigSource{tier, filePath, baseProcedure}
	}
}

// resolveFragmentPaths resolves fragment paths relative to config directory
func resolveFragmentPaths(configDir string, fragments []fragmentActionYAML) []FragmentAction {
	resolved := make([]FragmentAction, len(fragments))
	for i, frag := range fragments {
		resolved[i] = FragmentAction{
			Content:    frag.Content,
			Path:       frag.Path,
			Parameters: frag.Parameters,
		}
		// Only resolve path if not builtin: and path is specified
		if frag.Path != "" && !filepath.IsAbs(frag.Path) && frag.Path[:8] != "builtin:" {
			resolved[i].Path = filepath.Join(configDir, frag.Path)
		}
	}
	return resolved
}

// applyEnvVars applies environment variable overrides
func applyEnvVars(config *Config, provenance map[string]ConfigSource) {
	if v := os.Getenv("ROODA_LOOP_AI_CMD"); v != "" {
		config.Loop.AICmd = v
		provenance["loop.ai_cmd"] = ConfigSource{TierEnvVar, "", v}
	}
	if v := os.Getenv("ROODA_LOOP_AI_CMD_ALIAS"); v != "" {
		config.Loop.AICmdAlias = v
		provenance["loop.ai_cmd_alias"] = ConfigSource{TierEnvVar, "", v}
	}
	if v := os.Getenv("ROODA_LOOP_ITERATION_MODE"); v != "" {
		config.Loop.IterationMode = IterationMode(v)
		provenance["loop.iteration_mode"] = ConfigSource{TierEnvVar, "", v}
	}
	if v := os.Getenv("ROODA_LOOP_DEFAULT_MAX_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			config.Loop.DefaultMaxIterations = &n
			provenance["loop.default_max_iterations"] = ConfigSource{TierEnvVar, "", n}
		}
	}
	if v := os.Getenv("ROODA_LOOP_ITERATION_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			config.Loop.IterationTimeout = &n
			provenance["loop.iteration_timeout"] = ConfigSource{TierEnvVar, "", n}
		}
	}
	if v := os.Getenv("ROODA_LOOP_FAILURE_THRESHOLD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			config.Loop.FailureThreshold = n
			provenance["loop.failure_threshold"] = ConfigSource{TierEnvVar, "", n}
		}
	}
	if v := os.Getenv("ROODA_LOOP_LOG_LEVEL"); v != "" {
		config.Loop.LogLevel = LogLevel(v)
		provenance["loop.log_level"] = ConfigSource{TierEnvVar, "", v}
	}
	if v := os.Getenv("ROODA_LOOP_LOG_TIMESTAMP_FORMAT"); v != "" {
		config.Loop.LogTimestampFormat = TimestampFormat(v)
		provenance["loop.log_timestamp_format"] = ConfigSource{TierEnvVar, "", v}
	}
	if v := os.Getenv("ROODA_LOOP_SHOW_AI_OUTPUT"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			config.Loop.ShowAIOutput = b
			provenance["loop.show_ai_output"] = ConfigSource{TierEnvVar, "", b}
		}
	}
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
