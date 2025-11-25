package gorkflow_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/sicko7947/gorkflow"
	"github.com/sicko7947/gorkflow/engine"
	"github.com/sicko7947/gorkflow/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AppContext struct {
	UserID string
	Config map[string]string
}

func TestWorkflowWithCustomContext(t *testing.T) {
	// Define a step that uses the custom context
	stepHandler := func(ctx *gorkflow.StepContext, input string) (string, error) {
		appCtx, err := gorkflow.GetContext[*AppContext](ctx)
		if err != nil {
			return "", err
		}

		if appCtx.UserID != "user-123" {
			return "", fmt.Errorf("unexpected user ID: %s", appCtx.UserID)
		}

		val, ok := appCtx.Config["env"]
		if !ok || val != "production" {
			return "", fmt.Errorf("unexpected config value: %v", appCtx.Config)
		}

		return fmt.Sprintf("Processed for %s in %s", appCtx.UserID, val), nil
	}

	// Create workflow with custom context
	appCtx := &AppContext{
		UserID: "user-123",
		Config: map[string]string{"env": "production"},
	}

	wf := gorkflow.NewWorkflowInstance("test-wf", "Test Workflow", gorkflow.WithContext(appCtx))
	step := gorkflow.NewStep("step-1", "Step 1", stepHandler)
	wf.AddStep(step)
	wf.Graph().AddNode(step.GetID(), gorkflow.NodeTypeSequential)

	// Setup engine
	store := store.NewMemoryStore()
	eng := engine.NewEngine(store)

	// Run workflow
	runID, err := eng.StartWorkflow(context.Background(), wf, "start", gorkflow.WithSynchronousExecution())
	require.NoError(t, err)

	// Verify result
	run, err := eng.GetRun(context.Background(), runID)
	require.NoError(t, err)
	assert.Equal(t, gorkflow.RunStatusCompleted, run.Status)

	// Verify step output
	executions, err := eng.GetStepExecutions(context.Background(), runID)
	require.NoError(t, err)
	require.Len(t, executions, 1)
	assert.Equal(t, "\"Processed for user-123 in production\"", string(executions[0].Output))
}
