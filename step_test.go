package gorkflow

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test types
type TestInput struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

type TestOutput struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
}

// Test handler that doubles the input
func testHandler(ctx *StepContext, input TestInput) (TestOutput, error) {
	return TestOutput{
		Result:  input.Value * 2,
		Message: "Processed " + input.Name,
	}, nil
}

// Test handler that returns an error
func errorHandler(ctx *StepContext, input TestInput) (TestOutput, error) {
	return TestOutput{}, assert.AnError
}

func TestNewStep(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	assert.Equal(t, "test-step", step.GetID())
	assert.Equal(t, "Test Step", step.GetName())
	assert.NotNil(t, step.Handler)
	assert.Equal(t, DefaultExecutionConfig, step.GetConfig())
}

func TestNewStep_WithOptions(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler,
		WithRetries(5),
		WithTimeout(30*time.Second),
		WithBackoff(BackoffExponential),
		WithContinueOnError(true),
	)

	config := step.GetConfig()
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 30, config.TimeoutSeconds)
	assert.Equal(t, BackoffExponential, config.RetryBackoff)
	assert.True(t, config.ContinueOnError)
}

func TestStep_Execute_Success(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	// Create input
	input := TestInput{Value: 21, Name: "test"}
	inputBytes, err := json.Marshal(input)
	require.NoError(t, err)

	// Create step context
	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "test-step",
		Attempt: 0,
		Logger:  zerolog.Nop(),
	}

	// Execute
	outputBytes, err := step.Execute(ctx, inputBytes)
	require.NoError(t, err)

	// Verify output
	var output TestOutput
	err = json.Unmarshal(outputBytes, &output)
	require.NoError(t, err)

	assert.Equal(t, 42, output.Result)
	assert.Equal(t, "Processed test", output.Message)
}

func TestStep_Execute_InvalidInput(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	// Invalid JSON
	invalidInput := []byte("not valid json")

	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "test-step",
		Logger:  zerolog.Nop(),
	}

	_, err := step.Execute(ctx, invalidInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal input")
}

func TestStep_Execute_HandlerError(t *testing.T) {
	step := NewStep("error-step", "Error Step", errorHandler)

	input := TestInput{Value: 10, Name: "test"}
	inputBytes, _ := json.Marshal(input)

	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "error-step",
		Logger:  zerolog.Nop(),
	}

	_, err := step.Execute(ctx, inputBytes)
	assert.Error(t, err)
}

func TestStep_ValidateInput(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	// Valid input
	validInput := TestInput{Value: 10, Name: "test"}
	validBytes, _ := json.Marshal(validInput)
	err := step.ValidateInput(validBytes)
	assert.NoError(t, err)

	// Invalid input
	invalidBytes := []byte("not json")
	err = step.ValidateInput(invalidBytes)
	assert.Error(t, err)
}

func TestStep_ValidateOutput(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	// Valid output
	validOutput := TestOutput{Result: 42, Message: "test"}
	validBytes, _ := json.Marshal(validOutput)
	err := step.ValidateOutput(validBytes)
	assert.NoError(t, err)

	// Invalid output
	invalidBytes := []byte("not json")
	err = step.ValidateOutput(invalidBytes)
	assert.Error(t, err)
}

func TestStep_TypeInformation(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	inputType := step.InputType()
	outputType := step.OutputType()

	assert.NotNil(t, inputType)
	assert.NotNil(t, outputType)
	assert.Equal(t, "TestInput", inputType.Name())
	assert.Equal(t, "TestOutput", outputType.Name())
}

func TestConditionalStep_Execute_ConditionTrue(t *testing.T) {
	baseStep := NewStep("test-step", "Test Step", testHandler)

	condition := func(ctx *StepContext) (bool, error) {
		return true, nil
	}

	condStep := NewConditionalStep(baseStep, condition, nil)

	input := TestInput{Value: 21, Name: "test"}
	inputBytes, _ := json.Marshal(input)

	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "test-step",
		Logger:  zerolog.Nop(),
	}

	outputBytes, err := condStep.Execute(ctx, inputBytes)
	require.NoError(t, err)

	var output TestOutput
	json.Unmarshal(outputBytes, &output)
	assert.Equal(t, 42, output.Result)
}

func TestConditionalStep_Execute_ConditionFalse(t *testing.T) {
	baseStep := NewStep("test-step", "Test Step", testHandler)

	condition := func(ctx *StepContext) (bool, error) {
		return false, nil
	}

	defaultOutput := &TestOutput{Result: 0, Message: "skipped"}
	condStep := NewConditionalStep(baseStep, condition, defaultOutput)

	input := TestInput{Value: 21, Name: "test"}
	inputBytes, _ := json.Marshal(input)

	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "test-step",
		Logger:  zerolog.Nop(),
	}

	outputBytes, err := condStep.Execute(ctx, inputBytes)
	require.NoError(t, err)

	var output TestOutput
	json.Unmarshal(outputBytes, &output)
	assert.Equal(t, 0, output.Result)
	assert.Equal(t, "skipped", output.Message)
}

func TestConditionalStep_Execute_ConditionError(t *testing.T) {
	baseStep := NewStep("test-step", "Test Step", testHandler)

	condition := func(ctx *StepContext) (bool, error) {
		return false, assert.AnError
	}

	condStep := NewConditionalStep(baseStep, condition, nil)

	input := TestInput{Value: 21, Name: "test"}
	inputBytes, _ := json.Marshal(input)

	ctx := &StepContext{
		Context: context.Background(),
		RunID:   "test-run",
		StepID:  "test-step",
		Logger:  zerolog.Nop(),
	}

	_, err := condStep.Execute(ctx, inputBytes)
	assert.Error(t, err)
}

func TestStepOption_Chaining(t *testing.T) {
	step := NewStep("test-step", "Test Step", testHandler)

	// Apply options sequentially
	step.SetMaxRetries(3)
	step.SetTimeout(10) // SetTimeout takes int seconds
	step.SetBackoff(BackoffExponential)
	step.SetRetryDelay(500) // SetRetryDelay takes int ms
	step.SetContinueOnError(true)

	config := step.GetConfig()
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 10, config.TimeoutSeconds)
	assert.Equal(t, BackoffExponential, config.RetryBackoff)
	assert.Equal(t, 500, config.RetryDelayMs)
	assert.True(t, config.ContinueOnError)
}
