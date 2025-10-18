package domain

// CompletionRequest represents the input payload for a language model
// text generation request.
//
// The Prompt field contains the text instruction or query that the model
// should respond to. Implementations may enrich this with additional
// metadata (e.g. temperature, max tokens, or conversation context) in
// future extensions.
//
// Example:
//
//	req := CompletionRequest{
//	    Prompt: "Summarize the history of the Raspberry Pi in two sentences.",
//	}
type CompletionRequest struct {
	// Prompt is the text input provided to the language model.
	// It should clearly describe the desired output or question.
	Prompt string
}

// CompletionResult represents the output returned by a language model
// after processing a completion request.
//
// The Text field contains the generated response content, typically as
// plain text, but may include formatted or structured output depending
// on the model and implementation.
type CompletionResult struct {
	// Text is the generated text output from the completion model.
	Text string `json:"text"`
}
