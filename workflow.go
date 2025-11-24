package gorkflow

import (
	"fmt"
	"time"
)

// Workflow represents the complete workflow blueprint
type Workflow struct {
	id          string
	name        string
	description string
	version     string

	// Steps registered by ID
	steps map[string]StepExecutor

	// Execution graph
	graph *ExecutionGraph

	// Default config
	config ExecutionConfig

	// Metadata
	tags      map[string]string
	createdAt time.Time
}

// ID returns the workflow ID
func (w *Workflow) ID() string {
	return w.id
}

// Name returns the workflow name
func (w *Workflow) Name() string {
	return w.name
}

// Description returns the workflow description
func (w *Workflow) Description() string {
	return w.description
}

// Version returns the workflow version
func (w *Workflow) Version() string {
	return w.version
}

// Graph returns the execution graph
func (w *Workflow) Graph() *ExecutionGraph {
	return w.graph
}

// GetStep retrieves a step by ID
func (w *Workflow) GetStep(stepID string) (StepExecutor, error) {
	step, exists := w.steps[stepID]
	if !exists {
		return nil, fmt.Errorf("step %s not found in workflow", stepID)
	}
	return step, nil
}

// GetAllSteps returns all steps
func (w *Workflow) GetAllSteps() map[string]StepExecutor {
	return w.steps
}

// GetConfig returns the default execution config
func (w *Workflow) GetConfig() ExecutionConfig {
	return w.config
}

// NewWorkflowInstance creates a new workflow instance
func NewWorkflowInstance(id, name string) *Workflow {
	return &Workflow{
		id:        id,
		name:      name,
		version:   "1.0",
		steps:     make(map[string]StepExecutor),
		graph:     NewExecutionGraph(),
		config:    DefaultExecutionConfig,
		tags:      make(map[string]string),
		createdAt: time.Now(),
	}
}

// SetDescription sets the workflow description
func (w *Workflow) SetDescription(description string) {
	w.description = description
}

// SetVersion sets the workflow version
func (w *Workflow) SetVersion(version string) {
	w.version = version
}

// SetConfig sets the execution config
func (w *Workflow) SetConfig(config ExecutionConfig) {
	w.config = config
}

// SetTags sets the workflow tags
func (w *Workflow) SetTags(tags map[string]string) {
	w.tags = tags
}

// AddStep registers a step in the workflow
func (w *Workflow) AddStep(step StepExecutor) {
	w.steps[step.GetID()] = step
}
