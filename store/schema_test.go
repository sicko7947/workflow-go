package store

import (
	"strings"
	"testing"
)

func TestWorkflowRunPK(t *testing.T) {
	tests := []struct {
		name  string
		runID string
		want  string
	}{
		{
			name:  "simple run ID",
			runID: "test-run-1",
			want:  "RUN#test-run-1",
		},
		{
			name:  "UUID run ID",
			runID: "550e8400-e29b-41d4-a716-446655440000",
			want:  "RUN#550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:  "empty run ID",
			runID: "",
			want:  "RUN#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := workflowRunPK(tt.runID)
			if got != tt.want {
				t.Errorf("workflowRunPK(%s) = %s, want %s", tt.runID, got, tt.want)
			}
		})
	}
}

func TestWorkflowRunSK(t *testing.T) {
	got := workflowRunSK()
	want := "META"
	if got != want {
		t.Errorf("workflowRunSK() = %s, want %s", got, want)
	}
}

func TestWorkflowRunGSI1PK(t *testing.T) {
	tests := []struct {
		name       string
		workflowID string
		status     string
		want       string
	}{
		{
			name:       "pending status",
			workflowID: "email-workflow",
			status:     "PENDING",
			want:       "WF#email-workflow#STATUS#PENDING",
		},
		{
			name:       "running status",
			workflowID: "email-workflow",
			status:     "RUNNING",
			want:       "WF#email-workflow#STATUS#RUNNING",
		},
		{
			name:       "completed status",
			workflowID: "email-workflow",
			status:     "COMPLETED",
			want:       "WF#email-workflow#STATUS#COMPLETED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := workflowRunGSI1PK(tt.workflowID, tt.status)
			if got != tt.want {
				t.Errorf("workflowRunGSI1PK(%s, %s) = %s, want %s",
					tt.workflowID, tt.status, got, tt.want)
			}
		})
	}
}

func TestWorkflowRunGSI1SK(t *testing.T) {
	createdAt := "2023-10-15T10:30:00Z"
	got := workflowRunGSI1SK(createdAt)
	if got != createdAt {
		t.Errorf("workflowRunGSI1SK(%s) = %s, want %s", createdAt, got, createdAt)
	}
}

func TestWorkflowRunGSI2PK(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		status     string
		want       string
	}{
		{
			name:       "email resource with pending status",
			resourceID: "email-123",
			status:     "PENDING",
			want:       "RES#email-123#STATUS#PENDING",
		},
		{
			name:       "user resource with running status",
			resourceID: "user-456",
			status:     "RUNNING",
			want:       "RES#user-456#STATUS#RUNNING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := workflowRunGSI2PK(tt.resourceID, tt.status)
			if got != tt.want {
				t.Errorf("workflowRunGSI2PK(%s, %s) = %s, want %s", tt.resourceID, tt.status, got, tt.want)
			}
		})
	}
}

func TestWorkflowRunGSI2SK(t *testing.T) {
	tests := []struct {
		name      string
		createdAt string
		want      string
	}{
		{
			name:      "timestamp only",
			createdAt: "2023-10-15T10:30:00Z",
			want:      "2023-10-15T10:30:00Z",
		},
		{
			name:      "different timestamp",
			createdAt: "2023-10-15T11:00:00Z",
			want:      "2023-10-15T11:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := workflowRunGSI2SK(tt.createdAt)
			if got != tt.want {
				t.Errorf("workflowRunGSI2SK(%s) = %s, want %s",
					tt.createdAt, got, tt.want)
			}
		})
	}
}

func TestStepExecutionPK(t *testing.T) {
	tests := []struct {
		name  string
		runID string
		want  string
	}{
		{
			name:  "simple run ID",
			runID: "test-run-1",
			want:  "RUN#test-run-1",
		},
		{
			name:  "UUID run ID",
			runID: "550e8400-e29b-41d4-a716-446655440000",
			want:  "RUN#550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stepExecutionPK(tt.runID)
			if got != tt.want {
				t.Errorf("stepExecutionPK(%s) = %s, want %s", tt.runID, got, tt.want)
			}
		})
	}
}

