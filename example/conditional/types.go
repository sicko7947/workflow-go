package conditional

// ConditionalInput is input for conditional workflow example
type ConditionalInput struct {
	Value            int  `json:"value"`
	EnableDoubling   bool `json:"enable_doubling"`
	EnableFormatting bool `json:"enable_formatting"`
}

// DoubleInput for the doubling step
type DoubleInput struct {
	Value int `json:"value"`
}

// DoubleOutput from the doubling step
type DoubleOutput struct {
	Value   int    `json:"value"`
	Doubled bool   `json:"doubled"`
	Message string `json:"message,omitempty"`
}

// ConditionalFormatInput for the formatting step
type ConditionalFormatInput struct {
	Value   int    `json:"value"`
	Doubled bool   `json:"doubled"`
	Message string `json:"message,omitempty"`
}

// ConditionalFormatOutput final output
type ConditionalFormatOutput struct {
	Value     int    `json:"value"`
	Formatted string `json:"formatted"`
	Doubled   bool   `json:"doubled"`
}
