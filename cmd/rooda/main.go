package main

import (
	"fmt"
	"os"

	"github.com/jomadu/rooda/internal/config"
	"github.com/jomadu/rooda/internal/procedures"
)

func init() {
	// Register built-in procedures with config package
	config.BuiltInProceduresFunc = procedures.BuiltInProcedures
}

var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

const (
	ExitSuccess        = 0
	ExitUserError      = 1
	ExitConfigError    = 2
	ExitExecutionError = 3
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitUserError)
	}
}
