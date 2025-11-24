package workflow

import (
	"testing"
	"time"
)

func TestToPtr(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "int value",
			value: 42,
		},
		{
			name:  "string value",
			value: "test",
		},
		{
			name:  "bool value",
			value: true,
		},
		{
			name:  "struct value",
			value: struct{ name string }{name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				ptr := ToPtr(v)
				if ptr == nil {
					t.Fatal("ToPtr returned nil")
				}
				if *ptr != v {
					t.Errorf("ToPtr() = %v, want %v", *ptr, v)
				}
			case string:
				ptr := ToPtr(v)
				if ptr == nil {
					t.Fatal("ToPtr returned nil")
				}
				if *ptr != v {
					t.Errorf("ToPtr() = %v, want %v", *ptr, v)
				}
			case bool:
				ptr := ToPtr(v)
				if ptr == nil {
					t.Fatal("ToPtr returned nil")
				}
				if *ptr != v {
					t.Errorf("ToPtr() = %v, want %v", *ptr, v)
				}
			}
		})
	}
}

func TestToPtr_ModifyOriginal(t *testing.T) {
	// Verify that modifying the original doesn't affect the pointer
	original := 10
	ptr := ToPtr(original)
	original = 20

	if *ptr != 10 {
		t.Errorf("Pointer value changed unexpectedly: got %d, want 10", *ptr)
	}
}

func TestCalculateBackoff_FirstAttempt(t *testing.T) {
	strategies := []string{"EXPONENTIAL", "LINEAR", "NONE", "unknown"}

	for _, strategy := range strategies {
		t.Run(strategy, func(t *testing.T) {
			delay := CalculateBackoff(100, 0, strategy)
			if delay != 0 {
				t.Errorf("CalculateBackoff(100, 0, %s) = %v, want 0", strategy, delay)
			}
		})
	}
}

func TestCalculateBackoff_Exponential(t *testing.T) {
	tests := []struct {
		baseDelayMs int
		attempt     int
		want        time.Duration
	}{
		{100, 1, 100 * time.Millisecond},  // 100 * 2^0 = 100
		{100, 2, 200 * time.Millisecond},  // 100 * 2^1 = 200
		{100, 3, 400 * time.Millisecond},  // 100 * 2^2 = 400
		{100, 4, 800 * time.Millisecond},  // 100 * 2^3 = 800
		{100, 5, 1600 * time.Millisecond}, // 100 * 2^4 = 1600
		{50, 3, 200 * time.Millisecond},   // 50 * 2^2 = 200
		{200, 2, 400 * time.Millisecond},  // 200 * 2^1 = 400
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := CalculateBackoff(tt.baseDelayMs, tt.attempt, "EXPONENTIAL")
			if got != tt.want {
				t.Errorf("CalculateBackoff(%d, %d, EXPONENTIAL) = %v, want %v",
					tt.baseDelayMs, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestCalculateBackoff_Linear(t *testing.T) {
	tests := []struct {
		baseDelayMs int
		attempt     int
		want        time.Duration
	}{
		{100, 1, 100 * time.Millisecond}, // 100 * 1 = 100
		{100, 2, 200 * time.Millisecond}, // 100 * 2 = 200
		{100, 3, 300 * time.Millisecond}, // 100 * 3 = 300
		{100, 4, 400 * time.Millisecond}, // 100 * 4 = 400
		{100, 5, 500 * time.Millisecond}, // 100 * 5 = 500
		{50, 3, 150 * time.Millisecond},  // 50 * 3 = 150
		{200, 2, 400 * time.Millisecond}, // 200 * 2 = 400
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := CalculateBackoff(tt.baseDelayMs, tt.attempt, "LINEAR")
			if got != tt.want {
				t.Errorf("CalculateBackoff(%d, %d, LINEAR) = %v, want %v",
					tt.baseDelayMs, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestCalculateBackoff_None(t *testing.T) {
	tests := []struct {
		baseDelayMs int
		attempt     int
	}{
		{100, 1},
		{100, 2},
		{100, 5},
		{50, 3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := CalculateBackoff(tt.baseDelayMs, tt.attempt, "NONE")
			if got != 0 {
				t.Errorf("CalculateBackoff(%d, %d, NONE) = %v, want 0",
					tt.baseDelayMs, tt.attempt, got)
			}
		})
	}
}

func TestCalculateBackoff_UnknownStrategy(t *testing.T) {
	// Unknown strategy should default to LINEAR
	tests := []struct {
		baseDelayMs int
		attempt     int
		want        time.Duration
	}{
		{100, 1, 100 * time.Millisecond}, // 100 * 1 = 100
		{100, 2, 200 * time.Millisecond}, // 100 * 2 = 200
		{100, 3, 300 * time.Millisecond}, // 100 * 3 = 300
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := CalculateBackoff(tt.baseDelayMs, tt.attempt, "UNKNOWN")
			if got != tt.want {
				t.Errorf("CalculateBackoff(%d, %d, UNKNOWN) = %v, want %v (default LINEAR)",
					tt.baseDelayMs, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestCalculateBackoff_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		baseDelayMs int
		attempt     int
		strategy    string
		want        time.Duration
	}{
		{
			name:        "zero base delay",
			baseDelayMs: 0,
			attempt:     3,
			strategy:    "LINEAR",
			want:        0,
		},
		{
			name:        "zero base delay exponential",
			baseDelayMs: 0,
			attempt:     3,
			strategy:    "EXPONENTIAL",
			want:        0,
		},
		{
			name:        "large attempt number",
			baseDelayMs: 100,
			attempt:     10,
			strategy:    "LINEAR",
			want:        1000 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateBackoff(tt.baseDelayMs, tt.attempt, tt.strategy)
			if got != tt.want {
				t.Errorf("CalculateBackoff(%d, %d, %s) = %v, want %v",
					tt.baseDelayMs, tt.attempt, tt.strategy, got, tt.want)
			}
		})
	}
}
