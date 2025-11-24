package engine

import (
	"time"

	"github.com/sicko7947/workflow-go"
)

// calculateBackoff is a wrapper around the internal helper
func calculateBackoff(baseDelayMs int, attempt int, strategy string) time.Duration {
	return workflow.CalculateBackoff(baseDelayMs, attempt, strategy)
}
