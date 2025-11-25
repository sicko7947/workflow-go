package gorkflow

import (
	"encoding/json"
	"fmt"
	"time"
)

// Package internal provides shared utility functions for the workflow framework.
// This package contains generic helpers used across the workflow implementation.

// ToPtr returns a pointer to the given value.
// This is useful for creating pointers to literals or converting values to pointers.
func ToPtr[T any](v T) *T {
	return &v
}

// CalculateBackoff calculates the backoff delay for a retry attempt.
// It supports three strategies:
//   - EXPONENTIAL: baseDelay * 2^(attempt-1)
//   - LINEAR: baseDelay * attempt
//   - NONE: no backoff delay
//
// Parameters:
//   - baseDelayMs: the base delay in milliseconds
//   - attempt: the current retry attempt (0-based, where 0 = first attempt)
//   - strategy: the backoff strategy ("EXPONENTIAL", "LINEAR", "NONE")
//
// Returns the calculated delay duration. Returns 0 for attempt 0.
func CalculateBackoff(baseDelayMs int, attempt int, strategy string) time.Duration {
	if attempt == 0 {
		return 0
	}

	baseDelay := time.Duration(baseDelayMs) * time.Millisecond

	switch strategy {
	case "EXPONENTIAL":
		// Exponential: baseDelay * 2^(attempt-1)
		multiplier := 1 << (attempt - 1) // 2^(attempt-1)
		return baseDelay * time.Duration(multiplier)
	case "LINEAR":
		// Linear: baseDelay * attempt
		return baseDelay * time.Duration(attempt)
	case "NONE":
		// No backoff
		return 0
	default:
		// Default to linear
		return baseDelay * time.Duration(attempt)
	}
}

// GetRunContext retrieves and deserializes the custom context from a WorkflowRun
func GetRunContext[T any](run *WorkflowRun) (T, error) {
	var zero T
	if len(run.Context) == 0 {
		return zero, fmt.Errorf("workflow run has no context")
	}

	var result T
	if err := json.Unmarshal(run.Context, &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal context: %w", err)
	}
	return result, nil
}
