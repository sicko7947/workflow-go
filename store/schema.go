package store

import "fmt"

// DynamoDB schema constants for single-table design
const (
	// Table attributes
	AttrPK         = "PK"
	AttrSK         = "SK"
	AttrGSI1PK     = "GSI1PK"
	AttrGSI1SK     = "GSI1SK"
	AttrGSI2PK     = "GSI2PK"
	AttrGSI2SK     = "GSI2SK"
	AttrEntityType = "entity_type"
	AttrData       = "data"
	AttrTTL        = "ttl"

	// Entity types
	EntityTypeWorkflowRun   = "WorkflowRun"
	EntityTypeStepExecution = "StepExecution"
	EntityTypeStepOutput    = "StepOutput"
	EntityTypeState         = "State"

	// Index names
	IndexStatusIndex   = "GSI1"
	IndexResourceIndex = "GSI2"
)

// Key builders for single-table design

// WorkflowRun keys: PK=RUN#{runID}, SK=META
func workflowRunPK(runID string) string {
	return fmt.Sprintf("RUN#%s", runID)
}

func workflowRunSK() string {
	return "META"
}

func workflowRunGSI1PK(workflowID, status string) string {
	return fmt.Sprintf("WF#%s#STATUS#%s", workflowID, status)
}

func workflowRunGSI1SK(createdAt string) string {
	return createdAt
}

func workflowRunGSI2PK(resourceID, status string) string {
	return fmt.Sprintf("RES#%s#STATUS#%s", resourceID, status)
}

func workflowRunGSI2SK(createdAt string) string {
	return createdAt
}

// StepExecution keys: PK=RUN#{runID}, SK=STEP#{stepID}
func stepExecutionPK(runID string) string {
	return fmt.Sprintf("RUN#%s", runID)
}

func stepExecutionSK(stepID string) string {
	return fmt.Sprintf("STEP#%s", stepID)
}

// StepOutput keys: PK=RUN#{runID}, SK=OUTPUT#{stepID}
func stepOutputPK(runID string) string {
	return fmt.Sprintf("RUN#%s", runID)
}

func stepOutputSK(stepID string) string {
	return fmt.Sprintf("OUTPUT#%s", stepID)
}

// State keys: PK=RUN#{runID}, SK=STATE#{key}
func statePK(runID string) string {
	return fmt.Sprintf("RUN#%s", runID)
}

func stateSK(key string) string {
	return fmt.Sprintf("STATE#%s", key)
}

// Prefix for range queries
func statePrefix() string {
	return "STATE#"
}

func stepPrefix() string {
	return "STEP#"
}
