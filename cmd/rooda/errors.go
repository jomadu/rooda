package main

import "errors"

var (
	ErrInvalidMaxIterations = errors.New("--max-iterations must be >= 1")
	ErrEmptyContext         = errors.New("empty inline content not allowed for --context flag")
	ErrEmptyFragment        = errors.New("empty inline content not allowed for OODA phase flag")
)
