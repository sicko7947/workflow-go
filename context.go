package gorkflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

// StepContext provides rich context to step handlers
type StepContext struct {
	context.Context

	// Execution metadata
	RunID   string
	StepID  string
	Attempt int

	// Logger (enriched with step context)
	Logger zerolog.Logger

	// Access to other step outputs (type-safe)
	Outputs StepOutputAccessor

	// Access to workflow-level state
	State StateAccessor

	// Custom context (user-defined)
	CustomContext any
}

// GetContext retrieves the custom context from the step context
func GetContext[T any](ctx *StepContext) (T, error) {
	var zero T
	if ctx.CustomContext == nil {
		return zero, fmt.Errorf("custom context is nil")
	}

	val, ok := ctx.CustomContext.(T)
	if !ok {
		return zero, fmt.Errorf("custom context is not of type %T", zero)
	}
	return val, nil
}

// StepOutputAccessor provides type-safe access to other step outputs
type StepOutputAccessor interface {
	// GetOutput retrieves output from a specific step
	GetOutput(stepID string, target interface{}) error

	// HasOutput checks if a step has produced output
	HasOutput(stepID string) bool
}

// GetTypedOutput is a generic function for type-safe output retrieval
func GetTypedOutput[T any](accessor StepOutputAccessor, stepID string) (T, error) {
	var result T
	err := accessor.GetOutput(stepID, &result)
	return result, err
}

// StateAccessor provides type-safe access to workflow state
type StateAccessor interface {
	// Set stores a value in the workflow state
	Set(key string, value interface{}) error

	// Get retrieves a value from the workflow state
	Get(key string, target interface{}) error

	// Delete removes a key from the state
	Delete(key string) error

	// Has checks if a key exists
	Has(key string) bool

	// GetAll retrieves all state data
	GetAll() (map[string][]byte, error)
}

// SetTyped is a generic function for type-safe state setting
func SetTyped[T any](accessor StateAccessor, key string, value T) error {
	return accessor.Set(key, value)
}

// GetTyped is a generic function for type-safe state retrieval
func GetTyped[T any](accessor StateAccessor, key string) (T, error) {
	var result T
	err := accessor.Get(key, &result)
	return result, err
}

// stepOutputAccessor implements StepOutputAccessor
type stepOutputAccessor struct {
	runID string
	store WorkflowStore
	cache map[string][]byte
}

// newStepOutputAccessor creates a new output accessor
func newStepOutputAccessor(runID string, wfStore WorkflowStore) StepOutputAccessor {
	return &stepOutputAccessor{
		runID: runID,
		store: wfStore,
		cache: make(map[string][]byte),
	}
}

// NewStepOutputAccessor creates a new output accessor (exported)
func NewStepOutputAccessor(runID string, wfStore WorkflowStore) StepOutputAccessor {
	return newStepOutputAccessor(runID, wfStore)
}

func (a *stepOutputAccessor) GetOutput(stepID string, target interface{}) error {
	// Check cache first
	if data, ok := a.cache[stepID]; ok {
		return json.Unmarshal(data, target)
	}

	// Load from store
	data, err := a.store.LoadStepOutput(context.Background(), a.runID, stepID)
	if err != nil {
		return fmt.Errorf("failed to load output for step %s: %w", stepID, err)
	}

	// Cache it
	a.cache[stepID] = data

	// Unmarshal
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal output for step %s: %w", stepID, err)
	}

	return nil
}

func (a *stepOutputAccessor) HasOutput(stepID string) bool {
	// Check cache
	if _, ok := a.cache[stepID]; ok {
		return true
	}

	// Check store
	_, err := a.store.LoadStepOutput(context.Background(), a.runID, stepID)
	return err == nil
}

// stateAccessor implements StateAccessor
type stateAccessor struct {
	runID string
	store WorkflowStore
	cache map[string][]byte
}

// newStateAccessor creates a new state accessor
func newStateAccessor(runID string, wfStore WorkflowStore) StateAccessor {
	return &stateAccessor{
		runID: runID,
		store: wfStore,
		cache: make(map[string][]byte),
	}
}

// NewStateAccessor creates a new state accessor (exported)
func NewStateAccessor(runID string, wfStore WorkflowStore) StateAccessor {
	return newStateAccessor(runID, wfStore)
}

func (a *stateAccessor) Set(key string, value interface{}) error {
	// Marshal value
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal state value for key %s: %w", key, err)
	}

	// Update cache
	a.cache[key] = data

	// Persist to store
	if err := a.store.SaveState(context.Background(), a.runID, key, data); err != nil {
		return fmt.Errorf("failed to save state for key %s: %w", key, err)
	}

	return nil
}

func (a *stateAccessor) Get(key string, target interface{}) error {
	// Check cache first
	if data, ok := a.cache[key]; ok {
		return json.Unmarshal(data, target)
	}

	// Load from store
	data, err := a.store.LoadState(context.Background(), a.runID, key)
	if err != nil {
		return fmt.Errorf("failed to load state for key %s: %w", key, err)
	}

	// Cache it
	a.cache[key] = data

	// Unmarshal
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal state for key %s: %w", key, err)
	}

	return nil
}

func (a *stateAccessor) Delete(key string) error {
	// Remove from cache
	delete(a.cache, key)

	// Delete from store
	if err := a.store.DeleteState(context.Background(), a.runID, key); err != nil {
		return fmt.Errorf("failed to delete state for key %s: %w", key, err)
	}

	return nil
}

func (a *stateAccessor) Has(key string) bool {
	// Check cache
	if _, ok := a.cache[key]; ok {
		return true
	}

	// Check store
	_, err := a.store.LoadState(context.Background(), a.runID, key)
	return err == nil
}

func (a *stateAccessor) GetAll() (map[string][]byte, error) {
	// Get all from store
	data, err := a.store.GetAllState(context.Background(), a.runID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all state: %w", err)
	}

	// Update cache
	for k, v := range data {
		a.cache[k] = v
	}

	return data, nil
}