func TestStepExecutionSK(t *testing.T) {
	tests := []struct {
		name   string
		stepID string
		want   string
	}{
		{
			name:   "simple step ID",
			stepID: "send-email",
			want:   "STEP#send-email",
		},
		{
			name:   "numeric step ID",
			stepID: "step-1",
			want:   "STEP#step-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stepExecutionSK(tt.stepID)
			if got != tt.want {
				t.Errorf("stepExecutionSK(%s) = %s, want %s", tt.stepID, got, tt.want)
			}
		})
	}
}

func TestStepOutputPK(t *testing.T) {
	tests := []struct {
		name  string
		runID string
		want  string
	}{
		{
			name:  "simple run ID",
			runID: "test-run-1",
			want:  "RUN#test-run-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stepOutputPK(tt.runID)
			if got != tt.want {
				t.Errorf("stepOutputPK(%s) = %s, want %s", tt.runID, got, tt.want)
			}
		})
	}
}

func TestStepOutputSK(t *testing.T) {
	tests := []struct {
		name   string
		stepID string
		want   string
	}{
		{
			name:   "simple step ID",
			stepID: "send-email",
			want:   "OUTPUT#send-email",
		},
		{
			name:   "numeric step ID",
			stepID: "step-1",
			want:   "OUTPUT#step-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stepOutputSK(tt.stepID)
			if got != tt.want {
				t.Errorf("stepOutputSK(%s) = %s, want %s", tt.stepID, got, tt.want)
			}
		})
	}
}

func TestStatePK(t *testing.T) {
	tests := []struct {
		name  string
		runID string
		want  string
	}{
		{
			name:  "simple run ID",
			runID: "test-run-1",
			want:  "RUN#test-run-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statePK(tt.runID)
			if got != tt.want {
				t.Errorf("statePK(%s) = %s, want %s", tt.runID, got, tt.want)
			}
		})
	}
}

func TestStateSK(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "simple key",
			key:  "counter",
			want: "STATE#counter",
		},
		{
			name: "complex key",
			key:  "user.email.count",
			want: "STATE#user.email.count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stateSK(tt.key)
			if got != tt.want {
				t.Errorf("stateSK(%s) = %s, want %s", tt.key, got, tt.want)
			}
		})
	}
}

func TestStatePrefix(t *testing.T) {
	got := statePrefix()
	want := "STATE#"
	if got != want {
		t.Errorf("statePrefix() = %s, want %s", got, want)
	}
}

func TestStepPrefix(t *testing.T) {
	got := stepPrefix()
	want := "STEP#"
	if got != want {
		t.Errorf("stepPrefix() = %s, want %s", got, want)
	}
}

func TestKeyConsistency(t *testing.T) {
	// Test that PK functions are consistent across different entity types
	runID := "test-run-123"

	workflowPK := workflowRunPK(runID)
	stepPK := stepExecutionPK(runID)
	outputPK := stepOutputPK(runID)
	statePKVal := statePK(runID)

	// All should have the same PK for the same runID
	if workflowPK != stepPK {
		t.Errorf("Inconsistent PK: workflowRunPK=%s, stepExecutionPK=%s", workflowPK, stepPK)
	}
	if stepPK != outputPK {
		t.Errorf("Inconsistent PK: stepExecutionPK=%s, stepOutputPK=%s", stepPK, outputPK)
	}
	if outputPK != statePKVal {
		t.Errorf("Inconsistent PK: stepOutputPK=%s, statePK=%s", outputPK, statePKVal)
	}
}

func TestSKPrefixUniqueness(t *testing.T) {
	// Test that SK prefixes are unique to avoid collisions
	prefixes := map[string]string{
		"workflow": workflowRunSK(),
		"step":     stepPrefix(),
		"output":   "OUTPUT#",
		"state":    statePrefix(),
	}

	seen := make(map[string]string)
	for name, prefix := range prefixes {
		if existing, ok := seen[prefix]; ok {
			t.Errorf("Duplicate prefix %s used by %s and %s", prefix, name, existing)
		}
		seen[prefix] = name
	}
}

func TestStepExecutionSKPrefix(t *testing.T) {
	// Verify step execution SK has correct prefix
	stepID := "test-step"
	sk := stepExecutionSK(stepID)
	prefix := stepPrefix()

	if !strings.HasPrefix(sk, prefix) {
		t.Errorf("stepExecutionSK(%s) = %s does not start with prefix %s", stepID, sk, prefix)
	}
}

