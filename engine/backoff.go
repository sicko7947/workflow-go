package engine

import (
	"time"

	"github.com/sicko7947/gorkflow"
)

// calculateBackoff is a wrapper around the internal helper
func calculateBackoff(baseDelayMs int, attempt int, strategy string) time.Duration {
	return gorkflow.CalculateBackoff(baseDelayMs, attempt, strategy)
}
