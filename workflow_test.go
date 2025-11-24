package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflow_GetAllSteps(t *testing.T) {
	step1 := NewStep("step1", "Step 1", testHandler)
	step2 := NewStep("step2", "Step 2", testHandler)
	step3 := NewStep("step3", "Step 3", testHandler)

	wf := NewWorkflowInstance("test-workflow", "Test Workflow")
	wf.AddStep(step1)
	wf.AddStep(step2)
	wf.AddStep(step3)

	wf.Graph().AddNode("step1", NodeTypeSequential)
	wf.Graph().AddNode("step2", NodeTypeSequential)
	wf.Graph().AddNode("step3", NodeTypeSequential)

	wf.Graph().AddEdge("step1", "step2")
	wf.Graph().AddEdge("step2", "step3")
	wf.Graph().SetEntryPoint("step1")

	allSteps := wf.GetAllSteps()
	assert.Len(t, allSteps, 3)

	// Verify all steps are present
	stepIDs := make([]string, 0, len(allSteps))
	for _, step := range allSteps {
		stepIDs = append(stepIDs, step.GetID())
	}
	assert.Contains(t, stepIDs, "step1")
	assert.Contains(t, stepIDs, "step2")
	assert.Contains(t, stepIDs, "step3")
}

func TestWorkflow_ComplexGraph(t *testing.T) {
	// Create a diamond-shaped workflow:
	//     step1
	//     /   \
	//  step2  step3
	//     \   /
	//     step4

	step1 := NewStep("step1", "Step 1", testHandler)
	step2 := NewStep("step2", "Step 2", testHandler)
	step3 := NewStep("step3", "Step 3", testHandler)
	step4 := NewStep("step4", "Step 4", testHandler)

	wf := NewWorkflowInstance("test-workflow", "Test Workflow")
	wf.AddStep(step1)
	wf.AddStep(step2)
	wf.AddStep(step3)
	wf.AddStep(step4)

	wf.Graph().AddNode("step1", NodeTypeSequential)
	wf.Graph().AddNode("step2", NodeTypeSequential)
	wf.Graph().AddNode("step3", NodeTypeSequential)
	wf.Graph().AddNode("step4", NodeTypeSequential)

	wf.Graph().AddEdge("step1", "step2")
	wf.Graph().AddEdge("step1", "step3")
	wf.Graph().AddEdge("step2", "step4")
	wf.Graph().AddEdge("step3", "step4")
	wf.Graph().SetEntryPoint("step1")

	// Verify graph structure
	graph := wf.Graph()
	assert.Equal(t, "step1", graph.EntryPoint)

	// step1 should have two next steps
	nextFromStep1, err := graph.GetNextSteps("step1")
	require.NoError(t, err)
	assert.Len(t, nextFromStep1, 2)
	assert.Contains(t, nextFromStep1, "step2")
	assert.Contains(t, nextFromStep1, "step3")

	// step2 and step3 should both point to step4
	nextFromStep2, err := graph.GetNextSteps("step2")
	require.NoError(t, err)
	assert.Contains(t, nextFromStep2, "step4")

	nextFromStep3, err := graph.GetNextSteps("step3")
	require.NoError(t, err)
	assert.Contains(t, nextFromStep3, "step4")
}

func TestWorkflow_Accessors(t *testing.T) {
	step1 := NewStep("step1", "Step 1", testHandler)

	wf := NewWorkflowInstance("test-workflow", "Test Workflow")
	wf.SetDescription("Test Description")
	wf.SetVersion("2.0.0")
	wf.AddStep(step1)
	wf.Graph().AddNode("step1", NodeTypeSequential)
	wf.Graph().SetEntryPoint("step1")

	assert.Equal(t, "test-workflow", wf.ID())
	assert.Equal(t, "Test Workflow", wf.Name())
	assert.Equal(t, "Test Description", wf.Description())
	assert.Equal(t, "2.0.0", wf.Version())
	assert.NotNil(t, wf.Graph())
	assert.NotNil(t, wf.GetConfig())
}
