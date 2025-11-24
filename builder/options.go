package builder

import "github.com/sicko7947/gorkflow"

// WorkflowOption is a functional option for configuring workflows
type WorkflowOption func(*gorkflow.Workflow)

// WithDescription sets the workflow description
func WithDescription(description string) WorkflowOption {
	return func(w *gorkflow.Workflow) {
		w.SetDescription(description)
	}
}

// WithVersion sets the workflow version
func WithVersion(version string) WorkflowOption {
	return func(w *gorkflow.Workflow) {
		w.SetVersion(version)
	}
}

// WithDefaultConfig sets the default execution config
func WithDefaultConfig(config gorkflow.ExecutionConfig) WorkflowOption {
	return func(w *gorkflow.Workflow) {
		w.SetConfig(config)
	}
}

// WithTags sets workflow tags
func WithTags(tags map[string]string) WorkflowOption {
	return func(w *gorkflow.Workflow) {
		w.SetTags(tags)
	}
}

// ApplyOptions applies a list of options to a workflow
func ApplyOptions(w *gorkflow.Workflow, opts ...WorkflowOption) {
	for _, opt := range opts {
		opt(w)
	}
}