func TestStepOutputSKPrefix(t *testing.T) {
	// Verify step output SK has correct prefix
	stepID := "test-step"
	sk := stepOutputSK(stepID)
	prefix := "OUTPUT#"

	if !strings.HasPrefix(sk, prefix) {
		t.Errorf("stepOutputSK(%s) = %s does not start with prefix %s", stepID, sk, prefix)
	}
}

func TestStateSKPrefix(t *testing.T) {
	// Verify state SK has correct prefix
	key := "test-key"
	sk := stateSK(key)
	prefix := statePrefix()

	if !strings.HasPrefix(sk, prefix) {
		t.Errorf("stateSK(%s) = %s does not start with prefix %s", key, sk, prefix)
	}
}

func TestConstants(t *testing.T) {
	// Verify all constants are defined and non-empty
	constants := map[string]string{
		"AttrPK":                  AttrPK,
		"AttrSK":                  AttrSK,
		"AttrGSI1PK":              AttrGSI1PK,
		"AttrGSI1SK":              AttrGSI1SK,
		"AttrGSI2PK":              AttrGSI2PK,
		"AttrGSI2SK":              AttrGSI2SK,
		"AttrEntityType":          AttrEntityType,
		"AttrData":                AttrData,
		"AttrTTL":                 AttrTTL,
		"EntityTypeWorkflowRun":   EntityTypeWorkflowRun,
		"EntityTypeStepExecution": EntityTypeStepExecution,
		"EntityTypeStepOutput":    EntityTypeStepOutput,
		"EntityTypeState":         EntityTypeState,
		"IndexStatusIndex":        IndexStatusIndex,
		"IndexResourceIndex":      IndexResourceIndex,
	}

	for name, value := range constants {
		if value == "" {
			t.Errorf("Constant %s is empty", name)
		}
	}
}

func TestEntityTypeConstants(t *testing.T) {
	// Verify entity type constants are unique
	entityTypes := []string{
		EntityTypeWorkflowRun,
		EntityTypeStepExecution,
		EntityTypeStepOutput,
		EntityTypeState,
	}

	seen := make(map[string]bool)
	for _, et := range entityTypes {
		if seen[et] {
			t.Errorf("Duplicate entity type: %s", et)
		}
		seen[et] = true
	}
}

func TestIndexNameConstants(t *testing.T) {
	// Verify index name constants are unique
	indexes := []string{
		IndexStatusIndex,
		IndexResourceIndex,
	}

	seen := make(map[string]bool)
	for _, idx := range indexes {
		if seen[idx] {
			t.Errorf("Duplicate index name: %s", idx)
		}
		seen[idx] = true
	}
}

func TestGSI1KeyFormat(t *testing.T) {
	// Verify GSI1 keys follow expected format for queryability
	workflowID := "test-workflow"
	status := "RUNNING"
	createdAt := "2023-10-15T10:30:00Z"

	pk := workflowRunGSI1PK(workflowID, status)
	sk := workflowRunGSI1SK(createdAt)

	// PK should contain workflow ID and status
	if !strings.Contains(pk, workflowID) {
		t.Errorf("GSI1 PK %s does not contain workflow ID %s", pk, workflowID)
	}
	if !strings.Contains(pk, status) {
		t.Errorf("GSI1 PK %s does not contain status %s", pk, status)
	}

	// SK should be the timestamp for sorting
	if sk != createdAt {
		t.Errorf("GSI1 SK %s != createdAt %s", sk, createdAt)
	}
}

func TestGSI2KeyFormat(t *testing.T) {
	// Verify GSI2 keys follow expected format for queryability
	resourceID := "email-123"
	createdAt := "2023-10-15T10:30:00Z"
	status := "RUNNING"

	pk := workflowRunGSI2PK(resourceID, status)
	sk := workflowRunGSI2SK(createdAt)

	// PK should contain resource ID and status
	if !strings.Contains(pk, resourceID) {
		t.Errorf("GSI2 PK %s does not contain resource ID %s", pk, resourceID)
	}
	if !strings.Contains(pk, status) {
		t.Errorf("GSI2 PK %s does not contain status %s", pk, status)
	}

	// SK should be the timestamp for sorting
	if sk != createdAt {
		t.Errorf("GSI2 SK %s != createdAt %s", sk, createdAt)
	}
}
